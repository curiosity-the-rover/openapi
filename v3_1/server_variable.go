package v3_1

import "github.com/sv-tools/openapi/common"

// ServerVariable is an object representing a Server Variable for server URL template substitution.
//
// https://spec.openapis.org/oas/v3.1.1#server-variable-object
type ServerVariable struct {
	// REQUIRED.
	// The default value to use for substitution, which SHALL be sent if an alternate value is not supplied.
	// Note this behavior is different than the Schema Object’s treatment of default values,
	// because in those cases parameter values are optional.
	// If the enum is defined, the value MUST exist in the enum’s values.
	Default string `json:"default" yaml:"default"`
	// An optional description for the server variable.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// An enumeration of string values to be used if the substitution options are from a limited set.
	// The array MUST NOT be empty.
	Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
}

func (o *ServerVariable) validateSpec(location string, _ *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if o.Default == "" {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "default"), common.ErrRequired))
	}
	return errs
}

type ServerVariableBuilder struct {
	spec *common.Extendable[ServerVariable]
}

func NewServerVariableBuilder() *ServerVariableBuilder {
	return &ServerVariableBuilder{
		spec: common.NewExtendable[ServerVariable](&ServerVariable{}),
	}
}

func (b *ServerVariableBuilder) Build() *common.Extendable[ServerVariable] {
	return b.spec
}

func (b *ServerVariableBuilder) Extensions(v map[string]any) *ServerVariableBuilder {
	b.spec.Extensions = v
	return b
}

func (b *ServerVariableBuilder) AddExt(name string, value any) *ServerVariableBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *ServerVariableBuilder) Default(v string) *ServerVariableBuilder {
	b.spec.Spec.Default = v
	return b
}

func (b *ServerVariableBuilder) Description(v string) *ServerVariableBuilder {
	b.spec.Spec.Description = v
	return b
}

func (b *ServerVariableBuilder) Enum(v ...string) *ServerVariableBuilder {
	b.spec.Spec.Enum = v
	return b
}

func (b *ServerVariableBuilder) AddEnum(v ...string) *ServerVariableBuilder {
	b.spec.Spec.Enum = append(b.spec.Spec.Enum, v...)
	return b
}
