/*
Copyright 2021 The Crossplane Authors.

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

package resource

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tferrors "github.com/crossplane-contrib/terrajet/pkg/terraform/errors"
)

// Condition constants.
const (
	TypeAsyncOperation = "AsyncOperation"

	ReasonApplyFailure   xpv1.ConditionReason = "ApplyFailure"
	ReasonDestroyFailure xpv1.ConditionReason = "DestroyFailure"
	ReasonSuccess        xpv1.ConditionReason = "Success"
)

// LastOperationCondition returns the condition depending on the content
// of the error.
func LastOperationCondition(err error) xpv1.Condition {
	switch {
	case err == nil:
		return xpv1.Condition{
			Type:               TypeAsyncOperation,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             ReasonSuccess,
		}
	case tferrors.IsApplyFailed(err):
		return xpv1.Condition{
			Type:               TypeAsyncOperation,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             ReasonApplyFailure,
			Message:            err.Error(),
		}
	case tferrors.IsDestroyFailed(err):
		return xpv1.Condition{
			Type:               TypeAsyncOperation,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             ReasonDestroyFailure,
			Message:            err.Error(),
		}
	default:
		return xpv1.Condition{
			Type:               "Unknown",
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "Unknown",
			Message:            err.Error(),
		}
	}
}