# azurepreview_budget Resource

Creates an Azure budget.

## Example Usage

```hcl
resource "azurepreview_budget" "example" {
  name       = "example"
  scope      = "subscriptions/00000000-0000-0000-0000-000000000000"
  category   = "Cost"
  amount     = 1000
  time_grain = "BillingMonth"

  time_period {
    start_date = "2017-06-01T00:00:00Z"
    end_date   = "2035-06-01T00:00:00Z"
  }

  notification {
    name      = "notify-roles"
    operator  = "GreaterThan"
    threshold = 80
    contact_roles = [
      "Contributor",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the budget.

* `scope` - (Required) The scope of the budget. This includes `subscriptions/{subscriptionId}` for subscription scope, `subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}` for Resource Group scope, `providers/Microsoft.Billing/enrollmentAccounts/{enrollmentAccountId}` for Enrollment Account scope, `providers/Microsoft.Management/managementGroups/{managementGroupId}` for Management Group scope.

* `category` - (Required) The category of the budget, whether the budget tracks cost or usage. Possible values are: `Cost` and `Usage`.

* `amount` - (Required) The total amount of cost to track with the budget.

* `time_grain` - (Required) The time covered by a budget. Tracking of the amount will be reset based on the time grain. Possible values are: `Monthly`, `Quarterly`, `Annually`, `BillingMonth`, `BillingQuarter` and `BillingAnnual`.

* `time_period` - (Required) A `time_period` block as defined below. Has start and end date of the budget. The `start_date` must be first of the month and should be less than the `end_date`. Budget `start_date` must be on or after `June 1, 2017`. Future `start_date` should not be more than three months. Past `start_date` should be selected within the timegrain period. There are no restrictions on the `end_date`.

* `filters` - (Optional) A `filters` block as defined below. May be used to filter budgets by resource group, resource, or meter.

* `notification` - (Optional) A `notification` block as defined below. Notifications associated with the budget. Budget can have up to five notifications.

---

A `time_period` block supports the following:

* `start_date` - (Required) The start date for the budget.

* `end_date` - (Required) The end date for the budget.

---

A `filters` block supports the following:

* `resource_groups` - (Optional) The list of filters on resource groups, allowed at subscription level only.

* `resources` - (Optional) The list of filters on resources.

* `meters` - (Optional) The list of filters on meters (GUID), mandatory for budgets of usage category.

* `tag` - (Optional) A `tag` block as defined below.

---

A `tag` block supports the following:

* `name` - (Required) The name of the tag.

* `values` - (Required) List of values for the tag.

---

A `notification` block supports the following:

* `enabled` - (Optional) Whether the notification is enabled or not. Default is `true`.

* `operator` - (Required) The comparison operator. Possible values include: `EqualTo`, `GreaterThan`, `GreaterThanOrEqualTo`.

* `threshold` - (Required) Threshold value associated with a notification. Notification is sent when the cost exceeded the threshold. It is always percent and has to be between `0` and `1000`.

* `contact_emails` - (Optional) List of email addresses to send the budget notification to when the threshold is exceeded.

* `contact_roles` - (Optional) List of contact roles to send the budget notification to when the threshold is exceeded.

* `contact_groups` - (Optional) List of action groups to send the budget notification to when the threshold is exceeded.
