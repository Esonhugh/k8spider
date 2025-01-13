package metrics

import (
	"errors"

	"github.com/esonhugh/k8spider/define"
)

type MatchRules []*MetricMatcher

func (m MatchRules) Compile() error {
	var err error = nil
	for i := range m {
		e := m[i].Compile()
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

func (m MatchRules) Match(target string) (*MetricMatcher, error) {
	for _, r := range m {
		_, e := r.Match(target)
		if e != nil {
			continue
		} else {
			return r.CopyData(), nil
		}
	}
	return nil, errors.New("no match found")
}

type ResourceMergeHook func(m *MetricMatcher, res define.ResourceList) (r *define.Resource, addFlag bool)

var HookList []ResourceMergeHook

func ConvertToResource(r []*MetricMatcher, hooks ...ResourceMergeHook) []*define.Resource {
	var res []*define.Resource
	if len(hooks) == 0 {
		hooks = append(hooks, HookList...)
	}

	for _, m := range r {
		var resource *define.Resource
		var addFlag = true

		for _, hook := range hooks {
			resource, addFlag = hook(m, res)
			if resource != nil {
				break
			}
		}

		resourceType := m.Type
		if resource != nil {
			resourceType = resource.Type
		}
		if resource == nil && addFlag {
			resource = define.NewResource(resourceType)
		}

		resource.Namespace = m.FindLabel("namespace")
		resource.Name = m.FindLabel(m.LabelNameOfName())

		// merge endpoint_address and endpoint_port
		for _, l := range m.Labels {
			if l.Key != "namespace" && l.Key != m.LabelNameOfName() {
				resource.AddSpec(l.Key, l.Value)
			}
		}
		if addFlag {
			res = append(res, resource)
		}
	}
	return res
}
