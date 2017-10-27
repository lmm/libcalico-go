// Copyright (c) 2017 Tigera, Inc. All rights reserved.

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

package clientv2

import (
	"context"
	"fmt"

	apiv2 "github.com/projectcalico/libcalico-go/lib/apis/v2"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/names"
	"github.com/projectcalico/libcalico-go/lib/options"
	"github.com/projectcalico/libcalico-go/lib/watch"
)

// WorkloadEndpointInterface has methods to work with WorkloadEndpoint resources.
type WorkloadEndpointInterface interface {
	Create(ctx context.Context, res *apiv2.WorkloadEndpoint, opts options.SetOptions) (*apiv2.WorkloadEndpoint, error)
	Update(ctx context.Context, res *apiv2.WorkloadEndpoint, opts options.SetOptions) (*apiv2.WorkloadEndpoint, error)
	Delete(ctx context.Context, namespace, name string, opts options.DeleteOptions) (*apiv2.WorkloadEndpoint, error)
	Get(ctx context.Context, namespace, name string, opts options.GetOptions) (*apiv2.WorkloadEndpoint, error)
	List(ctx context.Context, opts options.ListOptions) (*apiv2.WorkloadEndpointList, error)
	Watch(ctx context.Context, opts options.ListOptions) (watch.Interface, error)
}

// workloadEndpoints implements WorkloadEndpointInterface
type workloadEndpoints struct {
	client client
}

// Create takes the representation of a WorkloadEndpoint and creates it.  Returns the stored
// representation of the WorkloadEndpoint, and an error, if there is any.
func (r workloadEndpoints) Create(ctx context.Context, res *apiv2.WorkloadEndpoint, opts options.SetOptions) (*apiv2.WorkloadEndpoint, error) {
	if err := r.validate(res); err != nil {
		return nil, err
	}
	r.updateLabelsForStorage(res)
	out, err := r.client.resources.Create(ctx, opts, apiv2.KindWorkloadEndpoint, res)
	if out != nil {
		return out.(*apiv2.WorkloadEndpoint), err
	}
	return nil, err
}

// Update takes the representation of a WorkloadEndpoint and updates it. Returns the stored
// representation of the WorkloadEndpoint, and an error, if there is any.
func (r workloadEndpoints) Update(ctx context.Context, res *apiv2.WorkloadEndpoint, opts options.SetOptions) (*apiv2.WorkloadEndpoint, error) {
	if err := r.validate(res); err != nil {
		return nil, err
	}
	r.updateLabelsForStorage(res)
	out, err := r.client.resources.Update(ctx, opts, apiv2.KindWorkloadEndpoint, res)
	if out != nil {
		return out.(*apiv2.WorkloadEndpoint), err
	}
	return nil, err
}

// Delete takes name of the WorkloadEndpoint and deletes it. Returns an error if one occurs.
func (r workloadEndpoints) Delete(ctx context.Context, namespace, name string, opts options.DeleteOptions) (*apiv2.WorkloadEndpoint, error) {
	out, err := r.client.resources.Delete(ctx, opts, apiv2.KindWorkloadEndpoint, namespace, name)
	if out != nil {
		return out.(*apiv2.WorkloadEndpoint), err
	}
	return nil, err
}

// Get takes name of the WorkloadEndpoint, and returns the corresponding WorkloadEndpoint object,
// and an error if there is any.
func (r workloadEndpoints) Get(ctx context.Context, namespace, name string, opts options.GetOptions) (*apiv2.WorkloadEndpoint, error) {
	out, err := r.client.resources.Get(ctx, opts, apiv2.KindWorkloadEndpoint, namespace, name)
	if out != nil {
		return out.(*apiv2.WorkloadEndpoint), err
	}
	return nil, err
}

// List returns the list of WorkloadEndpoint objects that match the supplied options.
func (r workloadEndpoints) List(ctx context.Context, opts options.ListOptions) (*apiv2.WorkloadEndpointList, error) {
	res := &apiv2.WorkloadEndpointList{}
	if err := r.client.resources.List(ctx, opts, apiv2.KindWorkloadEndpoint, apiv2.KindWorkloadEndpointList, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Watch returns a watch.Interface that watches the NetworkPolicies that match the
// supplied options.
func (r workloadEndpoints) Watch(ctx context.Context, opts options.ListOptions) (watch.Interface, error) {
	return r.client.resources.Watch(ctx, opts, apiv2.KindWorkloadEndpoint)
}

func (r workloadEndpoints) validate(res *apiv2.WorkloadEndpoint) error {
	// Validate the workload endpoint indices and the name match.
	wepids := names.WorkloadEndpointIdentifiers{
		Node:         res.Spec.Node,
		Orchestrator: res.Spec.Orchestrator,
		Endpoint:     res.Spec.Endpoint,
		Workload:     res.Spec.Workload,
		Pod:          res.Spec.Pod,
		ContainerID:  res.Spec.ContainerID,
	}
	expectedName, err := wepids.CalculateWorkloadEndpointName(false)
	if err != nil {
		return err
	}
	if len(res.Name) == 0 {
		// If a name was not specified then we will calculate it on behalf of the caller.
		res.Name = expectedName
		return nil
	}
	if res.Name != expectedName {
		return errors.ErrorValidation{
			ErroredFields: []errors.ErroredField{{
				Name:   "Name",
				Value:  res.Name,
				Reason: fmt.Sprintf("the WorkloadEndpoint name does not match the primary identifiers assigned in the Spec: expected name %s", expectedName),
			}},
		}
	}
	return nil
}

// updateLabelsForStorage updates the set of labels that we persist.  It adds/overrides
// the Namespace and Orchestrator labels which must be set to the correct values and are
// not user configurable.
func (r workloadEndpoints) updateLabelsForStorage(res *apiv2.WorkloadEndpoint) {
	if res.Labels == nil {
		res.Labels = make(map[string]string, 2)
	}
	res.Labels[apiv2.LabelNamespace] = res.Namespace
	res.Labels[apiv2.LabelOrchestrator] = res.Spec.Orchestrator
}