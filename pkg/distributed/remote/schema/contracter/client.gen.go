// Package contracter provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package contracter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// CancellationReason defines model for CancellationReason.
type CancellationReason struct {
	Reason string `json:"reason"`
}

// ConvertOrder defines model for ConvertOrder.
type ConvertOrder struct {
	Id     string        `json:"id"`
	Input  string        `json:"input"`
	Output string        `json:"output"`
	Params ConvertParams `json:"params"`
	State  string        `json:"state"`
	Type   string        `json:"type"`
}

// ConvertParams defines model for ConvertParams.
type ConvertParams struct {
	HwAccel          string `json:"hw_accel"`
	KeyframeInterval int    `json:"keyframe_interval"`
	Preset           string `json:"preset"`
	Scale            string `json:"scale"`
	VideoBitRate     string `json:"video_bit_rate"`
	VideoCodec       string `json:"video_codec"`
	VideoQuality     int    `json:"video_quality"`
}

// ConvertSegment defines model for ConvertSegment.
type ConvertSegment struct {
	Id       string        `json:"id"`
	Muxer    string        `json:"muxer"`
	OrderId  string        `json:"order_id"`
	Params   ConvertParams `json:"params"`
	Position int           `json:"position"`
	Type     string        `json:"type"`
}

// RFC 7807 Problem Details for HTTP APIs
type ProblemDetails struct {
	Detail *string                `json:"detail,omitempty"`
	Fields *ProblemDetails_Fields `json:"fields,omitempty"`
	Title  string                 `json:"title"`
	Type   *string                `json:"type,omitempty"`
}

// ProblemDetails_Fields defines model for ProblemDetails.Fields.
type ProblemDetails_Fields struct {
	AdditionalProperties map[string]string `json:"-"`
}

// OrderIDParam defines model for orderIDParam.
type OrderIDParam string

// SegmentIDParam defines model for segmentIDParam.
type SegmentIDParam string

// RFC 7807 Problem Details for HTTP APIs
type ResponseForbidden ProblemDetails

// RFC 7807 Problem Details for HTTP APIs
type ResponseNotFound ProblemDetails

// ResponseOrder defines model for ResponseOrder.
type ResponseOrder ConvertOrder

// ResponseOrders defines model for ResponseOrders.
type ResponseOrders []ConvertOrder

// ResponseSegment defines model for ResponseSegment.
type ResponseSegment ConvertSegment

// ResponseSegments defines model for ResponseSegments.
type ResponseSegments []ConvertSegment

// RFC 7807 Problem Details for HTTP APIs
type ResponseUnauthorized ProblemDetails

// CancelOrderByIDJSONBody defines parameters for CancelOrderByID.
type CancelOrderByIDJSONBody CancellationReason

// CancelOrderByIDJSONRequestBody defines body for CancelOrderByID for application/json ContentType.
type CancelOrderByIDJSONRequestBody CancelOrderByIDJSONBody

