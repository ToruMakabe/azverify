package cmd

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

type Client struct {
	OperationsClient resourcegraph.OperationsClient
}

func buildClient(b authentication.Builder) (*Client, error) {
	builder := &authentication.Builder{
		SubscriptionID: b.SubscriptionID,
		ClientID:       b.ClientID,
		ClientSecret:   b.ClientSecret,
		TenantID:       b.TenantID,
		Environment:    b.Environment,

		SupportsAuxiliaryTenants: false,

		SupportsClientCertAuth: true,
		ClientCertPath:         b.ClientCertPath,
		ClientCertPassword:     b.ClientCertPassword,

		SupportsClientSecretAuth:       true,
		SupportsManagedServiceIdentity: b.SupportsManagedServiceIdentity,
		SupportsAzureCliToken:          true,
	}

	authConfig, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("error building AzureRM Client: %s", err)
	}

	env, err := authentication.DetermineEnvironment(authConfig.Environment)
	if err != nil {
		return nil, err
	}

	oauthConfig, err := authConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	if oauthConfig == nil {
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s", authConfig.TenantID)
	}

	sender := autorest.DecorateSender(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	})

	auth, err := authConfig.GetAuthorizationToken(sender, oauthConfig, "https://management.azure.com/")
	if err != nil {
		return nil, err
	}

	operationsClient := resourcegraph.NewOperationsClient()
	operationsClient.Authorizer = auth

	result := &Client{
		OperationsClient: operationsClient,
	}

	return result, nil
}
