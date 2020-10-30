package azurepreview

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAzurePreviewResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzurePreviewResourcesRead,

		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"resource_group_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAzurePreviewResourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Resources

	tags := d.Get("tags").(map[string]interface{})

	if v, ok := d.GetOk("subscription_id"); ok {
		client.SubscriptionID = v.(string)
	}

	var filters []string

	if v, ok := d.GetOk("type"); ok {
		filters = append(filters, fmt.Sprintf("resourceType eq '%s'", v.(string)))
	}

	if v, ok := d.GetOk("name"); ok {
		filters = append(filters, fmt.Sprintf("name eq '%s'", v.(string)))
	}

	if v, ok := d.GetOk("resource_group_name"); ok {
		filters = append(filters, fmt.Sprintf("resourceGroup eq '%s'", v.(string)))
	}

	filter := strings.Join(filters, " and ")

	resp, err := client.ListComplete(ctx, filter, "", nil)
	if err != nil {
		return diag.Errorf("error reading resources: %+v", err)
	}

	resources := make([]map[string]interface{}, 0)

	for resp.NotDone() {
		resource := make(map[string]interface{})

		value := resp.Value()

		if v := value.ID; v != nil {
			resource["id"] = *v
		}

		if v := value.Name; v != nil {
			resource["name"] = *v
		}

		if v := value.Type; v != nil {
			resource["type"] = *v
		}

		if v := value.Location; v != nil {
			resource["location"] = *v
		}

		if err = resp.NextWithContext(ctx); err != nil {
			return diag.Errorf("error reading resources: %+v", err)
		}

		tagsFound := 0

		if value.Tags != nil {
			for requiredTagName, requiredTagValue := range tags {
				for tagName, tagValue := range value.Tags {
					if requiredTagName == tagName && requiredTagValue == *tagValue {
						tagsFound++
					}
				}
			}

			if tagsFound != len(tags) {
				continue
			}
		}

		resources = append(resources, resource)
	}

	id, _ := uuid.GenerateUUID()

	d.SetId(id)

	d.Set("resources", resources)

	return diags
}
