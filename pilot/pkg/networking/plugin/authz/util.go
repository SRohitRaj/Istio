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

package authz

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
)

// stringMatch checks if a string is in a list, it supports four types of string matches:
// 1. Exact match.
// 2. Wild character match. "*" matches any string.
// 3. Prefix match. For example, "book*" matches "bookstore", "bookshop", etc.
// 4. Suffix match. For example, "*/review" matches "/bookstore/review", "/products/review", etc.
func stringMatch(a string, list []string) bool {
	for _, s := range list {
		if a == s || s == "*" || prefixMatch(a, s) || suffixMatch(a, s) {
			return true
		}
	}
	return false
}

// prefixMatch checks if string "a" prefix matches "pattern".
func prefixMatch(a string, pattern string) bool {
	if !strings.HasSuffix(pattern, "*") {
		return false
	}
	pattern = strings.TrimSuffix(pattern, "*")
	return strings.HasPrefix(a, pattern)
}

// suffixMatch checks if string "a" prefix matches "pattern".
func suffixMatch(a string, pattern string) bool {
	if !strings.HasPrefix(pattern, "*") {
		return false
	}
	pattern = strings.TrimPrefix(pattern, "*")
	return strings.HasSuffix(a, pattern)
}

// convertToCidr converts a CIDR or a single IP string to a corresponding CidrRange. For a single IP
// string the converted CidrRange prefix is either 32 (for ipv4) or 128 (for ipv6).
func convertToCidr(v string) (*core.CidrRange, error) {
	var address string
	var prefixLen int

	if strings.Contains(v, "/") {
		if ip, ipnet, err := net.ParseCIDR(v); err == nil {
			address = ip.String()
			prefixLen, _ = ipnet.Mask.Size()
		} else {
			return nil, fmt.Errorf("invalid cidr range: %v", err)
		}
	} else {
		if ip := net.ParseIP(v); ip != nil {
			address = ip.String()
			if strings.Contains(v, ".") {
				// Set the prefixLen to 32 for ipv4 address.
				prefixLen = 32
			} else if strings.Contains(v, ":") {
				// Set the prefixLen to 128 for ipv6 address.
				prefixLen = 128
			} else {
				return nil, fmt.Errorf("invalid ip address: %s", v)
			}
		} else {
			return nil, fmt.Errorf("invalid ip address: %s", v)
		}
	}

	return &core.CidrRange{
		AddressPrefix: address,
		PrefixLen:     &types.UInt32Value{Value: uint32(prefixLen)},
	}, nil
}

// convertToPort converts a port string to a uint32.
func convertToPort(v string) (uint32, error) {
	p, err := strconv.ParseUint(v, 10, 32)
	if err != nil || p < 0 || p > 65535 {
		return 0, fmt.Errorf("invalid port %s: %v", v, err)
	}
	return uint32(p), nil
}

// convertToHeaderMatcher converts a key, value string pair to a corresponding HeaderMatcher.
func convertToHeaderMatcher(k, v string) *route.HeaderMatcher {
	//TODO(yangminzhu): Update the HeaderMatcher to support prefix and suffix match.
	if strings.Contains(v, "*") {
		return &route.HeaderMatcher{
			Name: k,
			HeaderMatchSpecifier: &route.HeaderMatcher_RegexMatch{
				RegexMatch: "^" + strings.Replace(v, "*", ".*", -1) + "$",
			},
		}
	}
	return &route.HeaderMatcher{
		Name: k,
		HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
			ExactMatch: v,
		},
	}
}
