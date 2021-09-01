// Copyright 2021 PingCAP, Inc.
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

package bgtask

import (
	"sync"
	"time"
)

// Stat summarizes a background task
type Stat struct {
	TaskName       string
	LastStartTime  time.Time
	LastFinishTime time.Time
	LastError      error
	LastErrorTime  time.Time
	ExecCount      int64
	SumErrors      int64
	SumLatency     int64
	MaxLatency     int64
	MinLatency     int64
}

// instanceStats provide background task stats for a tidb-server instance.
type instanceStats struct {
	sync.Mutex
	tasks map[string]*Stat
}

var stats instanceStats

// Stats summarizes statistics for infoschema
func Stats() map[string]*Stat {
	stats.Lock()
	defer stats.Unlock()
	// return a copy
	newMap := make(map[string]*Stat, len(stats.tasks))
	for k, v := range stats.tasks {
		newMap[k] = v // HACK, fixme!
	}
	return newMap
}

// Start records runtime statistics for a background task
// The task should be observed as in progress, so the finishtime and error code is nil
func Start(name string) {
	stats.Lock()
	defer stats.Unlock()
	if stats.tasks == nil {
		stats.tasks = make(map[string]*Stat)
	}
	if _, ok := stats.tasks[name]; !ok {
		stats.tasks[name] = &Stat{TaskName: name}
	}
	// ExecCount is added when the task finishes to that dividing SumLatency/ExecCount = AvgLatency.
	stats.tasks[name].LastStartTime = time.Now()
	// by convention set it to a zero time if it's in progress.
	stats.tasks[name].LastFinishTime = time.Time{}
}

// Finish records the finish time and the error code
// The mutex helps ensure the change is observed atomically
func Finish(name string, err error) {
	stats.Lock()
	defer stats.Unlock()
	s := stats.tasks[name] // reference
	s.LastFinishTime = time.Now()
	s.ExecCount++
	if err != nil {
		s.LastError = err
		s.LastErrorTime = s.LastFinishTime // TODO: should this be start time?
		s.SumErrors++
	}

	// bgtasks only supports serial execution
	// so this is expected to be safe
	duration := int64(s.LastFinishTime.Sub(s.LastStartTime))
	s.SumLatency += duration
	if duration > s.MaxLatency {
		s.MaxLatency = duration
	}
	if duration < s.MinLatency || s.MinLatency == 0 {
		s.MinLatency = duration
	}
}
