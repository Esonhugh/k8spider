package define

import (
	"encoding/json"
	"fmt"
	"io"
)

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value""`
}

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
	for _, v := range r.Spec[l.Key] { // check if the value already exists
		if v == l.Value {
			return
		}
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
