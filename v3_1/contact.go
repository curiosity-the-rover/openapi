package v3_1

import "github.com/sv-tools/openapi/common"

// Contact information for the exposed API.
//
// https://spec.openapis.org/oas/v3.1.1#contact-object
//
// Example:
//
//	name: API Support
//	url: https://www.example.com/support
//	email: support@example.com
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The URL pointing to the contact information.
	// This MUST be in the form of a URL.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
	// The email address of the contact person/organization.
	// This MUST be in the form of an email address.
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

func (o *Contact) validateSpec(location string, _ *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if err := common.checkURL(o.URL); err != nil {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "url"), err))
	}
	if err := common.checkEmail(o.Email); err != nil {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "email"), err))
	}
	return errs
}

type ContactBuilder struct {
	spec *common.Extendable[Contact]
}

func NewContactBuilder() *ContactBuilder {
	return &ContactBuilder{
		spec: common.NewExtendable(&Contact{}),
	}
}

func (b *ContactBuilder) Build() *common.Extendable[Contact] {
	return b.spec
}

func (b *ContactBuilder) Extensions(v map[string]any) *ContactBuilder {
	b.spec.Extensions = v
	return b
}

func (b *ContactBuilder) AddExt(name string, value any) *ContactBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *ContactBuilder) Name(v string) *ContactBuilder {
	b.spec.Spec.Name = v
	return b
}

func (b *ContactBuilder) URL(v string) *ContactBuilder {
	b.spec.Spec.URL = v
	return b
}

func (b *ContactBuilder) Email(v string) *ContactBuilder {
	b.spec.Spec.Email = v
	return b
}
