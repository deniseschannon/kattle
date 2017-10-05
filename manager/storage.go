package manager

import (
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/netes-agent/sync"
	"github.com/rancher/netes-agent/utils"
)

func (m *Manager) handleVolumeRemove(event *events.Event, rancherClient Client) (*client.Publish, error) {
	var volume client.Volume
	if err := utils.ConvertByJSON(event.Data["volume"], &volume); err != nil {
		return nil, err
	}

	clientset, watchClient, err := m.getCluster(rancherClient, volume.ClusterId)
	if err != nil {
		return emptyReply(event), nil
	}

	return emptyReply(event), sync.RemoveVolume(clientset, watchClient, volume)
}
