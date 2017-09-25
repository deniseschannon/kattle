package sync

import (
	"k8s.io/client-go/pkg/api/v1"

	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/netes-agent/labels"
)

func primary(d client.DeploymentSyncRequest) client.Container {
	if len(d.Containers) == 1 {
		return d.Containers[0]
	}
	for _, container := range d.Containers {
		if container.LaunchConfigName == labels.ServicePrimaryLaunchConfig {
			return container
		}
	}
	return client.Container{}
}

func primaryContainerNameFromPod(pod v1.Pod) string {
	return pod.Labels[labels.PrimaryContainerName]
}

func trimToLength(s string, size int) string {
	if len(s) > size {
		return s[:size]
	}
	return s
}
