package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"time"
)

func resourceSapBtpEntitlements() *schema.Resource {
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"plan_name": {
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
									"amount": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"sub_account_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
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

func resourceSapBtpEntitlementFixedAssignmentsCreate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}
	plans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
	return entitlementsUpdateSubAccountServicePlan(ctx, "created", plans, meta)
}

func resourceSapBtpEntitlementFixedAssignmentsRead(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceSapBtpEntitlementFixedAssignmentsUpdate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	plans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
	return entitlementsUpdateSubAccountServicePlan(ctx, "updated", plans, meta)
}

func resourceSapBtpEntitlementFixedAssignmentsDelete(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	plans := buildEntitlementsSubAccountServicePlan(d.Get("service"))
	for planIdx := range plans {
		assInfos := plans[planIdx].AssignmentInfo
		for assInfoIdx := range assInfos {
			if assInfos[assInfoIdx].Amount != nil {
				assInfos[assInfoIdx].Amount = sap.Uint(0)
			}
			if assInfos[assInfoIdx].Enable != nil {
				assInfos[assInfoIdx].Enable = nil
			}
		}
	}
	return entitlementsUpdateSubAccountServicePlan(ctx, "deleted", plans, meta)
}

func entitlementsUpdateSubAccountServicePlan(ctx context.Context, operation string,
	servicePlans []btpentitlements.SubAccountServicePlan, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	input := &btpentitlements.UpdateSubAccountServicePlanInput{
		SubAccountServicePlans: servicePlans,
	}
	if output, err := btpEntitlementsV1Client.UpdateSubAccountServicePlan(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Sub Account Entitlements can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Sub Account Entitlements can't be %s;  %v", operation, err)
		}
	} else if output.StatusCode != 202 {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Sub Account Entitlements can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Sub Account Entitlements can't be %s; Operation code %v",
				operation, output.StatusCode)
		}
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			jobInput := &btpentitlements.GetJobStatusInput{
				JobId: sap.StringValue(output.JobStatusId),
			}
			if jobOut, err := btpEntitlementsV1Client.GetJobStatus(ctx, jobInput); err != nil {
				return resource.RetryableError(err)
			} else {
				// IN_PROGRESS, COMPLETED, FAILED
				if jobOut.Status == "IN_PROGRESS" {
					return resource.RetryableError(
						fmt.Errorf("BTP Sub Account Entitlements in progress; %s", jobOut.Description))
				} else if jobOut.Status == "FAILED" {
					return resource.NonRetryableError(
						fmt.Errorf("BTP Sub Account Entitlements failed; %s", jobOut.Description))
				} else {
					return nil
				}
			}
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}
	}

	return nil
}

func buildEntitlementsSubAccountServicePlan(data interface{}) []btpentitlements.SubAccountServicePlan {
	if data == nil {
		return nil
	}

	datas, ok := data.([]interface{})
	if !ok {
		return nil
	}

	result := make([]btpentitlements.SubAccountServicePlan, 0)
	for idx := range datas {
		m, ok := datas[idx].(map[string]interface{})
		if !ok {
			continue
		}
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
	if data == nil {
		return nil
	}

	datas, ok := data.([]interface{})
	if !ok {
		return nil
	}

	result := make([]btpentitlements.AssignmentInfo, 0)
	for idx := range datas {
		m, ok := datas[idx].(map[string]interface{})
		if !ok {
			continue
		}

		elem := btpentitlements.AssignmentInfo{}
		if val, ok := m["amount"]; ok && val != nil {
			elem.Amount = sap.Uint(uint(val.(int)))
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
	if data == nil {
		return nil
	}

	datas, ok := data.([]interface{})
	if !ok {
		return nil
	}

	result := make([]btpentitlements.Resource, 0)
	for idx := range datas {
		m, ok := datas[idx].(map[string]interface{})
		if !ok {
			continue
		}

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
