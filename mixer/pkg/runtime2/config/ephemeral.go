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

package config

import (
	"istio.io/istio/mixer/pkg/adapter"
	configpb "istio.io/istio/mixer/pkg/config/proto"
	"istio.io/istio/mixer/pkg/config/store"
	"istio.io/istio/mixer/pkg/expr"
	"istio.io/istio/mixer/pkg/template"
	"istio.io/istio/pkg/log"
)

// Ephemeral configuration state that gets updated by incoming config change events. By itself, the data contained
// is not meaningful. BuildSnapshot must be called to create a new snapshot instance.
type Ephemeral struct {
	// Static information
	adapters  map[string]*adapter.Info // maps adapter shortName to Info.
	templates map[string]*template.Info

	// next snapshot id
	nextID int

	// whether the Attributes have changed since last snapshot
	attributesChanged bool

	// entries that are currently known.
	entries map[store.Key]*store.Resource

	// the latest snapshot.
	latest *Snapshot
}

// NewEphemeral returns a new Ephemeral instance.
func NewEphemeral(
	templates map[string]*template.Info,
	adapters map[string]*adapter.Info) *Ephemeral {

	e := &Ephemeral{
		templates: templates,
		adapters:  adapters,

		nextID: 0,

		attributesChanged: false,
		entries:           make(map[store.Key]*store.Resource, 0),
		latest:            nil,
	}

	// build the initial snapshot.
	_ = e.BuildSnapshot()

	return e
}

// SetState with the supplied state map. All existing ephemeral state is overrwritten.
func (e *Ephemeral) SetState(state map[store.Key]*store.Resource) {
	e.entries = state

	for k := range state {
		if k.Kind == AttributeManifestKind {
			e.attributesChanged = true
			break
		}
	}
}

// ApplyEvents to the internal ephemeral state. This gets called by an store event listener to relay store change
// events to this ephemeral config object.
func (e *Ephemeral) ApplyEvents(events []*store.Event) {

	if log.DebugEnabled() {
		log.Debugf("Incoming config change events: count='%d'", len(events))
	}

	for _, ev := range events {

		if ev.Kind == AttributeManifestKind {
			e.attributesChanged = true
			log.Debug("Received attribute manifest change event.")
		}

		switch ev.Type {
		case store.Update:
			e.entries[ev.Key] = ev.Value
		case store.Delete:
			delete(e.entries, ev.Key)
		}
	}
}

// BuildSnapshot builds a stable, fully-resolved snapshot view of the configuration.
func (e *Ephemeral) BuildSnapshot() *Snapshot {
	id := e.nextID
	e.nextID++

	log.Debugf("Building new config.Snapshot: id='%d'", id)

	attributes := e.processAttributeManifests()

	handlers := e.processHandlerConfigs()

	instances := e.processInstanceConfigs()

	rules := e.processRuleConfigs(handlers, instances)

	e.attributesChanged = false

	s := &Snapshot{
		ID:         id,
		Templates:  e.templates,
		Adapters:   e.adapters,
		Attributes: &attributeFinder{attrs: attributes},
		Handlers:   handlers,
		Instances:  instances,
		Rules:      rules,
	}

	e.latest = s

	log.Infof("Built new config.Snapshot: id='%d'", id)
	if log.DebugEnabled() {
		log.Debugf("config.Snapshot contents:\n%s", s.String())
	}
	return s
}

func (e *Ephemeral) processAttributeManifests() map[string]*configpb.AttributeManifest_AttributeInfo {
	if !e.attributesChanged && e.latest != nil {
		return e.latest.Attributes.attrs
	}

	attrs := make(map[string]*configpb.AttributeManifest_AttributeInfo)
	for k, obj := range e.entries {
		if k.Kind != AttributeManifestKind {
			continue
		}

		log.Debug("Start processing attributes from changed manifest...")

		cfg := obj.Spec
		for an, at := range cfg.(*configpb.AttributeManifest).Attributes {
			attrs[an] = at

			if log.DebugEnabled() {
				log.Debugf("Attribute '%s': '%s'.", an, at.ValueType.String())
			}
		}
	}

	// append all the well known attribute vocabulary from the templates.
	//
	// ATTRIBUTE_GENERATOR variety templates allows operators to write Attributes
	// using the $out.<field Name> convention, where $out refers to the output object from the attribute generating adapter.
	// The list of valid names for a given Template is available in the Template.Info.AttributeManifests object.
	for _, info := range e.templates {
		log.Debugf("Processing attributes from template: '%s'", info.Name)

		for _, v := range info.AttributeManifests {
			for an, at := range v.Attributes {
				attrs[an] = at

				if log.DebugEnabled() {
					log.Debugf("Attribute '%s': '%s'", an, at.ValueType.String())
				}
			}
		}
	}

	log.Debug("Completed processing attributes.")

	return attrs
}

