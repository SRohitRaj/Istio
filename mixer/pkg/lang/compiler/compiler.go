// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package compiler implements a compiler that converts Mixer's expression language into a
// Mixer IL-based program that can be executed via an interpreter.
package compiler

import (
	"fmt"
	"time"

	descriptor "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/il"
	"istio.io/istio/mixer/pkg/lang/ast"
	"istio.io/istio/pkg/log"
)

// Compiler is a stateful compiler that can be used to gradually build an il.Program out of multiple independent
// compilation of expressions.
type Compiler struct {
	program   *il.Program
	finder    ast.AttributeDescriptorFinder
	functions map[string]ast.FunctionMetadata

	nextFnID int
}

// New returns a new compiler instance.
func New(finder ast.AttributeDescriptorFinder, functions map[string]ast.FunctionMetadata) *Compiler {
	return &Compiler{
		finder:    finder,
		program:   il.NewProgram(),
		functions: functions,
		nextFnID:  0,
	}
}

// CompileExpression creates a new parameterless IL function, using the given expression text as its body. Upon success,
// it returns the id of the generated function.
func (c *Compiler) CompileExpression(text string) (uint32, descriptor.ValueType, error) {
	expression, err := ast.Parse(text)
	if err != nil {
		return 0, descriptor.VALUE_TYPE_UNSPECIFIED, err
	}

	exprType, err := expression.EvalType(c.finder, c.functions)
	if err != nil {
		return 0, descriptor.VALUE_TYPE_UNSPECIFIED, err
	}

	g := generator{
		program:   c.program,
		builder:   il.NewBuilder(c.program.Strings()),
		finder:    c.finder,
		functions: c.functions,
	}

	returnType := g.toIlType(exprType)
	g.generate(expression, 0, celMode, "")
	if g.err != nil {
		log.Warnf("compiler.Compile failed. expr:'%s', err:'%v'", text, g.err)
		return 0, descriptor.VALUE_TYPE_UNSPECIFIED, g.err
	}
	g.builder.Ret()

	body := g.builder.Build()

	name := fmt.Sprintf("$expression%d", c.nextFnID)
	c.nextFnID++

	if err = g.program.AddFunction(name, []il.Type{}, returnType, body); err != nil {
		return 0, descriptor.VALUE_TYPE_UNSPECIFIED, err
	}

	return g.program.Functions.IDOf(name), exprType, nil
}

// Program returns the program instance that is being built by this compiler.
func (c *Compiler) Program() *il.Program {
	return c.program
}

type generator struct {
	program   *il.Program
	builder   *il.Builder
	finder    ast.AttributeDescriptorFinder
	functions map[string]ast.FunctionMetadata
	err       error
}

// nilMode is an enum flag for specifying how the emitted code should be handling potential nils.
type nilMode int

const (
	// nmNone indicates that the emitted code should assume the expression will not yield any nils.
	// If a nil is encountered, then it will be treated as an error. The result of the expression
	// will be on top of the stack.
	nmNone nilMode = iota

	// nmJmpOnValue indicates that the emitted code should immediately jump to the label (that is
	// supplied on the side) when it resolves the expression to a non-nil value. If the value is nil,
	// the emitted code will "fallthrough" and allow the immediately-following piece of code to handle
	// the nil  case.
	nmJmpOnValue

	// CEL mode relies on CEL semantics to interpret expressions
	celMode
)

// Result is returned as the result of compilation.
type Result struct {
	Program    *il.Program
	Expression *ast.Expression
}

// Compile converts the given expression text, into an IL based program.
func Compile(text string, finder ast.AttributeDescriptorFinder, functions map[string]ast.FunctionMetadata) (*il.Program, error) {
	// TODO: This function should either be eliminated entirely, or use the Compiler struct, once we switch over to
	// to using compiled expressions. Keeping this here in its current form to avoid generation of excessive garbage
	// in the request path.

	p := il.NewProgram()

	expression, err := ast.Parse(text)
	if err != nil {
		return nil, err
	}

	exprType, err := expression.EvalType(finder, functions)
	if err != nil {
		return nil, err
	}

	g := generator{
		program:   p,
		builder:   il.NewBuilder(p.Strings()),
		finder:    finder,
		functions: functions,
	}

	returnType := g.toIlType(exprType)
	g.generate(expression, 0, celMode, "")
	if g.err != nil {
		log.Warnf("compiler.Compile failed. expr:'%s', err:'%v'", text, g.err)
		return nil, g.err
	}

	g.builder.Ret()
	body := g.builder.Build()
	if err = g.program.AddFunction("eval", []il.Type{}, returnType, body); err != nil {
		g.internalError(err.Error())
		return nil, err
	}

	return p, nil
}

