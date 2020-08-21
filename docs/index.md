# Azure Preview Provider

Use this provider to manage Azure resources.

## Example Usage

```hcl
provider "azure-preview" {}

resource "azurepreview_subscription" "example" {
  name               = "example"
  enrollment_account = "6d38255d-8321-4f17-8ddd-3bd94c57d988"
  offer_type         = "MS-AZR-0148P"
  provider           = azure-preview
}
```

## Argument Reference

* `subscription_id` - (Optional) The subscription ID. It can also be sourced from the `AZURE_SUBSCRIPTION_ID` environment variable.

* `client_id` - (Optional) The client ID. It can also be sourced from the `AZURE_CLIENT_ID` environment variable.

* `client_secret` - (Optional) The client secret. It can also be sourced from the `AZURE_CLIENT_SECRET` environment variable.

* `tenant_id` - (Optional) The tenant ID. It can also be sourced from the `AZURE_TENANT_ID` environment variable.

* `environment` - (Optional) The name of the Azure environment. It can also be sourced from the `AZURE_ENVIRONMENT` environment variable. Default is `AzurePublicCloud`.
