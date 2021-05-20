package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/service/btpsaasprovisioning"
	"time"
)

func resourceSapBtpTenantApplicationSubscriptions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpTenantApplicationSubscriptionsCreate,
		ReadContext:   resourceSapBtpTenantApplicationSubscriptionsRead,
		UpdateContext: resourceSapBtpTenantApplicationSubscriptionsUpdate,
		DeleteContext: resourceSapBtpTenantApplicationSubscriptionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"skip_unchanged_dependencies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skip_updating_dependencies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"update_application_url": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"custom_properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpTenantApplicationSubscriptionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpSaaSProvisioningV1Client := meta.(*SAPClient).btpSaaSProvisioningV1Client

	tenantId := d.Get("tenant_id")
	input := &btpsaasprovisioning.SubscribeTenantToApplicationInput{
		TenantId: tenantId.(string),
	}
	if output, err := btpSaaSProvisioningV1Client.SubscribeTenantToApplication(ctx, input); err != nil {
		if output != nil && output.Error != "" {
			return diag.Errorf("BTP SaaS Subscription to an application can't be done; Operation code %v; %s",
				output.StatusCode, output.Error)
		} else {
			return diag.Errorf("BTP SaaS Subscription to an application can't be done;  %v", err)
		}
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			jobInput := &btpsaasprovisioning.GetJobStatusInput{
				JobId: output.JobStatusId,
			}
			if jobOut, err := btpSaaSProvisioningV1Client.GetJobStatus(ctx, jobInput); err != nil {
				return resource.RetryableError(err)
			} else {
				// IN_PROGRESS, COMPLETED, FAILED
				if jobOut.Status == "IN_PROGRESS" {
					return resource.RetryableError(
						fmt.Errorf("BTP SaaS Subscription to an application in progress; %s", jobOut.Description))
				} else if jobOut.Status == "FAILED" {
					return resource.NonRetryableError(
						fmt.Errorf("BTP SaaS Subscription to an application failed; %s", jobOut.Description))
				} else {
					return nil
				}
			}
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}
	}

	d.SetId(input.TenantId)

	return nil
}

func resourceSapBtpTenantApplicationSubscriptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	return nil
}

func resourceSapBtpTenantApplicationSubscriptionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpSaaSProvisioningV1Client := meta.(*SAPClient).btpSaaSProvisioningV1Client

	tenantId := d.Get("tenant_id")
	input := &btpsaasprovisioning.UpdateSubscriptionDependenciesInput{
		TenantId: tenantId.(string),
	}
	if val, ok := d.GetOk("skip_unchanged_dependencies"); ok {
		input.SkipUnchangedDependencies = val.(bool)
	}
	if val, ok := d.GetOk("skip_updating_dependencies"); ok {
		input.SkipUpdatingDependencies = val.(bool)
	}
	if val, ok := d.GetOk("update_application_url"); ok {
		input.UpdateApplicationURL = val.(bool)
	}
	if val, ok := d.GetOk("custom_properties"); ok {
		input.UpdateApplicationDependencies = val.(map[string]interface{})
	}

	if output, err := btpSaaSProvisioningV1Client.UpdateSubscriptionDependencies(ctx, input); err != nil {
		if output != nil && output.Error != "" {
			return diag.Errorf("BTP SaaS Update Subscription can't be done; Operation code %v; %s",
				output.StatusCode, output.Error)
		} else {
			return diag.Errorf("BTP SaaS Update Subscription can't be done;  %v", err)
		}
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			jobInput := &btpsaasprovisioning.GetJobStatusInput{
				JobId: output.JobStatusId,
			}
			if jobOut, err := btpSaaSProvisioningV1Client.GetJobStatus(ctx, jobInput); err != nil {
				return resource.RetryableError(err)
			} else {
				// IN_PROGRESS, COMPLETED, FAILED
				if jobOut.Status == "IN_PROGRESS" {
					return resource.RetryableError(
						fmt.Errorf("BTP SaaS Update Subscription in progress; %s", jobOut.Description))
				} else if jobOut.Status == "FAILED" {
					return resource.NonRetryableError(
						fmt.Errorf("BTP SaaS Update Subscription failed; %s", jobOut.Description))
				} else {
					return nil
				}
			}
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}
	}

	d.SetId(input.TenantId)

	return nil
}

func resourceSapBtpTenantApplicationSubscriptionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpSaaSProvisioningV1Client := meta.(*SAPClient).btpSaaSProvisioningV1Client

	tenantId := d.Get("tenant_id")
	input := &btpsaasprovisioning.UnSubscribeTenantFromApplicationInput{
		TenantId: tenantId.(string),
	}
	if output, err := btpSaaSProvisioningV1Client.UnSubscribeTenantFromApplication(ctx, input); err != nil {
		if output != nil && output.Error != "" {
			return diag.Errorf("BTP SaaS UnSubscribing from an application can't be done; Operation code %v; %s",
				output.StatusCode, output.Error)
		} else {
			return diag.Errorf("BTP SaaS UnSubscribing from an application can't be done;  %v", err)
		}
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			jobInput := &btpsaasprovisioning.GetJobStatusInput{
				JobId: output.JobStatusId,
			}
			if jobOut, err := btpSaaSProvisioningV1Client.GetJobStatus(ctx, jobInput); err != nil {
				return resource.RetryableError(err)
			} else {
				// IN_PROGRESS, COMPLETED, FAILED
				if jobOut.Status == "IN_PROGRESS" {
					return resource.RetryableError(
						fmt.Errorf("BTP SaaS UnSubscribing from an application in progress; %s", jobOut.Description))
				} else if jobOut.Status == "FAILED" {
					return resource.NonRetryableError(
						fmt.Errorf("BTP SaaS UnSubscribing from an application failed; %s", jobOut.Description))
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