func (g *generator) toIlType(t descriptor.ValueType) il.Type {
	switch t {
	case descriptor.STRING:
		return il.String
	case descriptor.BOOL:
		return il.Bool
	case descriptor.INT64:
		return il.Integer
	case descriptor.DURATION:
		return il.Duration
	case descriptor.DOUBLE:
		return il.Double
	case descriptor.STRING_MAP:
		return il.Interface
	case descriptor.IP_ADDRESS:
		return il.Interface
	case descriptor.EMAIL_ADDRESS:
		return il.String
	case descriptor.DNS_NAME:
		return il.String
	case descriptor.URI:
		return il.String
	case descriptor.TIMESTAMP:
		return il.Interface
	default:
		g.internalError("unhandled expression type: '%v'", t)
		return il.Unknown
	}
}

func (g *generator) evalType(e *ast.Expression) il.Type {
	dvt, _ := e.EvalType(g.finder, g.functions)
	return g.toIlType(dvt)
}

func (g *generator) generate(e *ast.Expression, depth int, mode nilMode, valueJmpLabel string) {
	switch {
	case e.Const != nil:
		g.generateConstant(e.Const, mode, valueJmpLabel)
	case e.Var != nil:
		g.generateVariable(e.Var, mode, valueJmpLabel)
	case e.Fn != nil:
		g.generateFunction(e.Fn, depth, mode, valueJmpLabel)
	default:
		g.internalError("unexpected expression type encountered.")
	}
}

func (g *generator) generateVariable(v *ast.Variable, mode nilMode, valueJmpLabel string) {
	i := g.finder.GetAttribute(v.Name)
	ilType := g.toIlType(i.ValueType)
	switch ilType {
	case il.Integer, il.Duration:
		switch mode {
		case nmNone:
			g.builder.ResolveInt(v.Name)
		case nmJmpOnValue:
			g.builder.TResolveInt(v.Name)
			g.builder.Jnz(valueJmpLabel)
		case celMode:
			g.builder.TResolveInt(v.Name)
			end := g.builder.AllocateLabel()
			g.builder.Jnz(end)
			g.builder.APushInt(0)
			g.builder.SetLabelPos(end)
		}

	case il.String:
		switch mode {
		case nmNone:
			g.builder.ResolveString(v.Name)
		case nmJmpOnValue:
			g.builder.TResolveString(v.Name)
			g.builder.Jnz(valueJmpLabel)
		case celMode:
			g.builder.TResolveString(v.Name)
			end := g.builder.AllocateLabel()
			g.builder.Jnz(end)
			g.builder.APushStr("")
			g.builder.SetLabelPos(end)
		}

	case il.Bool:
		switch mode {
		case nmNone:
			g.builder.ResolveBool(v.Name)
		case nmJmpOnValue:
			g.builder.TResolveBool(v.Name)
			g.builder.Jnz(valueJmpLabel)
		case celMode:
			g.builder.TResolveBool(v.Name)
			end := g.builder.AllocateLabel()
			g.builder.Jnz(end)
			g.builder.APushBool(false)
			g.builder.SetLabelPos(end)
		}

	case il.Double:
		switch mode {
		case nmNone:
			g.builder.ResolveDouble(v.Name)
		case nmJmpOnValue:
			g.builder.TResolveDouble(v.Name)
			g.builder.Jnz(valueJmpLabel)
		case celMode:
			g.builder.TResolveDouble(v.Name)
			end := g.builder.AllocateLabel()
			g.builder.Jnz(end)
			g.builder.APushDouble(0)
			g.builder.SetLabelPos(end)
		}

	case il.Interface:
		switch mode {
		case nmNone:
			g.builder.ResolveInterface(v.Name)
		case nmJmpOnValue:
			g.builder.TResolveInterface(v.Name)
			g.builder.Jnz(valueJmpLabel)
		case celMode:
			g.builder.TResolveInterface(v.Name)
			end := g.builder.AllocateLabel()
			g.builder.Jnz(end)
			// null interface value
			g.builder.APushInt(0)
			g.builder.SetLabelPos(end)
		}
	default:
		g.internalError("unrecognized variable type: '%v'", i.ValueType)
	}
}

