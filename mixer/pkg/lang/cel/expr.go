// Copyright 2018 Istio Authors
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

package cel

import (
	"fmt"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/parser"
	"github.com/hashicorp/go-multierror"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

const (
	elvis = "getOrElse"
)

// Parse a CEL expression
func Parse(text string) (ex *exprpb.Expr, err error) {
	source := common.NewTextSource(text)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during CEL parsing of expression %q", text)
		}
	}()

	parsed, errors := parser.Parse(source)
	if errors != nil && len(errors.GetErrors()) > 0 {
		err = fmt.Errorf("parsing error: %v", errors.ToDisplayString())
		return
	}

	// TODO: compute max ID in an expression tree
	u := &unroller{id: 100000}
	ex = u.unroll(parsed.Expr)
	err = u.err
	return
}

// Check verifies a CEL expressions against an attribute manifest
func Check(ex *exprpb.Expr, env *checker.Env) (*exprpb.CheckedExpr, error) {
	checked, errors := checker.Check(&exprpb.ParsedExpr{Expr: ex}, nil, env)
	if errors != nil && len(errors.GetErrors()) > 0 {
		return nil, fmt.Errorf("type checking error: %s", errorString(errors))
	}
	return checked, nil
}

func errorString(errors *common.Errors) string {
	out := ""
	for _, err := range errors.GetErrors() {
		if out != "" {
			out = out + "\n"
		}
		out = out + err.Message
	}
	return out
}

type unroller struct {
	err error
	id  int64
}

func (u *unroller) nextID() int64 {
	u.id++
	return u.id
}

// unroll eliminates Elvis and conditional() macros
func (u *unroller) unroll(in *exprpb.Expr) *exprpb.Expr {
	switch v := in.ExprKind.(type) {
	case *exprpb.Expr_ConstExpr, *exprpb.Expr_IdentExpr:
		// do nothing
	case *exprpb.Expr_SelectExpr:
		// recurse
		return &exprpb.Expr{
			Id: u.nextID(),
			ExprKind: &exprpb.Expr_SelectExpr{
				SelectExpr: &exprpb.Expr_Select{
					Operand:  u.unroll(v.SelectExpr.Operand),
					Field:    v.SelectExpr.Field,
					TestOnly: v.SelectExpr.TestOnly,
				},
			},
		}
	case *exprpb.Expr_CallExpr:
		// recurse
		var target *exprpb.Expr
		if v.CallExpr.Target != nil {
			target = u.unroll(v.CallExpr.Target)
		}

		args := make([]*exprpb.Expr, len(v.CallExpr.Args))
		for i, arg := range v.CallExpr.Args {
			args[i] = u.unroll(arg)
		}

		switch v.CallExpr.Function {
		case "conditional":
			return &exprpb.Expr{
				Id: u.nextID(),
				ExprKind: &exprpb.Expr_CallExpr{
					CallExpr: &exprpb.Expr_Call{
						Function: operators.Conditional,
						Target:   target,
						Args:     args,
					},
				},
			}
		case elvis:
			if target != nil {
				u.err = multierror.Append(u.err, fmt.Errorf("unexpected target in expression %q", v))
				break
			}

			// step through arguments to construct a conditional chain
			out := args[len(args)-1]
			var selector *exprpb.Expr
			for i := len(args) - 2; i >= 0; i-- {
				selector = nil
				switch lhs := args[i].ExprKind.(type) {
				case *exprpb.Expr_SelectExpr:
					if !lhs.SelectExpr.TestOnly {
						// a.f | x --> has(a.f) ? a.f : x
						selector = &exprpb.Expr{
							Id: u.nextID(),
							ExprKind: &exprpb.Expr_SelectExpr{
								SelectExpr: &exprpb.Expr_Select{
									Operand:  lhs.SelectExpr.Operand,
									Field:    lhs.SelectExpr.Field,
									TestOnly: true,
								},
							},
						}
					}

				case *exprpb.Expr_CallExpr:
					if lhs.CallExpr.Function == operators.Index {
						// a["f"] | x --> "f" in a ? a["f"] : x
						selector = &exprpb.Expr{
							Id: u.nextID(),
							ExprKind: &exprpb.Expr_CallExpr{
								CallExpr: &exprpb.Expr_Call{
									Function: operators.In,
									Args:     []*exprpb.Expr{lhs.CallExpr.Args[1], lhs.CallExpr.Args[0]},
								},
							},
						}
					}
				}

				// otherwise, a | b --> a
				if selector == nil {
					out = args[i]
				} else {
					out = &exprpb.Expr{
						Id: u.nextID(),
						ExprKind: &exprpb.Expr_CallExpr{
							CallExpr: &exprpb.Expr_Call{
								Function: operators.Conditional,
								Args:     []*exprpb.Expr{selector, args[i], out},
							},
						},
					}
				}
			}
			return out
		}

	default:
		u.err = multierror.Append(u.err, fmt.Errorf("unsupported expression kind %q", v))
	}
	return in
}
