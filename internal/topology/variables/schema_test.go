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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Test_convertToAPIExtensionsJSONSchemaProps(t *testing.T) {
	basicSchema := clusterv1.JSONSchemaProps{
		Type: "integer",
	}
	schemaWithMinAndMax := clusterv1.JSONSchemaProps{
		Type:    "integer",
		Minimum: pointer.Int64(1),
		Maximum: pointer.Int64(43),
	}

	tests := []struct {
		name    string
		schema  *clusterv1.JSONSchemaProps
		want    *apiextensions.JSONSchemaProps
		wantErr bool
	}{
		{
			name: "pass for basic schema validation",
			schema: &clusterv1.JSONSchemaProps{
				Type: "integer",
			},
			want: &apiextensions.JSONSchemaProps{
				Type:             basicSchema.Type,
				Format:           basicSchema.Format,
				MaxLength:        basicSchema.MaxLength,
				MinLength:        basicSchema.MinLength,
				Pattern:          basicSchema.Pattern,
				ExclusiveMaximum: basicSchema.ExclusiveMaximum,
				ExclusiveMinimum: basicSchema.ExclusiveMinimum,
			},
		},
		{
			name:   "pass for schema with minimum and maximum",
			schema: &schemaWithMinAndMax,
			want: &apiextensions.JSONSchemaProps{
				Type:             schemaWithMinAndMax.Type,
				Format:           schemaWithMinAndMax.Format,
				MaxLength:        schemaWithMinAndMax.MaxLength,
				MinLength:        schemaWithMinAndMax.MinLength,
				Pattern:          schemaWithMinAndMax.Pattern,
				ExclusiveMaximum: schemaWithMinAndMax.ExclusiveMaximum,
				ExclusiveMinimum: schemaWithMinAndMax.ExclusiveMinimum,
				Minimum:          convertIntToFloatPointer(*schemaWithMinAndMax.Minimum),
				Maximum:          convertIntToFloatPointer(*schemaWithMinAndMax.Maximum),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToAPIExtensionsJSONSchemaProps(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToAPIExtensionsJSONSchemaProps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToAPIExtensionsJSONSchemaProps() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("pass for schema with default and enum", func(t *testing.T) {
		g := NewWithT(t)

		schema := &clusterv1.JSONSchemaProps{
			Default: &apiextensionsv1.JSON{
				Raw: []byte(`"defaultValue"`),
			},
			Enum: []apiextensionsv1.JSON{
				{Raw: []byte(`"enumValue1"`)},
				{Raw: []byte(`"enumValue2"`)},
			},
		}

		got, err := convertToAPIExtensionsJSONSchemaProps(schema)
		g.Expect(err).ToNot(HaveOccurred())

		g.Expect(*got.Default).To(Equal(apiextensions.JSON(`defaultValue`)))
		g.Expect(got.Enum).To(Equal([]apiextensions.JSON{`enumValue1`, `enumValue2`}))
	})
}
func convertIntToFloatPointer(i int64) *float64 {
	f := float64(i)
	return &f
}