func (g *generator) generateFunction(f *ast.Function, depth int, mode nilMode, valueJmpLabel string) {
	switch f.Name {
	case "EQ":
		g.generateEq(f, depth)
	case "NEQ":
		g.generateNeq(f, depth)
	case "LOR":
		g.generateLor(f, depth)
	case "LAND":
		g.generateLand(f, depth)
	case "INDEX":
		g.generateIndex(f, depth, mode, valueJmpLabel)
	case "OR":
		g.generateOr(f, depth, mode, valueJmpLabel)
	case "conditional":
		g.generateConditional(f, depth, mode, valueJmpLabel)
	case "ADD":
		g.generateAdd(f, depth)
	case "size":
		g.generateSize(f, depth)
	default:
		// The parameters to a function (and the function itself) is expected to exist, regardless of whether
		// we're in a nillable context. The call will either succeed or error out.
		if f.Target != nil {
			g.generate(f.Target, depth, celMode, "")
		}
		for _, arg := range f.Args {
			g.generate(arg, depth, celMode, "")
		}
		g.builder.Call(f.Name)
		// If we're in a nillable context, then simply short-circuit to the end. The function is either
		// guaranteed to succeed or error out.
		if mode == nmJmpOnValue {
			g.builder.Jmp(valueJmpLabel)
		}
	}
}

func (g *generator) generateEq(f *ast.Function, depth int) {
	exprType := g.evalType(f.Args[0])
	g.generate(f.Args[0], depth+1, celMode, "")

	var constArg1 interface{}
	if f.Args[1].Const != nil {
		constArg1 = f.Args[1].Const.Value
	} else {
		g.generate(f.Args[1], depth+1, celMode, "")
	}

	switch exprType {
	case il.Bool:
		if constArg1 != nil {
			g.builder.AEQBool(constArg1.(bool))
		} else {
			g.builder.EQBool()
		}

	case il.String:
		dvt, _ := f.Args[0].EvalType(g.finder, g.functions)
		// TODO: this does not handle constArg1 for extern calls -- not a problem since IL does not produce
		// constant DNS, email, or URI values
		switch dvt {
		case descriptor.DNS_NAME:
			g.builder.Call("dnsName_equal")
		case descriptor.EMAIL_ADDRESS:
			g.builder.Call("email_equal")
		case descriptor.URI:
			g.builder.Call("uri_equal")
		default:
			if constArg1 != nil {
				g.builder.AEQString(constArg1.(string))
			} else {
				g.builder.EQString()
			}
		}

	case il.Integer:
		if constArg1 != nil {
			g.builder.AEQInteger(constArg1.(int64))
		} else {
			g.builder.EQInteger()
		}

	case il.Double:
		if constArg1 != nil {
			g.builder.AEQDouble(constArg1.(float64))
		} else {
			g.builder.EQDouble()
		}

	case il.Interface:
		dvt, _ := f.Args[0].EvalType(g.finder, g.functions)
		switch dvt {
		case descriptor.IP_ADDRESS:
			g.builder.Call("ip_equal")
		case descriptor.TIMESTAMP:
			g.builder.Call("timestamp_equal")
		default:
			g.internalError("equality for type not yet implemented: %v", exprType)
		}

	default:
		g.internalError("equality for type not yet implemented: %v", exprType)
	}
}

func (g *generator) generateNeq(f *ast.Function, depth int) {
	g.generateEq(f, depth+1)
	g.builder.Not()
}

func (g *generator) generateLor(f *ast.Function, depth int) {
	g.generate(f.Args[0], depth+1, celMode, "")
	lr := g.builder.AllocateLabel()
	le := g.builder.AllocateLabel()
	g.builder.Jz(lr)
	g.builder.APushBool(true)
	if depth == 0 {
		g.builder.Ret()
	} else {
		g.builder.Jmp(le)
	}
	g.builder.SetLabelPos(lr)
	g.generate(f.Args[1], depth+1, celMode, "")

	if depth != 0 {
		g.builder.SetLabelPos(le)
	}
}

func (g *generator) generateLand(f *ast.Function, depth int) {
	// Short circuit jump point for arguments that evaluate to false.
	lfalse := g.builder.AllocateLabel()

	// Label for the end of the and block.
	lend := g.builder.AllocateLabel()

	for i, a := range f.Args {
		g.generate(a, depth+1, celMode, "")
		if i < len(f.Args)-1 {
			// if this is not the last argument, check and jump to the false label.
			g.builder.Jz(lfalse)
		} else {
			g.builder.Jmp(lend)
		}
	}

	g.builder.SetLabelPos(lfalse)
	g.builder.APushBool(false)

	g.builder.SetLabelPos(lend)
}

