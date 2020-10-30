package azurepreview

import (
	"context"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const TerraformProviderUserAgent = "terraform-provider-azurepreview"

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"AZURE_SUBSCRIPTION_ID", "ARM_SUBSCRIPTION_ID"}, nil),
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"client_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"AZURE_CLIENT_ID", "ARM_CLIENT_ID"}, nil),
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"client_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"AZURE_CLIENT_SECRET", "ARM_CLIENT_SECRET"}, nil),
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"tenant_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"AZURE_TENANT_ID", "ARM_TENANT_ID"}, nil),
				RequiredWith:     []string{"client_id", "client_secret"},
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"environment": {
				Type:             schema.TypeString,
				Required:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"AZURE_ENVIRONMENT", "ARM_ENVIRONMENT"}, azure.PublicCloud.Name),
				ValidateDiagFunc: stringIsNotEmpty,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"azurepreview_resources": dataSourceAzurePreviewResources(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"azurepreview_subscription": resourceAzurePreviewSubscription(),
			"azurepreview_budget":       resourceAzurePreviewBudget(),
		},
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := &Config{
			SubscriptionID: d.Get("subscription_id").(string),
			ClientID:       d.Get("client_id").(string),
			ClientSecret:   d.Get("client_secret").(string),
			TenantID:       d.Get("tenant_id").(string),
			Environment:    d.Get("environment").(string),
		}

		ua := p.UserAgent(TerraformProviderUserAgent, p.TerraformVersion)

		return config.Client(ua)
	}
}
