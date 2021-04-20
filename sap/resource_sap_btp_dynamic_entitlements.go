package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"time"
)

func resourceSapBtpDynamicEntitlements(plan string) *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpDynamicEntitlementCreate(plan),
		ReadContext:   resourceSapBtpDynamicEntitlementRead(plan),
		UpdateContext: resourceSapBtpDynamicEntitlementUpdate(plan),
		DeleteContext: resourceSapBtpDynamicEntitlementDelete(plan),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"assignment": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub_account_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"resource": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"provider": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"technical_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpDynamicEntitlementCreate(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		if uuidString, err := uuid.GenerateUUID(); err != nil {
			return diag.FromErr(err)
		} else {
			d.SetId(uuidString)
		}

		servicePlans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
		for spIdx := range servicePlans {
			service := servicePlans[spIdx]
			service.ServicePlanName = plan
			for infoIdx := range service.AssignmentInfo {
				service.AssignmentInfo[infoIdx].Amount = nil
				service.AssignmentInfo[infoIdx].Enable = sap.Bool(true)
			}
		}
		return entitlementsUpdateSubAccountServicePlan(ctx, "created", servicePlans, meta)
	}
}

func resourceSapBtpDynamicEntitlementRead(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

		return nil
	}
}

func resourceSapBtpDynamicEntitlementUpdate(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		servicePlans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
		for spIdx := range servicePlans {
			service := servicePlans[spIdx]
			service.ServicePlanName = plan
			for infoIdx := range service.AssignmentInfo {
				service.AssignmentInfo[infoIdx].Amount = nil
				service.AssignmentInfo[infoIdx].Enable = sap.Bool(true)
			}
		}
		return entitlementsUpdateSubAccountServicePlan(ctx, "updated", servicePlans, meta)
	}
}

func resourceSapBtpDynamicEntitlementDelete(plan string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		servicePlans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
		for spIdx := range servicePlans {
			service := servicePlans[spIdx]
			service.ServicePlanName = plan
			for infoIdx := range service.AssignmentInfo {
				service.AssignmentInfo[infoIdx].Amount = nil
				service.AssignmentInfo[infoIdx].Enable = sap.Bool(false)
			}
		}
		return entitlementsUpdateSubAccountServicePlan(ctx, "deleted", servicePlans, meta)
	}
}
