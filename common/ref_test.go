package common_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sv-tools/openapi/common"
	"github.com/sv-tools/openapi/v3_1"
	"go.yaml.in/yaml/v4"
)

type testRefOrSpec struct {
	A string `json:"a,omitempty" yaml:"a,omitempty"`
	B string `json:"b,omitempty" yaml:"b,omitempty"`
}

func TestNewRefOrSpec(t *testing.T) {
	for _, tt := range []struct {
		ref_or_spec any
		name        string
		nilRef      bool
		nilSpec     bool
	}{
		{
			name:    "empty",
			nilRef:  true,
			nilSpec: true,
		},
		{
			name:        "ref by reference",
			ref_or_spec: &common.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "ref by value",
			ref_or_spec: common.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "string",
			ref_or_spec: "foo",
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "spec by reference",
			ref_or_spec: &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
		{
			name:        "spec by value",
			ref_or_spec: testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o := common.NewRefOrSpec[testRefOrSpec](tt.ref_or_spec)
			require.NotNil(t, o)
			if tt.nilRef {
				require.Nil(t, o.Ref)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case string:
					require.Equal(t, v, o.Ref.Ref)
				case *common.Ref:
					require.Equal(t, v, o.Ref)
				case common.Ref:
					require.Equal(t, &v, o.Ref)
				default:
					t.Fatal("unexpected ref type")
				}
			}
			if tt.nilSpec {
				require.Nil(t, o.Spec)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case *testRefOrSpec:
					require.Equal(t, v, o.Spec)
				case testRefOrSpec:
					require.Equal(t, &v, o.Spec)
				default:
					t.Fatal("unexpected spec type")
				}
			}
		})
	}
}

func TestNewRefOrExtSpec(t *testing.T) {
	for _, tt := range []struct {
		ref_or_spec any
		name        string
		nilRef      bool
		nilSpec     bool
	}{
		{
			name:    "empty",
			nilRef:  true,
			nilSpec: true,
		},
		{
			name:        "ref by reference",
			ref_or_spec: &common.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "ref by value",
			ref_or_spec: common.Ref{Ref: "foo"},
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "string",
			ref_or_spec: "foo",
			nilRef:      false,
			nilSpec:     true,
		},
		{
			name:        "spec by reference",
			ref_or_spec: &testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
		{
			name:        "spec by value",
			ref_or_spec: testRefOrSpec{A: "foo", B: "bar"},
			nilRef:      true,
			nilSpec:     false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o := common.NewRefOrExtSpec[testRefOrSpec](tt.ref_or_spec)
			require.NotNil(t, o)
			if tt.nilRef {
				require.Nil(t, o.Ref)
			} else {
				switch v := tt.ref_or_spec.(type) {
				case string:
					require.Equal(t, v, o.Ref.Ref)
				case *common.Ref:
					require.Equal(t, v, o.Ref)
				case common.Ref:
					require.Equal(t, &v, o.Ref)
				default:
					t.Fatal("unexpected ref type")
				}
			}
			if tt.nilSpec {
				require.Nil(t, o.Spec)
			} else {
				require.IsType(t, &common.Extendable[testRefOrSpec]{}, o.Spec)
				switch v := tt.ref_or_spec.(type) {
				case *testRefOrSpec:
					require.Equal(t, v, o.Spec.Spec)
				case testRefOrSpec:
					require.Equal(t, &v, o.Spec.Spec)
				default:
					t.Fatal("unexpected spec type")
				}
			}
		})
	}
}

