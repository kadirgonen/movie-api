package model

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

type SignUp struct {

	// email
	// Required: true
	Email *string `json:"email"`

	// firstname
	// Required: true
	Firstname *string `json:"firstname"`

	// lastname
	// Required: true
	Lastname *string `json:"lastname"`

	// password
	// Required: true
	Password *string `json:"password"`
}

// Validate validates this sign up
func (m *SignUp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFirstname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SignUp) validateEmail(formats strfmt.Registry) error {

	if err := validate.Required("email", "body", m.Email); err != nil {
		return err
	}

	return nil
}

func (m *SignUp) validateFirstname(formats strfmt.Registry) error {

	if err := validate.Required("firstname", "body", m.Firstname); err != nil {
		return err
	}

	return nil
}

func (m *SignUp) validateLastname(formats strfmt.Registry) error {

	if err := validate.Required("lastname", "body", m.Lastname); err != nil {
		return err
	}

	return nil
}

func (m *SignUp) validatePassword(formats strfmt.Registry) error {

	if err := validate.Required("password", "body", m.Password); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this sign up based on context it is used
func (m *SignUp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SignUp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SignUp) UnmarshalBinary(b []byte) error {
	var res SignUp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
