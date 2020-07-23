# azurepreview_subscription Resource

Creates an Azure subscription.

## Example Usage

```hcl
resource "azurepreview_subscription" "example" {
  name               = "example"
  enrollment_account = "6d38255d-8321-4f17-8ddd-3bd94c57d988"
  offer_type         = "MS-AZR-0148P"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name of the subscription.

* `enrollment_account` - (Required) The name of the enrollment account to which the subscription will be billed.

* `owners` - (Optional) The list of principals that should be granted `Owner` access on the subscription. Principals should be of type `User`, `Service Principal` or `Security Group`.

* `offer_type` - (Optional) The offer type of the subscription. Only valid when creating a subscription in a enrollment account scope. Possible values include: `MS-AZR-0017P` (production use), `MS-AZR-0148P` (dev/test).

## Attributes Reference

* `id` - The fully qualified ID for the subscription. Example: `/subscriptions/00000000-0000-0000-0000-000000000000`.

* `subscription_id` - The subscription ID.
