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

func resourceSapBtpDirectoryEntitlements() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpDirectoryFixedEntitlementsCreate,
		ReadContext:   resourceSapBtpDirectoryFixedEntitlementsRead,
		UpdateContext: resourceSapBtpDirectoryFixedEntitlementsUpdate,
		DeleteContext: resourceSapBtpDirectoryFixedEntitlementsDelete,
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
						"amount": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"distribute": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						//"auto_assign": {
						//	Type:     schema.TypeBool,
						//	Optional: true,
						//},
						"auto_distribute_amount": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpDirectoryFixedEntitlementsCreate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}
	directoryId := d.Get("directory_id").(string)
	plans := buildDirectoryEntitlements(d.Get("assignment"))
	for idx := range plans {
		plans[idx].Enable = nil
		plans[idx].AutoAssign = true
	}
	return updateDirectoryEntitlements(ctx, "created", directoryId, plans, meta)
}

func resourceSapBtpDirectoryFixedEntitlementsRead(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceSapBtpDirectoryFixedEntitlementsUpdate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	directoryId := d.Get("directory_id").(string)
	plans := buildDirectoryEntitlements(d.Get("assignment"))
	for idx := range plans {
		plans[idx].Enable = nil
		plans[idx].AutoAssign = true
	}
	return updateDirectoryEntitlements(ctx, "updated", directoryId, plans, meta)
}

func resourceSapBtpDirectoryFixedEntitlementsDelete(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	directoryId := d.Get("directory_id").(string)
	plans := buildDirectoryEntitlements(d.Get("assignment"))
	for idx := range plans {
		plans[idx].Enable = nil
		plans[idx].Amount = sap.Uint(0)
		plans[idx].Distribute = false
		plans[idx].AutoAssign = true
		plans[idx].AutoDistributeAmount = sap.Uint(0)
	}
	return updateDirectoryEntitlements(ctx, "deleted", directoryId, plans, meta)
}

func updateDirectoryEntitlements(ctx context.Context, operation string,
	directoryId string, entitlements []btpentitlements.DirectoryEntitlement, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	input := &btpentitlements.UpdateDirectoryEntitlementsInput{
		DirectoryGuid:         directoryId,
		DirectoryEntitlements: entitlements,
	}
	if output, err := btpEntitlementsV1Client.UpdateDirectoryEntitlements(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Directory Entitlements can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Directory Entitlements can't be %s;  %v", operation, err)
		}
	} else if output.StatusCode != 200 {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Directory Entitlements can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Directory Entitlements can't be %s; Operation code %v", operation, output.StatusCode)
		}
	}
	//else {
	//	logDebug(output, "Directory Entitlements Output")
	//	if len(sap.StringValue(output.JobStatusId)) > 0 {
	//		jobInput := &btpentitlements.GetJobStatusInput{
	//			JobId: sap.StringValue(output.JobStatusId),
	//		}
	//		logDebug(jobInput, "Directory Entitlements GetJobStatus")
	//		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
	//			if jobOut, err := btpEntitlementsV1Client.GetJobStatus(ctx, jobInput); err != nil {
	//				return resource.RetryableError(err)
	//			} else {
	//				// IN_PROGRESS, COMPLETED, FAILED
	//				if jobOut.Status == "IN_PROGRESS" {
	//					return resource.RetryableError(
	//						fmt.Errorf("BTP Directory Entitlements in progress; %s", jobOut.Description))
	//				} else if jobOut.Status == "FAILED" {
	//					return resource.NonRetryableError(
	//						fmt.Errorf("BTP Directory Entitlements failed; %s", jobOut.Description))
	//				} else {
	//					return nil
	//				}
	//			}
	//		})
	//
	//		if retryErr != nil && isResourceTimeoutError(retryErr) {
	//			return diag.FromErr(retryErr)
	//		}
	//	}
	//}

	return nil
}

func buildDirectoryEntitlements(data interface{}) []btpentitlements.DirectoryEntitlement {
	if data == nil {
		return nil
	}

	datas, ok := data.([]interface{})
	if !ok {
		return nil
	}

	result := make([]btpentitlements.DirectoryEntitlement, 0)
	for idx := range datas {
		m, ok := datas[idx].(map[string]interface{})
		if !ok {
			continue
		}
		elem := btpentitlements.DirectoryEntitlement{}
		if val, ok := m["plan_name"]; ok && val != nil {
			elem.Plan = val.(string)
		}
		if val, ok := m["service_name"]; ok && val != nil {
			elem.Service = val.(string)
		}
		if val, ok := m["amount"]; ok && val != nil {
			elem.Amount = sap.Uint(uint(val.(int)))
		}
		if val, ok := m["distribute"]; ok && val != nil {
			elem.Distribute = val.(bool)
		}
		//if val, ok := m["auto_assign"]; ok && val != nil {
		//	elem.AutoAssign = val.(bool)
		//}
		if val, ok := m["auto_distribute_amount"]; ok && val != nil {
			elem.AutoDistributeAmount = sap.Uint(uint(val.(int)))
		}

		result = append(result, elem)
	}
	return result
}
