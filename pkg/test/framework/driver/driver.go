//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package driver

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"

	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/test/cluster"
	"istio.io/istio/pkg/test/dependency"
	"istio.io/istio/pkg/test/environment"
	"istio.io/istio/pkg/test/internal"
	"istio.io/istio/pkg/test/local"
)

const (
	maxTestIDLength = 30
)

var scope = log.RegisterScope("testframework", "General scope for the test framework", 0)
var lab = log.RegisterScope("testframework-lab", "Scope for normal log reporting to be used by the lab", 0)

// Driver for the test framework.
type Driver interface {
	// Run the tests by calling into testing.M. This must be called exactly once.
	Run(testID string, m *testing.M) (int, error)

	// SuiteRequires indicates that the whole suite requires particular dependencies.
	SuiteRequires(dependencies []dependency.Instance) error

	// Requires checks and initializes the supplied dependencies appropriately.
	Requires(t testing.TB, dependencies []dependency.Instance)

	// AcquireEnvironment resets and returns the environment. Once AcquireEnvironment should be called exactly
	// once per test.
	AcquireEnvironment(t testing.TB) environment.Interface
}

// internal state for the driver.
type state int

const (
	// Runtime states for the driver

	// The driver is freshly created. It hasn't run yet.
	created state = iota

	// The driver is currently running tests.
	running

	// the driver has completed running tests.
	completed
)

type driver struct {
	lock sync.Mutex

	args *args
	ctx  *internal.TestContext

	// The names of the tests that we've encountered so far.
	testNames map[string]struct{}

	state state
}

var _ Driver = &driver{}

// New returns a new driver instance.
func New() Driver {
	return &driver{
		testNames: make(map[string]struct{}),
		state:     created,
	}
}

//
//// Initialize implements same-named Driver method.
//func (d *driver) Initialize(a *args) error {
//	scope.Debugf("Enter: driver.Initialize (%s)", a.TestID)
//	d.lock.Lock()
//	defer d.lock.Unlock()
//
//	if d.ctx != nil {
//		return errors.New("test driver is already initialized")
//	}
//
//	// Make a copy of the args
//	args := &(*a)
//	if err := args.Validate(); err != nil {
//		return err
//	}
//
//	// Initialize the environment.
//	var env internal.Environment
//	switch a.Environment {
//	case EnvLocal:
//		env = local.NewEnvironment()
//	case EnvKube:
//		env = cluster.NewEnvironment()
//	default:
//		return fmt.Errorf("unrecognized environment: %s", a.Environment)
//	}
//
//	// Create the context, but do not attach it to the member fields before initializing.
//	ctx := internal.NewTestContext(
//		args.TestID,
//		generateRunID(args.TestID),
//		a.WorkDir,
//		a.Hub,
//		a.Tag,
//		a.KubeConfig,
//		env)
//
//	if err := env.Initialize(ctx); err != nil {
//		return err
//	}
//
//	d.ctx = ctx
//	d.args = args
//
//	return nil
//}

// Run implements same-named Driver method.
func (d *driver) Run(testID string, m *testing.M) (int, error) {
	d.lock.Lock()
	// do not defer unlock. We will explicitly unlock before starting test runs.
	if d.state != created {
		return -1, fmt.Errorf("driver.Run must be called only once")
	}
	d.state = running

	args, err := calculateArgs(testID, m)
	if err != nil {
		d.lock.Unlock()
		return -1, err
	}
	scope.Debugf("driver args: %+v", args)

	if err = d.initialize(args); err != nil {
		d.lock.Unlock()
		return -2, err
	}

	for _, dep := range d.args.SuiteDependencies {
		if err := d.ctx.Tracker().Initialize(d.ctx, d.ctx.Environment(), dep); err != nil {
			log.Errorf("driver.Run: Dependency error '%s': %v", dep, err)
			return -3, err
		}
	}

	// Call m.Run() while not holding the lock.
	d.lock.Unlock()
	lab.Infof(">>> Beginning test run for: '%s'", testID)
	rt := m.Run()
	lab.Infof("<<< Completing test run for: '%s'", testID)

	// Reacquire lock.
	d.lock.Lock()
	defer d.lock.Unlock()

	d.state = completed

	if !d.args.NoCleanup {
		d.ctx.Tracker().Cleanup()
	}

	return rt, nil
}

// AcquireEnvironment implementation
func (d *driver) AcquireEnvironment(t testing.TB) environment.Interface {
	t.Helper()
	scope.Debugf("Enter: driver.AcquireEnvionment (%s)", d.ctx.TestID())
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.state != running {
		t.Fatalf("Test driver is not running.")
	}

	// Double check the test name to see if this is a singleton call for this test?
	if _, ok := d.testNames[t.Name()]; ok {
		t.Fatalf("AcquireEnvironment should be called only once during a test session. (test='%s')", t.Name())
	}
	d.testNames[t.Name()] = struct{}{}

	if err := d.ctx.Environment().Reset(); err != nil {
		t.Fatalf("AcquireEnvironment failed to reset the environment state: %v", err)
	}

	// Reset all resettables, as we're going to be executing within the context of a new test.
	if err := d.ctx.Tracker().Reset(d.ctx); err != nil {
		t.Fatalf("driver.AcquireEnvironment failed to reset the resource state: %v", err)
	}

	return d.ctx.Environment()
}

// SuiteRequires indicates that the whole suite requires particular dependencies.
func (d *driver) SuiteRequires(dependencies []dependency.Instance) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.state != created {
		return fmt.Errorf("test driver is not in a valid state")
	}

	arguments.SuiteDependencies = append(arguments.SuiteDependencies, dependencies...)
	return nil
}

// Requires implements same-named Driver method.
func (d *driver) Requires(t testing.TB, dependencies []dependency.Instance) {
	t.Helper()
	scope.Debugf("Enter: driver.Requires (%s)", d.ctx.TestID())
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.state != running {
		t.Fatalf("Test driver is not running.")
	}

	// Initialize dependencies only once.
	for _, dep := range dependencies {
		if err := d.ctx.Tracker().Initialize(d.ctx, d.ctx.Environment(), dep); err != nil {
			log.Errorf("Failed to initialize dependency '%s': %v", dep, err)
			t.Fatalf("unable to satisfy dependency '%v': %v", dep, err)
		}
	}
}

func (d *driver) initialize(a *args) error {

	// Initialize the environment.
	var env internal.Environment
	switch a.Environment {
	case EnvLocal:
		env = local.NewEnvironment()
	case EnvKube:
		env = cluster.NewEnvironment()
	default:
		return fmt.Errorf("unrecognized environment: %s", a.Environment)
	}

	// Create the context, but do not attach it to the member fields before initializing.
	ctx := internal.NewTestContext(
		a.TestID,
		generateRunID(a.TestID),
		a.WorkDir,
		a.Hub,
		a.Tag,
		a.KubeConfig,
		env)

	if err := env.Initialize(ctx); err != nil {
		return err
	}

	d.ctx = ctx
	d.args = a

	return nil
}

func generateRunID(testID string) string {
	u := uuid.New().String()
	u = strings.Replace(u, "-", "", -1)
	testID = strings.Replace(testID, "_", "-", -1)
	// We want at least 6 characters of uuid padding
	padding := maxTestIDLength - len(testID)
	return fmt.Sprintf("%s-%s", testID, u[0:padding])
}

func calculateArgs(testID string, m *testing.M) (*args, error) {
	if err := processFlags(); err != nil {
		return nil, err
	}

	if err := arguments.Validate(); err != nil {
		return nil, err
	}

	arguments.TestID = testID
	arguments.M = m

	return arguments, nil
}
