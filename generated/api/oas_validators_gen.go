// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/validate"
)

func (s *Organization) Validate() error {
	if s == nil {
		return validate.ErrNilPointer
	}

	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    1,
			MinLengthSet: true,
			MaxLength:    255,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        regexMap["^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$"],
		}).Validate(string(s.Subdomain)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "subdomain",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}

func (s *OrganizationsPagedResponse) Validate() error {
	if s == nil {
		return validate.ErrNilPointer
	}

	var failures []validate.FieldError
	if err := func() error {
		var failures []validate.FieldError
		for i, elem := range s.Items {
			if err := func() error {
				if err := elem.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				failures = append(failures, validate.FieldError{
					Name:  fmt.Sprintf("[%d]", i),
					Error: err,
				})
			}
		}
		if len(failures) > 0 {
			return &validate.Error{Fields: failures}
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "items",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}

func (s *OrganizationsPagedResponseItemsItem) Validate() error {
	if s == nil {
		return validate.ErrNilPointer
	}

	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    1,
			MinLengthSet: true,
			MaxLength:    255,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        regexMap["^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$"],
		}).Validate(string(s.Subdomain)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "subdomain",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
