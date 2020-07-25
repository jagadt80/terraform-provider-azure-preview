package azurepreview

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-01-01/consumption"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

func resourceAzurePreviewBudget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzurePreviewBudgetCreate,
		ReadContext:   resourceAzurePreviewBudgetRead,
		UpdateContext: resourceAzurePreviewBudgetUpdate,
		DeleteContext: resourceAzurePreviewBudgetDelete,

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: stringInSlice([]string{
					string(consumption.Cost),
					string(consumption.Usage),
				}),
			},

			"amount": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"time_grain": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: stringInSlice([]string{
					string(consumption.TimeGrainTypeMonthly),
					string(consumption.TimeGrainTypeQuarterly),
					string(consumption.TimeGrainTypeAnnually),
					string(consumption.TimeGrainTypeBillingMonth),
					string(consumption.TimeGrainTypeBillingQuarter),
					string(consumption.TimeGrainTypeBillingAnnual),
				}),
			},

			"time_period": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_date": {
							Type:     schema.TypeString,
							Required: true,
						},

						"end_date": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"filters": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_groups": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"resources": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"meters": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: stringIsUUID,
							},
						},

						"tag": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"notification": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 5,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: stringInSlice([]string{
								string(consumption.EqualTo),
								string(consumption.GreaterThan),
								string(consumption.GreaterThanOrEqualTo),
							}),
						},

						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"contact_emails": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"contact_roles": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"contact_groups": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceAzurePreviewBudgetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Meta).Budgets

	scope := d.Get("scope").(string)
	budgetName := d.Get("name").(string)
	amount := decimal.NewFromInt(int64(d.Get("amount").(int)))

	params := consumption.Budget{
		BudgetProperties: &consumption.BudgetProperties{
			Category:      consumption.CategoryType(d.Get("category").(string)),
			Amount:        &amount,
			TimeGrain:     consumption.TimeGrainType(d.Get("time_grain").(string)),
			TimePeriod:    expandAzurePreviewBudgetTimePeriod(d.Get("time_period").([]interface{})),
			Filters:       expandAzurePreviewBudgetFilters(d.Get("filters").([]interface{})),
			Notifications: expandAzurePreviewBudgetNotifications(d.Get("notification").(*schema.Set).List()),
		},
	}

	resp, err := client.CreateOrUpdate(ctx, scope, budgetName, params)
	if err != nil {
		return diag.Errorf("error creating Budget %q (Scope %q): %+v", budgetName, scope, err)
	}

	d.SetId(*resp.ID)

	resourceAzurePreviewBudgetRead(ctx, d, meta)

	return diags
}

func resourceAzurePreviewBudgetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Budgets

	id, err := parseBudgetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.Get(ctx, id.Scope, id.BudgetName)
	if err != nil {
		if resp.IsHTTPStatus(404) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error reading Budget (ID %q): %+v", d.Id(), err)
	}

	d.Set("scope", id.Scope)
	d.Set("name", resp.Name)
	d.Set("amount", resp.Amount.IntPart())
	d.Set("time_grain", resp.TimeGrain)
	d.Set("time_period", flattenAzurePreviewBudgetTimePeriod(resp.TimePeriod))
	d.Set("filters", flattenAzurePreviewBudgetFilters(resp.Filters))
	d.Set("notification", flattenAzurePreviewBudgetNotifications(resp.Notifications))

	return diags
}

func resourceAzurePreviewBudgetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Budgets

	id, err := parseBudgetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	props := consumption.BudgetProperties{}

	if d.HasChange("amount") {
		amount := decimal.NewFromInt(int64(d.Get("amount").(int)))
		props.Amount = &amount
	}

	params := consumption.Budget{
		BudgetProperties: &props,
	}

	_, err = client.CreateOrUpdate(ctx, id.Scope, id.BudgetName, params)
	if err != nil {
		return diag.Errorf("error updating Budget %q (Scope %q): %+v", id.BudgetName, id.Scope, err)
	}

	resourceAzurePreviewBudgetRead(ctx, d, meta)

	return diags
}

func resourceAzurePreviewBudgetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Budgets

	id, err := parseBudgetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Delete(ctx, id.Scope, id.BudgetName)
	if err != nil {
		return diag.Errorf("error deleting Budget %q (Scope %q): %+v", id.BudgetName, id.Scope, err)
	}

	d.SetId("")

	return diags
}

func expandAzurePreviewBudgetTimePeriod(input []interface{}) *consumption.BudgetTimePeriod {
	if len(input) == 0 {
		return nil
	}

	values := input[0].(map[string]interface{})
	result := consumption.BudgetTimePeriod{}

	if v, ok := values["start_date"]; ok {
		startDate, _ := time.Parse(time.RFC3339, v.(string))
		result.StartDate = &date.Time{Time: startDate}
	}

	if v, ok := values["end_date"]; ok {
		endDate, _ := time.Parse(time.RFC3339, v.(string))
		result.EndDate = &date.Time{Time: endDate}
	}

	return &result
}

func expandAzurePreviewBudgetFilters(input []interface{}) *consumption.Filters {
	if len(input) == 0 {
		return nil
	}

	values := input[0].(map[string]interface{})
	result := consumption.Filters{}

	if v, ok := values["resource_groups"]; ok {
		result.ResourceGroups = expandStringSlice(v.([]interface{}))
	}

	if v, ok := values["resources"]; ok {
		result.Resources = expandStringSlice(v.([]interface{}))
	}

	if v, ok := values["meters"]; ok {
		ids := make([]uuid.UUID, 0)
		for _, item := range v.([]interface{}) {
			if item != nil {
				id, _ := uuid.FromString(item.(string))
				ids = append(ids, id)
			}
		}

		result.Meters = &ids
	}

	return &result
}

func expandAzurePreviewBudgetNotifications(input []interface{}) map[string]*consumption.Notification {
	if len(input) == 0 {
		return nil
	}

	results := make(map[string]*consumption.Notification)

	for _, item := range input {
		values := item.(map[string]interface{})
		result := consumption.Notification{}

		if v, ok := values["enabled"]; ok {
			result.Enabled = to.BoolPtr(v.(bool))
		}

		if v, ok := values["operator"]; ok {
			result.Operator = consumption.OperatorType(v.(string))
		}

		if v, ok := values["threshold"]; ok {
			threshold := decimal.NewFromInt(int64(v.(int)))
			result.Threshold = &threshold
		}

		if v, ok := values["contact_emails"]; ok {
			result.ContactEmails = expandStringSlice(v.([]interface{}))
		}

		if v, ok := values["contact_roles"]; ok {
			result.ContactRoles = expandStringSlice(v.([]interface{}))
		}

		if v, ok := values["contact_groups"]; ok {
			result.ContactGroups = expandStringSlice(v.([]interface{}))
		}

		if v, ok := values["name"]; ok {
			name := v.(string)
			results[name] = &result
		}
	}

	return results
}

func flattenAzurePreviewBudgetFilters(input *consumption.Filters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	values["resource_groups"] = input.ResourceGroups
	values["resources"] = input.Resources
	values["meters"] = input.Meters

	return []interface{}{values}
}

func flattenAzurePreviewBudgetTimePeriod(input *consumption.BudgetTimePeriod) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	values := make(map[string]interface{})

	values["start_date"] = input.StartDate.String()
	values["end_date"] = input.EndDate.String()

	return []interface{}{values}
}

func flattenAzurePreviewBudgetNotifications(input map[string]*consumption.Notification) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)

	for k, v := range input {
		values := make(map[string]interface{})

		values["name"] = k

		if v != nil {
			values["enabled"] = v.Enabled
			values["operator"] = string(v.Operator)
			values["threshold"] = v.Threshold.IntPart()
			values["contact_emails"] = v.ContactEmails
			values["contact_roles"] = v.ContactRoles
			values["contact_groups"] = v.ContactGroups
		}

		result = append(result, values)
	}

	return result
}
