package k8sgo

import (
	"fmt"
	opstreelabsinv1alpha1 "mongodb-operator/api/v1alpha1"
	"mongodb-operator/mongo"
)

// CreateMongoDBMonitoringUser is a method to create a monitoring user for MongoDB
func CreateMongoDBMonitoringUser(cr *opstreelabsinv1alpha1.MongoDB) error {
	logger := logGenerator(cr.ObjectMeta.Name, cr.Namespace, "MongoDB Monitoring User")
	serviceName := fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, "standalone")
	passwordParams := secretsParameters{Name: cr.ObjectMeta.Name, Namespace: cr.Namespace, SecretName: *cr.Spec.MongoDBSecurity.SecretRef.Name}
	password := getMongoDBPassword(passwordParams)
	monitoringPasswordParams := secretsParameters{Name: cr.ObjectMeta.Name, Namespace: cr.Namespace, SecretName: fmt.Sprintf("%s-%s", serviceName, "monitoring")}
	monitoringPassword := getMongoDBPassword(monitoringPasswordParams)
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:27017/", cr.Spec.MongoDBSecurity.MongoDBAdminUser, password, serviceName)
	mongoParams := mongogo.MongoDBParameters{
		MongoURL:  mongoURL,
		Namespace: cr.Namespace,
		Name:      cr.ObjectMeta.Name,
		Password:  monitoringPassword,
	}
	err := mongogo.CreateMonitoringUser(mongoParams)
	if err != nil {
		logger.Error(err, "Unable to create monitoring user in MongoDB")
		return err
	}
	logger.Info("Successfully created the monitoring user")
	return nil
}

// CheckMonitoringUser is a method to check if monitoring user exists in MongoDB
func CheckMonitoringUser(cr *opstreelabsinv1alpha1.MongoDB) bool {
	logger := logGenerator(cr.ObjectMeta.Name, cr.Namespace, "MongoDB Monitoring User")
	serviceName := fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, "standalone")
	passwordParams := secretsParameters{Name: cr.ObjectMeta.Name, Namespace: cr.Namespace, SecretName: *cr.Spec.MongoDBSecurity.SecretRef.Name}
	password := getMongoDBPassword(passwordParams)
	monitoringUser := "monitoring"
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:27017/", cr.Spec.MongoDBSecurity.MongoDBAdminUser, password, serviceName)
	mongoParams := mongogo.MongoDBParameters{
		MongoURL:  mongoURL,
		Namespace: cr.Namespace,
		Name:      cr.ObjectMeta.Name,
		UserName:  &monitoringUser,
	}
	output, err := mongogo.GetMongoDBUser(mongoParams)
	if err != nil {
		return false
	}
	logger.Info("Successfully executed the command to check monitoring user")
	return output
}