func (e *Ephemeral) processHandlerConfigs() map[string]*Handler {
	configs := make(map[string]*Handler)

	for key, resource := range e.entries {
		var info *adapter.Info
		var found bool
		if info, found = e.adapters[key.Kind]; !found {
			continue
		}

		adapterName := key.String()

		if log.DebugEnabled() {
			log.Debugf("Processing incoming handler config: name='%s'\n%s", adapterName, resource.Spec.String())
		}

		config := &Handler{
			Name:    adapterName,
			Adapter: info,
			Params:  resource.Spec,
		}

		configs[config.Name] = config
	}

	return configs
}

func (e *Ephemeral) processInstanceConfigs() map[string]*Instance {
	configs := make(map[string]*Instance)

	for key, resource := range e.entries {
		var info *template.Info
		var found bool
		if info, found = e.templates[key.Kind]; !found {
			continue
		}

		instanceName := key.String()

		if log.DebugEnabled() {
			log.Debugf("Processing incoming instance config: name='%s'\n%s", instanceName, resource.Spec.String())
		}

		config := &Instance{
			Name:     instanceName,
			Template: info,
			Params:   resource.Spec,
		}

		configs[config.Name] = config
	}

	return configs
}

func (e *Ephemeral) processRuleConfigs(
	handlers map[string]*Handler,
	instances map[string]*Instance) []*Rule {

	log.Debug("Begin processing rule configurations.")

	var configs []*Rule

	for ruleKey, resource := range e.entries {
		if ruleKey.Kind != RulesKind {
			continue
		}

		ruleName := ruleKey.String()

		cfg := resource.Spec.(*configpb.Rule)

		if log.DebugEnabled() {
			log.Debugf("Processing incoming rule: name='%s'\n%s", ruleName, cfg.String())
		}

		// resourceType is used for backwards compatibility with labels: [istio-protocol: tcp]
		rt := resourceType(resource.Metadata.Labels)
		if cfg.Match != "" {
			m, err := expr.ExtractEQMatches(cfg.Match)
			if err != nil {
				log.Warnf("ConfigWarning: Unable to extract resource type from rule: name='%s'", ruleName)
				continue
			}

			if ContextProtocolTCP == m[ContextProtocolAttributeName] {
				rt.protocol = protocolTCP
			}
		}

		var actions []*Action
		for i, a := range cfg.Actions {

			log.Debugf("Processing action: %s[%d]", ruleName, i)

			handlerName := canonicalize(a.Handler, ruleKey.Namespace)
			handler, found := handlers[handlerName]
			if !found {
				log.Warnf("ConfigWarning handler not found: handler='%s', action='%s[%d]'",
					handlerName, ruleName, i)
				continue
			}

			actionInstances := []*Instance{}
			for _, instanceName := range a.Instances {
				instanceName := canonicalize(instanceName, ruleKey.Namespace)
				instance, found := instances[instanceName]
				if !found {
					log.Warnf("ConfigWarning instance not found: instance='%s', action='%s[%d]'",
						instanceName, ruleName, i)
					continue
				}

				actionInstances = append(actionInstances, instance)
			}

			if len(actionInstances) == 0 {
				log.Warnf("ConfigWarning no valid instances found: action='%s[%d]'", ruleName, i)
				continue
			}

			action := &Action{
				Handler:   handler,
				Instances: actionInstances,
			}

			actions = append(actions, action)
		}

		if len(actions) == 0 {
			log.Warnf("ConfigWarning no valid actions found in rule: %s", ruleName)
			continue
		}

		rule := &Rule{
			Name:         ruleName,
			Namespace:    ruleKey.Namespace,
			Actions:      actions,
			ResourceType: rt,
			Match:        cfg.Match,
		}

		configs = append(configs, rule)
	}

	return configs
}

// resourceType maps labels to rule types.
func resourceType(labels map[string]string) ResourceType {
	rt := defaultResourcetype()
	if ContextProtocolTCP == labels[istioProtocol] {
		rt.protocol = protocolTCP
	}
	return rt
}
