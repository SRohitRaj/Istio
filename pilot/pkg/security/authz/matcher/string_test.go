// Copyright 2019 Istio Authors
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

package matcher

import (
	"reflect"
	"testing"

	matcherpb "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
)

type testCase struct {
	name string
	v    string
	want *matcherpb.StringMatcher
}

func TestStringMatcherWithPrefix(t *testing.T) {
	testCases := []testCase{
		{
			name: "wildcardAsRequired",
			v:    "*",
			want: StringMatcherRegex(".+"),
		},
		{
			name: "prefix",
			v:    "-prefix-*",
			want: &matcherpb.StringMatcher{
				MatchPattern: &matcherpb.StringMatcher_Prefix{
					Prefix: "abc-prefix-",
				},
			},
		},
		{
			name: "suffix",
			v:    "*-suffix",
			want: &matcherpb.StringMatcher{
				MatchPattern: &matcherpb.StringMatcher_Suffix{
					Suffix: "abc-suffix",
				},
			},
		},
		{
			name: "exact",
			v:    "-exact",
			want: &matcherpb.StringMatcher{
				MatchPattern: &matcherpb.StringMatcher_Exact{
					Exact: "abc-exact",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := StringMatcherWithPrefix(tc.v, "abc")
			if !reflect.DeepEqual(*actual, *tc.want) {
				t.Errorf("want %s but got %s", tc.want.String(), actual.String())
			}
		})
	}
}

func TestStringMatcherRegex(t *testing.T) {
	testCases := []testCase{
		{
			name: "wildcardAsRequired",
			v:    "*",
			want: &matcherpb.StringMatcher{
				MatchPattern: &matcherpb.StringMatcher_SafeRegex{
					SafeRegex: &matcherpb.RegexMatcher{
						Regex: "*",
						EngineType: &matcherpb.RegexMatcher_GoogleRe2{
							GoogleRe2: &matcherpb.RegexMatcher_GoogleRE2{},
						},
					},
				},
			},
		},
		{
			name: "regexExpression",
			v:    "+?",
			want: &matcherpb.StringMatcher{
				MatchPattern: &matcherpb.StringMatcher_SafeRegex{
					SafeRegex: &matcherpb.RegexMatcher{
						Regex: "+?",
						EngineType: &matcherpb.RegexMatcher_GoogleRe2{
							GoogleRe2: &matcherpb.RegexMatcher_GoogleRE2{},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := StringMatcherRegex(tc.v); !reflect.DeepEqual(*actual, *tc.want) {
				t.Errorf("want %s but got %s", tc.want.String(), actual.String())
			}
		})
	}
}
