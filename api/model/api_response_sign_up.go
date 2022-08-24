package model

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type APIResponseSignUp struct {

	// code
	Code int32 `json:"code,omitempty"`

	// token
	Token *Token `json:"token,omitempty"`
}

// Validate validates this Api response sign up
func (m *APIResponseSignUp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIResponseSignUp) validateToken(formats strfmt.Registry) error {
	if swag.IsZero(m.Token) { // not required
		return nil
	}

	if m.Token != nil {
		if err := m.Token.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("token")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("token")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this Api response sign up based on the context it is used
func (m *APIResponseSignUp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateToken(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIResponseSignUp) contextValidateToken(ctx context.Context, formats strfmt.Registry) error {

	if m.Token != nil {
		if err := m.Token.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("token")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("token")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *APIResponseSignUp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIResponseSignUp) UnmarshalBinary(b []byte) error {
	var res APIResponseSignUp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
