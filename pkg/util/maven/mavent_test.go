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

package maven

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedPom = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <groupId>org.apache.camel.k.integration</groupId>
  <artifactId>camel-k-integration</artifactId>
  <version>1.0.0</version>
  <dependencyManagement>
    <dependencies>
      <dependency>
        <groupId>org.apache.camel</groupId>
        <artifactId>camel-bom</artifactId>
        <version>2.22.1</version>
        <type>pom</type>
        <scope>import</scope>
      </dependency>
    </dependencies>
  </dependencyManagement>
  <dependencies>
    <dependency>
      <groupId>org.apache.camel.k</groupId>
      <artifactId>camel-k-runtime-jvm</artifactId>
      <version>1.0.0</version>
    </dependency>
  </dependencies>
</project>`

func TestPomGeneration(t *testing.T) {
	project := Project{
		XMLName:           xml.Name{Local: "project"},
		XmlNs:             "http://maven.apache.org/POM/4.0.0",
		XmlNsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLocation: "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd",
		ModelVersion:      "4.0.0",
		GroupId:           "org.apache.camel.k.integration",
		ArtifactId:        "camel-k-integration",
		Version:           "1.0.0",
		DependencyManagement: DependencyManagement{
			Dependencies: Dependencies{
				Dependencies: []Dependency{
					{
						GroupId:    "org.apache.camel",
						ArtifactId: "camel-bom",
						Version:    "2.22.1",
						Type:       "pom",
						Scope:      "import",
					},
				},
			},
		},
		Dependencies: Dependencies{
			Dependencies: []Dependency{
				{
					GroupId:    "org.apache.camel.k",
					ArtifactId: "camel-k-runtime-jvm",
					Version:    "1.0.0",
				},
			},
		},
	}

	pom, err := pomFileContent(project)

	assert.Nil(t, err)
	assert.NotNil(t, pom)

	assert.Equal(t, pom, expectedPom)
}

func TestParseSimpleGAV(t *testing.T) {
	dep, err := ParseGAV("org.apache.camel:camel-core:2.21.1")

	assert.Nil(t, err)
	assert.Equal(t, dep.GroupId, "org.apache.camel")
	assert.Equal(t, dep.ArtifactId, "camel-core")
	assert.Equal(t, dep.Version, "2.21.1")
	assert.Equal(t, dep.Type, "jar")
	assert.Equal(t, dep.Classifier, "")
}

func TestParseGAVWithType(t *testing.T) {
	dep, err := ParseGAV("org.apache.camel:camel-core:war:2.21.1")

	assert.Nil(t, err)
	assert.Equal(t, dep.GroupId, "org.apache.camel")
	assert.Equal(t, dep.ArtifactId, "camel-core")
	assert.Equal(t, dep.Version, "2.21.1")
	assert.Equal(t, dep.Type, "war")
	assert.Equal(t, dep.Classifier, "")
}

func TestParseGAVWithClassifierAndType(t *testing.T) {
	dep, err := ParseGAV("org.apache.camel:camel-core:war:test:2.21.1")

	assert.Nil(t, err)
	assert.Equal(t, dep.GroupId, "org.apache.camel")
	assert.Equal(t, dep.ArtifactId, "camel-core")
	assert.Equal(t, dep.Version, "2.21.1")
	assert.Equal(t, dep.Type, "war")
	assert.Equal(t, dep.Classifier, "test")
}
