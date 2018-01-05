// Copyright 2017 Istio Authors.
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

package papertrail

import (
	"fmt"
	"net"
	"strings"

	"istio.io/istio/mixer/pkg/adapter"
)

// LoggerImpl is an implementation of adapter.Logger for testing purposes
type LoggerImpl struct{}

// VerbosityLevel verifies the current verbosity level of the logger
func (l *LoggerImpl) VerbosityLevel(level adapter.VerbosityLevel) bool {
	return true
}

// Infof is for logging info level messages
func (l *LoggerImpl) Infof(format string, args ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", args...)
}

// Warningf is for logging warn level messages
func (l *LoggerImpl) Warningf(format string, args ...interface{}) {
	fmt.Printf("WARN: "+format+"\n", args...)
}

// Errorf is for logging error level messages
func (l *LoggerImpl) Errorf(format string, args ...interface{}) error {
	err := fmt.Errorf("Error: "+format+"\n", args...)
	fmt.Printf("%v", err)
	return err
}

// RunUDPServer is for running a test udp server
func RunUDPServer(port int, logger adapter.Logger, stopChan chan struct{}, trackChan chan struct{}) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Errorf("Error: %v", err)
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	go func() {
		for {
			noOfBytes, remoteAddr, err := conn.ReadFromUDP(buf)
			logger.Infof("udp server - data received: %s from %v", strings.TrimSpace(string(buf[0:noOfBytes])), remoteAddr)
			if err != nil {
				logger.Errorf("Error: %v", err)
				return
			}
			trackChan <- struct{}{}
		}
	}()
	var tobrk bool
	for {
		select {
		case <-stopChan:
			tobrk = true
		default:
		}
		if tobrk {
			break
		}
	}
}
