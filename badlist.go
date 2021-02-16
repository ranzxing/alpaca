// Copyright 2021 The Alpaca Authors
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
	"fmt"
	"sync"
	"time"
)

// This duration was chosen to match Chrome's behaviour (see "Evaluating proxy lists" in
// https://crsrc.org/net/docs/proxy.md).
const blockDuration = 5 * time.Minute

type badList struct {
	entries []string             // Slice of entries ordered by time added
	times   map[string]time.Time // Map containing the time that each bad entry was added
	now     func() time.Time
	mux     sync.Mutex
}

func newBadList() *badList {
	return &badList{
		entries: []string{},
		times:   make(map[string]time.Time),
		now:     time.Now,
	}
}

func (b *badList) add(entry string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.sweep()
	b.times[entry] = b.now()
	b.entries = append(b.entries, entry)
}

func (b *badList) contains(entry string) bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.sweep()
	_, ok := b.times[entry]
	return ok
}

func (b *badList) sweep() {
	// Delete any stale entries from both the slice and the map. This function is *not*
	// reentrant; `mux` should be locked before calling this function!
	count := 0
	for _, entry := range b.entries {
		then, ok := b.times[entry]
		if !ok {
			// This should never happen.
			panic(fmt.Sprintf("%q is in the bad list with an entry but no time", entry))
		}
		if b.now().Sub(then) < blockDuration {
			break
		}
		delete(b.times, entry)
		count++
	}
	b.entries = b.entries[count:]
}
