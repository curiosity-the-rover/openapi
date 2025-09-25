package common

import (
	"encoding/json"

	"go.yaml.in/yaml/v4"
)

// BoolOrSchema handles Boolean or Any type.
//
// It MUST be used as a pointer,
// otherwise the `false` can be omitted by json or yaml encoders in case of `omitempty` tag is set.
type BoolOrSchema[T any] struct {
	Schema  *RefOrSpec[T]
	Allowed bool
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *BoolOrSchema[T]) UnmarshalJSON(data []byte) error {
	if json.Unmarshal(data, &o.Allowed) == nil {
		o.Schema = nil
		return nil
	}
	if err := json.Unmarshal(data, &o.Schema); err != nil {
		return err
	}
	o.Allowed = true
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (o *BoolOrSchema[T]) MarshalJSON() ([]byte, error) {
	var v any
	if o.Schema != nil {
		v = o.Schema
	} else {
		v = o.Allowed
	}
	return json.Marshal(&v)
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (o *BoolOrSchema[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Decode(&o.Allowed) == nil {
		o.Schema = nil
		return nil
	}
	if err := node.Decode(&o.Schema); err != nil {
		return err
	}
	o.Allowed = true
	return nil
}

// MarshalYAML implements yaml.Marshaler interface.
func (o *BoolOrSchema[T]) MarshalYAML() (any, error) {
	var v any
	if o.Schema != nil {
		v = o.Schema
	} else {
		v = o.Allowed
	}

	return v, nil
}

func (o *BoolOrSchema[T]) validateSpec(path string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.Schema != nil {
		errs = append(errs, o.Schema.validateSpec(path, validator)...)
	}
	return errs
}

func NewBoolOrSchema[T any](v any) *BoolOrSchema[T] {
	switch v := v.(type) {
	case bool:
		return &BoolOrSchema[T]{Allowed: v}
	case *RefOrSpec[T]:
		return &BoolOrSchema[T]{Schema: v}
	case Builder[RefOrSpec[T]]:
		return &BoolOrSchema[T]{Schema: v.Build()}
	default:
		return nil
	}
}
