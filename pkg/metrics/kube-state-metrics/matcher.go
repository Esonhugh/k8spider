package kube_state_metrics

import (
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/metrics"
)

func DefaultMatchRules() metrics.MatchRules {
	return []*metrics.MetricMatcher{
		metrics.NewMetricMatcher("configmap").AddLabel("namespace").AddLabel("configmap"),
		metrics.NewMetricMatcher("secret").AddLabel("namespace").AddLabel("secret"),

		metrics.NewMetricMatcher("node").AddLabel("node").AddLabel("kernel_version").
			AddLabel("os_image").AddLabel("container_runtime_version").
			AddLabel("provider_id").AddLabel("internal_ip"),
		metrics.NewMetricMatcher("node_role").SetHeader("kube_node_role").
			AddLabel("node").AddLabel("role").SetNameLabel("node"),

		metrics.NewMetricMatcher("pod").AddLabel("namespace").AddLabel("pod").
			AddLabel("node").AddLabel("pod_ip").
			AddLabel("host_ip").AddLabel("host_network"),
		metrics.NewMetricMatcher("container").SetHeader("kube_pod_container_info").
			AddLabel("namespace").AddLabel("pod").
			SetNameLabel("container").AddLabel("container").
			AddLabel("image").AddLabel("image_spec"),
		metrics.NewMetricMatcher("init_container").SetHeader("kube_pod_init_container_info").
			AddLabel("namespace").AddLabel("pod").
			SetNameLabel("container").AddLabel("container").
			AddLabel("image").AddLabel("image_spec"),

		metrics.NewMetricMatcher("cronjob").AddLabel("namespace").AddLabel("cronjob").
			AddLabel("schedule").AddLabel("concurrency_policy"),

		metrics.NewMetricMatcher("service_account").SetHeader("kube_pod_service_account").
			AddLabel("namespace").AddLabel("pod").AddLabel("service_account"),

		metrics.NewMetricMatcher("service").AddLabel("namespace").AddLabel("service").
			AddLabel("cluster_ip").AddLabel("external_name").AddLabel("load_balancer_ip"),
		metrics.NewMetricMatcher("endpoint_address").SetHeader("kube_endpoint_address").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("ip"),
		metrics.NewMetricMatcher("endpoint_port").SetHeader("kube_endpoint_ports").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("port_number"),

		metrics.NewMetricMatcher("persistentvolume").AddLabel("persistentvolume").AddLabel("storageclass").
			AddLabel("gce_persistent_disk_name").AddLabel("ebs_volume_id").AddLabel("azure_disk_name").
			AddLabel("nfs_server").AddLabel("nfs_path").AddLabel("csi_driver").AddLabel("csi_volume_handle").
			AddLabel("local_path").AddLabel("local_fs").AddLabel("host_path").AddLabel("host_path_type"),

		metrics.NewMetricMatcher("validating_webhook").SetNameLabel("webhook_name").
			SetHeader("kube_validatingwebhookconfiguration_webhook_clientconfig_service").
			AddLabel("namespace").AddLabel("webhook_name").
			AddLabel("service_name").AddLabel("service_namespace"),

		metrics.NewMetricMatcher("mutating_webhook").SetNameLabel("webhook_name").
			SetHeader("kube_mutatingwebhookconfiguration_webhook_clientconfig_service").
			AddLabel("namespace").AddLabel("webhook_name").
			AddLabel("service_name").AddLabel("service_namespace"),
	}
}

var NodeMergeHook metrics.ResourceMergeHook = func(m *metrics.MetricMatcher, res define.ResourceList) (r *define.Resource, addFlag bool) {
	if m.Type == "node" || m.Type == "node_role" {
		for i := len(res) - 1; i >= 0; i-- {
			c := res[i]
			if m.FindLabel("node") == c.Name {
				r = res[i]
				return r, false
			}
		}
		return define.NewResource("node"), true
	}
	return nil, true
}

var EndpointMergeHook metrics.ResourceMergeHook = func(m *metrics.MetricMatcher, res define.ResourceList) (r *define.Resource, addFlag bool) {
	if m.Type == "endpoint_address" || m.Type == "endpoint_port" {
		for i := len(res) - 1; i >= 0; i-- {
			c := res[i]
			if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
				r = res[i]
				return r, false
			}
		}
		return define.NewResource("endpoint"), true
		/*
			for i, c := range res {
				if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Type {
					r = res[i]
					return r, false
				}
			}
			return Newmetrics.Resource("endpoint"), true
		*/
	}
	return nil, true
}

func init() {
	metrics.HookList = append(metrics.HookList, NodeMergeHook, EndpointMergeHook)
}
