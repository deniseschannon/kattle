package sync

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/netes-agent/utils"
	"github.com/rancher/netes-agent/watch"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	toBeDeletedAnnotation = "to-be-deleted"
)

func reconcileVolumes(clientset *kubernetes.Clientset, watchClient *watch.Client, deploymentUnit client.DeploymentSyncRequest) error {
	for _, volume := range deploymentUnit.Volumes {
		if volume.StorageClass == "" {
			var pv v1.PersistentVolume
			var err error
			if volume.Local {
				pv, err = pvFromLocalVolume(volume, deploymentUnit.NodeName)
				if err != nil {
					return err
				}
			} else {
				pv, err = pvFromVolume(volume)
				if err != nil {
					return err
				}
			}
			if err := createPv(clientset, pv); err != nil {
				return err
			}
		}

		pvc, err := pvcFromVolume(volume, deploymentUnit.Namespace)
		if err != nil {
			return err
		}
		if err := createPvc(clientset, pvc); err != nil {
			return err
		}
	}
	return nil
}

func createPv(clientset *kubernetes.Clientset, pv v1.PersistentVolume) error {
	log.Infof("Creating persistent volume %s", pv.Name)
	_, err := clientset.PersistentVolumes().Create(&pv)
	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func createPvc(clientset *kubernetes.Clientset, pvc v1.PersistentVolumeClaim) error {
	log.Infof("Creating persistent volume claim %s", pvc.Name)
	_, err := clientset.PersistentVolumeClaims(v1.NamespaceDefault).Create(&pvc)
	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func deletePersistentVolume(clientset *kubernetes.Clientset, volume client.Volume) error {
	name := utils.Hash(volume.Id)
	log.Infof("Deleting persistent volume %s", name)
	if volume.Local {
		pv, err := clientset.PersistentVolumes().Get(name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return err
		}
		pv.Annotations[toBeDeletedAnnotation] = "true"
		_, err = clientset.PersistentVolumes().Update(pv)
		return err
	}
	return clientset.PersistentVolumes().Delete(name, &metav1.DeleteOptions{})
}
