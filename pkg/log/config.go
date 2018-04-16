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

// Package log provides the canonical logging functionality used by Go-based
// Istio components.
//
// Istio's logging subsystem is built on top of the [Zap](https://godoc.org/go.uber.org/zap) package.
// High performance scenarios should use the Error, Warn, Info, and Debug methods. Lower perf
// scenarios can use the more expensive convenience methods such as Debugf and Warnw.
//
// The package provides direct integration with the Cobra command-line processor which makes it
// easy to build programs that use a consistent interface for logging. Here's an example
// of a simple Cobra-based program using this log package:
//
//		func main() {
//			// get the default logging options
//			options := log.DefaultOptions()
//
//			rootCmd := &cobra.Command{
//				Run: func(cmd *cobra.Command, args []string) {
//
//					// configure the logging system
//					if err := log.Configure(options); err != nil {
//                      // print an error and quit
//                  }
//
//					// output some logs
//					log.Info("Hello")
//					log.Sync()
//				},
//			}
//
//			// add logging-specific flags to the cobra command
//			options.AttachCobraFlags(rootCmd)
//			rootCmd.SetArgs(os.Args[1:])
//			rootCmd.Execute()
//		}
//
// Once configured, this package intercepts the output of the standard golang "log" package as well as anything
// sent to the global zap logger (zap.L()).
package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

// This type exists so we can have a bit of logic sitting in front of all the real core calls for the
// case where we are intercepting other log packages. This lets us subject the log output to the
// default scope's log level, as is these other log packages were in fact directly outputting to the
// default scope.
type interceptor struct {
	zapcore.Core
}

func (in interceptor) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	switch ent.Level {
	case zapcore.ErrorLevel:
		if !defaultScope.ErrorEnabled() {
			return nil
		}
	case zapcore.WarnLevel:
		if !defaultScope.WarnEnabled() {
			return nil
		}
	case zapcore.InfoLevel:
		if !defaultScope.InfoEnabled() {
			return nil
		}
	case zapcore.DebugLevel:
		if !defaultScope.DebugEnabled() {
			return nil
		}
	}

	return in.Core.Check(ent, ce)
}

func (in interceptor) With(fields []zapcore.Field) zapcore.Core {
	return &interceptor{Core: in.Core.With(fields)}
}

func (in interceptor) Write(ent zapcore.Entry, f []zapcore.Field) error {
	return in.Core.Write(ent, f)
}

// none is used to disable logging output as well as to disable stack tracing.
const none zapcore.Level = 100

var levelToZap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	NoneLevel:  none,
}

func init() {
	// use our defaults for starters so that logging works even before everything is fully configured
	_ = Configure(DefaultOptions())
}

// prepZap is a utility function used by the Configure function.
func prepZap(options *Options) (zapcore.Core, zapcore.WriteSyncer, error) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "scope",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     formatDate,
	}

	var enc zapcore.Encoder
	if options.JSONEncoding {
		enc = zapcore.NewJSONEncoder(encCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encCfg)
	}

	var rotaterSink zapcore.WriteSyncer
	if options.RotateOutputPath != "" {
		rotaterSink = zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.RotateOutputPath,
			MaxSize:    options.RotationMaxSize,
			MaxBackups: options.RotationMaxAge,
			MaxAge:     options.RotationMaxBackups,
		})
	}

	errSink, closeErrorSink, err := zap.Open(options.ErrorOutputPaths...)
	if err != nil {
		return nil, nil, err
	}

	var outputSink zapcore.WriteSyncer
	if len(options.OutputPaths) > 0 {
		outputSink, _, err = zap.Open(options.OutputPaths...)
		if err != nil {
			closeErrorSink()
			return nil, nil, err
		}
	}

	var sink zapcore.WriteSyncer
	if rotaterSink != nil && outputSink != nil {
		sink = zapcore.NewMultiWriteSyncer(outputSink, rotaterSink)
	} else if rotaterSink != nil {
		sink = rotaterSink
	} else {
		sink = outputSink
	}

	return zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(zapcore.DebugLevel)), errSink, nil
}

func formatDate(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 27)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'

	enc.AppendString(string(buf))
}

func updateScopes(options *Options, core zapcore.Core, errSink zapcore.WriteSyncer) error {
	// init the global I/O funcs
	writeFn = core.Write
	syncFn = core.Sync
	errorSink = errSink

	// update the output levels of all scopes
	levels := strings.Split(options.outputLevels, ",")
	for _, sl := range levels {
		s, l, err := convertScopedLevel(sl)
		if err != nil {
			return err
		}

		if scope, ok := scopes[s]; ok {
			scope.Level = l
		} else {
			return fmt.Errorf("unknown scope '%s' specified", s)
		}
	}

	// update the stack tracing levels of all scopes
	levels = strings.Split(options.stackTraceLevels, ",")
	for _, sl := range levels {
		s, l, err := convertScopedLevel(sl)
		if err != nil {
			return err
		}

		if scope, ok := scopes[s]; ok {
			scope.StackTraceLevel = l
		} else {
			return fmt.Errorf("unknown scope '%s' specified", s)
		}
	}

	// update the caller location setting of all scopes
	sc := strings.Split(options.logCallers, ",")
	for _, s := range sc {
		if s == "" {
			continue
		}

		if scope, ok := scopes[s]; ok {
			scope.LogCallers = true
		} else {
			return fmt.Errorf("unknown scope '%s' specified", s)
		}
	}

	return nil
}

// Configure initializes Istio's logging subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the logging system is ready to accept data.
func Configure(options *Options) error {
	core, errSink, err := prepZap(options)
	if err != nil {
		return err
	}

	if err = updateScopes(options, core, errSink); err != nil {
		return err
	}

	opts := []zap.Option{
		zap.ErrorOutput(errSink),
		zap.AddCallerSkip(1),
	}

	if defaultScope.LogCallers {
		opts = append(opts, zap.AddCaller())
	}

	if defaultScope.StackTraceLevel != NoneLevel {
		opts = append(opts, zap.AddStacktrace(levelToZap[defaultScope.StackTraceLevel]))
	}

	captureLogger := zap.New(interceptor{core}, opts...)

	// capture global zap logging and force it through our logger
	_ = zap.ReplaceGlobals(captureLogger)

	// capture standard golang "log" package output and force it through our logger
	_ = zap.RedirectStdLog(captureLogger)

	// capture gRPC logging
	if options.LogGrpc {
		grpclog.SetLogger(zapgrpc.NewLogger(captureLogger.WithOptions(zap.AddCallerSkip(2))))
	}

	return nil
}

// reset by the Configure method
var syncFn = nopSync

func nopSync() error {
	return nil
}

// Sync flushes any buffered log entries.
// Processes should normally take care to call Sync before exiting.
func Sync() error {
	return syncFn()
}
