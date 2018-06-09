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

package checker

import (
	"fmt"
	"go/token"
)

// Report populates lint report.
type Report struct {
	items []string
}

// NewLintReport creates and returns a Report object.
func NewLintReport() *Report {
	return &Report{}
}

// LReport returns report from Linter
func (lt *Report) Items() []string {
	return lt.items
}

// AddItem creates a new lint error report.
func (lr *Report) AddItem(pos token.Pos, fs *token.FileSet, msg string) {
	item := fmt.Sprintf("%v:%v:%v:%s",
		fs.Position(pos).Filename,
		fs.Position(pos).Line,
		fs.Position(pos).Column,
		msg)
	lr.items = append(lr.items, item)
}

// AddString creates a new string line in report.
func (lr *Report) AddString(msg string) {
	lr.items = append(lr.items, msg)
}
