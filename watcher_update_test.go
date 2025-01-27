// Copyright 2020 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casbin

import (
	"testing"
)

type SampleWatcherUpdatable struct {
	SampleWatcher
}

func (w SampleWatcherUpdatable) UpdateForUpdatePolicy(sec, ptype string, oldRule, newRule []string) error {
	return nil
}

func (w SampleWatcherUpdatable) UpdateForUpdatePolicies(sec, ptype string, oldRules, newRules [][]string) error {
	return nil
}

func TestSetWatcherUpdatable(t *testing.T) {
	e, _ := NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	sampleWatcherEx := SampleWatcherUpdatable{}
	err := e.SetWatcher(sampleWatcherEx)
	if err != nil {
		t.Fatal(err)
	}

	_ = e.SavePolicy()                                                                                                                                                // calls watcherEx.UpdateForSavePolicy()
	_, _ = e.UpdatePolicy([]string{"admin", "data1", "read"}, []string{"admin", "data2", "read"})                                                                     // calls watcherEx.UpdateForUpdatePolicy()
	_, _ = e.UpdatePolicies([][]string{{"alice", "data1", "read"}, {"alice", "data2", "read"}}, [][]string{{"alice", "data1", "write"}, {"alice", "data2", "write"}}) // calls watcherEx.UpdateForUpdatePolicies()

}
