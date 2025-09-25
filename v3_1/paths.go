package v3_1

import (
	"encoding/json"
	"strings"

	"github.com/sv-tools/openapi/common"
	"go.yaml.in/yaml/v4"
)

// Paths holds the relative paths to the individual endpoints and their operations.
// The path is appended to the URL from the Server Object in order to construct the full URL.
// The Paths MAY be empty, due to Access Control List (ACL) constraints.
//
// https://spec.openapis.org/oas/v3.1.1#paths-object
//
// Example:
//
//	/pets:
//	  get:
//	    description: Returns all pets from the system that the user has access to
//	    responses:
//	      '200':
//	        description: A list of pets.
//	        content:
//	          application/json:
//	            schema:
//	              type: array
//	              items:
//	                $ref: '#/components/schemas/pet'
type Paths struct {
	// A relative path to an individual endpoint.
	// The field name MUST begin with a forward slash (/).
	// The path is appended (no relative URL resolution) to the expanded URL
	// from the Server Object’s url field in order to construct the full URL.
	// Path templating is allowed.
	// When matching URLs, concrete (non-templated) paths would be matched before their templated counterparts.
	// Templated paths with the same hierarchy but different templated names MUST NOT exist as they are identical.
	// In case of ambiguous matching, it’s up to the tooling to decide which one to use.
	Paths map[string]*common.RefOrSpec[common.Extendable[PathItem]] `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler interface.
func (o *Paths) MarshalJSON() ([]byte, error) {
	return json.Marshal(&o.Paths)
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *Paths) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&o.Paths)
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Paths) MarshalYAML() (any, error) {
	return o.Paths, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Paths) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.Paths)
}

func (o *Paths) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	for k, v := range o.Paths {
		if !strings.HasPrefix(k, "/") {
			errs = append(errs, common.newValidationError(common.joinLoc(location, k), "path must start with a forward slash (`/`)"))
		}
		if v == nil {
			errs = append(errs, common.newValidationError(common.joinLoc(location, k), "path item cannot be empty"))
		} else {
			errs = append(errs, v.validateSpec(common.joinLoc(location, k), validator)...)
		}
	}
	return errs
}

func (o *Paths) Add(path string, item *common.RefOrSpec[common.Extendable[PathItem]]) *Paths {
	if item == nil {
		return o
	}
	if o.Paths == nil {
		o.Paths = make(map[string]*common.RefOrSpec[common.Extendable[PathItem]])
	}
	o.Paths[path] = item
	return o
}

func NewPaths() *common.Extendable[Paths] {
	return common.NewExtendable[Paths](&Paths{})
}
