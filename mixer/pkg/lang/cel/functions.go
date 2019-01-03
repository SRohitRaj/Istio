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
	"github.com/google/cel-go/checker/decls"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func standardFunctions() []*exprpb.Decl {
	return []*exprpb.Decl{
		decls.NewFunction("startsWith",
			decls.NewInstanceOverload("startsWith",
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool)),
		decls.NewFunction("endsWith",
			decls.NewInstanceOverload("endsWith",
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool)),
		decls.NewFunction("match",
			decls.NewOverload("match",
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool)),
		decls.NewFunction("reverse",
			decls.NewInstanceOverload("reverse",
				[]*exprpb.Type{decls.String}, decls.String)),
		decls.NewFunction("reverse",
			decls.NewOverload("reverse",
				[]*exprpb.Type{decls.String}, decls.String)),
		decls.NewFunction("toLower",
			decls.NewOverload("toLower",
				[]*exprpb.Type{decls.String}, decls.String)),
		decls.NewFunction("email",
			decls.NewOverload("email",
				[]*exprpb.Type{decls.String}, decls.NewObjectType(emailAddressType))),
		decls.NewFunction("dnsName",
			decls.NewOverload("dnsName",
				[]*exprpb.Type{decls.String}, decls.NewObjectType(dnsType))),
		decls.NewFunction("uri",
			decls.NewOverload("uri",
				[]*exprpb.Type{decls.String}, decls.NewObjectType(uriType))),
		decls.NewFunction("ip",
			decls.NewOverload("ip",
				[]*exprpb.Type{decls.String}, decls.NewObjectType(ipAddressType))),
		decls.NewFunction("emptyStringMap",
			decls.NewOverload("emptyStringMap",
				[]*exprpb.Type{}, stringMapType)),
	}
}
