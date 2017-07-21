// Copyright 2017 Istio Authors
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

package fortio

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// LogLevel is the level of logging (0 Verbose -> 5 Fatal).
type LogLevel int

// Log levels
const (
	D LogLevel = iota // Debug
	V LogLevel = iota // Verbose
	I LogLevel = iota // Info
	W LogLevel = iota // Warning
	E LogLevel = iota // Error
	C LogLevel = iota // Critical
	F LogLevel = iota // Fatal
)

var level = I // default is Info and up
var levelToStrA []string
var levelToStrM map[string]LogLevel

func init() {
	levelToStrA = []string{
		"Debug",
		"Verbose",
		"Info",
		"Warning",
		"Error",
		"Critical",
		"Fatal",
	}
	levelToStrM = make(map[string]LogLevel, 2*len(levelToStrA))
	for l, name := range levelToStrA {
		// Allow both -loglevel Verbose and -loglevel verbose ...
		levelToStrM[name] = LogLevel(l)
		levelToStrM[strings.ToLower(name)] = LogLevel(l)
	}
	flag.Var(&level, "loglevel", fmt.Sprintf("loglevel, one of %v", levelToStrA))
	log.SetFlags(log.Flags() | log.Lshortfile)
}

// String returns the string representation of the level.
// Needed for flag Var interface.
func (l *LogLevel) String() string {
	return (*l).ToString()
}

// ToString returns the string representation of the level.
// (somehow this can't be the same name as the pointer receiver version)
func (l LogLevel) ToString() string {
	return levelToStrA[l]
}

// Set is called by the flags.
func (l *LogLevel) Set(str string) error {
	var lvl LogLevel
	var ok bool
	if lvl, ok = levelToStrM[str]; !ok {
		return fmt.Errorf("Invalid value for loglevel '%s' should be one of %v", str, levelToStrA)
	}
	SetLogLevel(lvl)
	return nil
}

// SetLogLevel sets the log level and returns the previous one.
func SetLogLevel(lvl LogLevel) LogLevel {
	prev := level
	if lvl < D {
		log.Printf("SetLogLevel called with level %d lower than Debug!", lvl)
		return -1
	}
	if lvl > F {
		log.Printf("SetLogLevel called with level %d higher than Fatal!", lvl)
		return -1
	}
	Log(I, "Log level is now %d %s (was %d %s)\n", lvl, lvl.ToString(), prev, prev.ToString())
	level = lvl
	return prev
}

// GetLogLevel returns the current level
func GetLogLevel() LogLevel {
	return level
}

// IsLogging returns true if a given level is currently logged.
func IsLogging(lvl LogLevel) bool {
	return lvl >= level
}

// LogLevelByName returns the LogLevel by its name.
func LogLevelByName(str string) LogLevel {
	return levelToStrM[str]
}

// Log at the given level.
func Log(lvl LogLevel, format string, rest ...interface{}) {
	if !IsLogging(lvl) {
		return
	}
	log.Print(levelToStrA[lvl][0:1], " ", fmt.Sprintf(format, rest...))
	if lvl == F {
		panic("aborting...")
	}
}

// -- would be nice to be able to create those in a loop instead of copypasta:

// LogDebug logs if Debug level is on.
func LogDebug(format string, rest ...interface{}) {
	Log(D, format, rest...)
}

// LogV logs if Verbose level is on.
func LogV(format string, rest ...interface{}) {
	Log(V, format, rest...)
}

// LogInfo logs if Info level is on.
func LogInfo(format string, rest ...interface{}) {
	Log(I, format, rest...)
}

// LogWarning logs if Warning level is on.
func LogWarning(format string, rest ...interface{}) {
	Log(W, format, rest...)
}

// LogError logs if Warning level is on.
func LogError(format string, rest ...interface{}) {
	Log(E, format, rest...)
}

// LogCritical logs if Warning level is on.
func LogCritical(format string, rest ...interface{}) {
	Log(C, format, rest...)
}

// LogFatal logs if Warning level is on.
func LogFatal(format string, rest ...interface{}) {
	Log(F, format, rest...)
}

// IsLoggingDebug is a shortcut for IsLogging(D).
func IsLoggingDebug() bool {
	return IsLogging(D)
}

// IsLoggingV is a shortcut for IsLogging(V).
func IsLoggingV() bool {
	return IsLogging(V)
}
