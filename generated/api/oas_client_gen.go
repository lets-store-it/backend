// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

func trimTrailingSlashes(u *url.URL) {
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
}

// Invoker invokes operations described by OpenAPI v3 specification.
type Invoker interface {
	// CreateOrganization invokes createOrganization operation.
	//
	// Create Organization.
	//
	// POST /orgs
	CreateOrganization(ctx context.Context, request *CreateOrganizationRequest) (*CreateOrganizationResponse, error)
	// CreateUnit invokes createUnit operation.
	//
	// Create Organization Unit.
	//
	// POST /units
	CreateUnit(ctx context.Context, request *CreateOrganizationUnitRequest) (*CreateOrganizationUnitResponse, error)
	// DeleteOrganization invokes deleteOrganization operation.
	//
	// Delete Organization.
	//
	// DELETE /orgs/{id}
	DeleteOrganization(ctx context.Context, params DeleteOrganizationParams) error
	// DeleteOrganizationUnit invokes deleteOrganizationUnit operation.
	//
	// Delete Organization Unit.
	//
	// DELETE /units/{id}
	DeleteOrganizationUnit(ctx context.Context, params DeleteOrganizationUnitParams) error
	// GetOrganizationById invokes getOrganizationById operation.
	//
	// Get Organization by ID.
	//
	// GET /orgs/{id}
	GetOrganizationById(ctx context.Context, params GetOrganizationByIdParams) (*GetOrganizationByIdResponse, error)
	// GetOrganizationUnitById invokes getOrganizationUnitById operation.
	//
	// Get Unit by ID with Spaces.
	//
	// GET /units/{id}
	GetOrganizationUnitById(ctx context.Context, params GetOrganizationUnitByIdParams) (*GetOrganizationUnitByIdResponse, error)
	// GetOrganizationUnits invokes getOrganizationUnits operation.
	//
	// Get list of Organization Units.
	//
	// GET /units
	GetOrganizationUnits(ctx context.Context) (*GetOrganizationUnitsResponse, error)
	// GetOrganizations invokes getOrganizations operation.
	//
	// Get list of Organizations.
	//
	// GET /orgs
	GetOrganizations(ctx context.Context) (*GetOrganizationsResponse, error)
	// PatchOrganization invokes patchOrganization operation.
	//
	// Update Organization.
	//
	// PATCH /orgs/{id}
	PatchOrganization(ctx context.Context, request *PatchOrganizationRequest, params PatchOrganizationParams) (*PatchOrganizationResponse, error)
	// PatchOrganizationUnit invokes patchOrganizationUnit operation.
	//
	// Patch Organization Unit.
	//
	// PATCH /units/{id}
	PatchOrganizationUnit(ctx context.Context, request *PatchOrganizationUnitRequest, params PatchOrganizationUnitParams) (*PatchOrganizationUnitResponse, error)
	// UpdateOrganization invokes updateOrganization operation.
	//
	// Update Organization.
	//
	// PUT /orgs/{id}
	UpdateOrganization(ctx context.Context, request *UpdateOrganizationRequest, params UpdateOrganizationParams) (*UpdateOrganizationResponse, error)
	// UpdateOrganizationUnit invokes updateOrganizationUnit operation.
	//
	// Update Organization Unit.
	//
	// PUT /units/{id}
	UpdateOrganizationUnit(ctx context.Context, request *UpdateOrganizationUnitRequest, params UpdateOrganizationUnitParams) (*UpdateOrganizationUnitResponse, error)
}

// Client implements OAS client.
type Client struct {
	serverURL *url.URL
	baseClient
}
type errorHandler interface {
	NewError(ctx context.Context, err error) *ErrorStatusCode
}

var _ Handler = struct {
	errorHandler
	*Client
}{}

// NewClient initializes new Client defined by OAS.
func NewClient(serverURL string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	trimTrailingSlashes(u)

	c, err := newClientConfig(opts...).baseClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		serverURL:  u,
		baseClient: c,
	}, nil
}

type serverURLKey struct{}

// WithServerURL sets context key to override server URL.
func WithServerURL(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, serverURLKey{}, u)
}

func (c *Client) requestURL(ctx context.Context) *url.URL {
	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
	if !ok {
		return c.serverURL
	}
	return u
}

// CreateOrganization invokes createOrganization operation.
//
// Create Organization.
//
// POST /orgs
func (c *Client) CreateOrganization(ctx context.Context, request *CreateOrganizationRequest) (*CreateOrganizationResponse, error) {
	res, err := c.sendCreateOrganization(ctx, request)
	return res, err
}

