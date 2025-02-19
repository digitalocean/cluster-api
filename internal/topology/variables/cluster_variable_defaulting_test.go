/*
Copyright 2021 The Kubernetes Authors.

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

package variables

import (
	"testing"

	. "github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Test_DefaultClusterVariables(t *testing.T) {
	tests := []struct {
		name                  string
		clusterClassVariables []clusterv1.ClusterClassVariable
		clusterVariables      []clusterv1.ClusterVariable
		want                  []clusterv1.ClusterVariable
		wantErr               bool
	}{
		{
			name: "Default one variable of each valid type",
			clusterClassVariables: []clusterv1.ClusterClassVariable{
				{
					Name:     "cpu",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "integer",
							Default: &apiextensionsv1.JSON{Raw: []byte(`1`)},
						},
					},
				},
				{
					Name:     "location",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "string",
							Default: &apiextensionsv1.JSON{Raw: []byte(`"us-east"`)},
						},
					},
				},

				{
					Name:     "count",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "number",
							Default: &apiextensionsv1.JSON{Raw: []byte(`0.1`)},
						},
					},
				},
				{
					Name:     "correct",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "boolean",
							Default: &apiextensionsv1.JSON{Raw: []byte(`true`)},
						},
					},
				},
			},
			clusterVariables: []clusterv1.ClusterVariable{},
			want: []clusterv1.ClusterVariable{
				{
					Name: "cpu",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`1`),
					},
				},
				{
					Name: "location",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`"us-east"`),
					},
				},
				{
					Name: "count",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`0.1`),
					},
				},
				{
					Name: "correct",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`true`),
					},
				},
			},
		},
		{
			name: "Don't default variables that are set",
			clusterClassVariables: []clusterv1.ClusterClassVariable{
				{
					Name:     "cpu",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "integer",
							Default: &apiextensionsv1.JSON{Raw: []byte(`1`)},
						},
					},
				},
				{
					Name:     "correct",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "boolean",
							Default: &apiextensionsv1.JSON{Raw: []byte(`true`)},
						},
					},
				},
			},
			clusterVariables: []clusterv1.ClusterVariable{
				{
					Name: "correct",

					// Value is set here and shouldn't be defaulted.
					Value: apiextensionsv1.JSON{
						Raw: []byte(`false`),
					},
				},
			},
			want: []clusterv1.ClusterVariable{
				{
					Name: "cpu",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`1`),
					},
				},
				{
					Name: "correct",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`false`),
					},
				},
			},
		},
		{
			name: "Don't add variables that have no default schema",
			clusterClassVariables: []clusterv1.ClusterClassVariable{
				{
					Name:     "cpu",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type:    "integer",
							Default: &apiextensionsv1.JSON{Raw: []byte(`1`)},
						},
					},
				},
				{
					Name:     "correct",
					Required: true,
					Schema: clusterv1.VariableSchema{
						OpenAPIV3Schema: clusterv1.JSONSchemaProps{
							Type: "boolean",
						},
					},
				},
			},
			clusterVariables: []clusterv1.ClusterVariable{},
			want: []clusterv1.ClusterVariable{
				{
					Name: "cpu",
					Value: apiextensionsv1.JSON{
						Raw: []byte(`1`),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			vars, errList := DefaultClusterVariables(tt.clusterVariables, tt.clusterClassVariables,
				field.NewPath("spec", "topology", "variables"))

			if tt.wantErr {
				g.Expect(errList).NotTo(BeEmpty())
				return
			}
			g.Expect(errList).To(BeEmpty())
			g.Expect(vars).To(ConsistOf(tt.want))
		})
	}
}

func Test_DefaultClusterVariable(t *testing.T) {
	tests := []struct {
		name                 string
		clusterClassVariable *clusterv1.ClusterClassVariable
		clusterVariable      *clusterv1.ClusterVariable
		want                 *clusterv1.ClusterVariable
		wantErr              bool
	}{
		{
			name: "Default integer",
			clusterClassVariable: &clusterv1.ClusterClassVariable{
				Name:     "cpu",
				Required: true,
				Schema: clusterv1.VariableSchema{
					OpenAPIV3Schema: clusterv1.JSONSchemaProps{
						Type:    "integer",
						Default: &apiextensionsv1.JSON{Raw: []byte(`1`)},
					},
				},
			},
			clusterVariable: &clusterv1.ClusterVariable{Name: "cpu"},
			want: &clusterv1.ClusterVariable{
				Name: "cpu",
				Value: apiextensionsv1.JSON{
					Raw: []byte(`1`),
				},
			},
		},
		{
			name: "Default string",
			clusterClassVariable: &clusterv1.ClusterClassVariable{
				Name:     "location",
				Required: true,
				Schema: clusterv1.VariableSchema{
					OpenAPIV3Schema: clusterv1.JSONSchemaProps{
						Type:    "string",
						Default: &apiextensionsv1.JSON{Raw: []byte(`"us-east"`)},
					},
				},
			},
			clusterVariable: &clusterv1.ClusterVariable{Name: "location"},
			want: &clusterv1.ClusterVariable{
				Name: "location",
				Value: apiextensionsv1.JSON{
					Raw: []byte(`"us-east"`),
				},
			},
		},
		{
			name: "Default number",
			clusterClassVariable: &clusterv1.ClusterClassVariable{
				Name:     "location",
				Required: true,
				Schema: clusterv1.VariableSchema{
					OpenAPIV3Schema: clusterv1.JSONSchemaProps{
						Type:    "number",
						Default: &apiextensionsv1.JSON{Raw: []byte(`0.1`)},
					},
				},
			},
			clusterVariable: &clusterv1.ClusterVariable{Name: "location"},
			want: &clusterv1.ClusterVariable{
				Name: "location",
				Value: apiextensionsv1.JSON{
					Raw: []byte(`0.1`),
				},
			},
		},
		{
			name: "Default boolean",
			clusterClassVariable: &clusterv1.ClusterClassVariable{
				Name:     "correct",
				Required: true,
				Schema: clusterv1.VariableSchema{
					OpenAPIV3Schema: clusterv1.JSONSchemaProps{
						Type:    "boolean",
						Default: &apiextensionsv1.JSON{Raw: []byte(`true`)},
					},
				},
			},
			clusterVariable: &clusterv1.ClusterVariable{Name: "location"},
			want: &clusterv1.ClusterVariable{
				Name: "correct",
				Value: apiextensionsv1.JSON{
					Raw: []byte(`true`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			clusterVariable := &clusterv1.ClusterVariable{
				Name: tt.clusterClassVariable.Name,
			}
			errList := defaultClusterVariable(clusterVariable, tt.clusterClassVariable,
				field.NewPath("spec", "topology", "variables"))

			if tt.wantErr {
				g.Expect(errList).NotTo(BeEmpty())
				return
			}
			g.Expect(errList).To(BeEmpty())

			g.Expect(clusterVariable).To(Equal(tt.want))
		})
	}
}
