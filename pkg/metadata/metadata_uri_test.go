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

package metadata

import (
	"testing"

	"github.com/apache/camel-k/pkg/util/camel"
	"github.com/apache/camel-k/pkg/util/defaults"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestJava1(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			import org.apache.camel.builder.RouteBuilder;

			public class Sample extends RouteBuilder {
  				@Override
  				public void configure() throws Exception {
		  			from("timer:tick")
		    			.setBody(constant("-\n             r\n             o\n             c\nHello! Camel K\n             s\n             !\n"))
							.to("log:info?skipBodyLineSeparator=false");
  				}
			}
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Contains(t, metadata.FromURIs, "timer:tick")
	assert.Len(t, metadata.FromURIs, 1)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Len(t, metadata.ToURIs, 1)
}

func TestJava2(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			import org.apache.camel.builder.RouteBuilder;

			public class Sample extends RouteBuilder {
  				@Override
  				public void configure() throws Exception {
		  			from("timer:tick")
		    			.setBody(constant("!\n"))
							.to           (

											"log:info?skipBodyLineSeparator=false"

							               )
							.toD("uri:2")
							.toF("uri:%s", "3");
  				}
			}
		`,
		},
		Language: v1alpha1.LanguageJavaSource,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Contains(t, metadata.FromURIs, "timer:tick")
	assert.Len(t, metadata.FromURIs, 1)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Contains(t, metadata.ToURIs, "uri:2")
	assert.Contains(t, metadata.ToURIs, "uri:%s") // resolution not supported yet
	assert.Len(t, metadata.ToURIs, 3)
}

func TestGroovy1(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			
		  	from( "timer:tick")
		    	.setBody().constant("aa")
				.to   ('log:info?skipBodyLineSeparator=false').to(
											'http://url' )
  			
			from("uri:2")
		    	.setBody().constant("aa")
				.to('uri:3')
		`,
		},
		Language: v1alpha1.LanguageGroovy,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Contains(t, metadata.FromURIs, "timer:tick")
	assert.Contains(t, metadata.FromURIs, "uri:2")
	assert.Len(t, metadata.FromURIs, 2)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Contains(t, metadata.ToURIs, "http://url")
	assert.Contains(t, metadata.ToURIs, "uri:3")
	assert.Len(t, metadata.ToURIs, 3)
}

func TestGroovy2(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `			
			rest().get("/")
				.to   ('log:info?skipBodyLineSeparator=false').to( 'http://url' )
						.toD('dyn:1')
						.tony('thisisnot:anuri')
						.toD( "dyn:2")
						.toF( "f:%s", "2")
		`,
		},
		Language: v1alpha1.LanguageGroovy,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Empty(t, metadata.FromURIs)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Contains(t, metadata.ToURIs, "http://url")
	assert.Contains(t, metadata.ToURIs, "dyn:1")
	assert.Contains(t, metadata.ToURIs, "dyn:2")
	assert.Contains(t, metadata.ToURIs, "f:%s") // resolution not supported yet
	assert.Len(t, metadata.ToURIs, 5)
}

func TestXml1(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			<routes>
			    <route id="hello">
        			<from uri="timer:hello?period=3s"/>
        			<setBody>
            			<constant>Hello World!!!</constant>
        			</setBody>
        			<to uri="log:info"/>
					<to uri="log:info2"/>
					<toD uri="log:info3"/>
    			</route>
			</routes>
		`,
		},
		Language: v1alpha1.LanguageXML,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Contains(t, metadata.FromURIs, "timer:hello?period=3s")
	assert.Len(t, metadata.FromURIs, 1)
	assert.Contains(t, metadata.ToURIs, "log:info")
	assert.Contains(t, metadata.ToURIs, "log:info2")
	assert.Contains(t, metadata.ToURIs, "log:info3")
	assert.Len(t, metadata.ToURIs, 3)
}

func TestKotlin1(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			
		  	from( "timer:tick")
		    	.setBody().constant("aa")
				.to   ("log:info?skipBodyLineSeparator=false").to(
											"http://url" )
  			
			from("uri:2")
		    	.setBody().constant("aa")
				.to("uri:3")
				.toD("uri:4")
				.toF("uri:%s", 5)
		`,
		},
		Language: v1alpha1.LanguageKotlin,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Contains(t, metadata.FromURIs, "timer:tick")
	assert.Contains(t, metadata.FromURIs, "uri:2")
	assert.Len(t, metadata.FromURIs, 2)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Contains(t, metadata.ToURIs, "http://url")
	assert.Contains(t, metadata.ToURIs, "uri:3")
	assert.Contains(t, metadata.ToURIs, "uri:4")
	assert.Contains(t, metadata.ToURIs, "uri:%s") // resolution not supported yet
	assert.Len(t, metadata.ToURIs, 5)
}

func TestJavascript1(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name: "test",
			Content: `
			
			rest().get("/")
				.to   ('log:info?skipBodyLineSeparator=false').to( 'http://url' )
				.toD("uri:2")
				.toF("uri:%s", "3") 
		`,
		},
		Language: v1alpha1.LanguageJavaScript,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.Empty(t, metadata.FromURIs)
	assert.Contains(t, metadata.ToURIs, "log:info?skipBodyLineSeparator=false")
	assert.Contains(t, metadata.ToURIs, "http://url")
	assert.Contains(t, metadata.ToURIs, "uri:2")
	assert.Contains(t, metadata.ToURIs, "uri:%s") // resolution not supported yet
	assert.Len(t, metadata.ToURIs, 4)
}

const yamlFlow = `
- steps:
  - kind: "endpoint"
    uri: "timer:tick"
  - kind: "endpoint"
    uri: "log:info"
`

func TestJYamlFlow(t *testing.T) {
	source := v1alpha1.SourceSpec{
		DataSpec: v1alpha1.DataSpec{
			Name:    "test",
			Content: yamlFlow,
		},
		Language: v1alpha1.LanguageYamlFlow,
	}

	catalog := camel.Catalog(defaults.CamelVersion)
	metadata := Extract(catalog, source)

	assert.NotEmpty(t, metadata.FromURIs)
	assert.Contains(t, metadata.FromURIs, "timer:tick")
	assert.Len(t, metadata.FromURIs, 1)

	assert.NotEmpty(t, metadata.ToURIs)
	assert.Contains(t, metadata.ToURIs, "log:info")
	assert.Len(t, metadata.ToURIs, 1)
}
