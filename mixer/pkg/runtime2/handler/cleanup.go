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

package handler

import (
	"istio.io/istio/pkg/log"
)

// Cleanup the old table by selectively closing handlers that are not used in the given table.
func (t *Table) Cleanup(current *Table) {
	toCleanup := []entry{}

	for name, oldEntry := range t.entries {
		if currentEntry, found := current.entries[name]; found && currentEntry.signature.Equals(oldEntry.signature) {
			// this entry is still in use. Skip it.
			continue
		}

		if oldEntry.startupError != nil {
			log.Debugf("skipping cleanup of handler with startup error: %s: '%s'", oldEntry.name, oldEntry.startupError)
			continue
		}

		// schedule for cleanup
		toCleanup = append(toCleanup, oldEntry)
	}

	for _, entry := range toCleanup {
		log.Debugf("closing adapter %s/%v", entry.name, entry.handler)
		err := entry.handler.Close()
		if err != nil {
			log.Warnf("error closing adapter: %s/%v: '%v'", entry.name, entry.handler, err)
		}
	}
}
