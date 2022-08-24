package model

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type APIResponse struct {

	// code
	Code int64 `json:"code,omitempty"`

	// a (key, value) map.
	Details interface{} `json:"details,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this Api response
func (m *APIResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this Api response based on context it is used
func (m *APIResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *APIResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIResponse) UnmarshalBinary(b []byte) error {
	var res APIResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
