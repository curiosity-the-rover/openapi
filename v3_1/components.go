package v3_1

import (
	"regexp"

	"github.com/sv-tools/openapi/common"
)

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced
// from properties outside the components object.
//
// https://spec.openapis.org/oas/v3.1.1#components-object
//
// Example:
//
//	components:
//	  schemas:
//	    GeneralError:
//	      type: object
//	      properties:
//	        code:
//	          type: integer
//	          format: int32
//	        message:
//	          type: string
//	    Category:
//	      type: object
//	      properties:
//	        id:
//	          type: integer
//	          format: int64
//	        name:
//	          type: string
//	    Tag:
//	      type: object
//	      properties:
//	        id:
//	          type: integer
//	          format: int64
//	        name:
//	          type: string
//	  parameters:
//	    skipParam:
//	      name: skip
//	      in: query
//	      description: number of items to skip
//	      required: true
//	      schema:
//	        type: integer
//	        format: int32
//	    limitParam:
//	      name: limit
//	      in: query
//	      description: max records to return
//	      required: true
//	      schema:
//	        type: integer
//	        format: int32
//	  responses:
//	    NotFound:
//	      description: Entity not found.
//	    IllegalInput:
//	      description: Illegal input for operation.
//	    GeneralError:
//	      description: General Error
//	      content:
//	        application/json:
//	          schema:
//	            $ref: '#/components/schemas/GeneralError'
//	  securitySchemes:
//	    api_key:
//	      type: apiKey
//	      name: api_key
//	      in: header
//	    petstore_auth:
//	      type: oauth2
//	      flows:
//	        implicit:
//	          authorizationUrl: https://example.org/api/oauth/dialog
//	          scopes:
//	            write:pets: modify pets in your account
//	            read:pets: read your pets
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas map[string]*common.RefOrSpec[Schema] `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses map[string]*common.RefOrSpec[common.Extendable[Response]] `json:"responses,omitempty" yaml:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters map[string]*common.RefOrSpec[common.Extendable[Parameter]] `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples map[string]*common.RefOrSpec[common.Extendable[Example]] `json:"examples,omitempty" yaml:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies map[string]*common.RefOrSpec[common.Extendable[RequestBody]] `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers map[string]*common.RefOrSpec[common.Extendable[Header]] `json:"headers,omitempty" yaml:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes map[string]*common.RefOrSpec[common.Extendable[SecurityScheme]] `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links map[string]*common.RefOrSpec[common.Extendable[Link]] `json:"links,omitempty" yaml:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks map[string]*common.RefOrSpec[common.Extendable[Callback]] `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	Paths map[string]*common.RefOrSpec[common.Extendable[PathItem]] `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// Add adds the given object to the appropriate list based on a type and returns the current object (self|this).
func (o *Components) Add(name string, v any) *Components {
	if v == nil {
		return o
	}
	switch spec := v.(type) {
	case *common.RefOrSpec[Schema]:
		if o.Schemas == nil {
			o.Schemas = make(map[string]*common.RefOrSpec[Schema], 1)
		}
		o.Schemas[name] = spec
	case *common.RefOrSpec[common.Extendable[Response]]:
		if o.Responses == nil {
			o.Responses = make(map[string]*common.RefOrSpec[common.Extendable[Response]], 1)
		}
		o.Responses[name] = spec
	case *common.RefOrSpec[common.Extendable[Parameter]]:
		if o.Parameters == nil {
			o.Parameters = make(map[string]*common.RefOrSpec[common.Extendable[Parameter]], 1)
		}
		o.Parameters[name] = spec
	case *common.RefOrSpec[common.Extendable[Example]]:
		if o.Examples == nil {
			o.Examples = make(map[string]*common.RefOrSpec[common.Extendable[Example]], 1)
		}
		o.Examples[name] = spec
	case *common.RefOrSpec[common.Extendable[RequestBody]]:
		if o.RequestBodies == nil {
			o.RequestBodies = make(map[string]*common.RefOrSpec[common.Extendable[RequestBody]], 1)
		}
		o.RequestBodies[name] = spec
	case *common.RefOrSpec[common.Extendable[Header]]:
		if o.Headers == nil {
			o.Headers = make(map[string]*common.RefOrSpec[common.Extendable[Header]], 1)
		}
		o.Headers[name] = spec
	case *common.RefOrSpec[common.Extendable[SecurityScheme]]:
		if o.SecuritySchemes == nil {
			o.SecuritySchemes = make(map[string]*common.RefOrSpec[common.Extendable[SecurityScheme]], 1)
		}
		o.SecuritySchemes[name] = spec
	case *common.RefOrSpec[common.Extendable[Link]]:
		if o.Links == nil {
			o.Links = make(map[string]*common.RefOrSpec[common.Extendable[Link]], 1)
		}
		o.Links[name] = spec
	case *common.RefOrSpec[common.Extendable[Callback]]:
		if o.Callbacks == nil {
			o.Callbacks = make(map[string]*common.RefOrSpec[common.Extendable[Callback]], 1)
		}
		o.Callbacks[name] = spec
	case *common.RefOrSpec[common.Extendable[PathItem]]:
		if o.Paths == nil {
			o.Paths = make(map[string]*common.RefOrSpec[common.Extendable[PathItem]], 1)
		}
		o.Paths[name] = spec
	default:
		// ignore to avoid panic
	}
	return o
}

var namePattern = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)

func (o *Components) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	for k, v := range o.Schemas {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "schemas", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "schemas", k), validator)...)
	}

	for k, v := range o.Responses {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "responses", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "responses", k), validator)...)
	}

	for k, v := range o.Parameters {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "parameters", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "parameters", k), validator)...)
	}

	for k, v := range o.Examples {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "examples", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "examples", k), validator)...)
	}

	for k, v := range o.RequestBodies {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "requestBodies", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "requestBodies", k), validator)...)
	}

	for k, v := range o.Headers {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "headers", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "headers", k), validator)...)
	}

	for k, v := range o.SecuritySchemes {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "securitySchemes", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "securitySchemes", k), validator)...)
	}

	for k, v := range o.Links {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "links", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "links", k), validator)...)
	}

	for k, v := range o.Callbacks {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "callbacks", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "callbacks", k), validator)...)
	}

	for k, v := range o.Paths {
		if !namePattern.MatchString(k) {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "paths", k), "invalid name %q, must match %q", k, namePattern.String()))
		}
		errs = append(errs, v.validateSpec(common.joinLoc(location, "paths", k), validator)...)
	}

	return errs
}

func NewComponents() *common.Extendable[Components] {
	return common.NewExtendable[Components](&Components{})
}
