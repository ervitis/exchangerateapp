// Code generated by go-swagger; DO NOT EDIT.

package rates

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new rates API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for rates API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetByDate(params *GetByDateParams, opts ...ClientOption) (*GetByDateOK, *GetByDateCreated, *GetByDateAccepted, error)

	GetLatest(params *GetLatestParams, opts ...ClientOption) (*GetLatestOK, *GetLatestCreated, *GetLatestAccepted, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetByDate gets historical rates from a gived date

  Returns the foreign exchange reference rates for an historical date. Rates are quoted against the Euro by default. Specify the symbols returned (default = all)
*/
func (a *Client) GetByDate(params *GetByDateParams, opts ...ClientOption) (*GetByDateOK, *GetByDateCreated, *GetByDateAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetByDateParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getByDate",
		Method:             "GET",
		PathPattern:        "/{date}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetByDateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, nil, nil, err
	}
	switch value := result.(type) {
	case *GetByDateOK:
		return value, nil, nil, nil
	case *GetByDateCreated:
		return nil, value, nil, nil
	case *GetByDateAccepted:
		return nil, nil, value, nil
	}
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for rates: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetLatest gets the latest foreign exchange reference rates

  Returns the latest foreign exchange reference rates. Rates are quoted against the Euro by default. Specify the symbols returned (default = all)
*/
func (a *Client) GetLatest(params *GetLatestParams, opts ...ClientOption) (*GetLatestOK, *GetLatestCreated, *GetLatestAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetLatestParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getLatest",
		Method:             "GET",
		PathPattern:        "/latest",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetLatestReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, nil, nil, err
	}
	switch value := result.(type) {
	case *GetLatestOK:
		return value, nil, nil, nil
	case *GetLatestCreated:
		return nil, value, nil, nil
	case *GetLatestAccepted:
		return nil, nil, value, nil
	}
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for rates: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
