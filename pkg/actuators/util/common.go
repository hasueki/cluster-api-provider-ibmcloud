/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"fmt"

	ibmcloudproviderv1 "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"
	machoneapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	apicorev1 "k8s.io/api/core/v1"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	controllerRuntimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	credentialsSecretKey = "ibmcloud_api_key"
)

// GetCredentialsSecret returns base64 encoded credential secret data
func GetCredentialsSecret(coreClient controllerRuntimeClient.Client, namespace string, spec ibmcloudproviderv1.IBMCloudMachineProviderSpec) (string, error) {
	if spec.CredentialsSecret == nil {
		return "", nil
	}
	var credentialsSecret apicorev1.Secret

	if err := coreClient.Get(context.Background(), controllerRuntimeClient.ObjectKey{Namespace: namespace, Name: spec.CredentialsSecret.Name}, &credentialsSecret); err != nil {
		if apimachineryerrors.IsNotFound(err) {
			machoneapierrors.InvalidMachineConfiguration("credentials secret %q in namespace %q not found: %v", spec.CredentialsSecret.Name, namespace, err.Error())
		}
		return "", fmt.Errorf("error getting credentials secret %q in namespace %q: %v", spec.CredentialsSecret.Name, namespace, err)
	}
	data, exists := credentialsSecret.Data[credentialsSecretKey]
	if !exists {
		return "", machoneapierrors.InvalidMachineConfiguration("secret %v/%v does not have %q field set. Thus, no credentials applied when creating an instance", namespace, spec.CredentialsSecret.Name, credentialsSecretKey)
	}

	return string(data), nil
}

// UpdateConditionFailed returns provider condition obj for failed machine creation
func UpdateConditionFailed() ibmcloudproviderv1.IBMCloudMachineProviderCondition {
	return ibmcloudproviderv1.IBMCloudMachineProviderCondition{
		Type:   ibmcloudproviderv1.MachineCreated,
		Status: apicorev1.ConditionFalse,
		Reason: ibmcloudproviderv1.MachineCreationFailed,
	}
}

// UpdateConditionSuccess returns provider condition obj for successful machine creation
func UpdateConditionSuccess() ibmcloudproviderv1.IBMCloudMachineProviderCondition {
	return ibmcloudproviderv1.IBMCloudMachineProviderCondition{
		Type:   ibmcloudproviderv1.MachineCreated,
		Status: apicorev1.ConditionTrue,
		Reason: ibmcloudproviderv1.MachineCreationSucceeded,
	}
}
