// Copyright 2017 The casbin Authors. All Rights Reserved.
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
	"fmt"

	Err "github.com/casbin/casbin/v2/errors"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

const (
	notImplemented = "not implemented"
)

func (e *Enforcer) shouldPersist() bool {
	return e.adapter != nil && e.autoSave
}

// addPolicy adds a rule to the current policy.
func (e *Enforcer) addPolicy(sec string, ptype string, rule []string) (bool, error) {
	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.AddPolicies(sec, ptype, [][]string{rule})
	}

	if e.model.HasPolicy(sec, ptype, rule) {
		return false, nil
	}

	if e.shouldPersist() {
		if err := e.adapter.AddPolicy(sec, ptype, rule); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	err := e.operator.AddPolicy(sec, ptype, rule)
	if err != nil {
		return true, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForAddPolicy(sec, ptype, rule...)
		} else {
			err = e.watcher.Update()
		}
		return true, err
	}

	return true, nil
}

// addPolicies adds rules to the current policy.
func (e *Enforcer) addPolicies(sec string, ptype string, rules [][]string) (bool, error) {
	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.AddPolicies(sec, ptype, rules)
	}

	if e.model.HasPolicies(sec, ptype, rules) {
		return false, nil
	}

	if e.shouldPersist() {
		if err := e.adapter.(persist.BatchAdapter).AddPolicies(sec, ptype, rules); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	err := e.operator.AddPolicies(sec, ptype, rules)
	if err != nil {
		return true, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForAddPolicies(sec, ptype, rules...)
		} else {
			err = e.watcher.Update()
		}
		return true, err
	}

	return true, nil
}

// removePolicy removes a rule from the current policy.
func (e *Enforcer) removePolicy(sec string, ptype string, rule []string) (bool, error) {
	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.RemovePolicies(sec, ptype, [][]string{rule})
	}

	if e.shouldPersist() {
		if err := e.adapter.RemovePolicy(sec, ptype, rule); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	ruleRemoved, err := e.operator.RemovePolicy(sec, ptype, rule)
	if err != nil {
		return ruleRemoved, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForRemovePolicy(sec, ptype, rule...)
		} else {
			err = e.watcher.Update()
		}
		return ruleRemoved, err

	}

	return ruleRemoved, nil
}

func (e *Enforcer) updatePolicy(sec string, ptype string, oldRule []string, newRule []string) (bool, error) {
	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.UpdatePolicy(sec, ptype, oldRule, newRule)
	}

	if e.shouldPersist() {
		if err := e.adapter.(persist.UpdatableAdapter).UpdatePolicy(sec, ptype, oldRule, newRule); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	ruleUpdated, err := e.operator.UpdatePolicy(sec, ptype, oldRule, newRule)
	if err != nil {
		return ruleUpdated, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherUpdatable); ok {
			err = watcher.UpdateForUpdatePolicy(sec, ptype, oldRule, newRule)
		} else {
			err = e.watcher.Update()
		}
		return ruleUpdated, err
	}

	return ruleUpdated, nil
}

func (e *Enforcer) updatePolicies(sec string, ptype string, oldRules [][]string, newRules [][]string) (bool, error) {
	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.UpdatePolicies(sec, ptype, oldRules, newRules)
	}

	if e.shouldPersist() {
		if err := e.adapter.(persist.UpdatableAdapter).UpdatePolicies(sec, ptype, oldRules, newRules); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	ruleUpdated, err := e.operator.UpdatePolicies(sec, ptype, oldRules, newRules)
	if err != nil {
		return ruleUpdated, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherUpdatable); ok {
			err = watcher.UpdateForUpdatePolicies(sec, ptype, oldRules, newRules)
		} else {
			err = e.watcher.Update()
		}
		return ruleUpdated, err
	}

	return ruleUpdated, nil
}

// removePolicies removes rules from the current policy.
func (e *Enforcer) removePolicies(sec string, ptype string, rules [][]string) (bool, error) {
	if !e.model.HasPolicies(sec, ptype, rules) {
		return false, nil
	}

	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.RemovePolicies(sec, ptype, rules)
	}

	if e.shouldPersist() {
		if err := e.adapter.(persist.BatchAdapter).RemovePolicies(sec, ptype, rules); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	rulesRemoved, err := e.operator.RemovePolicies(sec, ptype, rules)
	if err != nil {
		return rulesRemoved, err
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, rules)
		if err != nil {
			return rulesRemoved, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForRemovePolicies(sec, ptype, rules...)
		} else {
			err = e.watcher.Update()
		}
		return true, err
	}

	return rulesRemoved, nil
}

