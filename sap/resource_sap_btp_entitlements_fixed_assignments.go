package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"time"
)

func resourceSapBtpEntitlementFixedAssignments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpEntitlementFixedAssignmentsCreate,
		ReadContext:   resourceSapBtpEntitlementFixedAssignmentsRead,
		UpdateContext: resourceSapBtpEntitlementFixedAssignmentsUpdate,
		DeleteContext: resourceSapBtpEntitlementFixedAssignmentsDelete,
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
						"plan_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"assignment": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
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

func resourceSapBtpEntitlementFixedAssignmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	input := &btpentitlements.UpdateSubAccountServicePlanInput{
		SubAccountServicePlans: buildEntitlementsSubAccountServicePlan(d.Get("service")),
	}
	if output, err := btpEntitlementsV1Client.UpdateSubAccountServicePlan(ctx, input); err != nil {
		return diag.Errorf("BTP Sub Account assignment can't be created:  %v", err)
	} else if output.StatusCode != 202 {
		return diag.Errorf("BTP Sub Account assignment can't be created; Operation code %v; %v", output.StatusCode, output.Error.Message)
	}

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}

	return nil
}

func resourceSapBtpEntitlementFixedAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client
	d.SetId(d.Id())

	return nil
}

func resourceSapBtpEntitlementFixedAssignmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	input := &btpentitlements.UpdateSubAccountServicePlanInput{
		SubAccountServicePlans: buildEntitlementsSubAccountServicePlan(d.Get("service")),
	}
	if output, err := btpEntitlementsV1Client.UpdateSubAccountServicePlan(ctx, input); err != nil {
		return diag.Errorf("BTP Sub Account assignment can't be updated:  %v", err)
	} else if output.StatusCode != 202 {
		return diag.Errorf("BTP Sub Account assignment can't be updated; Operation code %v; %v", output.StatusCode, output.Error.Message)
	}

	d.SetId(d.Id())

	return nil
}

func resourceSapBtpEntitlementFixedAssignmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	plans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
	for _, plan := range plans {
		for _, assignmentInfo := range plan.AssignmentInfo {
			if assignmentInfo.Amount != nil {
				assignmentInfo.Amount = sap.Float32(0)
			}
			if assignmentInfo.Amount != nil {
				assignmentInfo.Enable = sap.Bool(false)
			}
		}
	}
	input := &btpentitlements.UpdateSubAccountServicePlanInput{
		SubAccountServicePlans: plans,
	}
	if output, err := btpEntitlementsV1Client.UpdateSubAccountServicePlan(ctx, input); err != nil {
		return diag.Errorf("BTP Sub Account assignment can't be delete:  %v", err)
	} else if output.StatusCode != 202 {
		return diag.Errorf("BTP Sub Account assignment can't be delete; Operation code %v; %v", output.StatusCode, output.Error.Message)
	}

	d.SetId("")

	return nil
}

func buildEntitlementsSubAccountServicePlan(data interface{}) []btpentitlements.SubAccountServicePlan {
	result := make([]btpentitlements.SubAccountServicePlan, 0)

	if data == nil {
		return result
	}

	maps, ok := data.([]map[string]interface{})
	if !ok {
		return result
	}

	for _, m := range maps {
		elem := btpentitlements.SubAccountServicePlan{}
		if val, ok := m["name"]; ok && val != nil {
			elem.ServiceName = val.(string)
		}
		if val, ok := m["plan_name"]; ok && val != nil {
			elem.ServicePlanName = val.(string)
		}
		if val, ok := m["assignment"]; ok && val != nil {
			elem.AssignmentInfo = buildEntitlementsAssignments(val)
		}

		result = append(result, elem)
	}
	return result
}

func buildEntitlementsAssignments(data interface{}) []btpentitlements.AssignmentInfo {
	result := make([]btpentitlements.AssignmentInfo, 0)

	if data == nil {
		return result
	}

	maps, ok := data.([]map[string]interface{})
	if !ok {
		return result
	}

	for _, m := range maps {
		elem := btpentitlements.AssignmentInfo{}
		if val, ok := m["amount"]; ok && val != nil {
			elem.Amount = sap.Float32(val.(float32))
		}
		if val, ok := m["sub_account_id"]; ok && val != nil {
			elem.SubAccountGuid = val.(string)
		}
		if val, ok := m["resource"]; ok && val != nil {
			elem.Resources = buildEntitlementsResources(val)
		}

		result = append(result, elem)
	}
	return result
}

func buildEntitlementsResources(data interface{}) []btpentitlements.Resource {
	result := make([]btpentitlements.Resource, 0)

	if data == nil {
		return result
	}

	maps, ok := data.([]map[string]interface{})
	if !ok {
		return result
	}

	for _, m := range maps {
		elem := btpentitlements.Resource{}
		if val, ok := m["name"]; ok && val != nil {
			elem.Name = val.(string)
		}
		if val, ok := m["provider"]; ok && val != nil {
			elem.Provider = val.(string)
		}
		if val, ok := m["technical_name"]; ok && val != nil {
			elem.TechnicalName = val.(string)
		}
		if val, ok := m["type"]; ok && val != nil {
			elem.Type = val.(string)
		}
		if val, ok := m["data"]; ok && val != nil {
			elem.Data = val.(string)
		}
		result = append(result, elem)
	}
	return result
}
