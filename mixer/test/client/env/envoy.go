// Copyright 2017 Istio Authors. All Rights Reserved.
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

package env

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Envoy struct {
	cmd   *exec.Cmd
	ports *Ports
}

// Run command and return the merged output from stderr and stdout, error code
func Run(name string, args ...string) (s string, err error) {
	log.Println(">", name, strings.Join(args, " "))
	c := exec.Command(name, args...)
	bytes, err := c.CombinedOutput()
	s = string(bytes)
	for _, line := range strings.Split(s, "\n") {
		log.Println(line)
	}
	if err != nil {
		log.Println(err)
	}
	return
}

func NewEnvoy(conf, flags string, stress, faultInject bool, v2 *V2Conf, ports *Ports) (*Envoy, error) {
	// Asssume test environment has copied latest envoy to $HOME/go/bin in bin/init.sh
	envoy_path := "envoy"
	conf_path := fmt.Sprintf("/tmp/config.conf.%v", ports.AdminPort)
	log.Printf("Envoy config: in %v\n%v\n", conf_path, conf)
	if err := CreateEnvoyConf(conf_path, conf, flags, stress, faultInject, v2, ports); err != nil {
		return nil, err
	}

	args := []string{"-c", conf_path, "--base-id", strconv.Itoa(int(ports.AdminPort))}
	if stress {
		args = append(args, "--concurrency", "10")
	} else {
		args = append(args, "-l", "debug", "--concurrency", "1")
	}

	cmd := exec.Command(envoy_path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return &Envoy{
		cmd:   cmd,
		ports: ports,
	}, nil
}

func (s *Envoy) Start() error {
	err := s.cmd.Start()
	if err == nil {
		url := fmt.Sprintf("http://localhost:%v/server_info", s.ports.AdminPort)
		WaitForHttpServer(url)
		WaitForPort(s.ports.ClientProxyPort)
		WaitForPort(s.ports.ServerProxyPort)
	}
	return err
}

func (s *Envoy) Stop() error {
	log.Printf("Kill Envoy ...\n")
	err := s.cmd.Process.Kill()
	log.Printf("Kill Envoy ... Done\n")
	return err
}
