package v3_1

import (
	"encoding/json"

	"github.com/sv-tools/openapi/common"
	"go.yaml.in/yaml/v4"
)

// Callback is a map of possible out-of band callbacks related to the parent operation.
// Each value in the map is a Path Item Object that describes a set of requests that may be initiated by
// the API provider and the expected responses.
// The key value used to identify the path item object is an expression, evaluated at runtime,
// that identifies a URL to use for the callback operation.
// To describe incoming requests from the API provider independent from another API call, use the webhooks field.
//
// https://spec.openapis.org/oas/v3.1.1#callback-object
//
// Example:
//
//	myCallback:
//	  '{$request.query.queryUrl}':
//	    post:
//	      requestBody:
//	        description: Callback payload
//	        content:
//	          'application/json':
//	            schema:
//	              $ref: '#/components/schemas/SomePayload'
//	      responses:
//	        '200':
//	          description: callback successfully processed
type Callback struct {
	Paths map[string]*common.RefOrSpec[common.Extendable[PathItem]]
}

// MarshalJSON implements json.Marshaler interface.
func (o *Callback) MarshalJSON() ([]byte, error) {
	return json.Marshal(&o.Paths)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Callback) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.Paths)
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *Callback) MarshalYAML() (any, error) {
	return o.Paths, nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *Callback) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&o.Paths)
}

func (o *Callback) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	for k, v := range o.Paths {
		errs = append(errs, v.validateSpec(common.joinLoc(location, k), validator)...)
	}
	return errs
}

func (o *Callback) Add(expression string, item *common.RefOrSpec[common.Extendable[PathItem]]) *Callback {
	if o.Paths == nil {
		o.Paths = make(map[string]*common.RefOrSpec[common.Extendable[PathItem]], 1)
	}
	o.Paths[expression] = item
	return o
}

type CallbackBuilder struct {
	spec *common.RefOrSpec[common.Extendable[Callback]]
}

func NewCallbackBuilder() *CallbackBuilder {
	return &CallbackBuilder{
		spec: common.NewRefOrExtSpec[Callback](&Callback{
			Paths: make(map[string]*common.RefOrSpec[common.Extendable[PathItem]]),
		}),
	}
}

func (b *CallbackBuilder) Build() *common.RefOrSpec[common.Extendable[Callback]] {
	return b.spec
}

func (b *CallbackBuilder) Extensions(v map[string]any) *CallbackBuilder {
	b.spec.Spec.Extensions = v
	return b
}

func (b *CallbackBuilder) AddExt(name string, value any) *CallbackBuilder {
	b.spec.Spec.AddExt(name, value)
	return b
}

func (b *CallbackBuilder) Paths(paths map[string]*common.RefOrSpec[common.Extendable[PathItem]]) *CallbackBuilder {
	b.spec.Spec.Spec.Paths = paths
	return b
}

func (b *CallbackBuilder) AddPathItem(expression string, item *common.RefOrSpec[common.Extendable[PathItem]]) *CallbackBuilder {
	b.spec.Spec.Spec.Add(expression, item)
	return b
}
