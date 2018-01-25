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

package handler

import (
	"testing"

	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime2/config"
	"istio.io/istio/mixer/pkg/runtime2/testing/data"
	"istio.io/istio/mixer/pkg/runtime2/testing/util"
)

// Create a standard global config with Handler H1, Instance I1 and rule R1 referencing I1 and H1.
// Most of the tests use this config.
var globalCfg = data.JoinConfigs(data.HandlerACheck1, data.InstanceCheck1, data.RuleCheck1)

//  Alternate configuration where R2 references H1, but uses instances I1 and I2.
var globalCfgI2 = data.JoinConfigs(data.HandlerACheck1, data.InstanceCheck1, data.InstanceCheck2, data.RuleCheck1WithInstance1And2)

func TestNew_EmptyConfig(t *testing.T) {
	s := config.Empty()

	table := NewTable(Empty(), s, nil)
	e, found := table.Get(data.FqnACheck1)
	if found {
		t.Fatal("found")
	}

	if e.name != "" {
		t.Fatal("non-empty handler")
	}
}

func TestNew_EmptyOldTable(t *testing.T) {
	adapters := data.BuildAdapters()
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)
	e, found := table.Get(data.FqnACheck1)
	if !found {
		t.Fatal("not found")
	}

	if e.name == "" {
		t.Fatal("empty handler")
	}
}

func TestNew_Reuse(t *testing.T) {
	adapters := data.BuildAdapters()
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)

	// NewTable again using the same config, but add fault to the adapter to detect change.
	adapters = data.BuildAdapters(data.FakeAdapterSettings{Name: "tcheck", ErrorAtBuild: true})
	s = util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table2 := NewTable(table, s, nil)

	if len(table2.entries) != 1 {
		t.Fatal("size")
	}

	if newEntry, found := table2.entries[data.FqnACheck1]; !found || newEntry != table.entries[data.FqnACheck1] {
		t.Fail()
	}
}

func TestNew_NoReuse_DifferentConfig(t *testing.T) {
	adapters := data.BuildAdapters()
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)

	// NewTable again using the slightly different config
	s = util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfgI2)

	table2 := NewTable(table, s, nil)

	if len(table2.entries) != 1 {
		t.Fatal("size")
	}

	if table2.entries[data.FqnACheck1] == table.entries[data.FqnACheck1] {
		t.Fail()
	}
}

func TestTable_Get(t *testing.T) {
	table := &table{
		entries: make(map[string]entry),
	}

	table.entries["h1"] = entry{
		name:    "h1",
		handler: &data.FakeHandler{},
	}

	e, found := table.Get("h1")
	if !found {
		t.Fail()
	}

	if e.name == "" {
		t.Fail()
	}
}

func TestTable_Get_Empty(t *testing.T) {
	table := &table{
		entries: make(map[string]entry),
	}

	_, found := table.Get("h1")
	if found {
		t.Fail()
	}
}

func TestEmpty(t *testing.T) {
	table := Empty()

	if len(table.entries) > 0 {
		t.Fail()
	}
}

func TestCleanup_Basic(t *testing.T) {
	closeCalled := false
	adapters := data.BuildAdapters(data.FakeAdapterSettings{Name: "acheck", CloseCalled: &closeCalled})
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)

	s = config.Empty()

	table2 := NewTable(table, s, nil)

	table.Cleanup(table2)

	if !closeCalled {
		t.Fail()
	}
}

func TestCleanup_WorkerNotClosed(t *testing.T) {
	tests := []struct {
		SpawnWorker      bool
		SpawnDaemon      bool
		CloseGoRoutines  bool
		wantStrayRoutine float64
	}{

		{
			SpawnWorker:      true,
			SpawnDaemon:      true,
			CloseGoRoutines:  false,
			wantStrayRoutine: 1,
		},
		{
			SpawnWorker:      true,
			SpawnDaemon:      true,
			CloseGoRoutines:  true,
			wantStrayRoutine: 0,
		},
	}
	for idx, tt := range tests {
		t.Run(fmt.Sprintf("[%d]", idx), func(t *testing.T) {
			adapters := data.BuildAdapters(data.FakeAdapterSettings{Name: "acheck", SpawnWorker: tt.SpawnWorker, SpawnDaemon: tt.SpawnDaemon, CloseGoRoutines: tt.CloseGoRoutines})
			templates := data.BuildTemplates()

			s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)
			s.ID = int64(idx * 2)

			oldTable := NewTable(Empty(), s, pool.NewGoroutinePool(5, false))

			s = config.Empty()
			s.ID = int64(idx*2 + 1)

			newTable := NewTable(oldTable, s, nil)

			oldTable.Cleanup(newTable)

			var c prometheus.Metric = oldTable.entries["hcheck1.acheck.istio-system"].env.counters.workers
			m := new(dto.Metric)
			_ = c.Write(m)
			if *m.GetGauge().Value != tt.wantStrayRoutine {
				t.Fatalf("expected %v worked; got %v", tt.wantStrayRoutine, *m.GetGauge().Value)
			}
		})
	}
}

func TestCleanup_NoChange(t *testing.T) {
	closeCalled := false
	adapters := data.BuildAdapters(data.FakeAdapterSettings{Name: "acheck", CloseCalled: &closeCalled})
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)

	// use same config again.
	table2 := NewTable(table, s, nil)

	table.Cleanup(table2)

	if closeCalled {
		t.Fail()
	}
}

func TestCleanup_EmptyNewTable(t *testing.T) {
	adapters := data.BuildAdapters()
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)

	table := NewTable(Empty(), s, nil)

	// Use an empty table as current.
	table.Cleanup(Empty())
}

func TestCleanup_WithStartupError(t *testing.T) {
	adapters := data.BuildAdapters()

	templates := data.BuildTemplates(data.FakeTemplateSettings{Name: "tcheck", HandlerDoesNotSupportTemplate: true})

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)
	table := NewTable(Empty(), s, nil)

	if _, found := table.Get(data.FqnACheck1); found {
		t.Fail()
	}

	// use different config to force cleanup
	table2 := NewTable(table, config.Empty(), nil)

	table.Cleanup(table2)

	// Close shouldn't be called.
	if len(table2.entries) != 0 {
		t.Fail()
	}
}

func TestCleanup_CloseError(t *testing.T) {
	closeCalled := false
	adapters := data.BuildAdapters(data.FakeAdapterSettings{Name: "acheck", ErrorAtHandlerClose: true, CloseCalled: &closeCalled})
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)
	table := NewTable(Empty(), s, nil)

	// use different config to force cleanup
	table2 := NewTable(table, config.Empty(), nil)

	table.Cleanup(table2)

	if len(table2.entries) != 0 {
		t.Fail()
	}

	if !closeCalled {
		t.Fail()
	}
}

func TestCleanup_ClosePanic(t *testing.T) {
	closeCalled := false
	adapters := data.BuildAdapters(data.FakeAdapterSettings{Name: "acheck", PanicAtHandlerClose: true, CloseCalled: &closeCalled})
	templates := data.BuildTemplates()

	s := util.GetSnapshot(templates, adapters, data.ServiceConfig, globalCfg)
	table := NewTable(Empty(), s, nil)

	// use different config to force cleanup
	table2 := NewTable(table, config.Empty(), nil)

	table.Cleanup(table2)

	if len(table2.entries) != 0 {
		t.Fail()
	}

	if !closeCalled {
		t.Fail()
	}
}
