package v3_1

import "github.com/sv-tools/openapi/common"

// OAuthFlows allows configuration of the supported OAuth Flows.
//
// https://spec.openapis.org/oas/v3.1.1#oauth-flows-object
//
// Example:
//
//	type: oauth2
//	flows:
//	  implicit:
//	    authorizationUrl: https://example.com/api/oauth/dialog
//	    scopes:
//	      write:pets: modify pets in your account
//	      read:pets: read your pets
//	  authorizationCode:
//	    authorizationUrl: https://example.com/api/oauth/dialog
//	    tokenUrl: https://example.com/api/oauth/token
//	    scopes:
//	      write:pets: modify pets in your account
//	      read:pets: read your pets
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow.
	Implicit *common.Extendable[OAuthFlow] `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow.
	Password *common.Extendable[OAuthFlow] `json:"password,omitempty" yaml:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow.
	// Previously called application in OpenAPI 2.0.
	ClientCredentials *common.Extendable[OAuthFlow] `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow.
	// Previously called accessCode in OpenAPI 2.0.
	AuthorizationCode *common.Extendable[OAuthFlow] `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

func (o *OAuthFlows) validateSpec(location string, validator *common.Validator) []*common.validationError {
	var errs []*common.validationError
	if o.Implicit != nil {
		errs = append(errs, o.Implicit.validateSpec(common.joinLoc(location, "implicit"), validator)...)
		if o.Implicit.Spec.AuthorizationURL == "" {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "implicit", "authorizationUrl"), common.ErrRequired))
		}
	}
	if o.Password != nil {
		errs = append(errs, o.Password.validateSpec(common.joinLoc(location, "password"), validator)...)
		if o.Password.Spec.TokenURL == "" {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "password", "tokenUrl"), common.ErrRequired))
		}
	}
	if o.ClientCredentials != nil {
		errs = append(errs, o.ClientCredentials.validateSpec(common.joinLoc(location, "clientCredentials"), validator)...)
		if o.ClientCredentials.Spec.TokenURL == "" {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "clientCredentials", "tokenUrl"), common.ErrRequired))
		}
	}
	if o.AuthorizationCode != nil {
		errs = append(errs, o.AuthorizationCode.validateSpec(common.joinLoc(location, "authorizationCode"), validator)...)
		if o.AuthorizationCode.Spec.AuthorizationURL == "" {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "authorizationCode", "authorizationUrl"), common.ErrRequired))
		}
		if o.AuthorizationCode.Spec.TokenURL == "" {
			errs = append(errs, common.newValidationError(common.joinLoc(location, "authorizationCode", "tokenUrl"), common.ErrRequired))
		}
	}

	return errs
}

type OAuthFlowsBuilder struct {
	spec *common.Extendable[OAuthFlows]
}

func NewOAuthFlowsBuilder() *OAuthFlowsBuilder {
	return &OAuthFlowsBuilder{
		spec: common.NewExtendable[OAuthFlows](&OAuthFlows{}),
	}
}

func (b *OAuthFlowsBuilder) Build() *common.Extendable[OAuthFlows] {
	return b.spec
}

func (b *OAuthFlowsBuilder) Extensions(v map[string]any) *OAuthFlowsBuilder {
	b.spec.Extensions = v
	return b
}

func (b *OAuthFlowsBuilder) AddExt(name string, value any) *OAuthFlowsBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *OAuthFlowsBuilder) Implicit(v *common.Extendable[OAuthFlow]) *OAuthFlowsBuilder {
	b.spec.Spec.Implicit = v
	return b
}

func (b *OAuthFlowsBuilder) Password(v *common.Extendable[OAuthFlow]) *OAuthFlowsBuilder {
	b.spec.Spec.Password = v
	return b
}

func (b *OAuthFlowsBuilder) ClientCredentials(v *common.Extendable[OAuthFlow]) *OAuthFlowsBuilder {
	b.spec.Spec.ClientCredentials = v
	return b
}

func (b *OAuthFlowsBuilder) AuthorizationCode(v *common.Extendable[OAuthFlow]) *OAuthFlowsBuilder {
	b.spec.Spec.AuthorizationCode = v
	return b
}
