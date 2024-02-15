// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AccessToken access token
//
// swagger:model AccessToken
type AccessToken struct {

	// Access Token
	// Required: true
	AccessToken *string `json:"accessToken"`
}

// Validate validates this access token
func (m *AccessToken) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccessToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AccessToken) validateAccessToken(formats strfmt.Registry) error {

	if err := validate.Required("accessToken", "body", m.AccessToken); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this access token based on context it is used
func (m *AccessToken) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AccessToken) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccessToken) UnmarshalBinary(b []byte) error {
	var res AccessToken
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
