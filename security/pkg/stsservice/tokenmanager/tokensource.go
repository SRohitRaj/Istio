// Copyright Istio Authors
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

package tokenmanager

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/oauth2"

	"istio.io/istio/security/pkg/stsservice"
)

// TokenSource specifies a oauth token source based on STS token exchange.
type TokenSource struct {
	tm           stsservice.TokenManager
	subjectToken string
	authScope    string
}

var _ oauth2.TokenSource = &TokenSource{}

// NewTokenSource creates a token source based on STS token exchange.
func NewTokenSource(trustDomain, subjectToken, authScope string) (*TokenSource, error) {
	return &TokenSource{
		tm:           CreateTokenManager(GoogleTokenExchange, Config{TrustDomain: trustDomain}),
		subjectToken: subjectToken,
		authScope:    authScope,
	}, nil
}

func (ts *TokenSource) setTokenManager(tm stsservice.TokenManager) {
	ts.tm = tm
}

// Token returns Oauth token received from sts token exchange.
func (ts *TokenSource) Token() (*oauth2.Token, error) {
	params := stsservice.StsRequestParameters{
		GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
		Scope:            ts.authScope,
		SubjectToken:     ts.subjectToken,
		SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
	}
	body, err := ts.tm.GenerateToken(params)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange access token: %v", err)
	}
	respData := &stsservice.StsResponseParameters{}
	if err := json.Unmarshal(body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access token response data: %v", err)
	}

	return &oauth2.Token{
		AccessToken:  respData.AccessToken,
		TokenType:    respData.TokenType,
		RefreshToken: respData.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(respData.ExpiresIn)),
	}, nil
}
