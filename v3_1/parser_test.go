package v3_1_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sv-tools/openapi/common"
	"github.com/sv-tools/openapi/v3_1"
)

type Simple struct {
	Fs        string `json:"fs,omitempty" openapi:",format:password" yaml:"FS,omitempty"` // json name should be used
	Fi        int    `yaml:"FI,omitempty"`                                                // yaml name shuld be used
	Fb        *bool
	Fbs       []byte            `json:"fbs,omitempty"                                  openapi:"fBS" yaml:"FS,omitempty"` // openapi name should be used
	Fm        map[string]string `openapi:",required,title:Map of strings,addtype:null"`                                   // default field name should be used
	Excluded1 map[string]string `openapi:"-"`
	Excluded2 map[string]string `json:"-"`
	Excluded3 map[string]string `yaml:"-"`
	Fa        any               `openapi:",deprecated"`

	fp string
}

type Complex struct {
	Simple          // anonymous field
	Next   *Complex `json:"Next"` // circular references
}

type SimpleByRef struct {
	S Simple `json:"s" openapi:"s,required,title:Simple By Ref,ref:#/components/schemas/github.com.sv-tools.openapi_test.Simple" yaml:"s"`
}

func TestParseObject(t *testing.T) {
	trueVar := true
	strVar := "foo"
	var nilBool *bool

	for _, tt := range []struct {
		name               string
		obj                any
		expected           *common.RefOrSpec[v3_1.Schema]
		expectedComponents *v3_1.Components
		err                string
	}{
		{
			name: "nil",
			obj:  nil,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.NullType).
				GoType("nil").Build(),
		},
		{
			name: "bool true",
			obj:  true,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.BooleanType).
				GoType("bool").Build(),
		},
		{
			name: "bool false",
			obj:  false,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.BooleanType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to bool true",
			obj:  &trueVar,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.BooleanType, v3_1.NullType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to bool false",
			obj:  &trueVar,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.BooleanType, v3_1.NullType).
				GoType("bool").Build(),
		},
		{
			name: "ptr to *bool nil",
			obj:  nilBool,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.BooleanType, v3_1.NullType).
				GoType("bool").Build(),
		},
		{
			name: "int",
			obj:  42,
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int64Format).
				GoType("int").Build(),
		},
		{
			name: "int8",
			obj:  int8(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int32Format).
				GoType("int8").Build(),
		},
		{
			name: "int32",
			obj:  int32(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int32Format).
				GoType("int32").Build(),
		},
		{
			name: "int64",
			obj:  int64(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int64Format).
				GoType("int64").Build(),
		},
		{
			name: "uint",
			obj:  uint(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int64Format).
				GoType("uint").Build(),
		},
		{
			name: "uint8",
			obj:  uint8(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int32Format).
				GoType("uint8").Build(),
		},
		{
			name: "uint32",
			obj:  uint32(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int32Format).
				GoType("uint32").Build(),
		},
		{
			name: "uint64",
			obj:  uint64(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.IntegerType).
				Format(v3_1.Int64Format).
				GoType("uint64").Build(),
		},
		{
			name: "float32",
			obj:  float32(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.NumberType).
				Format(v3_1.FloatFormat).
				GoType("float32").Build(),
		},
		{
			name: "float64",
			obj:  float64(42),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.NumberType).
				Format(v3_1.DoubleFormat).
				GoType("float64").Build(),
		},
		{
			name: "string",
			obj:  "foo",
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.StringType).
				GoType("string").Build(),
		},
		{
			name: "bytes",
			obj:  []byte("foo"),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.StringType).
				ContentEncoding(v3_1.Base64Encoding).
				GoType("[]byte").Build(),
		},
		{
			name: "map string",
			obj:  map[string]string{"foo": "bar"},
			expected: v3_1.NewSchemaBuilder().Type(v3_1.ObjectType).
				AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.StringType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "map string ref string",
			obj:  map[string]*string{"foo": &strVar},
			expected: v3_1.NewSchemaBuilder().Type(v3_1.ObjectType).
				AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.StringType, v3_1.NullType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "map string int",
			obj:  map[string]int{"foo": 42},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ObjectType).
				AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.IntegerType).
					Format(v3_1.Int64Format).
					GoType("int").Build(),
				)).Build(),
		},
		{
			name: "map string any",
			obj:  map[string]any{"foo": 42, "bar": "baz"},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ObjectType).
				AdditionalProperties(common.NewBoolOrSchema(true)).
				Build(),
		},
		{
			name: "slice int",
			obj:  []int{42},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.IntegerType).
					Format(v3_1.Int64Format).
					GoType("int").Build(),
				)).Build(),
		},
		{
			name: "slice string",
			obj:  []string{"foo"},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.StringType).
					GoType("string").Build(),
				)).Build(),
		},
		{
			name: "slice any",
			obj:  []any{"foo", 42},
			expected: v3_1.NewSchemaBuilder().Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(true)).
				Build(),
		},
		{
			name: "double slice any",
			obj:  [][]any{{"foo", 42}},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.ArrayType).
					Items(common.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "triple slice any",
			obj:  [][][]any{{{"foo", 42}}},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.ArrayType).
					Items(common.NewBoolOrSchema(
						v3_1.NewSchemaBuilder().
							Type(v3_1.ArrayType).
							Items(common.NewBoolOrSchema(true)).
							Build(),
					)).Build(),
				)).Build(),
		},
		{
			name: "map string any",
			obj:  map[string]map[string]any{"xyz": {"foo": 42, "bar": "baz"}},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ObjectType).
				AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AdditionalProperties(common.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "slice map string any",
			obj:  []map[string]any{{"foo": 42, "bar": "baz"}},
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.ArrayType).
				Items(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AdditionalProperties(common.NewBoolOrSchema(true)).
					Build(),
				)).Build(),
		},
		{
			name: "json number",
			obj:  json.Number("42"),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.NumberType).
				GoPackage("encoding/json").GoType("json.Number").Build(),
		},
		{
			name: "json raw",
			obj:  json.RawMessage(`"foo"`),
			expected: v3_1.NewSchemaBuilder().
				Type(v3_1.StringType).
				ContentMediaType("application/json").
				GoPackage("encoding/json").GoType("json.RawMessage").Build(),
		},
		{
			name: "simple struct",
			obj: Simple{
				Fs:  "foo",
				Fi:  42,
				Fb:  &trueVar,
				Fbs: []byte("bar"),
				Fm:  map[string]string{"baz": "qux"},
				Fa:  []any{"435", 42, false},
				fp:  "baz",
			},
			expected: v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Build(),
			expectedComponents: v3_1.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.Simple",
				v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AddProperty("fs", v3_1.NewSchemaBuilder().Type(v3_1.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", v3_1.NewSchemaBuilder().Type(v3_1.IntegerType).Format(v3_1.Int64Format).GoType("int").Build()).
					AddProperty("Fb", v3_1.NewSchemaBuilder().Type(v3_1.BooleanType, v3_1.NullType).GoType("bool").Build()).
					AddProperty("fBS", v3_1.NewSchemaBuilder().Type(v3_1.StringType).ContentEncoding(v3_1.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", v3_1.NewSchemaBuilder().
						Type(v3_1.ObjectType, v3_1.NullType).
						Title("Map of strings").
						AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
							Type(v3_1.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", v3_1.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
		{
			name: "complex struct",
			obj: Complex{
				Simple: Simple{
					Fs:  "foo",
					Fi:  42,
					Fb:  &trueVar,
					Fbs: []byte("bar"),
					Fm:  map[string]string{"baz": "qux"},
					Fa:  []any{"435", 42, false},
					fp:  "baz",
				},
				Next: &Complex{},
			},
			expected: v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Complex").Build(),
			expectedComponents: v3_1.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.Complex",
				v3_1.NewSchemaBuilder().
					AllOf(
						v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Build(),
						v3_1.NewSchemaBuilder().
							Type(v3_1.ObjectType).
							AddProperty("Next", v3_1.NewSchemaBuilder().
								OneOf(
									v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Complex").Build(),
									v3_1.NewSchemaBuilder().Type(v3_1.NullType).Build(),
								).
								Build(),
							).
							Build(),
					).
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Complex").Build(),
			).Add(
				"github.com.sv-tools.openapi_test.Simple",
				v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AddProperty("fs", v3_1.NewSchemaBuilder().Type(v3_1.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", v3_1.NewSchemaBuilder().Type(v3_1.IntegerType).Format(v3_1.Int64Format).GoType("int").Build()).
					AddProperty("Fb", v3_1.NewSchemaBuilder().Type(v3_1.BooleanType, v3_1.NullType).GoType("bool").Build()).
					AddProperty("fBS", v3_1.NewSchemaBuilder().Type(v3_1.StringType).ContentEncoding(v3_1.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", v3_1.NewSchemaBuilder().
						Type(v3_1.ObjectType, v3_1.NullType).
						Title("Map of strings").
						AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
							Type(v3_1.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", v3_1.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
		{
			name: "simple by ref",
			obj: SimpleByRef{
				S: Simple{},
			},
			expected: v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.SimpleByRef").Build(),
			expectedComponents: v3_1.NewComponents().Spec.Add(
				"github.com.sv-tools.openapi_test.SimpleByRef",
				v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AddProperty("s", v3_1.NewSchemaBuilder().Ref("#/components/schemas/github.com.sv-tools.openapi_test.Simple").Title("Simple By Ref").Build()).
					Required("s").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.SimpleByRef").Build(),
			).Add(
				"github.com.sv-tools.openapi_test.Simple",
				v3_1.NewSchemaBuilder().
					Type(v3_1.ObjectType).
					AddProperty("fs", v3_1.NewSchemaBuilder().Type(v3_1.StringType).GoType("string").Format("password").Build()).
					AddProperty("FI", v3_1.NewSchemaBuilder().Type(v3_1.IntegerType).Format(v3_1.Int64Format).GoType("int").Build()).
					AddProperty("Fb", v3_1.NewSchemaBuilder().Type(v3_1.BooleanType, v3_1.NullType).GoType("bool").Build()).
					AddProperty("fBS", v3_1.NewSchemaBuilder().Type(v3_1.StringType).ContentEncoding(v3_1.Base64Encoding).GoType("[]byte").Build()).
					AddProperty("Fm", v3_1.NewSchemaBuilder().
						Type(v3_1.ObjectType, v3_1.NullType).
						Title("Map of strings").
						AdditionalProperties(common.NewBoolOrSchema(v3_1.NewSchemaBuilder().
							Type(v3_1.StringType).
							GoType("string").Build(),
						)).Build(),
					).
					AddProperty("Fa", v3_1.NewSchemaBuilder().
						Deprecated(true).
						GoType("any").Build(),
					).
					AddRequired("Fm").
					GoPackage("github.com/sv-tools/openapi_test").GoType("openapi_test.Simple").Build(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			spec := v3_1.NewOpenAPIBuilder().
				Info(v3_1.NewInfoBuilder().Title("Test").Version("1.0").Build()).
				Components(v3_1.NewComponents()).
				Build()
			schema, err := v3_1.ParseObject(tt.obj, spec.Spec.Components)
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, schema)

			actual, err := schema.Build().MarshalJSON()
			require.NoError(t, err)

			expected, err := tt.expected.MarshalJSON()
			require.NoError(t, err)

			require.JSONEq(t, string(expected), string(actual))

			if tt.expectedComponents != nil {
				actualComponents, err := spec.Spec.Components.MarshalJSON()
				require.NoError(t, err)

				expectedComponents, err := json.Marshal(tt.expectedComponents)
				require.NoError(t, err)

				require.JSONEq(t, string(expectedComponents), string(actualComponents))
			}

			spec.Spec.Components.Spec.Add("test", schema.Build())
			validator, err := common.NewValidator(
				spec,
				common.AllowUnusedComponents(),
			)
			require.NoError(t, err)

			require.NoError(t, validator.ValidateSpec())

			value, err := common.ConvertToJSON(tt.obj)
			require.NoError(t, err)

			pretty, _ := json.MarshalIndent(tt.obj, "", "  ")
			t.Logf("obj: %s", pretty)

			require.NoError(t, validator.ValidateData("#/components/schemas/test", value))
		})
	}
}
