package sync

import (
	"fmt"
	"path"

	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/netes-agent/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	nodeNameAnnotation  = "node"
	localVolumeLocation = "/var/lib/rancher/volumes"

	accessModeMultiHostRW      = "multiHostRW"
	accessModeSingleHostRW     = "singleHostRW"
	accessModeSingleInstanceRW = "singleInstanceRW"
)

func pvFromVolume(volume client.Volume) (v1.PersistentVolume, error) {
	name := utils.Hash(volume.Id)
	source, err := getVolumeSource(volume)
	if err != nil {
		return v1.PersistentVolume{}, err
	}

	pv := v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PersistentVolumeSpec{
			StorageClassName: name,
			AccessModes: []v1.PersistentVolumeAccessMode{
				getAccessMode(volume),
			},
			PersistentVolumeSource:        source,
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimDelete,
		},
	}

	if volume.SizeMb != 0 {
		size, err := getSize(volume)
		if err != nil {
			return v1.PersistentVolume{}, err
		}
		pv.Spec.Capacity = v1.ResourceList{
			"storage": size,
		}
	}

	return pv, nil
}

func pvFromLocalVolume(volume client.Volume, nodeName string) (v1.PersistentVolume, error) {
	name := utils.Hash(volume.Id)

	pv := v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				nodeNameAnnotation: nodeName,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			StorageClassName: name,
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: path.Join(localVolumeLocation, name),
				},
			},
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimDelete,
		},
	}

	if volume.SizeMb != 0 {
		size, err := getSize(volume)
		if err != nil {
			return v1.PersistentVolume{}, err
		}
		pv.Spec.Capacity = v1.ResourceList{
			"storage": size,
		}
	}

	return pv, nil
}

func pvcFromVolume(volume client.Volume, namespace string) (v1.PersistentVolumeClaim, error) {
	name := utils.Hash(volume.Id)

	claim := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: &name,
			AccessModes: []v1.PersistentVolumeAccessMode{
				getAccessMode(volume),
			},
		},
	}

	if volume.SizeMb != 0 {
		size, err := getSize(volume)
		if err != nil {
			return v1.PersistentVolumeClaim{}, err
		}
		claim.Spec.Resources = v1.ResourceRequirements{
			Requests: v1.ResourceList{
				"storage": size,
			},
		}
	}

	if volume.StorageClass == "" {
		claim.Spec.StorageClassName = &name
	} else {
		claim.Spec.StorageClassName = &volume.StorageClass
	}

	return claim, nil
}

func getSize(volume client.Volume) (resource.Quantity, error) {
	return resource.ParseQuantity(fmt.Sprintf("%dMi", volume.SizeMb))
}

func getVolumeSource(volume client.Volume) (v1.PersistentVolumeSource, error) {
	var source v1.PersistentVolumeSource
	err := utils.ConvertByJSON(volume.PvConfig, &source)
	return source, err
}

func getAccessMode(volume client.Volume) v1.PersistentVolumeAccessMode {
	switch volume.AccessMode {
	case accessModeMultiHostRW:
		return v1.ReadWriteMany
	case accessModeSingleHostRW:
		fallthrough
	case accessModeSingleInstanceRW:
		return v1.ReadWriteOnce
	}
	return v1.PersistentVolumeAccessMode("")
}
