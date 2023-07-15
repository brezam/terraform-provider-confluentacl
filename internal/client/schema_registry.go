package client

import (
	"terraform-provider-confluentacl/internal/client/request"
)

type SchemaCluster struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type SchemaReadResponse struct {
	SchemaClusters []SchemaCluster `json:"clusters"`
}

const (
	schemaRegistryReadEndpoint = "schema_registries"
)

func (c *Client) GetFirstSchemaRegistry(environmentId string) (*SchemaCluster, error) {
	schemaReadResponse := &SchemaReadResponse{}
	response, err := c.RequestBuilder().
		Endpoint(schemaRegistryReadEndpoint).
		SetQueryParams(map[string]string{"account_id": environmentId}).
		Get().
		ExecuteAndRetryOn429()
	if err != nil {
		return nil, err
	}
	err = request.UnpackJSONResponse(response, &schemaReadResponse)
	if err != nil {
		return nil, err
	}
	return &schemaReadResponse.SchemaClusters[0], nil
}
