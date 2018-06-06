// Copyright 2018 Istio Authors. All Rights Reserved.
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

package linter

import (
	"istio.io/istio/tests/util/golinter/rules"
)

// UnitTestRules are rules which should apply to unit test file
var UnitTestRules = []rules.LintRule{
	rules.NewSkipByIssueRule(),
	rules.NewNoGoroutineRule(),
	rules.NewNoSleepRule(),
	rules.NewNoShortRule(),
}

// IntegTestRules are rules which should apply to integration test file
var IntegTestRules = []rules.LintRule{
	rules.NewSkipByIssueRule(),
	rules.NewSkipByShortRule(),
}

// E2eTestRules are rules which should apply to e2e test file
var E2eTestRules = []rules.LintRule{
	rules.NewSkipByIssueRule(),
	rules.NewSkipByShortRule(),
}