func (g *generator) generateIndex(f *ast.Function, depth int, mode nilMode, valueJmpLabel string) {

	switch mode {
	case celMode:
		// should throw an error
		g.generate(f.Args[0], depth+1, celMode, "")
		g.generate(f.Args[1], depth+1, celMode, "")
		g.builder.Lookup()

	case nmNone:
		// Assume both the indexing target (arg[0]) and the index variable (arg[1]) are non-nil.
		g.generate(f.Args[0], depth+1, nmNone, "")

		if f.Args[1].Const != nil {
			str := f.Args[1].Const.Value.(string)
			g.builder.ANLookup(str)
		} else {
			g.generate(f.Args[1], depth+1, nmNone, "")
			g.builder.NLookup()
		}

	case nmJmpOnValue:
		// Assume both the indexing target (arg[0]) and the index variable (arg[1]) can be nil.
		// Stack the operations such that first the target is resolved, if not nil, then the
		// indexer is resolved, and finally if the indexer is not nil, the value is looked up.
		//
		//   As an example, consider the expression (ar["c"] | "foo") and we were supplied "LOuterLabel"
		//   as valueJmpLabel
		//   ...
		//                                     // Begin Args[0] eval (i.e. ar)
		//   tresolve_m "ar"                   // Try to resolve "ar"
		//   jnz LTargetResolved               // If resolved, then jump to LTargetResolved.
		//                                     // End Args[0] eval (i.e. ar)
		//
		//   jmp LEnd                          // This means Args[0] evaluated to nil, so jump to end
		//                                     // and fallthrough.
		// LTargetResolved:
		//                                     // Begin Args[1] eval (i.e. "c")
		//   apush_s "c"                       // Constant is directly evaluated.
		//                                     // End Args[1] eval (i.e. "c")
		//
		//   tlookup                           // With Args[0] and Arg[1] on stack, perform conditional lookup
		//   jnz LOuterLabel                   // If successful, jump to LOuterLabel
		// LEnd:                               // Otherwise fall through with value evaluated to nil.
		//                                     // From here on, evaluate the `... | "foo"` part, as "..." evaluated to nil.
		//   apush_s "foo"
		// LOuterLabel:
		//
		lEnd := g.builder.AllocateLabel()
		lTargetResolved := g.builder.AllocateLabel()
		g.generate(f.Args[0], depth+1, nmJmpOnValue, lTargetResolved)
		g.builder.Jmp(lEnd)
		g.builder.SetLabelPos(lTargetResolved)

		if f.Args[1].Const != nil {
			str := f.Args[1].Const.Value.(string)
			g.builder.APushStr(str)
		} else {
			lArgumentResolved := g.builder.AllocateLabel()
			g.generate(f.Args[1], depth+1, nmJmpOnValue, lArgumentResolved)
			g.builder.Jmp(lEnd)
			g.builder.SetLabelPos(lArgumentResolved)
		}
		g.builder.TLookup()
		g.builder.Jnz(valueJmpLabel)
		g.builder.SetLabelPos(lEnd)
	}
}