// Getter for additional properties for ProblemDetails_Fields. Returns the specified
// element and whether it was found
func (a ProblemDetails_Fields) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ProblemDetails_Fields
func (a *ProblemDetails_Fields) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ProblemDetails_Fields to handle AdditionalProperties
func (a *ProblemDetails_Fields) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ProblemDetails_Fields to handle AdditionalProperties
func (a ProblemDetails_Fields) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SearchOrders request
	SearchOrders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetOrderByID request
	GetOrderByID(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CancelOrderByID request  with any body
	CancelOrderByIDWithBody(ctx context.Context, orderID OrderIDParam, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CancelOrderByID(ctx context.Context, orderID OrderIDParam, body CancelOrderByIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SearchSegmentsByOrderID request
	SearchSegmentsByOrderID(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SearchSegments request
	SearchSegments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSegmentByID request
	GetSegmentByID(ctx context.Context, segmentID SegmentIDParam, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SearchOrders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchOrdersRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetOrderByID(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetOrderByIDRequest(c.Server, orderID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelOrderByIDWithBody(ctx context.Context, orderID OrderIDParam, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelOrderByIDRequestWithBody(c.Server, orderID, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelOrderByID(ctx context.Context, orderID OrderIDParam, body CancelOrderByIDJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelOrderByIDRequest(c.Server, orderID, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SearchSegmentsByOrderID(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchSegmentsByOrderIDRequest(c.Server, orderID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SearchSegments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchSegmentsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetSegmentByID(ctx context.Context, segmentID SegmentIDParam, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSegmentByIDRequest(c.Server, segmentID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSearchOrdersRequest generates requests for SearchOrders
func NewSearchOrdersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/orders")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetOrderByIDRequest generates requests for GetOrderByID
func NewGetOrderByIDRequest(server string, orderID OrderIDParam) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/orders/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCancelOrderByIDRequest calls the generic CancelOrderByID builder with application/json body
func NewCancelOrderByIDRequest(server string, orderID OrderIDParam, body CancelOrderByIDJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCancelOrderByIDRequestWithBody(server, orderID, "application/json", bodyReader)
}

// NewCancelOrderByIDRequestWithBody generates requests for CancelOrderByID with any type of body
func NewCancelOrderByIDRequestWithBody(server string, orderID OrderIDParam, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/orders/%s/cancel", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewSearchSegmentsByOrderIDRequest generates requests for SearchSegmentsByOrderID
func NewSearchSegmentsByOrderIDRequest(server string, orderID OrderIDParam) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "orderID", runtime.ParamLocationPath, orderID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/orders/%s/segments", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSearchSegmentsRequest generates requests for SearchSegments
func NewSearchSegmentsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/segments")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetSegmentByIDRequest generates requests for GetSegmentByID
func NewGetSegmentByIDRequest(server string, segmentID SegmentIDParam) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "segmentID", runtime.ParamLocationPath, segmentID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/segments/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SearchOrders request
	SearchOrdersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*SearchOrdersResponse, error)

	// GetOrderByID request
	GetOrderByIDWithResponse(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*GetOrderByIDResponse, error)

	// CancelOrderByID request  with any body
	CancelOrderByIDWithBodyWithResponse(ctx context.Context, orderID OrderIDParam, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelOrderByIDResponse, error)

	CancelOrderByIDWithResponse(ctx context.Context, orderID OrderIDParam, body CancelOrderByIDJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelOrderByIDResponse, error)

	// SearchSegmentsByOrderID request
	SearchSegmentsByOrderIDWithResponse(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*SearchSegmentsByOrderIDResponse, error)

	// SearchSegments request
	SearchSegmentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*SearchSegmentsResponse, error)

	// GetSegmentByID request
	GetSegmentByIDWithResponse(ctx context.Context, segmentID SegmentIDParam, reqEditors ...RequestEditorFn) (*GetSegmentByIDResponse, error)
}

type SearchOrdersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ConvertOrder
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r SearchOrdersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchOrdersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetOrderByIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ConvertOrder
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
	JSON404      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r GetOrderByIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOrderByIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CancelOrderByIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
	JSON404      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r CancelOrderByIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CancelOrderByIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SearchSegmentsByOrderIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ConvertSegment
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
	JSON404      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r SearchSegmentsByOrderIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchSegmentsByOrderIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SearchSegmentsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ConvertSegment
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r SearchSegmentsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchSegmentsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSegmentByIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ConvertSegment
	JSON401      *ProblemDetails
	JSON403      *ProblemDetails
	JSON404      *ProblemDetails
}

// Status returns HTTPResponse.Status
func (r GetSegmentByIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSegmentByIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SearchOrdersWithResponse request returning *SearchOrdersResponse
func (c *ClientWithResponses) SearchOrdersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*SearchOrdersResponse, error) {
	rsp, err := c.SearchOrders(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchOrdersResponse(rsp)
}

// GetOrderByIDWithResponse request returning *GetOrderByIDResponse
func (c *ClientWithResponses) GetOrderByIDWithResponse(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*GetOrderByIDResponse, error) {
	rsp, err := c.GetOrderByID(ctx, orderID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetOrderByIDResponse(rsp)
}

// CancelOrderByIDWithBodyWithResponse request with arbitrary body returning *CancelOrderByIDResponse
func (c *ClientWithResponses) CancelOrderByIDWithBodyWithResponse(ctx context.Context, orderID OrderIDParam, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelOrderByIDResponse, error) {
	rsp, err := c.CancelOrderByIDWithBody(ctx, orderID, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelOrderByIDResponse(rsp)
}

func (c *ClientWithResponses) CancelOrderByIDWithResponse(ctx context.Context, orderID OrderIDParam, body CancelOrderByIDJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelOrderByIDResponse, error) {
	rsp, err := c.CancelOrderByID(ctx, orderID, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelOrderByIDResponse(rsp)
}

// SearchSegmentsByOrderIDWithResponse request returning *SearchSegmentsByOrderIDResponse
func (c *ClientWithResponses) SearchSegmentsByOrderIDWithResponse(ctx context.Context, orderID OrderIDParam, reqEditors ...RequestEditorFn) (*SearchSegmentsByOrderIDResponse, error) {
	rsp, err := c.SearchSegmentsByOrderID(ctx, orderID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchSegmentsByOrderIDResponse(rsp)
}

// SearchSegmentsWithResponse request returning *SearchSegmentsResponse
func (c *ClientWithResponses) SearchSegmentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*SearchSegmentsResponse, error) {
	rsp, err := c.SearchSegments(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchSegmentsResponse(rsp)
}

// GetSegmentByIDWithResponse request returning *GetSegmentByIDResponse
func (c *ClientWithResponses) GetSegmentByIDWithResponse(ctx context.Context, segmentID SegmentIDParam, reqEditors ...RequestEditorFn) (*GetSegmentByIDResponse, error) {
	rsp, err := c.GetSegmentByID(ctx, segmentID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSegmentByIDResponse(rsp)
}

// ParseSearchOrdersResponse parses an HTTP response from a SearchOrdersWithResponse call
func ParseSearchOrdersResponse(rsp *http.Response) (*SearchOrdersResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &SearchOrdersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ConvertOrder
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	}

	return response, nil
}

// ParseGetOrderByIDResponse parses an HTTP response from a GetOrderByIDWithResponse call
func ParseGetOrderByIDResponse(rsp *http.Response) (*GetOrderByIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetOrderByIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ConvertOrder
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseCancelOrderByIDResponse parses an HTTP response from a CancelOrderByIDWithResponse call
func ParseCancelOrderByIDResponse(rsp *http.Response) (*CancelOrderByIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CancelOrderByIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseSearchSegmentsByOrderIDResponse parses an HTTP response from a SearchSegmentsByOrderIDWithResponse call
func ParseSearchSegmentsByOrderIDResponse(rsp *http.Response) (*SearchSegmentsByOrderIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &SearchSegmentsByOrderIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ConvertSegment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseSearchSegmentsResponse parses an HTTP response from a SearchSegmentsWithResponse call
func ParseSearchSegmentsResponse(rsp *http.Response) (*SearchSegmentsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &SearchSegmentsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ConvertSegment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	}

	return response, nil
}

// ParseGetSegmentByIDResponse parses an HTTP response from a GetSegmentByIDWithResponse call
func ParseGetSegmentByIDResponse(rsp *http.Response) (*GetSegmentByIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetSegmentByIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ConvertSegment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ProblemDetails
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /orders)
	SearchOrders(ctx echo.Context) error

	// (GET /orders/{orderID})
	GetOrderByID(ctx echo.Context, orderID OrderIDParam) error

	// (GET /orders/{orderID}/cancel)
	CancelOrderByID(ctx echo.Context, orderID OrderIDParam) error

	// (GET /orders/{orderID}/segments)
	SearchSegmentsByOrderID(ctx echo.Context, orderID OrderIDParam) error

	// (GET /segments)
	SearchSegments(ctx echo.Context) error

	// (GET /segments/{segmentID})
	GetSegmentByID(ctx echo.Context, segmentID SegmentIDParam) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SearchOrders converts echo context to params.
func (w *ServerInterfaceWrapper) SearchOrders(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SearchOrders(ctx)
	return err
}

// GetOrderByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetOrderByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "orderID" -------------
	var orderID OrderIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "orderID", runtime.ParamLocationPath, ctx.Param("orderID"), &orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter orderID: %s", err))
	}

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetOrderByID(ctx, orderID)
	return err
}

// CancelOrderByID converts echo context to params.
func (w *ServerInterfaceWrapper) CancelOrderByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "orderID" -------------
	var orderID OrderIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "orderID", runtime.ParamLocationPath, ctx.Param("orderID"), &orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter orderID: %s", err))
	}

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CancelOrderByID(ctx, orderID)
	return err
}

// SearchSegmentsByOrderID converts echo context to params.
func (w *ServerInterfaceWrapper) SearchSegmentsByOrderID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "orderID" -------------
	var orderID OrderIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "orderID", runtime.ParamLocationPath, ctx.Param("orderID"), &orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter orderID: %s", err))
	}

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SearchSegmentsByOrderID(ctx, orderID)
	return err
}

// SearchSegments converts echo context to params.
func (w *ServerInterfaceWrapper) SearchSegments(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SearchSegments(ctx)
	return err
}

// GetSegmentByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetSegmentByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "segmentID" -------------
	var segmentID SegmentIDParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "segmentID", runtime.ParamLocationPath, ctx.Param("segmentID"), &segmentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter segmentID: %s", err))
	}

	ctx.Set(BasicAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSegmentByID(ctx, segmentID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/orders", wrapper.SearchOrders)
	router.GET(baseURL+"/orders/:orderID", wrapper.GetOrderByID)
	router.GET(baseURL+"/orders/:orderID/cancel", wrapper.CancelOrderByID)
	router.GET(baseURL+"/orders/:orderID/segments", wrapper.SearchSegmentsByOrderID)
	router.GET(baseURL+"/segments", wrapper.SearchSegments)
	router.GET(baseURL+"/segments/:segmentID", wrapper.GetSegmentByID)

}
