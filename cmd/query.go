package cmd

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph"
)

func rgQuery(c *Client, sub, query string) (resourcegraph.QueryResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RequestOptions := resourcegraph.QueryRequestOptions{
		ResultFormat: "objectArray",
	}

	subList := []string{sub}

	Request := resourcegraph.QueryRequest{
		Subscriptions: &subList,
		Query:         &query,
		Options:       &RequestOptions,
	}

	results, err := c.OperationsClient.Resources(ctx, Request)
	if err != nil {
		return results, err
	}

	return results, nil
}
