package labels

import (
	"fmt"
	"strings"
)

const (
	RancherPrefix = "io.rancher"

	Revision       = "io.rancher.revision"
	DeploymentUUID = "io.rancher.deployment.uuid"

	SchedulingPrefix               = "io.rancher.scheduler"
	Global                         = "io.rancher.scheduler.global"
	HostAffinity                   = "io.rancher.scheduler.affinity:host_label"
	HostAntiAffinity               = "io.rancher.scheduler.affinity:host_label_ne"
	HostSoftAffinity               = "io.rancher.scheduler.affinity:host_label_soft"
	HostSoftAntiAffinity           = "io.rancher.scheduler.affinity:host_label_soft_ne"
	ContainerLabelAffinity         = "io.rancher.scheduler.affinity:container_label"
	ContainerLabelAntiAffinity     = "io.rancher.scheduler.affinity:container_label_ne"
	ContainerLabelSoftAffinity     = "io.rancher.scheduler.affinity:container_label_soft"
	ContainerLabelSoftAntiAffinity = "io.rancher.scheduler.affinity:container_label_soft_ne"
	ContainerNameAffinity          = "io.rancher.scheduler.affinity:container"
	ContainerNameAntiAffinity      = "io.rancher.scheduler.affinity:container_ne"
	ContainerNameSoftAffinity      = "io.rancher.scheduler.affinity:container_soft"
	ContainerNameSoftAntiAffinity  = "io.rancher.scheduler.affinity:container_soft_ne"

	ServiceLaunchConfig        = "io.rancher.service.launch.config"
	ServicePrimaryLaunchConfig = "io.rancher.service.primary.launch.config"
	RancherDNS                 = "io.rancher.container.dns"
	RancherDNSPriority         = "io.rancher.container.dns.priority"
	RancherDNSSearch           = "io.rancher.container.dnssearch"

	ContainerName        = "io.rancher.container.name"
	PrimaryContainerName = "io.rancher.container.primary"

	ServiceAccount = "io.rancher.kubernetes.service_account"
)

func ParseMap(label interface{}) map[string]string {
	labelMap := map[string]string{}
	kvPairs := strings.Split(fmt.Sprint(label), ",")
	for _, kvPair := range kvPairs {
		kv := strings.SplitN(kvPair, "=", 2)
		if len(kv) > 1 {
			labelMap[kv[0]] = kv[1]
		}
	}
	return labelMap
}

func ParseSlice(label interface{}) []string {
	return strings.Split(fmt.Sprint(label), ",")
}
