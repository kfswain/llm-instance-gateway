package framework

import (
	"fmt"
	"reflect"
)

type DataSourceType string

const (
	PoolLevel       DataSourceType = "pool-level"
	RequestSpecific DataSourceType = "request-specific"
)

type BasePlugin interface {
	Consumes() map[string]any
	Name() string
}

type DataSourcePlugin interface {
	BasePlugin
	DataSourceType() DataSourceType
	Produces() map[string]any
}

type MetricsDataSource struct {
}

func (m *MetricsDataSource) Name() string {
	return "metrics-data-source"
}

func (m *MetricsDataSource) DataSourceType() DataSourceType {
	return PoolLevel
}

func (m *MetricsDataSource) Consumes() map[string]any {
	return nil
}

func (m *MetricsDataSource) Produces() map[string]any {
	return map[string]any{
		"lora-info":            LoraInfo{},
		"kv-cache-utilization": float32(0),
		"queued-requests":      int(0),
	}
}

type LoraInfo struct {
}

type DataConsumingPlugin struct {
}

func (s *DataConsumingPlugin) Name() string {
	return "scoring-plugin"
}
func (s *DataConsumingPlugin) Consumes() map[string]any {
	return map[string]any{
		"kv-cache-utilization": float32(0),
	}
}

func produceDataGraph(dataProducingPlugins []DataSourcePlugin, dataConsumingPlugins []DataConsumingPlugin) ([]*DependencyNode, error) {
	// Can probably call this function as a sibling to this one, but for this initial implementation, keeping it here.
	err := validateKeysAndTypes(dataProducingPlugins, dataConsumingPlugins)
	if err != nil {
		return nil, err
	}

	baseSources := []*DependencyNode{}
	sourceAdded := make(map[string]*DependencyNode)

	for len(sourceAdded) < len(dataProducingPlugins) {
		oldSourceCount := len(sourceAdded)
		for _, dsp := range dataProducingPlugins {
			if _, ok := sourceAdded[dsp.Name()]; ok {
				// Already added, skip.
				continue
			} else if dsp.Consumes() == nil || len(dsp.Consumes()) == 0 {
				// No dependencies, can add as root node.
				node := &DependencyNode{
					Consumes:   nil,
					Produces:   dsp.Produces(),
					Dependents: make(map[string]*DependencyNode),
				}
				baseSources = append(baseSources, node)
				sourceAdded[dsp.Name()] = node
			} else {
				// Has dependencies, check if they are all satisfied.
				depsSatisfied := true
				node := DependencyNode{
					Consumes:   make(map[string]*DependencyNode),
					Produces:   dsp.Produces(),
					Dependents: make(map[string]*DependencyNode),
				}
				for dep := range dsp.Consumes() {
					for _, source := range sourceAdded {
						if _, ok := source.Produces[dep]; !ok {
							depsSatisfied = false
							break
						} else {
							node.Consumes[dep] = source
						}
					}
				}
				if depsSatisfied {
					// All dependencies satisfied, apply changes to graph
					for _, source := range node.Consumes {
						source.Dependents[dsp.Name()] = &node
					}
					baseSources = append(baseSources, &node)
					sourceAdded[dsp.Name()] = &node
				}
			}
		}
		if len(sourceAdded) == oldSourceCount {
			return nil, fmt.Errorf("circular or unsatisfiable dependencies detected among data source plugins")
		}
	}

	// Add scoring plugins as dependents to the graph
	for len(sourceAdded) < (len(dataProducingPlugins) + len(dataConsumingPlugins)) {
		oldSourceCount := len(sourceAdded)
		for _, dcp := range dataConsumingPlugins {
			if _, ok := sourceAdded[dcp.Name()]; ok {
				// Already added, skip.
				continue
			} else {
				depsSatisfied := true
				node := DependencyNode{
					Consumes:   make(map[string]*DependencyNode),
					Produces:   nil,
					Dependents: make(map[string]*DependencyNode),
				}
				for dep := range dcp.Consumes() {
					for _, source := range sourceAdded {
						if _, ok := source.Produces[dep]; !ok {
							depsSatisfied = false
							break
						} else {
							node.Consumes[dep] = source
						}
					}
				}
				if depsSatisfied {
					// All dependencies satisfied, apply changes to graph
					for _, source := range node.Consumes {
						source.Dependents[dcp.Name()] = &node
					}
					baseSources = append(baseSources, &node)
					sourceAdded[dcp.Name()] = &node
				}
			}
		}
		if len(sourceAdded) == oldSourceCount {
			return nil, fmt.Errorf("circular or unsatisfiable dependencies detected among data consuming plugins")
		}
	}

	return baseSources, nil
}

func validateKeysAndTypes(dataProducingPlugins []DataSourcePlugin, scoringPlugins []DataConsumingPlugin) error {
	dataProduced := make(map[string]any)
	dataConsumed := make(map[string]any)
	// First pass, validate that there are no conflicts in data types and all consumes have a corresponding producer
	for _, sp := range scoringPlugins {
		for k, v := range sp.Consumes() {
			if _, ok := dataConsumed[k]; !ok {
				dataConsumed[k] = v
			} else if reflect.TypeOf(dataConsumed[k]) != reflect.TypeOf(v) {
				// Error: two dependents on different data types. data type mismatch, bail out with an error message
				return fmt.Errorf("plugins expecting different types for the same key. Key: %s Types: %T , %T", k, dataProduced[k], v)
			}
		}
	}

	for _, dsp := range dataProducingPlugins {
		for k, v := range dsp.Consumes() {
			if _, ok := dataConsumed[k]; !ok {
				dataConsumed[k] = v
			} else if reflect.TypeOf(dataConsumed[k]) != reflect.TypeOf(v) {
				// Error: two dependents on different data types. data type mismatch, bail out with an error message
				return fmt.Errorf("plugins expecting different types for the same key. Key: %s Types: %T , %T", k, dataProduced[k], v)
			}
		}
		for k, v := range dsp.Produces() {
			if _, ok := dataProduced[k]; ok {
				// Error: two data sources producing the same key, bail out with an error message
				return fmt.Errorf("two data sources producing the same key: %s", k)
			}
			dataProduced[k] = v
		}
	}

	for k, v := range dataConsumed {
		if _, ok := dataProduced[k]; !ok {
			// Error: a data source is consuming data that is not produced by any data source. Bail out with an error message
			return fmt.Errorf("a data source is consuming data that is not produced by any data source. Key: %s", k)
		} else if reflect.TypeOf(dataProduced[k]) != reflect.TypeOf(v) {
			// Error: two dependents on different data types. data type mismatch, bail out with an error message
			return fmt.Errorf("Data type produced does not match the expected data type consumed. Key: %s Types: %T , %T", k, dataProduced[k], v)
		}
	}

	return nil
}

type DependencyNode struct {
	Consumes   map[string]*DependencyNode
	Produces   map[string]any
	Dependents map[string]*DependencyNode
}
