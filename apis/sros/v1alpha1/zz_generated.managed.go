// +build !ignore_autogenerated
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
// Code generated by ndd-gen. DO NOT EDIT.

package v1alpha1

import nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"

// GetActive of this Registration.
func (mg *Registration) GetActive() bool {
	return mg.Spec.Active
}

// GetCondition of this Registration.
func (mg *Registration) GetCondition(ck nddv1.ConditionKind) nddv1.Condition {
	return mg.Status.GetCondition(ck)
}

// GetDeletionPolicy of this Registration.
func (mg *Registration) GetDeletionPolicy() nddv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetExternalLeafRefs of this Registration.
func (mg *Registration) GetExternalLeafRefs() []string {
	return mg.Status.ExternalLeafRefs
}

// GetNetworkNodeReference of this Registration.
func (mg *Registration) GetNetworkNodeReference() *nddv1.Reference {
	return mg.Spec.NetworkNodeReference
}

// GetResourceIndexes of this Registration.
func (mg *Registration) GetResourceIndexes() map[string]string {
	return mg.Status.ResourceIndexes
}

// GetTarget of this Registration.
func (mg *Registration) GetTarget() []string {
	return mg.Status.Target
}

// SetActive of this Registration.
func (mg *Registration) SetActive(b bool) {
	mg.Spec.Active = b
}

// SetConditions of this Registration.
func (mg *Registration) SetConditions(c ...nddv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Registration.
func (mg *Registration) SetDeletionPolicy(r nddv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetExternalLeafRefs of this Registration.
func (mg *Registration) SetExternalLeafRefs(n []string) {
	mg.Status.ExternalLeafRefs = n
}

// SetNetworkNodeReference of this Registration.
func (mg *Registration) SetNetworkNodeReference(r *nddv1.Reference) {
	mg.Spec.NetworkNodeReference = r
}

// SetResourceIndexes of this Registration.
func (mg *Registration) SetResourceIndexes(n map[string]string) {
	mg.Status.ResourceIndexes = n
}

// SetTarget of this Registration.
func (mg *Registration) SetTarget(t []string) {
	mg.Status.Target = t
}

// GetActive of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetActive() bool {
	return mg.Spec.Active
}

// GetCondition of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetCondition(ck nddv1.ConditionKind) nddv1.Condition {
	return mg.Status.GetCondition(ck)
}

// GetDeletionPolicy of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetDeletionPolicy() nddv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetExternalLeafRefs of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetExternalLeafRefs() []string {
	return mg.Status.ExternalLeafRefs
}

// GetNetworkNodeReference of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetNetworkNodeReference() *nddv1.Reference {
	return mg.Spec.NetworkNodeReference
}

// GetResourceIndexes of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetResourceIndexes() map[string]string {
	return mg.Status.ResourceIndexes
}

// GetTarget of this SrosConfigurePort.
func (mg *SrosConfigurePort) GetTarget() []string {
	return mg.Status.Target
}

// SetActive of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetActive(b bool) {
	mg.Spec.Active = b
}

// SetConditions of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetConditions(c ...nddv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetDeletionPolicy(r nddv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetExternalLeafRefs of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetExternalLeafRefs(n []string) {
	mg.Status.ExternalLeafRefs = n
}

// SetNetworkNodeReference of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetNetworkNodeReference(r *nddv1.Reference) {
	mg.Spec.NetworkNodeReference = r
}

// SetResourceIndexes of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetResourceIndexes(n map[string]string) {
	mg.Status.ResourceIndexes = n
}

// SetTarget of this SrosConfigurePort.
func (mg *SrosConfigurePort) SetTarget(t []string) {
	mg.Status.Target = t
}