func (c *Client) sendCreateOrganization(ctx context.Context, request *CreateOrganizationRequest) (res *CreateOrganizationResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/orgs"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeCreateOrganizationRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeCreateOrganizationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// CreateUnit invokes createUnit operation.
//
// Create Organization Unit.
//
// POST /units
func (c *Client) CreateUnit(ctx context.Context, request *CreateOrganizationUnitRequest) (*CreateOrganizationUnitResponse, error) {
	res, err := c.sendCreateUnit(ctx, request)
	return res, err
}

func (c *Client) sendCreateUnit(ctx context.Context, request *CreateOrganizationUnitRequest) (res *CreateOrganizationUnitResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/units"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeCreateUnitRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeCreateUnitResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// DeleteOrganization invokes deleteOrganization operation.
//
// Delete Organization.
//
// DELETE /orgs/{id}
func (c *Client) DeleteOrganization(ctx context.Context, params DeleteOrganizationParams) error {
	_, err := c.sendDeleteOrganization(ctx, params)
	return err
}

func (c *Client) sendDeleteOrganization(ctx context.Context, params DeleteOrganizationParams) (res *DeleteOrganizationOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/orgs/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "DELETE", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeDeleteOrganizationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// DeleteOrganizationUnit invokes deleteOrganizationUnit operation.
//
// Delete Organization Unit.
//
// DELETE /units/{id}
func (c *Client) DeleteOrganizationUnit(ctx context.Context, params DeleteOrganizationUnitParams) error {
	_, err := c.sendDeleteOrganizationUnit(ctx, params)
	return err
}

func (c *Client) sendDeleteOrganizationUnit(ctx context.Context, params DeleteOrganizationUnitParams) (res *DeleteOrganizationUnitOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/units/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "DELETE", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeDeleteOrganizationUnitResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// GetOrganizationById invokes getOrganizationById operation.
//
// Get Organization by ID.
//
// GET /orgs/{id}
func (c *Client) GetOrganizationById(ctx context.Context, params GetOrganizationByIdParams) (*GetOrganizationByIdResponse, error) {
	res, err := c.sendGetOrganizationById(ctx, params)
	return res, err
}

func (c *Client) sendGetOrganizationById(ctx context.Context, params GetOrganizationByIdParams) (res *GetOrganizationByIdResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/orgs/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeGetOrganizationByIdResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// GetOrganizationUnitById invokes getOrganizationUnitById operation.
//
// Get Unit by ID with Spaces.
//
// GET /units/{id}
func (c *Client) GetOrganizationUnitById(ctx context.Context, params GetOrganizationUnitByIdParams) (*GetOrganizationUnitByIdResponse, error) {
	res, err := c.sendGetOrganizationUnitById(ctx, params)
	return res, err
}

func (c *Client) sendGetOrganizationUnitById(ctx context.Context, params GetOrganizationUnitByIdParams) (res *GetOrganizationUnitByIdResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/units/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeGetOrganizationUnitByIdResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// GetOrganizationUnits invokes getOrganizationUnits operation.
//
// Get list of Organization Units.
//
// GET /units
func (c *Client) GetOrganizationUnits(ctx context.Context) (*GetOrganizationUnitsResponse, error) {
	res, err := c.sendGetOrganizationUnits(ctx)
	return res, err
}

func (c *Client) sendGetOrganizationUnits(ctx context.Context) (res *GetOrganizationUnitsResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/units"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeGetOrganizationUnitsResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// GetOrganizations invokes getOrganizations operation.
//
// Get list of Organizations.
//
// GET /orgs
func (c *Client) GetOrganizations(ctx context.Context) (*GetOrganizationsResponse, error) {
	res, err := c.sendGetOrganizations(ctx)
	return res, err
}

func (c *Client) sendGetOrganizations(ctx context.Context) (res *GetOrganizationsResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/orgs"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeGetOrganizationsResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// PatchOrganization invokes patchOrganization operation.
//
// Update Organization.
//
// PATCH /orgs/{id}
func (c *Client) PatchOrganization(ctx context.Context, request *PatchOrganizationRequest, params PatchOrganizationParams) (*PatchOrganizationResponse, error) {
	res, err := c.sendPatchOrganization(ctx, request, params)
	return res, err
}

func (c *Client) sendPatchOrganization(ctx context.Context, request *PatchOrganizationRequest, params PatchOrganizationParams) (res *PatchOrganizationResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/orgs/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "PATCH", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodePatchOrganizationRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodePatchOrganizationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// PatchOrganizationUnit invokes patchOrganizationUnit operation.
//
// Patch Organization Unit.
//
// PATCH /units/{id}
func (c *Client) PatchOrganizationUnit(ctx context.Context, request *PatchOrganizationUnitRequest, params PatchOrganizationUnitParams) (*PatchOrganizationUnitResponse, error) {
	res, err := c.sendPatchOrganizationUnit(ctx, request, params)
	return res, err
}

func (c *Client) sendPatchOrganizationUnit(ctx context.Context, request *PatchOrganizationUnitRequest, params PatchOrganizationUnitParams) (res *PatchOrganizationUnitResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/units/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "PATCH", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodePatchOrganizationUnitRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodePatchOrganizationUnitResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// UpdateOrganization invokes updateOrganization operation.
//
// Update Organization.
//
// PUT /orgs/{id}
func (c *Client) UpdateOrganization(ctx context.Context, request *UpdateOrganizationRequest, params UpdateOrganizationParams) (*UpdateOrganizationResponse, error) {
	res, err := c.sendUpdateOrganization(ctx, request, params)
	return res, err
}

func (c *Client) sendUpdateOrganization(ctx context.Context, request *UpdateOrganizationRequest, params UpdateOrganizationParams) (res *UpdateOrganizationResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/orgs/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "PUT", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeUpdateOrganizationRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeUpdateOrganizationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// UpdateOrganizationUnit invokes updateOrganizationUnit operation.
//
// Update Organization Unit.
//
// PUT /units/{id}
func (c *Client) UpdateOrganizationUnit(ctx context.Context, request *UpdateOrganizationUnitRequest, params UpdateOrganizationUnitParams) (*UpdateOrganizationUnitResponse, error) {
	res, err := c.sendUpdateOrganizationUnit(ctx, request, params)
	return res, err
}

func (c *Client) sendUpdateOrganizationUnit(ctx context.Context, request *UpdateOrganizationUnitRequest, params UpdateOrganizationUnitParams) (res *UpdateOrganizationUnitResponse, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/units/"
	{
		// Encode "id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.ID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "PUT", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeUpdateOrganizationUnitRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeUpdateOrganizationUnitResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}
