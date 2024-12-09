package metrics

import "errors"

type MatchRules []*MetricMatcher

func GenerateMatchRules() MatchRules {
	return make(MatchRules, 0)
}

func DefaultMatchRules() MatchRules {
	return []*MetricMatcher{
		NewMetricMatcher("configmap").AddLabel("namespace").AddLabel("configmap"),
		NewMetricMatcher("secret").AddLabel("namespace").AddLabel("secret"),

		NewMetricMatcher("node").AddLabel("node").AddLabel("kernel_version").
			AddLabel("os_image").AddLabel("container_runtime_version").
			AddLabel("provider_id").AddLabel("internal_ip"),
		NewMetricMatcher("node_role").SetHeader("kube_node_role").
			AddLabel("node").AddLabel("role").SetNameLabel("node"),

		NewMetricMatcher("pod").AddLabel("namespace").AddLabel("pod").
			AddLabel("node").AddLabel("pod_ip").
			AddLabel("host_ip").AddLabel("host_network"),
		NewMetricMatcher("container").SetHeader("kube_pod_container_info").
			AddLabel("namespace").AddLabel("pod").
			SetNameLabel("container").AddLabel("container").
			AddLabel("image").AddLabel("image_spec"),
		NewMetricMatcher("init_container").SetHeader("kube_pod_init_container_info").
			AddLabel("namespace").AddLabel("pod").
			SetNameLabel("container").AddLabel("container").
			AddLabel("image").AddLabel("image_spec"),

		NewMetricMatcher("cronjob").AddLabel("namespace").AddLabel("cronjob").
			AddLabel("schedule").AddLabel("concurrency_policy"),

		NewMetricMatcher("service_account").SetHeader("kube_pod_service_account").
			AddLabel("namespace").AddLabel("pod").AddLabel("service_account"),

		NewMetricMatcher("service").AddLabel("namespace").AddLabel("service").
			AddLabel("cluster_ip").AddLabel("external_name").AddLabel("load_balancer_ip"),
		NewMetricMatcher("endpoint_address").SetHeader("kube_endpoint_address").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("ip"),
		NewMetricMatcher("endpoint_port").SetHeader("kube_endpoint_ports").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("port_number"),

		NewMetricMatcher("persistentvolume").AddLabel("persistentvolume").AddLabel("storageclass").
			AddLabel("gce_persistent_disk_name").AddLabel("ebs_volume_id").AddLabel("azure_disk_name").
			AddLabel("nfs_server").AddLabel("nfs_path").AddLabel("csi_driver").AddLabel("csi_volume_handle").
			AddLabel("local_path").AddLabel("local_fs").AddLabel("host_path").AddLabel("host_path_type"),

		NewMetricMatcher("validating_webhook").SetNameLabel("webhook_name").
			SetHeader("kube_validatingwebhookconfiguration_webhook_clientconfig_service").
			AddLabel("namespace").AddLabel("webhook_name").
			AddLabel("service_name").AddLabel("service_namespace"),

		NewMetricMatcher("mutating_webhook").SetNameLabel("webhook_name").
			SetHeader("kube_mutatingwebhookconfiguration_webhook_clientconfig_service").
			AddLabel("namespace").AddLabel("webhook_name").
			AddLabel("service_name").AddLabel("service_namespace"),
	}
}

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
