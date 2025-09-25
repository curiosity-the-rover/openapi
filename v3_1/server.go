package v3_1

import (
	"strings"

	"github.com/sv-tools/openapi/common"
)

// Server is an object representing a Server.
//
// https://spec.openapis.org/oas/v3.1.1#server-object
//
// Example:
//
//	servers:
//	- url: https://development.gigantic-server.com/v1
//	  description: Development server
//	- url: https://staging.gigantic-server.com/v1
//	  description: Staging server
//	- url: https://api.gigantic-server.com/v1
//	  description: Production server
type Server struct {
	// A map between a variable name and its value.
	// The value is used for substitution in the server’s URL template.
	Variables map[string]*common.Extendable[ServerVariable] `json:"variables,omitempty" yaml:"variables,omitempty"`
	// REQUIRED.
	// A URL to the target host.
	// This URL supports Server Variables and MAY be relative, to indicate that the host location is relative
	// to the location where the OpenAPI document is being served.
	// Variable substitutions will be made when a variable is named in {brackets}.
	URL string `json:"url" yaml:"url"`
	// An optional string describing the host designated by the URL.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (o *Server) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if o.URL == "" {
		errs = append(errs, common.newValidationError(common.joinLoc(location, "url"), common.ErrRequired))
	}
	if l := len(o.Variables); l == 0 {
		if err := common.checkURL(o.URL); err != nil {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "url"), err))
		}
	} else {
		oldnew := make([]string, 0, l*2)
		for k, v := range o.Variables {
			errs = append(errs, v.validateSpec(common.joinLoc(location, "variables", k), validator)...)
			oldnew = append(oldnew, "{"+k+"}", v.Spec.Default)
		}
		u := strings.NewReplacer(oldnew...).Replace(o.URL)
		if err := common.checkURL(u); err != nil {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "url"), err))
		}
	}
	return errs
}

type ServerBuilder struct {
	spec *common.Extendable[Server]
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		spec: common.NewExtendable[Server](&Server{}),
	}
}

func (b *ServerBuilder) Build() *common.Extendable[Server] {
	return b.spec
}

func (b *ServerBuilder) Extensions(v map[string]any) *ServerBuilder {
	b.spec.Extensions = v
	return b
}

func (b *ServerBuilder) AddExt(name string, value any) *ServerBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *ServerBuilder) Variables(v map[string]*common.Extendable[ServerVariable]) *ServerBuilder {
	b.spec.Spec.Variables = v
	return b
}

func (b *ServerBuilder) AddVariable(name string, value *common.Extendable[ServerVariable]) *ServerBuilder {
	if b.spec.Spec.Variables == nil {
		b.spec.Spec.Variables = make(map[string]*common.Extendable[ServerVariable], 1)
	}
	b.spec.Spec.Variables[name] = value
	return b
}

func (b *ServerBuilder) URL(v string) *ServerBuilder {
	b.spec.Spec.URL = v
	return b
}

func (b *ServerBuilder) Description(v string) *ServerBuilder {
	b.spec.Spec.Description = v
	return b
}
