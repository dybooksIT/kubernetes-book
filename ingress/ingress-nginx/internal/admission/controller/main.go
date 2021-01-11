/*
Copyright 2019 The Kubernetes Authors.

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

package controller

import (
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	networking "k8s.io/api/networking/v1beta1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/klog/v2"
)

// Checker must return an error if the ingress provided as argument
// contains invalid instructions
type Checker interface {
	CheckIngress(ing *networking.Ingress) error
}

// IngressAdmission implements the AdmissionController interface
// to handle Admission Reviews and deny requests that are not validated
type IngressAdmission struct {
	Checker Checker
}

var (
	ingressResource = metav1.GroupVersionKind{
		Group:   networking.GroupName,
		Version: "v1beta1",
		Kind:    "Ingress",
	}
)

// HandleAdmission populates the admission Response
// with Allowed=false if the Object is an ingress that would prevent nginx to reload the configuration
// with Allowed=true otherwise
func (ia *IngressAdmission) HandleAdmission(obj runtime.Object) (runtime.Object, error) {
	outputVersion := admissionv1.SchemeGroupVersion

	review, isV1 := obj.(*admissionv1.AdmissionReview)

	if !isV1 {
		outputVersion = admissionv1beta1.SchemeGroupVersion
		reviewv1beta1, isv1beta1 := obj.(*admissionv1beta1.AdmissionReview)
		if !isv1beta1 {
			return nil, fmt.Errorf("request is not of type AdmissionReview v1 or v1beta1")
		}

		review = &admissionv1.AdmissionReview{}
		convertV1beta1AdmissionReviewToAdmissionAdmissionReview(reviewv1beta1, review)
	}

	if !apiequality.Semantic.DeepEqual(review.Request.Kind, ingressResource) {
		return nil, fmt.Errorf("rejecting admission review because the request does not contain an Ingress resource but %s with name %s in namespace %s",
			review.Request.Kind.String(), review.Request.Name, review.Request.Namespace)
	}

	status := &admissionv1.AdmissionResponse{}
	status.UID = review.Request.UID

	ingress := networking.Ingress{}

	codec := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{
		Pretty: true,
	})
	codec.Decode(review.Request.Object.Raw, nil, nil)
	_, _, err := codec.Decode(review.Request.Object.Raw, nil, &ingress)
	if err != nil {
		klog.ErrorS(err, "failed to decode ingress")
		status.Allowed = false
		status.Result = &metav1.Status{
			Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
			Message: err.Error(),
		}

		review.Response = status
		return convertResponse(review, outputVersion), nil
	}

	if err := ia.Checker.CheckIngress(&ingress); err != nil {
		klog.ErrorS(err, "invalid ingress configuration", "ingress", fmt.Sprintf("%v/%v", review.Request.Name, review.Request.Namespace))
		status.Allowed = false
		status.Result = &metav1.Status{
			Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
			Message: err.Error(),
		}

		review.Response = status
		return convertResponse(review, outputVersion), nil
	}

	klog.InfoS("successfully validated configuration, accepting", "ingress", fmt.Sprintf("%v/%v", review.Request.Name, review.Request.Namespace))
	status.Allowed = true
	review.Response = status

	return convertResponse(review, outputVersion), nil
}

func convertResponse(review *admissionv1.AdmissionReview, outputVersion schema.GroupVersion) runtime.Object {
	// reply v1
	if outputVersion.Version == admissionv1.SchemeGroupVersion.Version {
		return review
	}

	// reply v1beta1
	reviewv1beta1 := &admissionv1beta1.AdmissionReview{}
	convertAdmissionAdmissionReviewToV1beta1AdmissionReview(review, reviewv1beta1)
	return review
}
