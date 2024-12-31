package coredns

import "github.com/esonhugh/k8spider/pkg/metrics"

func init() {
	metrics.HookList = append(metrics.HookList, MergeCoreDnsPlugin)
}

func CoreDNSMatchRules() metrics.MatchRules {
	return []*metrics.MetricMatcher{
		metrics.NewMetricMatcher("coredns_plugin").SetHeader("coredns_plugin_enabled").
			SetNameLabel("zone").AddLabel("name").AddLabel("zone"),
	}
}

var MergeCoreDnsPlugin metrics.ResourceMergeHook = func(m *metrics.MetricMatcher, res metrics.ResourceList) (r *metrics.Resource, addFlag bool) {
	if m.Type == "coredns_plugin" {
		for i := len(res) - 1; i >= 0; i-- {
			if m.Type == "coredns_plugin" {
				r = res[i]
				return r, false
			}
		}
		return metrics.NewResource("coredns_plugin"), true
	}
	return nil, true
}
