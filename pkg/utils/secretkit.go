// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"

	dplv1alpha1 "github.com/IBM/multicloud-operators-deployable/pkg/apis/app/v1alpha1"
	"github.com/IBM/multicloud-operators-deployable/pkg/utils"
	appv1alpha1 "github.com/IBM/multicloud-operators-subscription/pkg/apis/app/v1alpha1"
)

//PackageSecert put the secret to the deployable template
func PackageSecert(s v1.Secret) *dplv1alpha1.Deployable {
	if klog.V(QuiteLogLel) {
		fnName := GetFnName()
		klog.Infof("Entering: %v()", fnName)

		defer klog.Infof("Exiting: %v()", fnName)
	}

	dpl := &dplv1alpha1.Deployable{}
	dpl.Name = s.GetName()
	dpl.Namespace = s.GetNamespace()
	dpl.Spec.Template = &runtime.RawExtension{}

	sRaw, err := json.Marshal(s)
	dpl.Spec.Template.Raw = sRaw

	if err != nil {
		klog.Error("Failed to unmashall ", s.GetNamespace(), "/", s.GetName(), " err:", err)
	}

	klog.V(10).Infof("Retived Dpl: %v", dpl)

	return dpl
}

//ApplyFilters will apply the subscription level filters to the secret
func ApplyFilters(secret v1.Secret, sub *appv1alpha1.Subscription) (v1.Secret, bool) {
	if klog.V(utils.QuiteLogLel) {
		fnName := utils.GetFnName()
		klog.Infof("Entering: %v()", fnName)

		defer klog.Infof("Exiting: %v()", fnName)
	}

	secret = CleanUpObject(secret)

	if sub.Spec.PackageFilter != nil {
		if sub.Spec.Package != "" && sub.Spec.Package != secret.GetName() {
			klog.Info("Name does not match, skiping:", sub.Spec.Package, "|", secret.GetName())
			return secret, false
		}

		subAnno := sub.GetAnnotations()
		klog.V(10).Info("checking annotations filter:", subAnno)

		if subAnno != nil {
			secretsAnno := secret.GetAnnotations()
			for k, v := range subAnno {
				if secretsAnno[k] != v {
					klog.Info("Annotation filter does not match:", k, "|", v, "|", secretsAnno[k])
					return secret, false
				}
			}
		}
	}

	return secret, true
}

//CleanUpObject is used to reset the sercet fields in order to put the secret into deployable template
func CleanUpObject(s v1.Secret) v1.Secret {
	s.SetResourceVersion("")

	t := types.UID("")

	s.SetUID(t)

	s.SetSelfLink("")

	gvk := schema.GroupVersionKind{
		Kind:    "Secret",
		Version: "v1",
	}
	s.SetGroupVersionKind(gvk)

	return s
}
