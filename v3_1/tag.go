package v3_1

import "github.com/sv-tools/openapi/common"

// Tag adds metadata to a single tag that is used by the Operation Object.
// It is not mandatory to have a Tag Object per tag defined in the Operation Object instances.
//
// https://spec.openapis.org/oas/v3.1.1#tag-object
//
// Example:
//
//	name: pet
//	description: Pets operations
type Tag struct {
	// Additional external documentation for this tag.
	ExternalDocs *common.Extendable[ExternalDocs] `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// REQUIRED.
	// The name of the tag.
	Name string `json:"name" yaml:"name"`
	// A description for the tag.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (o *Tag) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if o.Name == "" {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "name"), common.ErrRequired))
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.ExternalDocs.validateSpec(common.joinLoc(location, "externalDocs"), validator)...)
	}
	validator.visited[common.joinLoc("tags", o.Name)] = true
	return errs
}

type TagBuilder struct {
	spec *common.Extendable[Tag]
}

func NewTagBuilder() *TagBuilder {
	return &TagBuilder{
		spec: common.NewExtendable[Tag](&Tag{}),
	}
}

func (b *TagBuilder) Build() *common.Extendable[Tag] {
	return b.spec
}

func (b *TagBuilder) Extensions(v map[string]any) *TagBuilder {
	b.spec.Extensions = v
	return b
}

func (b *TagBuilder) AddExt(name string, value any) *TagBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *TagBuilder) ExternalDocs(v *common.Extendable[ExternalDocs]) *TagBuilder {
	b.spec.Spec.ExternalDocs = v
	return b
}

func (b *TagBuilder) Name(v string) *TagBuilder {
	b.spec.Spec.Name = v
	return b
}

func (b *TagBuilder) Description(v string) *TagBuilder {
	b.spec.Spec.Description = v
	return b
}
