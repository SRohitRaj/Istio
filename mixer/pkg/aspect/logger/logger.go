// Copyright 2017 Google Inc.
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

package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"istio.io/mixer/pkg/aspect"
)

type (
	// Aspect is the interface for adapters that will handle logs data
	// within the mixer.
	Aspect interface {
		aspect.Aspect

		// Log directs a backend adapter to process a batch of
		// log entries derived from potentially several Report() calls.
		Log([]Entry) error
	}

	// Entry is the set of data that together constitutes a log entry.
	Entry struct {
		// LogName is the name of the log in which to record this entry.
		LogName string `json:"logName,omitempty"`
		// Labels are a set of metadata associated with this entry.
		// For instance, Labels can be used to describe the monitored
		// resource that corresponds to the log entry.
		Labels map[string]interface{} `json:"labels,omitempty"`
		// Timestamp is the time value for the log entry
		Timestamp string `json:"timestamp,omitempty"`
		// Severity indicates the log level for the log entry.
		Severity Severity `json:"severity,omitempty"`
		// TextPayload is textual logs data for which the entry is being
		// generated.
		TextPayload string `json:"textPayload,omitempty"`
		// StructPayload is a structured set of data, extracted from
		// a serialized JSON Object, that represent the actual data
		// for which the entry is being generated.
		StructPayload map[string]interface{} `json:"structPayload,omitempty"`
	}

	// Severity provides a set of logging levels for logger.Entry.
	Severity int

	// Adapter is the interface for building Aspect instances for mixer
	// logging backends.
	Adapter interface {
		aspect.Adapter

		// NewLogger returns a new Logger implementation, based on the
		// supplied Aspect configuration for the backend.
		NewLogger(env aspect.Env, c aspect.Config) (Aspect, error)
	}
)

const (
	// Default indicates that the log entry has no assigned severity level.
	Default Severity = iota
	// Debug indicates that the log entry has debug or trace information.
	Debug
	// Info indicates that the log entry has routine info, including status or performance data
	Info
	// Notice indicates that the log entry has normal, but significant events, such as config changes.
	Notice
	// Warning indicates that the log entry has data on events that might cause issues.
	Warning
	// Error indicates that the log entry has data on events that are likely to cause issues.
	Error
	// Critical indicates that the log entry has data on events related to severe problems or outages.
	Critical
	// Alert indicates that the log entry has data on which human action should be taken immediately.
	Alert
	// Emergency indicates that the log entry has data on one (or more) systems that are unusable.
	Emergency
)

var (
	severityMap = map[Severity]string{
		Default:   "DEFAULT",
		Debug:     "DEBUG",
		Info:      "INFO",
		Notice:    "NOTICE",
		Warning:   "WARNING",
		Error:     "ERROR",
		Critical:  "CRITICAL",
		Alert:     "ALERT",
		Emergency: "EMERGENCY",
	}
)

func (s Severity) String() string {
	if str, ok := severityMap[s]; ok {
		return str
	}
	return fmt.Sprintf("%d", s)
}

// SeverityByName returns a logger.Severity based on the supplied string. It
// attempts to look up the logger.Severity based on a map from enum to string.
// Default (with false) is returned if the name is not found.
func SeverityByName(s string) (Severity, bool) {
	s = strings.ToUpper(s)
	for severity, name := range severityMap {
		if name == s {
			return severity, true
		}
	}
	return Default, false
}

// MarshalJSON implements the json.Marshaler interface.
func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
