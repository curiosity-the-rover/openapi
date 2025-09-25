package v3_1_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sv-tools/openapi/common"
	"github.com/sv-tools/openapi/v3_1"
	"go.yaml.in/yaml/v4"
)

func TestSchema_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name            string
		data            string
		expected        string
		emptyExtensions bool
	}{
		{
			name:            "spec only",
			data:            `{"title": "foo"}`,
			emptyExtensions: true,
		},
		{
			name:            "spec with extension field",
			data:            `{"title": "foo", "b": "bar"}`,
			emptyExtensions: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v *v3_1.Schema
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})
			t.Run("yaml", func(t *testing.T) {
				var v *v3_1.Schema
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				if tt.emptyExtensions {
					require.Empty(t, v.Extensions)
				} else {
					require.NotEmpty(t, v.Extensions)
				}
				data, err := yaml.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.YAMLEq(t, tt.expected, string(data))
			})
		})
	}
}

func TestSchema_AddExt(t *testing.T) {
	for _, tt := range []struct {
		name     string
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:  "without prefix",
			key:   "foo",
			value: 42,
			expected: map[string]any{
				"foo": 42,
			},
		},
		{
			name:  "with prefix",
			key:   "x-foo",
			value: 43,
			expected: map[string]any{
				"x-foo": 43,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ext := v3_1.Schema{}
			ext.AddExt(tt.key, tt.value)
			require.Equal(t, tt.expected, ext.Extensions)
		})
	}
}

func TestSchemaConstWithMatchingDefaultInteger(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("Number", v3_1.NewSchemaBuilder().
			AddType("integer").
			Const(5).
			Default(5).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	require.NoError(t, v.ValidateSpec())
}

func TestSchemaConstWithNonMatchingDefaultError(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("Number", v3_1.NewSchemaBuilder().
			AddType("integer").
			Const(5).
			Default(6).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "expected to be equal to const value")
}

func TestSchemaEnumDuplicateValueError(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("DuplicateEnum", v3_1.NewSchemaBuilder().
			AddType("integer").
			Enum(1, 1).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "duplicate value found in enum")
}

func TestSchemaConstWithEnumConflict(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("ConstEnum", v3_1.NewSchemaBuilder().
			AddType("integer").
			Const(1).
			Enum(1).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "cannot be used together with enum")
}

func TestSchemaFractionalMultipleOfAllowed(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("Price", v3_1.NewSchemaBuilder().
			AddType("number").
			MultipleOf(0.01).
			Minimum(0.0).
			Maximum(100.0).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	require.NoError(t, v.ValidateSpec())
}

func TestSchemaInvalidMaximumLessThanMinimumError(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("Range", v3_1.NewSchemaBuilder().
			AddType("number").
			Minimum(10.0).
			Maximum(5.0).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "must be greater than or equal to minimum")
}

func TestSchemaExclusiveMaximumLessThanExclusiveMinimumError(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("ExclusiveRange", v3_1.NewSchemaBuilder().
			AddType("number").
			ExclusiveMinimum(10.0).
			ExclusiveMaximum(5.0).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "exclusiveMaximum")
}

func TestSchemaDollarSchemaAbsoluteURIOk(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("DialectSchema", v3_1.NewSchemaBuilder().
			Schema("https://example.com/my-dialect").
			AddType("string").
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	require.NoError(t, v.ValidateSpec())
}

func TestSchemaDollarSchemaNonAbsoluteURIError(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("BadDialect", v3_1.NewSchemaBuilder().
			Schema("not-absolute").
			AddType("string").
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "must be an absolute URI")
}

func TestSchema_MaxContainsWithoutContains_Error(t *testing.T) {
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("ArrayNoContainsMax", v3_1.NewSchemaBuilder().
			AddType("array").
			MaxContains(2).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "'contains' keyword is required")
}

func TestOperationID_Unique_OK(t *testing.T) {
	path1 := v3_1.NewPathItemBuilder().Get(v3_1.NewOperationBuilder().OperationID("opA").Build()).Build()
	path2 := v3_1.NewPathItemBuilder().Get(v3_1.NewOperationBuilder().OperationID("opB").Build()).Build()
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Paths(v3_1.NewPaths()).
		Build()
	spec.Spec.Paths.Spec.Add("/a", path1)
	spec.Spec.Paths.Spec.Add("/b", path2)
	v, err := common.NewValidator(spec)
	require.NoError(t, err)
	require.NoError(t, v.ValidateSpec())
}

func TestSchema_EnumUniqueness_SliceDuplicate(t *testing.T) {
	dup := []any{"a", 1}
	spec := v3_1.NewOpenAPIBuilder().
		Info(v3_1.NewInfoBuilder().Title("Spec").Version("1.0.0").Build()).
		Components(v3_1.NewComponents()).
		AddComponent("SliceEnum", v3_1.NewSchemaBuilder().
			AddType("array"). // type itself not critical
			Enum(dup, dup).
			Build()).
		Build()
	v, err := common.NewValidator(spec, common.AllowUnusedComponents())
	require.NoError(t, err)
	err = v.ValidateSpec()
	require.ErrorContains(t, err, "duplicate value found in enum")
}
