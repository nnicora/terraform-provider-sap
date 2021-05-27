package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpsaasmanager"
	"github.com/pkg/errors"
)

func dataSourceSapBtpApplicationRegistration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpApplicationRegistrationRead,
		Schema: map[string]*schema.Schema{
			"saas_manager_service": {
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

			"service_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"xsappname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"commercial_app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_urls": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"formation_solution_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpApplicationRegistrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpSaasManagerV1Client := meta.(*SAPClient).btpSaasManagerV1Client
	session := meta.(*SAPClient).session
	serviceList := d.Get("saas_manager_service").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("SaaS manager service is required")
	}

	err := session.AddEndpointWithReplace(btpsaasmanager.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP SaaS manager service OAuth2;  %v", err))
	}
	btpSaasManagerV1Client := btpsaasmanager.New(session)

	input := &btpsaasmanager.GetApplicationRegistrationInput{}
	if output, err := btpSaasManagerV1Client.GetApplicationRegistration(ctx, input); err != nil {
		if output != nil && output.Error != "" {
			return diag.Errorf("BTP SaaS Application Registration can't be read; Operation code %v; %s",
				output.StatusCode, output.Error)
		} else {
			return diag.Errorf("BTP SaaS Subscription to an application can't be done;  %v", err)
		}
	} else {
		d.Set("service_instance_id", output.ServiceInstanceId)
		d.Set("organization_id", output.OrganizationGuid)
		d.Set("space_id", output.SpaceGuid)
		d.Set("xsappname", output.XSAppName)
		d.Set("app_id", output.AppId)
		d.Set("app_name", output.AppName)
		d.Set("commercial_app_name", output.CommercialAppName)
		d.Set("app_urls", output.AppUrls)
		d.Set("provider_tenant_id", output.ProviderTenantId)
		d.Set("app_type", output.AppType)
		d.Set("display_name", output.DisplayName)
		d.Set("description", output.Description)
		d.Set("category", output.Category)
		d.Set("global_account_id", output.GlobalAccountId)
		d.Set("formation_solution_name", output.FormationSolutionName)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
