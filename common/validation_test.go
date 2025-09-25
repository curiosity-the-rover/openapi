package common_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sv-tools/openapi/common"
	"github.com/sv-tools/openapi/v3_1"
	"go.yaml.in/yaml/v4"
)

func TestValidator_ValidateSpec(t *testing.T) {
	info, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, f := range info {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(path.Join("testdata", name))
			require.NoError(t, err)
			var o *common.Extendable[v3_1.OpenAPI]
			switch path.Ext(name) {
			case ".yaml":
				require.NoError(t, yaml.Unmarshal(data, &o))
				newData, err := yaml.Marshal(&o)
				require.NoError(t, err)
				require.YAMLEq(t, string(data), string(newData))
			case ".json":
				require.NoError(t, json.Unmarshal(data, &o))
				newData, err := json.Marshal(&o)
				require.NoError(t, err)
				require.JSONEq(t, string(data), string(newData))
			default:
				t.Fatal("wrong file")
			}
			v, err := common.NewValidator(
				o,
				common.AllowUndefinedTagsInOperation(),
				common.ValidateStringDataAsJSON(),
			)
			require.NoError(t, err)
			require.NoError(t, v.ValidateSpec())
		})
	}
}

func TestValidator_ValidateSpec_ManuallyCreated(t *testing.T) {
	for _, tt := range []struct {
		name string
		spec *common.Extendable[v3_1.OpenAPI]
		opts []common.ValidationOption
		err  string
	}{
		{
			name: "info required",
			spec: v3_1.NewOpenAPIBuilder().Build(),
			err:  "/info: required",
		},
		{
			name: "any of path or webhooks or components required",
			spec: v3_1.NewOpenAPIBuilder().Build(),
			err:  "/paths||webhooks||components: required",
		},
		{
			name: "minimal valid with empty paths",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).Paths(v3_1.NewPaths()).Build(),
		},
		{
			name: "minimal valid with empty components",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).Components(v3_1.NewComponents()).Build(),
		},
		{
			name: "minimal valid with empty webhooks",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).WebHooks(v3_1.NewWebhooks()).Build(),
		},
		{
			name: "xml component",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", v3_1.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", v3_1.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					XML(v3_1.NewXMLBuilder().Attribute(true).Build()).
					Build(),
				).
				AddProperty("name", v3_1.NewSchemaBuilder().
					AddType("string").
					XML(v3_1.NewXMLBuilder().
						Namespace("https://example.com/schema/sample").
						Prefix("sample").
						Build(),
					).
					Build(),
				).
				Build(),
			).Build(),
			opts: []common.ValidationOption{common.AllowUnusedComponents()},
		},
		{
			name: "properties examples",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", v3_1.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", v3_1.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Build(),
				).
				AddProperty("name", v3_1.NewSchemaBuilder().
					AddType("string").
					Build(),
				).
				AddExamples(
					map[string]any{
						"id":   123,
						"name": "John Doe 1",
					},
					struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{
						ID:   124,
						Name: "John Doe 2",
					},
				).Build(),
			).Build(),
			opts: []common.ValidationOption{common.AllowUnusedComponents()},
		},
		{
			name: "properties examples error",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", v3_1.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", v3_1.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Build(),
				).
				AddProperty("name", v3_1.NewSchemaBuilder().
					AddType("string").
					Build(),
				).
				AddExamples(
					map[string]any{
						"id":   "123",
						"name": false,
					},
				).Build(),
			).Build(),
			opts: []common.ValidationOption{common.AllowUnusedComponents()},
			err:  "at '/id': got string, want integer",
		},
		{
			name: "properties default",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", v3_1.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", v3_1.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Default(42).
					Build(),
				).
				AddProperty("name", v3_1.NewSchemaBuilder().
					AddType("string").
					Default("John Doe").
					Build(),
				).Build(),
			).Build(),
			opts: []common.ValidationOption{common.AllowUnusedComponents()},
		},
		{
			name: "properties default error",
			spec: v3_1.NewOpenAPIBuilder().Info(
				v3_1.NewInfoBuilder().
					Title("Minimal Valid Spec").
					Version("1.0.0").
					Build(),
			).AddComponent("Person", v3_1.NewSchemaBuilder().
				AddType("object").
				AddProperty("id", v3_1.NewSchemaBuilder().
					AddType("integer").
					Format("int32").
					Default("42").
					Build(),
				).
				AddProperty("name", v3_1.NewSchemaBuilder().
					AddType("string").
					Default(false).
					Build(),
				).Build(),
			).Build(),
			opts: []common.ValidationOption{common.AllowUnusedComponents()},
			err:  "at '': got string, want integer",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, err := common.NewValidator(tt.spec, tt.opts...)
			require.NoError(t, err)

			err = v.ValidateSpec()
			t.Log("error: ", err)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestNewValidator(t *testing.T) {
	data, err := os.ReadFile(path.Join("testdata", "petstore.json"))
	require.NoError(t, err)
	var petStore common.Extendable[v3_1.OpenAPI]
	require.NoError(t, json.Unmarshal(data, &petStore))

	for _, tt := range []struct {
		name string
		spec *common.Extendable[v3_1.OpenAPI]
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			spec: common.NewExtendable(&v3_1.OpenAPI{}),
		},
		{
			name: "petstore",
			spec: &petStore,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := common.NewValidator(tt.spec)
			require.NoError(t, err)
		})
	}
}

