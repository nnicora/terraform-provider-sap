package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpprovisioning"
	"github.com/pkg/errors"
)

func dataSourceSapBtpProvisioningAvailableEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpProvisioningAvailableEnvironmentsRead,
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
						"availability_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"create_schema": {
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
						"landscape_label": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"plan_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"plan_updatable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"service_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_display_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_documentation_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_image_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_long_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_support_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"technical_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"update_schema": {
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

func dataSourceSapBtpProvisioningAvailableEnvironmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if output, err := btpProvisioningV1Client.GetAvailableEnvironments(ctx); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Provisioning Environments can't be read; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Provisioning Environments can't be read;  %v", err)
		}
	} else {
		result := make([]map[string]interface{}, 0, len(output.Environments))
		for _, outEnv := range output.Environments {
			m := map[string]interface{}{
				"availability_level":        outEnv.AvailabilityLevel,
				"create_schema":             outEnv.CreateSchema,
				"description":               outEnv.Description,
				"environment_type":          outEnv.EnvironmentType,
				"landscape_label":           outEnv.LandscapeLabel,
				"plan_name":                 outEnv.PlanName,
				"plan_updatable":            outEnv.PlanUpdatable,
				"service_description":       outEnv.ServiceDescription,
				"service_display_name":      outEnv.ServiceDisplayName,
				"service_documentation_url": outEnv.ServiceDocumentationUrl,
				"service_image_url":         outEnv.ServiceImageUrl,
				"service_long_description":  outEnv.ServiceLongDescription,
				"service_name":              outEnv.ServiceName,
				"service_support_url":       outEnv.ServiceSupportUrl,
				"technical_key":             outEnv.TechnicalKey,
				"update_schema":             outEnv.UpdateSchema,
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
