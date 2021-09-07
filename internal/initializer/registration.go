/*
Copyright 2021 NDD.

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

package initializer

import (
	"context"

	"github.com/pkg/errors"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	srosv1alpha1 "github.com/yndd/ndd-provider-sros/apis/sros/v1alpha1"
)

const (
	errApplyRegistration = "cannot apply Registration object"
)

// NewRegistrationObject returns a new *RegistrationObject initializer.
func NewRegistrationObject() *RegistrationObject {
	return &RegistrationObject{}
}

// RegistrationObject has the initializer for creating the Registration object.
type RegistrationObject struct{}

// Run makes sure Registration object exists.
func (lo *RegistrationObject) Run(ctx context.Context, kube client.Client) error {
	l := &srosv1alpha1.Registration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sros-registrations",
			//OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(&d.ObjectMeta, appsv1.SchemeGroupVersion.WithKind("Deployment")))},
		},
		Spec: srosv1alpha1.RegistrationSpec{
			ResourceSpec: nddv1.ResourceSpec{
				NetworkNodeReference: &nddv1.Reference{
					Name: "all",
				},
				DeletionPolicy: nddv1.DeletionDelete,
				Active:         true,
			},
			ForNetworkNode: srosv1alpha1.RegistrationParameters{
				Subscriptions: []string{
					"/acl",
					"/bfd",
					"/interface",
					"/network-instance",
					"/platform",
					"/qos",
					"/routing-policy",
					"/tunnel",
					"/tunnel-interface",
					"/system/snmp",
					"/system/sflow",
					"/system/ntp",
					"/system/network-instance",
					"/system/name",
					"/system/mtu",
					"/system/maintenance",
					"/system/lldp",
					"/system/lacp",
					"/system/authentication",
					"/system/banner",
					"/system/bridge-table",
					"/system/ftp-server",
					"/system/ip-load-balancing",
					"/system/json-rpc-server",
				},
				ExceptionPaths: []string{
					"/interface[name=mgmt0]",
					"/network-instance[name=mgmt]",
					"/system/gnmi-server",
					"/system/tls",
					"/system/ssh-server",
					"/system/aaa",
					"/system/logging",
					"/acl/cpm-filter",
				},
				ExplicitExceptionPaths: []string{
					"/acl",
					"/bfd",
					"/platform",
					"/qos",
					"/routing-policy",
					"/system",
					"/tunnel",
				},
			},
		},
	}
	return errors.Wrap(resource.NewAPIPatchingApplicator(kube).Apply(ctx, l), errApplyRegistration)
}
