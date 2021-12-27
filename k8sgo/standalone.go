package k8sgo

import (
	"fmt"
	opstreelabsinv1alpha1 "mongodb-operator/api/v1alpha1"
)

// CreateMongoStandaloneSetup
func CreateMongoStandaloneSetup(cr *opstreelabsinv1alpha1.MongoDB) error {
	logger := logGenerator(cr.ObjectMeta.Name, cr.Namespace, "StatefulSet")
	err := CreateOrUpdateStateFul(getMongoDBStandaloneParams(cr))
	if err != nil {
		logger.Error(err, "Cannot create standalone StatefulSet for MongoDB")
		return err
	}
	return nil
}

func getMongoDBStandaloneParams(cr *opstreelabsinv1alpha1.MongoDB) statefulSetParameters {
	replicas := int32(1)
	trueProperty := true
	falseProperty := false
	appName := fmt.Sprintf("cr.ObjectMeta.Name-%s", "standalone")
	labels := map[string]string{
		"app":           appName,
		"mongodb_setup": "standalone",
		"role":          "standalone",
	}
	params := statefulSetParameters{
		StatefulSetMeta: generateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, generateAnnotations()),
		OwnerDef:        mongoAsOwner(cr),
		Namespace:       cr.Namespace,
		ContainerParams: containerParameters{
			Image:           cr.Spec.KubernetesConfig.Image,
			ImagePullPolicy: cr.Spec.KubernetesConfig.ImagePullPolicy,
		},
		Replicas:    &replicas,
		Labels:      labels,
		Annotations: generateAnnotations(),
	}

	if cr.Spec.Storage != nil {
		params.ContainerParams.PersistenceEnabled = &trueProperty
		params.PVCParameters = pvcParameters{
			Name:             appName,
			Namespace:        cr.Namespace,
			Labels:           labels,
			Annotations:      generateAnnotations(),
			StorageSize:      cr.Spec.Storage.StorageSize,
			StorageClassName: cr.Spec.Storage.StorageClassName,
			AccessModes:      cr.Spec.Storage.AccessModes,
		}
	} else {
		params.ContainerParams.PersistenceEnabled = &falseProperty
	}
	return params
}
