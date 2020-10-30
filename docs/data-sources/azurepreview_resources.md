# azurepreview_resources Data Source

Use this data source to get information about Azure resources.

## Example Usage

```hcl
data "azurepreview_resources" "example" {
  type = "Microsoft.Network/virtualNetworks"
  tags = {
    is_spoke = true
  }
}
```

## Argument Reference

* `subscription_id` - (Optional) The ID of the subscription.

* `name` - (Optional) The name of the resource.

* `resource_group_name` (Optional) The name of the resource group.

* `type` - (Optional) The type of resource. Example: `Microsoft.Network/virtualNetworks`.

* `tags` - (Optional) A mapping of tags used to filter the resources.

## Attribute Reference

* `resources` - One or more `resource` blocks as defined below.

The `resource` block contains:

* `id` - The ID of the resource.

* `name` - The name of the resource.

* `type` - The type of resource.

* `location` - The Azure region where the resource exists.
