// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

func encodeCreateItemResponse(response *CreateItemResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeCreateOrganizationResponse(response *CreateOrganizationResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeCreateStorageGroupResponse(response *CreateStorageGroupResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeCreateUnitResponse(response *CreateOrganizationUnitResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeDeleteItemResponse(response *DeleteItemOK, w http.ResponseWriter) error {
	w.WriteHeader(200)

	return nil
}

func encodeDeleteOrganizationResponse(response *DeleteOrganizationOK, w http.ResponseWriter) error {
	w.WriteHeader(200)

	return nil
}

func encodeDeleteOrganizationUnitResponse(response *DeleteOrganizationUnitOK, w http.ResponseWriter) error {
	w.WriteHeader(200)

	return nil
}

func encodeDeleteStorageGroupResponse(response *DeleteStorageGroupOK, w http.ResponseWriter) error {
	w.WriteHeader(200)

	return nil
}

func encodeExchangeYandexAccessTokenResponse(response *AuthResponse, w http.ResponseWriter) error {
	// Encoding response headers.
	{
		h := uri.NewHeaderEncoder(w.Header())
		// Encode "Set-Cookie" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Set-Cookie",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				return e.EncodeValue(conv.StringToString(response.SetCookie))
			}); err != nil {
				return errors.Wrap(err, "encode Set-Cookie header")
			}
		}
	}
	w.WriteHeader(200)

	return nil
}

func encodeGetAuthCookieByEmailResponse(response *GetAuthCookieByEmailOK, w http.ResponseWriter) error {
	// Encoding response headers.
	{
		h := uri.NewHeaderEncoder(w.Header())
		// Encode "Set-Cookie" header.
		{
			cfg := uri.HeaderParameterEncodingConfig{
				Name:    "Set-Cookie",
				Explode: false,
			}
			if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
				return e.EncodeValue(conv.StringToString(response.SetCookie))
			}); err != nil {
				return errors.Wrap(err, "encode Set-Cookie header")
			}
		}
	}
	w.WriteHeader(200)

	return nil
}

func encodeGetCurrentUserResponse(response *GetCurrentUserResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetItemByIdResponse(response *GetItemByIdResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetItemsResponse(response *GetItemsResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetOrganizationByIdResponse(response *GetOrganizationByIdResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetOrganizationUnitByIdResponse(response *GetOrganizationUnitByIdResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetOrganizationUnitsResponse(response *GetOrganizationUnitsResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetOrganizationsResponse(response *GetOrganizationsResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetStorageGroupByIdResponse(response *GetStorageGroupByIdResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetStorageGroupsResponse(response *GetStorageGroupsResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodePatchItemResponse(response *PatchItemResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodePatchOrganizationResponse(response *PatchOrganizationResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodePatchOrganizationUnitResponse(response *PatchOrganizationUnitResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodePatchStorageGroupResponse(response *PatchStorageGroupResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeUpdateItemResponse(response *UpdateItemResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeUpdateOrganizationResponse(response *UpdateOrganizationResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeUpdateOrganizationUnitResponse(response *UpdateOrganizationUnitResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeUpdateStorageGroupResponse(response *UpdateStorageGroupResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeErrorResponse(response *DefaultErrorStatusCode, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := response.StatusCode
	if code == 0 {
		// Set default status code.
		code = http.StatusOK
	}
	w.WriteHeader(code)

	e := new(jx.Encoder)
	response.Response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	if code >= http.StatusInternalServerError {
		return errors.Wrapf(ht.ErrInternalServerErrorResponse, "code: %d, message: %s", code, http.StatusText(code))
	}
	return nil

}
