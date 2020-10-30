package azurepreview

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2019-10-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Config struct {
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	TenantID       string
	Environment    string
}

type Meta struct {
	Budgets       consumption.BudgetsClient
	Resources     resources.Client
	Subscription  subscription.Client
	Subscriptions subscriptions.Client
	StopContext   context.Context
}

func (c *Config) Client(userAgent string) (*Meta, diag.Diagnostics) {
	meta := Meta{
		StopContext: context.Background(),
	}

	authorizer, err := c.getAuthorizer()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	meta.Budgets = consumption.NewBudgetsClient(c.SubscriptionID)
	configureClient(&meta.Budgets.Client, userAgent, authorizer)

	meta.Resources = resources.NewClient(c.SubscriptionID)
	configureClient(&meta.Resources.Client, userAgent, authorizer)

	meta.Subscription = subscription.NewClient()
	configureClient(&meta.Subscription.Client, userAgent, authorizer)

	meta.Subscriptions = subscriptions.NewClient()
	configureClient(&meta.Subscriptions.Client, userAgent, authorizer)

	return &meta, nil
}

func configureClient(client *autorest.Client, userAgent string, authorizer autorest.Authorizer) {
	client.Authorizer = authorizer
	client.UserAgent = userAgent
}

func (c *Config) getAuthorizer() (autorest.Authorizer, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(token), nil
}

func (c *Config) getToken() (*adal.Token, error) {
	env, err := azure.EnvironmentFromName(c.Environment)
	if err != nil {
		return nil, err
	}

	if c.ClientID != "" && c.ClientSecret != "" && c.TenantID != "" {
		oauthConfig, err := adal.NewOAuthConfigWithAPIVersion(
			env.ActiveDirectoryEndpoint,
			c.TenantID,
			nil,
		)
		if err != nil {
			return nil, err
		}

		spToken, err := adal.NewServicePrincipalToken(
			*oauthConfig,
			c.ClientID,
			c.ClientSecret,
			env.ResourceManagerEndpoint)
		if err != nil {
			return nil, err
		}

		err = spToken.Refresh()
		if err != nil {
			return nil, err
		}

		adalToken := spToken.Token()

		return &adalToken, nil
	}

	cliToken, err := cli.GetTokenFromCLI(env.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	adalToken, err := cliToken.ToADALToken()
	if err != nil {
		return nil, err
	}

	return &adalToken, nil
}
