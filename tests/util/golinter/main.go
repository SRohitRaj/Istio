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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"istio.io/istio/tests/util/golinter/linter"
)

var exitCode int

func main() {
	flag.Parse()
	for _, r := range getReport(flag.Args()) {
		reportErr(r)
	}
	os.Exit(exitCode)
}

func getReport(args []string) []string {
	var reports []string
	if len(args) == 0 {
		reports = doAllDirs([]string{"."})
	} else {
		reports = doAllDirs(args)
	}
	return reports
}

func doAllDirs(args []string) []string {
	rpts := make([]string, 0)
	pFilter := linter.NewPathFilter()
	for _, path := range args {
		if !filepath.IsAbs(path) {
			path, _ = filepath.Abs(path)
		}
		err := filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
			if err != nil {
				reportErr(fmt.Sprintf("pervent panic by handling failure accessing a path %q: %v", fpath, err))
				return err
			}
			if testType, sRuleMaps := pFilter.GetTestType(fpath, info); testType != linter.NonTest {
				lt := linter.NewLinter(fpath, testType, sRuleMaps)
				lt.Run()
				rpts = append(rpts, lt.LReport()...)
			}
			return nil
		})
		if err != nil {
			reportErr(fmt.Sprintf("error visiting the path %q: %v", path, err))
		}
	}
	return rpts
}

func reportErr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	exitCode = 2
}
