package service

import (
	"net/url"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

// ApplyQueries sends the list of Queries to be applied (upserted) to the
// Fleet instance.
func (c *Client) ApplyQueries(specs []*fleet.QuerySpec) error {
	req := applyQuerySpecsRequest{Specs: specs}
	verb, path := "POST", "/api/latest/fleet/spec/queries"
	var responseBody applyQuerySpecsResponse
	return c.authenticatedRequest(req, verb, path, &responseBody)
}

// GetQuery retrieves the list of all Queries.
func (c *Client) GetQuery(name string) (*fleet.QuerySpec, error) {
	verb, path := "GET", "/api/latest/fleet/spec/queries/"+url.PathEscape(name)
	var responseBody getQuerySpecResponse
	err := c.authenticatedRequest(nil, verb, path, &responseBody)
	return responseBody.Spec, err
}

// GetQueries retrieves the list of all Queries.
func (c *Client) GetQueries() ([]fleet.Query, error) {
	verb, path := "GET", "/api/latest/fleet/queries"
	var responseBody listQueriesResponse
	err := c.authenticatedRequest(nil, verb, path, &responseBody)
	return responseBody.Queries, err
}

// DeleteQuery deletes the query with the matching name.
func (c *Client) DeleteQuery(name string) error {
	verb, path := "DELETE", "/api/latest/fleet/queries/"+url.PathEscape(name)
	var responseBody deleteQueryResponse
	return c.authenticatedRequest(nil, verb, path, &responseBody)
}