// removeFilteredPolicy removes rules based on field filters from the current policy.
func (e *Enforcer) removeFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	if len(fieldValues) == 0 {
		return false, Err.INVALID_FIELDVAULES_PARAMETER
	}

	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
	}

	if e.shouldPersist() {
		if err := e.adapter.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	ruleRemoved, err := e.operator.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
	if err != nil {
		return ruleRemoved, nil
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForRemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
		} else {
			err = e.watcher.Update()
		}
		return ruleRemoved, err
	}

	return ruleRemoved, nil
}

func (e *Enforcer) updateFilteredPolicies(sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) (bool, error) {
	var (
		oldRules [][]string
		err      error
	)

	if e.shouldPersist() {
		if oldRules, err = e.adapter.(persist.UpdatableAdapter).UpdateFilteredPolicies(sec, ptype, newRules, fieldIndex, fieldValues...); err != nil {
			if err.Error() != notImplemented {
				return false, err
			}
		}
	}

	if e.dispatcher != nil && e.autoNotifyDispatcher {
		return true, e.dispatcher.UpdateFilteredPolicies(sec, ptype, oldRules, newRules)
	}

	ruleChanged, err := e.operator.UpdatePolicies(sec, ptype, oldRules, newRules)
	if err != nil {
		return ruleChanged, err
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherUpdatable); ok {
			err = watcher.UpdateForUpdatePolicies(sec, ptype, oldRules, newRules)
		} else {
			err = e.watcher.Update()
		}
		return ruleChanged, err
	}

	return ruleChanged, nil
}

func (e *Enforcer) getDomainIndex(ptype string) int {
	p := e.model["p"][ptype]
	pattern := fmt.Sprintf("%s_dom", ptype)
	index := len(p.Tokens)
	for i, token := range p.Tokens {
		if token == pattern {
			index = i
			break
		}
	}
	return index
}

// PolicyOperator encapsulates policy mutation logic
// This can be used in Watcher callback implementors to share common policy mutation logic which is otherwise internal to the Enforcer.
// PolicyOperator is internal to Casbin, it should be used only with Watcher callbacks.
type PolicyOperator struct {
	e *Enforcer
}

// AddPolicy adds a rule to the current policy.
func (p *PolicyOperator) AddPolicy(sec string, ptype string, rule []string) error {
	p.e.model.AddPolicy(sec, ptype, rule)

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, [][]string{rule})
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPolicies adds rules to the current policy.
func (p *PolicyOperator) AddPolicies(sec string, ptype string, rules [][]string) error {
	p.e.model.AddPolicies(sec, ptype, rules)

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, rules)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemovePolicy removes a rule from the current policy.
func (p *PolicyOperator) RemovePolicy(sec string, ptype string, rule []string) (bool, error) {
	ruleRemoved := p.e.model.RemovePolicy(sec, ptype, rule)
	if !ruleRemoved {
		return ruleRemoved, nil
	}

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, [][]string{rule})
		if err != nil {
			return ruleRemoved, err
		}
	}
	return ruleRemoved, nil
}

// UpdatePolicy updates a rule in the current polic
func (p *PolicyOperator) UpdatePolicy(sec string, ptype string, oldRule []string, newRule []string) (bool, error) {
	ruleUpdated := p.e.model.UpdatePolicy(sec, ptype, oldRule, newRule)
	if !ruleUpdated {
		return ruleUpdated, nil
	}

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, [][]string{oldRule}) // remove the old rule
		if err != nil {
			return ruleUpdated, err
		}
		err = p.e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, [][]string{newRule}) // add the new rule
		if err != nil {
			return ruleUpdated, err
		}
	}
	return ruleUpdated, nil
}

// UpdatePolicies updates rules in the current policy
func (p *PolicyOperator) UpdatePolicies(sec string, ptype string, oldRules [][]string, newRules [][]string) (bool, error) {
	ruleUpdated := p.e.model.UpdatePolicies(sec, ptype, oldRules, newRules)
	if !ruleUpdated {
		return ruleUpdated, nil
	}

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, oldRules) // remove the old rules
		if err != nil {
			return ruleUpdated, err
		}
		err = p.e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, newRules) // add the new rules
		if err != nil {
			return ruleUpdated, err
		}
	}
	return ruleUpdated, nil
}

// RemovePolicies removes rules from the current policy.
func (p *PolicyOperator) RemovePolicies(sec string, ptype string, rules [][]string) (bool, error) {
	rulesRemoved := p.e.model.RemovePolicies(sec, ptype, rules)
	if !rulesRemoved {
		return rulesRemoved, nil
	}

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, rules)
		if err != nil {
			return rulesRemoved, err
		}
	}
	return rulesRemoved, nil
}

// RemoveFilteredPolicy removes rules based on field filters from the current policy.
func (p *PolicyOperator) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	ruleRemoved, effects := p.e.model.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
	if !ruleRemoved {
		return ruleRemoved, nil
	}

	if sec == "g" {
		err := p.e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, effects)
		if err != nil {
			return ruleRemoved, err
		}
	}
	return ruleRemoved, nil
}
