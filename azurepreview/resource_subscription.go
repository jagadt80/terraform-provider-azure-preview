package azurepreview

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2019-10-01-preview/subscription"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzurePreviewSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzurePreviewSubscriptionCreate,
		ReadContext:   resourceAzurePreviewSubscriptionRead,
		UpdateContext: resourceAzurePreviewSubscriptionUpdate,
		DeleteContext: resourceAzurePreviewSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: stringLengthBetween(1, 60),
			},

			"enrollment_account": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: stringIsNotEmpty,
			},

			"owners": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"offer_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: stringInSlice([]string{
					string(subscription.MSAZR0017P),
					string(subscription.MSAZR0148P),
				}),
			},

			"additional_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAzurePreviewSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Meta).Subscription

	name := d.Get("name").(string)
	enrollmentAccount := d.Get("enrollment_account").(string)
	owners := d.Get("owners").([]interface{})

	adPrincipals := make([]subscription.AdPrincipal, 0)
	for _, owner := range owners {
		adPrincipal := subscription.AdPrincipal{
			ObjectID: to.StringPtr(owner.(string)),
		}

		adPrincipals = append(adPrincipals, adPrincipal)
	}

	params := subscription.CreationParameters{
		DisplayName: &name,
		Owners:      &adPrincipals,
		OfferType:   subscription.OfferType(d.Get("offer_type").(string)),
	}

	future, err := client.CreateSubscriptionInEnrollmentAccount(ctx, enrollmentAccount, params)
	if err != nil {
		return diag.Errorf("error creating Subscription %q in Enrollment Account %q: %+v", name, enrollmentAccount, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return diag.Errorf("error waiting for Subscription %q in Enrollment Account %q to finish creating: %+v", name, enrollmentAccount, err)
	}

	resp, err := future.Result(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.SubscriptionLink)

	resourceAzurePreviewSubscriptionRead(ctx, d, meta)

	return diags
}

func resourceAzurePreviewSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Subscriptions

	subscriptionID, err := parseSubscriptionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.Get(ctx, subscriptionID)
	if err != nil {
		if resp.IsHTTPStatus(404) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error reading Subscription (ID %q): %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("subscription_id", resp.SubscriptionID)
	d.Set("tenant_id", resp.TenantID)

	return diags
}

func resourceAzurePreviewSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Subscription

	subscriptionID, err := parseSubscriptionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	subscriptionName := subscription.Name{
		SubscriptionName: to.StringPtr(d.Get("name").(string)),
	}

	if d.HasChange("name") {
		_, err := client.Rename(ctx, subscriptionID, subscriptionName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resourceAzurePreviewSubscriptionRead(ctx, d, meta)

	return diags
}

func resourceAzurePreviewSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*Meta).Subscription

	subscriptionID, err := parseSubscriptionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Cancel(ctx, subscriptionID)
	if err != nil {
		return diag.Errorf("error cancelling Subscription (ID %q): %+v", d.Id(), err)
	}

	d.SetId("")

	return diags
}
