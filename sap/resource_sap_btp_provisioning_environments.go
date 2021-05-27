package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpprovisioning"
	"github.com/pkg/errors"
	"time"
)

func resourceSapBtpProvisioningEnvironments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpProvisioningEnvironmentsCreate,
		ReadContext:   resourceSapBtpProvisioningEnvironmentsRead,
		//UpdateContext: resourceSapBtpProvisioningEnvironmentsUpdate,
		UpdateContext: resourceSapBtpProvisioningEnvironmentsRead,
		DeleteContext: resourceSapBtpProvisioningEnvironmentsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"provisioning_service": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"oauth2": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"grant_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "client_credentials",
										Description: "SAP OAuth2 Grant Type.",
									},
									"client_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SAP OAuth2 Client Id.",
									},
									"client_secret": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SAP OAuth2 Client Secret.",
									},
									"token_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SAP OAuth2 Token Url.",
									},
									"authorization_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP OAuth2 Authorization Url.",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP OAuth2 Redirect Url.",
									},

									"username": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP OAuth2 Username. Used in case if 'grant_type=password'.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP OAuth2 Password. Used in case if 'grant_type=password'.",
									},

									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     60,
										Description: "SAP OAuth2 HTTP Client timeout.",
									},
								},
							},
						},
					},
				},
			},

			"environment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"plan_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"technical_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"landscape_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// computed
			"broker_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"commercial_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"dashboard_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"global_account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"modified_date": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"operation": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"platform_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"state_message": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"sub_account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpProvisioningEnvironmentsCreate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpProvisioningV1Client := meta.(*SAPClient).btpProvisioningV1Client
	session := meta.(*SAPClient).session

	serviceList := d.Get("provisioning_service").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("SaaS manager service is required")
	}

	err := session.AddEndpointWithReplace(btpprovisioning.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Provisioning Service OAuth2;  %v", err))
	}

	btpProvisioningV1Client := btpprovisioning.New(session)

	input := &btpprovisioning.CreateEnvironmentInstanceInput{
		EnvironmentType: d.Get("environment_type").(string),
		PlanName:        d.Get("plan_name").(string),
		TechnicalKey:    d.Get("technical_key").(string),
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}
	if val, ok := d.GetOk("landscape_label"); ok {
		input.LandscapeLabel = val.(string)
	}
	if val, ok := d.GetOk("name"); ok {
		input.Name = val.(string)
	}
	if val, ok := d.GetOk("origin"); ok {
		input.Origin = val.(string)
	}
	if val, ok := d.GetOk("service_name"); ok {
		input.ServiceName = val.(string)
	}
	if val, ok := d.GetOk("user"); ok {
		input.User = val.(string)
	}
	if val, ok := d.GetOk("parameters"); ok {
		if m, isMap := val.(map[string]interface{}); isMap {
			input.Parameters = m
		}
	}

	if output, err := btpProvisioningV1Client.CreateEnvironmentInstance(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Provisioning Environment can't be created; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Provisioning Environment can't be created;  %v", err)
		}
	} else {
		d.SetId(output.Id)
	}

	return nil
}

func resourceSapBtpProvisioningEnvironmentsRead(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpProvisioningV1Client := meta.(*SAPClient).btpProvisioningV1Client
	session := meta.(*SAPClient).session
	serviceList := d.Get("provisioning_service").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("Provisioning service is required")
	}

	err := session.AddEndpointWithReplace(btpprovisioning.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Provisioning Service OAuth2;  %v", err))
	}
	btpProvisioningV1Client := btpprovisioning.New(session)

	input := &btpprovisioning.GetEnvironmentInstanceInput{
		EnvironmentInstanceId: d.Id(),
	}

	if output, err := btpProvisioningV1Client.GetEnvironmentInstance(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Provisioning Environment can't be read; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Provisioning Environment can't be read;  %v", err)
		}
	} else {
		d.SetId(output.Id)
		d.Set("environment_type", output.EnvironmentType)
		d.Set("plan_name", output.PlanName)
		d.Set("description", output.Description)
		d.Set("landscape_label", output.LandscapeLabel)
		d.Set("name", output.Name)
		d.Set("service_name", output.ServiceName)

		// computed
		d.Set("broker_id", output.BrokerId)
		d.Set("commercial_type", output.CommercialType)
		d.Set("created_date", output.CreatedDate)
		d.Set("dashboard_url", output.DashboardUrl)
		d.Set("global_account_id", output.GlobalAccountGuid)
		d.Set("labels", output.Labels)
		d.Set("modified_date", output.ModifiedDate)
		d.Set("operation", output.Operation)
		//d.Set("parameters", output.Parameters)
		d.Set("plan_id", output.PlanId)
		d.Set("platform_id", output.PlatformId)
		d.Set("service_id", output.ServiceId)
		d.Set("state", output.State)
		d.Set("state_message", output.StateMessage)
		d.Set("sub_account_id", output.SubAccountGuid)
		d.Set("tenant_id", output.TenantId)
		d.Set("type", output.Type)
	}

	return nil
}

func resourceSapBtpProvisioningEnvironmentsUpdate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpProvisioningV1Client := meta.(*SAPClient).btpProvisioningV1Client
	//session := meta.(*SAPClient).session
	//serviceList := d.Get("provisioning_service").([]interface{})
	//if len(serviceList) < 1 {
	//	return diag.Errorf("Provisioning service is required")
	//}
	//
	//	err := session.AddEndpointWithReplace(btpprovisioning.EndpointsID, extractEndpointConfig(serviceList))
	//	if err != nil {
	//		return diag.FromErr(errors.Errorf("BTP Provisioning Service OAuth2;  %v", err))
	//	}
	//btpProvisioningV1Client := btpprovisioning.New(session)
	//
	//input := &btpprovisioning.UpdateEnvironmentInstanceInput{
	//	EnvironmentInstanceId: d.Id(),
	//	PlanName:              d.Get("plan_name").(string),
	//}
	//
	//if val, ok := d.GetOk("parameters"); ok {
	//	if m, isMap := val.(map[string]interface{}); isMap {
	//		input.Parameters = m
	//	}
	//}
	//if output, err := btpProvisioningV1Client.UpdateEnvironmentInstance(ctx, input); err != nil {
	//	if output != nil && output.Error != nil {
	//		return diag.Errorf("BTP Provisioning Environment can't be updated; Operation code %v; %s",
	//			output.StatusCode, sap.StringValue(output.Error.Message))
	//	} else {
	//		return diag.Errorf("BTP Provisioning Environment can't be updated;  %v", err)
	//	}
	//}
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

func resourceSapBtpProvisioningEnvironmentsDelete(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpProvisioningV1Client := meta.(*SAPClient).btpProvisioningV1Client
	session := meta.(*SAPClient).session
	serviceList := d.Get("provisioning_service").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("Provisioning service is required")
	}

	err := session.AddEndpointWithReplace(btpprovisioning.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Provisioning Service OAuth2;  %v", err))
	}

	btpProvisioningV1Client := btpprovisioning.New(session)

	input := &btpprovisioning.DeleteEnvironmentInstanceInput{
		EnvironmentInstanceId: d.Id(),
	}
	if output, err := btpProvisioningV1Client.DeleteEnvironmentInstance(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Provisioning Environment can't be deleted; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Provisioning Environment can't be deleted;  %v", err)
		}
	}
	return nil
}