func TestValidator_ValidateData(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(path.Join("testdata", "petstore.json"))
	require.NoError(t, err)
	var spec common.Extendable[v3_1.OpenAPI]
	require.NoError(t, json.Unmarshal(data, &spec))
	validator, err := common.NewValidator(&spec)
	require.NoError(t, err)

	for _, tt := range []struct {
		name          string
		ref           string
		data          string
		compileError  string
		validateError string
	}{
		{
			name: "by component",
			ref:  "#/components/schemas/Pet",
			data: `{"id": 123, "name": "foo", "tag": "bar"}`,
		},
		{
			name:          "by component failed",
			ref:           "/components/schemas/Pet",
			data:          `{"id": "123", "name": "foo", "tag": "bar"}`,
			validateError: "got string, want integer",
		},
		{
			name: "by route",
			ref:  "/paths/~1pets~1{petId}/get/responses/200/content/application~1json/schema",
			data: `{"id": 123, "name": "foo", "tag": "bar"}`,
		},
		{
			name:          "by route failed",
			ref:           "/paths/~1pets~1{petId}/get/responses/200/content/application~1json/schema",
			data:          `{"id": "123", "name": "foo", "tag": "bar"}`,
			validateError: "got string, want integer",
		},
		{
			name:         "component not found",
			ref:          "/components/schemas/Fake",
			data:         `{}`,
			compileError: "not found",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var data any
			require.NoError(t, json.Unmarshal([]byte(tt.data), &data))
			err := validator.ValidateData(tt.ref, data)

			if tt.compileError != "" {
				require.ErrorContains(t, err, tt.compileError)
				return
			}

			if tt.validateError != "" {
				require.ErrorContains(t, err, tt.validateError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_UnusedTagsOption(t *testing.T) {
	// Spec with one declared tag but unused
	base := func() *common.Extendable[v3_1.OpenAPI] {
		return v3_1.NewOpenAPIBuilder().
			Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
			Paths(v3_1.NewPaths()).
			Tags(v3_1.NewTagBuilder().Name("unused").Build()).
			Build()
	}

	t.Run("unused tag error by default", func(t *testing.T) {
		v, err := common.NewValidator(base())
		require.NoError(t, err)
		err = v.ValidateSpec()
		require.ErrorContains(t, err, "unused")
	})

	t.Run("unused tag allowed when option set", func(t *testing.T) {
		v, err := common.NewValidator(base(), common.AllowUnusedTags())
		require.NoError(t, err)
		require.NoError(t, v.ValidateSpec())
	})
}