func TestRefOrSpec_Marshal_Unmarshal(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     string
		expected string
		nilRef   bool
		nilSpec  bool
	}{
		{
			name:    "ref",
			data:    `{"$ref": "foo"}`,
			nilRef:  false,
			nilSpec: true,
		},
		{
			name:    "spec",
			data:    `{"a": "foo", "b": "bar"}`,
			nilRef:  true,
			nilSpec: false,
		},
		{
			name:     "both",
			data:     `{"$ref": "foo", "a": "foo", "b": "bar"}`,
			expected: `{"$ref": "foo"}`,
			nilRef:   false,
			nilSpec:  true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				var v *common.RefOrSpec[testRefOrSpec]
				require.NoError(t, json.Unmarshal([]byte(tt.data), &v))
				if tt.nilRef {
					require.Nil(t, v.Ref)
				} else {
					require.NotNil(t, v.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, v.Spec)
				} else {
					require.NotNil(t, v.Spec)
				}
				data, err := json.Marshal(&v)
				require.NoError(t, err)
				if tt.expected == "" {
					tt.expected = tt.data
				}
				require.JSONEq(t, tt.expected, string(data))
			})

			t.Run("yaml", func(t *testing.T) {
				var v *common.RefOrSpec[testRefOrSpec]
				require.NoError(t, yaml.Unmarshal([]byte(tt.data), &v))
				if tt.nilRef {
					require.Nil(t, v.Ref)
				} else {
					require.NotNil(t, v.Ref)
				}
				if tt.nilSpec {
					require.Nil(t, v.Spec)
				} else {
					require.NotNil(t, v.Spec)
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

func TestRefOrSpec_GetSpec(t *testing.T) {
	for _, tt := range []struct {
		name   string
		ref    any
		c      *common.Extendable[v3_1.Components]
		exp    any
		expErr string
	}{
		{
			name: "with spec",
			ref:  common.NewRefOrSpec[testRefOrSpec](&testRefOrSpec{A: "foo"}),
			exp:  &testRefOrSpec{A: "foo"},
		},
		{
			name:   "empty",
			ref:    common.NewRefOrSpec[testRefOrSpec](nil),
			expErr: "not found",
		},
		{
			name:   "no components prefix if ref",
			ref:    common.NewRefOrSpec[testRefOrSpec]("fooo"),
			expErr: "is not implemented",
		},
		{
			name:   "no components but with correct ref",
			ref:    common.NewRefOrSpec[testRefOrSpec]("#/components/schemas/Pet"),
			expErr: "components is required",
		},
		{
			name: "correct ref and components",
			ref:  common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet"),
			c: common.NewExtendable((&v3_1.Components{}).
				Add("Pet", common.NewRefOrSpec[v3_1.Schema](&v3_1.Schema{Title: "foo"})),
			),
			exp: &v3_1.Schema{Title: "foo"},
		},
		{
			name: "ref to ref",
			ref:  common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet"),
			c: common.NewExtendable((&v3_1.Components{}).
				Add("Pet", common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet2")).
				Add("Pet2", common.NewRefOrSpec[v3_1.Schema](&v3_1.Schema{Title: "foo"})),
			),
			exp: &v3_1.Schema{Title: "foo"},
		},
		{
			name: "ref to incorrect ref",
			ref:  common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet"),
			c: common.NewExtendable((&v3_1.Components{}).
				Add("Pet", common.NewRefOrSpec[v3_1.Schema]("fooo")),
			),
			expErr: "is not implemented",
		},
		{
			name: "cycle ref",
			ref:  common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet"),
			c: common.NewExtendable((&v3_1.Components{}).
				Add("Pet", common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet2")).
				Add("Pet2", common.NewRefOrSpec[v3_1.Schema]("#/components/schemas/Pet")),
			),
			expErr: "cycle ref",
		},
		{
			name:   "ref to unexpected component",
			ref:    common.NewRefOrSpec[testRefOrSpec]("#/components/test/Pet"),
			c:      common.NewExtendable(&v3_1.Components{}),
			expErr: "unexpected component",
		},
		{
			name: "ref to unexpected component",
			ref:  common.NewRefOrSpec[v3_1.Operation]("#/components/schemas/Pet"),
			c: common.NewExtendable((&v3_1.Components{}).
				Add("Pet", common.NewRefOrSpec[v3_1.Schema](&v3_1.Schema{Title: "foo"})),
			),
			expErr: "expected spec of type",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.ref).MethodByName("GetSpec").Call([]reflect.Value{reflect.ValueOf(tt.c)})
			require.Len(t, val, 2)
			if tt.expErr != "" {
				err, ok := val[1].Interface().(error)
				require.Truef(t, ok, "not error: %+v", val[1].Interface())
				require.ErrorContains(t, err, tt.expErr)
				require.Nil(t, val[0].Interface())
				return
			}
			require.Nil(t, val[1].Interface())
			require.Equal(t, tt.exp, val[0].Interface())
		})
	}
}
