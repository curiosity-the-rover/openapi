package v3_1

import "github.com/sv-tools/openapi/common"

// License information for the exposed API.
//
// https://spec.openapis.org/oas/v3.1.1#license-object
//
// Example:
//
//	name: Apache 2.0
//	identifier: Apache-2.0
type License struct {
	// REQUIRED.
	// The license name used for the API.
	Name string `json:"name" yaml:"name"`
	// An SPDX license expression for the API.
	// The identifier field is mutually exclusive of the url field.
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	// A URL to the license used for the API.
	// This MUST be in the form of a URL.
	// The url field is mutually exclusive of the identifier field.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (o *License) validateSpec(location string, _ *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if o.Name == "" {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "name"), common.ErrRequired))
	}
	if o.Identifier != "" && o.URL != "" {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "identifier&url"), common.ErrMutuallyExclusive))
	}
	if err := common.checkURL(o.URL); err != nil {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "url"), err))
	}
	return errs
}

type LicenseBuilder struct {
	spec *common.Extendable[License]
}

func NewLicenseBuilder() *LicenseBuilder {
	return &LicenseBuilder{
		spec: common.NewExtendable[License](&License{}),
	}
}

func (b *LicenseBuilder) Build() *common.Extendable[License] {
	return b.spec
}

func (b *LicenseBuilder) Extensions(v map[string]any) *LicenseBuilder {
	b.spec.Extensions = v
	return b
}

func (b *LicenseBuilder) AddExt(name string, value any) *LicenseBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *LicenseBuilder) Name(v string) *LicenseBuilder {
	b.spec.Spec.Name = v
	return b
}

func (b *LicenseBuilder) Identifier(v string) *LicenseBuilder {
	b.spec.Spec.Identifier = v
	return b
}

func (b *LicenseBuilder) URL(v string) *LicenseBuilder {
	b.spec.Spec.URL = v
	return b
}
