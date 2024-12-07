package metrics

import (
	"encoding/json"
	"fmt"
	"io"
)

type Resource struct {
	Namespace string              `json:"namespace"`
	Type      string              `json:"type"`
	Name      string              `json:"name"`
	Spec      map[string][]string `json:"spec"`
}

func NewResource(t string) *Resource {
	return &Resource{
		Type: t,
		Spec: make(map[string][]string, 4),
	}
}

func (r *Resource) AddLabelSpec(l Label) {
	if l.Value == "" {
		return
	}
	if _, ok := r.Spec[l.Key]; !ok {
		r.Spec[l.Key] = make([]string, 0)
	}
	r.Spec[l.Key] = append(r.Spec[l.Key], l.Value)
}

func (r *Resource) AddSpec(key string, value string) {
	r.AddLabelSpec(Label{Key: key, Value: value})
}

type ResourceList []*Resource

func (r *Resource) JSON() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (rl *ResourceList) Print(writer ...io.Writer) {
	var w io.Writer = io.MultiWriter(writer...)
	for _, r := range *rl {
		_, _ = fmt.Fprintf(w, "%v\n", r.JSON())
	}
}

type ResourceMergeHook func(m *MetricMatcher, resource ResourceList) (res *Resource, addFlag bool)

var EndpointMergeHook ResourceMergeHook = func(m *MetricMatcher, res ResourceList) (r *Resource, addFlag bool) {
	if m.Name == "endpoint_address" || m.Name == "endpoint_port" {
		for i := len(res) - 1; i >= 0; i-- {
			c := res[i]
			if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
				r = res[i]
				return r, false
			}
		}
		return NewResource("endpoint"), true
		/*
			for i, c := range res {
				if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
					r = res[i]
					return r, false
				}
			}
			return NewResource("endpoint"), true
		*/
	}
	return nil, true
}

func ConvertToResource(r []*MetricMatcher, hooks ...ResourceMergeHook) []*Resource {
	var res []*Resource
	if len(hooks) == 0 {
		hooks = append(hooks, EndpointMergeHook)
	}

	for _, m := range r {
		var resource *Resource
		var addFlag = true

		for _, hook := range hooks {
			resource, addFlag = hook(m, res)
			if resource != nil {
				break
			}
		}

		resourceType := m.Name
		if resource != nil {
			resourceType = resource.Type
		}
		if resource == nil && addFlag {
			resource = NewResource(resourceType)
		}

		resource.Namespace = m.FindLabel("namespace")
		resource.Name = m.FindLabel(resourceType)

		// merge endpoint_address and endpoint_port
		for _, l := range m.Labels {
			if l.Key != "namespace" && l.Key != resourceType {
				resource.AddLabelSpec(l)
			}
		}
		if addFlag {
			res = append(res, resource)
		}
	}
	return res
}
