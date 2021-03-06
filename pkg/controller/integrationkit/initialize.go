/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integrationkit

import (
	"context"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/platform"
	"github.com/apache/camel-k/pkg/trait"
	"github.com/apache/camel-k/pkg/util/digest"
)

// NewInitializeAction creates a new initialization handling action for the kit
func NewInitializeAction() Action {
	return &initializeAction{}
}

type initializeAction struct {
	baseAction
}

func (action *initializeAction) Name() string {
	return "initialize"
}

func (action *initializeAction) CanHandle(kit *v1alpha1.IntegrationKit) bool {
	return kit.Status.Phase == ""
}

func (action *initializeAction) Handle(ctx context.Context, kit *v1alpha1.IntegrationKit) error {
	// The integration platform needs to be initialized before starting to create kits
	if _, err := platform.GetCurrentPlatform(ctx, action.client, kit.Namespace); err != nil {
		action.L.Info("Waiting for the integration platform to be initialized")
		return nil
	}

	target := kit.DeepCopy()

	_, err := trait.Apply(ctx, action.client, nil, target)
	if err != nil {
		return err
	}

	// Updating the whole integration kit as it may have changed
	action.L.Info("Updating IntegrationKit")
	if err := action.client.Update(ctx, target); err != nil {
		return err
	}

	if target.Spec.Image == "" {
		// by default the kit should be built
		target.Status.Phase = v1alpha1.IntegrationKitPhaseBuildSubmitted
	} else {
		// but in case it has been created from an image, mark the
		// kit as ready
		target.Status.Phase = v1alpha1.IntegrationKitPhaseReady

		// and set the image to be used
		target.Status.Image = target.Spec.Image
	}

	dgst, err := digest.ComputeForIntegrationKit(target)
	if err != nil {
		return err
	}
	target.Status.Digest = dgst

	action.L.Info("IntegrationKit state transition", "phase", target.Status.Phase)
	return action.client.Status().Update(ctx, target)
}
