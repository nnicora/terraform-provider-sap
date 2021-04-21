package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/sap"
	"time"
)

func resourceSapBtpDirectoryDynamicEntitlements(plan string) *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpDirectoryDynamicEntitlementCreate(plan),
		ReadContext:   resourceSapBtpDirectoryDynamicEntitlementRead(plan),
		UpdateContext: resourceSapBtpDirectoryDynamicEntitlementUpdate(plan),
		DeleteContext: resourceSapBtpDirectoryDynamicEntitlementDelete(plan),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"assignment": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"service_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"distribute": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpDirectoryDynamicEntitlementCreate(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		if uuidString, err := uuid.GenerateUUID(); err != nil {
			return diag.FromErr(err)
		} else {
			d.SetId(uuidString)
		}

		directoryId := d.Get("directory_id").(string)
		plans := buildDirectoryEntitlements(d.Get("assignment"))
		for idx := range plans {
			plans[idx].Amount = nil
			plans[idx].AutoDistributeAmount = nil
			plans[idx].Enable = sap.Bool(true)
			plans[idx].AutoAssign = true
		}
		return updateDirectoryEntitlements(ctx, "created", directoryId, plans, meta)
	}
}

func resourceSapBtpDirectoryDynamicEntitlementRead(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

		return nil
	}
}

func resourceSapBtpDirectoryDynamicEntitlementUpdate(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		directoryId := d.Get("directory_id").(string)
		plans := buildDirectoryEntitlements(d.Get("assignment"))
		for idx := range plans {
			plans[idx].Amount = nil
			plans[idx].AutoDistributeAmount = nil
			plans[idx].Enable = sap.Bool(true)
			plans[idx].AutoAssign = true
		}
		return updateDirectoryEntitlements(ctx, "updated", directoryId, plans, meta)
	}
}

func resourceSapBtpDirectoryDynamicEntitlementDelete(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		directoryId := d.Get("directory_id").(string)
		plans := buildDirectoryEntitlements(d.Get("assignment"))
		for idx := range plans {
			plans[idx].Amount = nil
			plans[idx].AutoDistributeAmount = nil
			plans[idx].Enable = sap.Bool(false)
		}
		return updateDirectoryEntitlements(ctx, "deleted", directoryId, plans, meta)
	}
}
