// Copyright 2016 Google Inc.
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
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	tracer "github.com/opentracing/basictracer-go"
	ot "github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"

	"istio.io/mixer/pkg/api"
	"istio.io/mixer/pkg/attribute"
)

type serverArgs struct {
	port                 uint
	maxMessageSize       uint
	maxConcurrentStreams uint
	compressedPayload    bool
	serverCertFile       string
	serverKeyFile        string
	clientCertFiles      string
}

func serverCmd(errorf errorFn) *cobra.Command {
	sa := &serverArgs{}
	serverCmd := cobra.Command{
		Use:   "server",
		Short: "Starts the mixer as a server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting gRPC server on port %v, type 'exit' to stop gracefully.\n", sa.port)

			// The tracer will keep track of spans in memory and print them to stdout when the server finishes.
			recorder := tracer.NewInMemoryRecorder()
			t := tracer.New(recorder)
			defer func() {
				for _, span := range recorder.GetSpans() {
					fmt.Printf("%v\n", span)
				}
			}()

			srv, err := newServer(sa, t)
			if err != nil {
				errorf("%v", err)
				os.Exit(1)
			}

			go func() {
				if err := srv.Start(); err != nil {
					errorf("%v", err)
					os.Exit(1)
				}
			}()

			loopTillStopped(srv)
		},
	}
	serverCmd.PersistentFlags().UintVarP(&sa.port, "port", "p", 9091, "TCP port to use for the mixer's gRPC API")
	serverCmd.PersistentFlags().UintVarP(&sa.maxMessageSize, "maxMessageSize", "", 1024*1024, "Maximum size of individual gRPC messages")
	serverCmd.PersistentFlags().UintVarP(&sa.maxConcurrentStreams, "maxConcurrentStreams", "", 32, "Maximum supported number of concurrent gRPC streams")
	serverCmd.PersistentFlags().BoolVarP(&sa.compressedPayload, "compressedPayload", "", false, "Whether to compress gRPC messages")
	serverCmd.PersistentFlags().StringVarP(&sa.serverCertFile, "serverCertFile", "", "", "The TLS cert file")
	serverCmd.PersistentFlags().StringVarP(&sa.serverKeyFile, "serverKeyFile", "", "", "The TLS key file")
	serverCmd.PersistentFlags().StringVarP(&sa.clientCertFiles, "clientCertFiles", "", "", "A set of comma-separated client X509 cert files")

	return &serverCmd
}

func newServer(sa *serverArgs, tracer ot.Tracer) (*api.GRPCServer, error) {
	var serverCert *tls.Certificate
	var clientCerts *x509.CertPool

	if sa.serverCertFile != "" && sa.serverKeyFile != "" {
		sc, err := tls.LoadX509KeyPair(sa.serverCertFile, sa.serverKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load server certificate and server key: %v", err)
		}
		serverCert = &sc
	}

	if sa.clientCertFiles != "" {
		clientCerts = x509.NewCertPool()
		for _, clientCertFile := range strings.Split(sa.clientCertFiles, ",") {
			pem, err := ioutil.ReadFile(clientCertFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load client certificate: %v", err)
			}
			clientCerts.AppendCertsFromPEM(pem)
		}
	}

	attrMgr := attribute.NewManager()

	grpcServerOptions := api.GRPCServerOptions{
		Port:                 uint16(sa.port),
		MaxMessageSize:       sa.maxMessageSize,
		MaxConcurrentStreams: sa.maxConcurrentStreams,
		CompressedPayload:    sa.compressedPayload,
		ServerCertificate:    serverCert,
		ClientCertificates:   clientCerts,
		Handlers:             api.NewMethodHandlers(),
		AttributeManager:     attrMgr,
		EnableTracing:        true,
		Tracer:               tracer,
	}

	var grpcServer *api.GRPCServer
	var err error
	if grpcServer, err = api.NewGRPCServer(&grpcServerOptions); err != nil {
		return nil, fmt.Errorf("unable to initialize gRPC server " + err.Error())
	}

	return grpcServer, nil
}

func loopTillStopped(srv *api.GRPCServer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "exit" {
			srv.Stop()
			return
		}
	}
}
