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

package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
)

const (
	statCdsVersion    = "cluster_manager.cds.version"
	statLdsVersion    = "listener_manager.lds.version"
	statServerState   = "server.state"
	versionStatsRegex = "^(cluster_manager.cds|listener_manager.lds).version$"
)

type stat struct {
	name  string
	value *uint64
	found bool
}

// Stats contains values of interest from a poll of Envoy stats.
type Stats struct {
	// Hash of the contents from the last successful cluster update.
	CDSVersion uint64
	// Hash of the contents from the last successful listener update.
	LDSVersion uint64
	// Server State of Envoy.
	ServerState uint64
}

// String representation of the Stats.
func (s *Stats) String() string {
	cdsStatus := "Not Received"
	ldsStatus := "Not Received"
	if s.CDSVersion > 0 {
		cdsStatus = "Received"
	}
	if s.LDSVersion > 0 {
		ldsStatus = "Received"
	}
	return fmt.Sprintf("cds update: %s , lds update: %s", cdsStatus, ldsStatus)
}

// GetServerState returns the current Envoy state by checking the "server.state" stat.
func GetServerState(localHostAddr string, adminPort uint16) (*uint64, error) {
	stats, err := doHTTPGet(fmt.Sprintf("http://%s:%d/stats?usedonly&filter=%s", localHostAddr, adminPort, statServerState))
	if err != nil {
		return nil, multierror.Prefix(err, "failed retrieving Envoy server state stat:")
	}

	s := &Stats{}
	allStats := []*stat{
		{name: statServerState, value: &s.ServerState},
	}
	if err := parseStats(stats, allStats); err != nil {
		return nil, err
	}

	return &s.ServerState, nil
}

// GetVersionStats returns the version stats for CDS and LDS.
func GetVersionStats(localHostAddr string, adminPort uint16) (*Stats, error) {
	stats, err := doHTTPGet(fmt.Sprintf("http://%s:%d/stats?usedonly&filter=%s", localHostAddr, adminPort, versionStatsRegex))
	if err != nil {
		return nil, multierror.Prefix(err, "failed retrieving Envoy stats:")
	}

	s := &Stats{}
	allStats := []*stat{
		{name: statCdsVersion, value: &s.CDSVersion},
		{name: statLdsVersion, value: &s.LDSVersion},
	}
	if err := parseStats(stats, allStats); err != nil {
		return nil, err
	}

	return s, nil
}

func parseStats(input *bytes.Buffer, stats []*stat) (err error) {
	for input.Len() > 0 {
		line, _ := input.ReadString('\n')
		for _, stat := range stats {
			if e := stat.processLine(line); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, stat := range stats {
		if !stat.found {
			*stat.value = 0
		}
	}
	return
}

func (s *stat) processLine(line string) error {
	if !s.found && strings.HasPrefix(line, s.name) {
		s.found = true

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return fmt.Errorf("envoy stat %s missing separator. line:%s", s.name, line)
		}

		val, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			return fmt.Errorf("failed parsing Envoy stat %s (error: %s) line: %s", s.name, err.Error(), line)
		}

		*s.value = val
	}

	return nil
}
