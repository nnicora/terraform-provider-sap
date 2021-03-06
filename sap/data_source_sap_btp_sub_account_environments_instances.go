package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpprovisioning"
	"github.com/pkg/errors"
)

func dataSourceSapBtpSubAccountEnvironmentsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpSubAccountEnvironmentsInstancesRead,
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

			"environments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"commercial_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dashboard_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"environment_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"global_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"landscape_label": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"modified_date": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"plan_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"plan_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"platform_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"state_message": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpSubAccountEnvironmentsInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if output, err := btpProvisioningV1Client.GetEnvironmentInstances(ctx); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Sub Account Environments can't be read; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Sub Account Environments can't be read;  %v", err)
		}
	} else {
		result := make([]map[string]interface{}, 0, len(output.Environments))
		for _, outEnv := range output.Environments {
			m := map[string]interface{}{
				"broker_id":         outEnv.BrokerId,
				"commercial_type":   outEnv.CommercialType,
				"created_date":      outEnv.CreatedDate,
				"dashboard_url":     outEnv.DashboardUrl,
				"description":       outEnv.Description,
				"environment_type":  outEnv.EnvironmentType,
				"global_account_id": outEnv.GlobalAccountGuid,
				"id":                outEnv.Id,
				"labels":            outEnv.Labels,
				"landscape_label":   outEnv.LandscapeLabel,
				"modified_date":     outEnv.ModifiedDate,
				"name":              outEnv.Name,
				"operation":         outEnv.Operation,
				"parameters":        outEnv.Parameters,
				"plan_id":           outEnv.PlanId,
				"plan_name":         outEnv.PlanName,
				"platform_id":       outEnv.PlatformId,
				"service_id":        outEnv.ServiceId,
				"service_name":      outEnv.ServiceName,
				"state":             outEnv.State,
				"state_message":     outEnv.StateMessage,
				"sub_account_id":    outEnv.SubAccountGuid,
				"tenant_id":         outEnv.TenantId,
				"type":              outEnv.Type,
			}
			result = append(result, m)
		}
		d.Set("environments", result)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
