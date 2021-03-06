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

package trait

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/builder"
	"github.com/apache/camel-k/pkg/builder/kaniko"
	"github.com/apache/camel-k/pkg/builder/s2i"
	"github.com/apache/camel-k/pkg/util/defaults"
	"github.com/apache/camel-k/pkg/util/kubernetes"
	"github.com/apache/camel-k/pkg/util/test"

	"github.com/stretchr/testify/assert"

	"github.com/scylladb/go-set/strset"
)

func TestBuilderTraitNotAppliedBecauseOfNilKit(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I),
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko),
	}

	for _, e := range environments {
		e := e // pin
		e.IntegrationKit = nil

		t.Run(string(e.Platform.Spec.Cluster), func(t *testing.T) {
			err := NewBuilderTestCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.Nil(t, e.GetTrait(ID("builder")))
			assert.Empty(t, e.Steps)
		})
	}
}

func TestBuilderTraitNotAppliedBecauseOfNilPhase(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I),
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko),
	}

	for _, e := range environments {
		e := e // pin
		e.IntegrationKit.Status.Phase = ""

		t.Run(string(e.Platform.Spec.Cluster), func(t *testing.T) {
			err := NewBuilderTestCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.Nil(t, e.GetTrait(ID("builder")))
			assert.Empty(t, e.Steps)
		})
	}
}

func TestS2IBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I)
	err := NewBuilderTestCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.NotNil(t, env.GetTrait(ID("builder")))
	assert.NotEmpty(t, env.Steps)
	assert.Len(t, env.Steps, 7)
	assert.Condition(t, func() bool {
		for _, s := range env.Steps {
			if s == s2i.Steps.Publisher && s.Phase() == builder.ApplicationPublishPhase {
				return true
			}
		}

		return false
	})
}

func TestKanikoBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko)
	err := NewBuilderTestCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.NotNil(t, env.GetTrait(ID("builder")))
	assert.NotEmpty(t, env.Steps)
	assert.Len(t, env.Steps, 7)
	assert.Condition(t, func() bool {
		for _, s := range env.Steps {
			if s == kaniko.Steps.Publisher && s.Phase() == builder.ApplicationPublishPhase {
				return true
			}
		}

		return false
	})
}

func createBuilderTestEnv(cluster v1alpha1.IntegrationPlatformCluster, strategy v1alpha1.IntegrationPlatformBuildPublishStrategy) *Environment {
	c, err := test.DefaultCatalog()
	if err != nil {
		panic(err)
	}

	return &Environment{
		C:            context.TODO(),
		CamelCatalog: c,
		Catalog:      NewCatalog(context.TODO(), nil),
		Integration: &v1alpha1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "ns",
			},
			Status: v1alpha1.IntegrationStatus{
				Phase: v1alpha1.IntegrationPhaseDeploying,
			},
		},
		IntegrationKit: &v1alpha1.IntegrationKit{
			Status: v1alpha1.IntegrationKitStatus{
				Phase: v1alpha1.IntegrationKitPhaseBuildSubmitted,
			},
		},
		Platform: &v1alpha1.IntegrationPlatform{
			Spec: v1alpha1.IntegrationPlatformSpec{
				Cluster: cluster,
				Build: v1alpha1.IntegrationPlatformBuildSpec{
					PublishStrategy: strategy,
					Registry:        v1alpha1.IntegrationPlatformRegistrySpec{Address: "registry"},
					CamelVersion:    defaults.CamelVersionConstraint,
				},
			},
		},
		EnvVars:        make([]corev1.EnvVar, 0),
		ExecutedTraits: make([]Trait, 0),
		Resources:      kubernetes.NewCollection(),
		Classpath:      strset.New(),
	}
}

func NewBuilderTestCatalog() *Catalog {
	return NewCatalog(context.TODO(), nil)
}