func (g *generator) generateOr(f *ast.Function, depth int, mode nilMode, valueJmpLabel string) {
	// TODO: Optimize code generation to check whether the first argument can be guaranteed to return a value (or error)
	// If so, we can elide the right-hand-side entirely.
	switch mode {
	case celMode:

		lEnd := g.builder.AllocateLabel()
		switch {
		case f.Args[0].Fn != nil && f.Args[0].Fn.Name == "OR":
			// rebalance expression (a | b) | c into a | (b | c)
			g.generateOr(&ast.Function{
				Name: "OR",
				Args: []*ast.Expression{f.Args[0].Fn.Args[0], &ast.Expression{
					Fn: &ast.Function{
						Name: "OR",
						Args: []*ast.Expression{f.Args[0].Fn.Args[1], f.Args[1]},
					},
				}},
			}, depth, celMode, valueJmpLabel)
			return

		case f.Args[0].Fn != nil && f.Args[0].Fn.Name == "INDEX":
			lhs := f.Args[0].Fn
			g.generate(lhs.Args[0], depth+1, celMode, "")
			g.generate(lhs.Args[1], depth+1, celMode, "")
			g.builder.TLookup()
			g.builder.Jnz(lEnd)
		case f.Args[0].Var != nil:
			g.generate(f.Args[0], depth+1, nmJmpOnValue, lEnd)
		case f.Args[0].Const != nil:
			g.generate(f.Args[0], depth+1, nmJmpOnValue, lEnd)
		default:
			g.internalError("unexpected LHS expression: %#v", f.Args[0].Fn)
			break
		}

		g.generate(f.Args[1], depth+1, celMode, "")
		g.builder.SetLabelPos(lEnd)
	case nmNone:
		// If the caller expects non-null result, evaluate Args[1] as non-null, and jump to end if
		// it resolves to a value.
		lEnd := g.builder.AllocateLabel()
		g.generate(f.Args[0], depth+1, nmJmpOnValue, lEnd)

		// Continue calculation of the right part of or, assuming left part evaluated to nil.
		// If this is a chain of "OR"s, we can chain through the ORs (i.e. the inner Or's Arg[0]
		// is evaluated as nmJmpOnValue, this allows short-circuiting all the way to the end.
		if f.Args[1].Fn != nil && f.Args[1].Fn.Name == "OR" {
			g.generate(f.Args[1], depth+1, nmJmpOnValue, lEnd)
		} else {
			g.generate(f.Args[1], depth+1, nmNone, "")
		}
		g.builder.SetLabelPos(lEnd)

	case nmJmpOnValue:
		g.generate(f.Args[0], depth+1, nmJmpOnValue, valueJmpLabel)
		g.generate(f.Args[1], depth+1, nmJmpOnValue, valueJmpLabel)
	}
}

func (g *generator) generateConditional(f *ast.Function, depth int, mode nilMode, valueJmplLabel string) {
	g.generate(f.Args[0], depth+1, celMode, "")

	labelFalse := g.builder.AllocateLabel()
	g.builder.Jz(labelFalse)
	g.generate(f.Args[1], depth+1, celMode, "")

	labelEnd := g.builder.AllocateLabel()
	g.builder.Jmp(labelEnd)

	g.builder.SetLabelPos(labelFalse)
	g.generate(f.Args[2], depth+1, celMode, "")

	g.builder.SetLabelPos(labelEnd)

	if mode == nmJmpOnValue {
		g.builder.Jmp(valueJmplLabel)
	}
}

func (g *generator) generateAdd(f *ast.Function, depth int) {
	exprType := g.evalType(f.Args[0])
	g.generate(f.Args[0], depth+1, celMode, "")
	g.generate(f.Args[1], depth+1, celMode, "")

	switch exprType {
	case il.String:
		g.builder.AddString()
	case il.Double:
		g.builder.AddDouble()
	case il.Integer:
		g.builder.AddInteger()
	default:
		g.internalError("Add for type not yet implemented: %v", exprType)
	}
}

func (g *generator) generateSize(f *ast.Function, depth int) {
	// eliminate polymorphic size function in bytecode
	exprType := g.evalType(f.Args[0])

	switch exprType {
	case il.String:
		if f.Args[0].Const != nil {
			constArg0 := f.Args[0].Const.Value.(string)
			g.builder.APushInt(int64(len(constArg0)))
		} else {
			g.generate(f.Args[0], depth+1, celMode, "")
			g.builder.SizeString()
		}
	default:
		g.internalError("Size for type not yet implemented: %v", exprType)
	}
}

func (g *generator) generateConstant(c *ast.Constant, mode nilMode, valueJmpLabel string) {
	switch c.Type {
	case descriptor.STRING:
		s := c.Value.(string)
		g.builder.APushStr(s)
	case descriptor.BOOL:
		b := c.Value.(bool)
		g.builder.APushBool(b)
	case descriptor.INT64:
		i := c.Value.(int64)
		g.builder.APushInt(i)
	case descriptor.DOUBLE:
		d := c.Value.(float64)
		g.builder.APushDouble(d)
	case descriptor.DURATION:
		u := c.Value.(time.Duration)
		g.builder.APushInt(int64(u))
	default:
		g.internalError("unhandled constant type: %v", c.Type)
	}

	switch mode {
	case nmNone:
		break

	case nmJmpOnValue:
		g.builder.Jmp(valueJmpLabel)

	case celMode:
		break

	default:
		g.internalError("unhandled nil mode: %v", mode)
	}
}

func (g *generator) internalError(format string, args ...interface{}) {
	log.Warnf("internal compiler error -- "+format, args)
	if g.err == nil {
		g.err = fmt.Errorf("internal compiler error -- %s", fmt.Sprintf(format, args...))
	}
}
