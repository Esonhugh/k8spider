package metrics

import "errors"

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
