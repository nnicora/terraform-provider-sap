package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpmanagment"
	"github.com/pkg/errors"
	"log"
	"time"
)

func resourceSapBtpSubAccountServiceManagementInstances() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountServiceManagementInstancesCreate,
		ReadContext:   resourceSapBtpSubAccountServiceManagementInstancesRead,
		//UpdateContext: resourceSapBtpSubAccountServiceManagementInstancesUpdate,
		UpdateContext: resourceSapBtpSubAccountServiceManagementInstancesRead,
		DeleteContext: resourceSapBtpSubAccountServiceManagementInstancesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_management_url": {
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_plan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_offering_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_plan_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeList},
			},

			"ready": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"platform_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dashboard_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"context": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"maintenance_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"usable": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpSubAccountServiceManagementInstancesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.CreateServiceInstanceInput{
		Async:      false,
		Name:       d.Get("name").(string),
		Parameters: expandMapString(d.Get("parameters")),
		Labels:     expandMapListString(d.Get("labels")),
	}
	if val, ok := d.GetOk("service_plan_id"); ok {
		input.ServicePlanId = val.(string)
	}
	if val, ok := d.GetOk("service_offering_name"); ok {
		input.ServiceOfferingName = val.(string)
	}
	if val, ok := d.GetOk("service_plan_name"); ok {
		input.ServicePlanName = val.(string)
	}
	if output, err := btpServiceManagementV1Client.CreateServiceInstance(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Instances can't be created; %s", output.ErrorMessage))
		}
	} else {
		d.SetId(output.Id)
		d.Set("service_plan_id", output.ServicePlanId)
		d.Set("ready", output.Ready)
		d.Set("platform_id", output.PlatformId)
		d.Set("dashboard_url", output.DashboardUrl)
		d.Set("context", output.Context)
		d.Set("maintenance_info", output.MaintenanceInfo)
		d.Set("usable", output.Usable)
	}

	return nil
}

func resourceSapBtpSubAccountServiceManagementInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.GetServiceInstanceInput{
		ServiceInstanceID: d.Id(),
	}
	if output, err := btpServiceManagementV1Client.GetServiceInstance(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Instances can't be created; %s", output.ErrorMessage))
		}
	} else {
		d.Set("service_plan_id", output.ServicePlanId)
		d.Set("ready", output.Ready)
		d.Set("platform_id", output.PlatformId)
		d.Set("dashboard_url", output.DashboardUrl)
		d.Set("context", output.Context)
		d.Set("maintenance_info", output.MaintenanceInfo)
		d.Set("usable", output.Usable)
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementInstancesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//session := meta.(*SAPClient).session
	//
	//oauth2Map := mapFrom(d.Get("oauth2"))
	//log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)
	//
	//serviceManagementUrl := d.Get("service_management_url").(string)
	//
	//endpointConfig := &sap.EndpointConfig {
	//	Host: serviceManagementUrl,
	//	OAuth2: oauth2ConfigFrom(oauth2Map),
	//}
	//if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
	//	return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	//}
	//
	//btpServiceManagementV1Client := btpmanagment.New(session)
	//
	//input := &btpmanagment.UpdateServiceInstanceInput {
	//	ServiceInstanceID: d.Id(),
	//	Async: false,
	//	Name: d.Get("name").(string),
	//	Parameters: expandMapString(d.Get("parameters")),
	//	Labels: expandMapListString(d.Get("labels")),
	//}
	//if val, ok := d.GetOk("service_plan_id"); ok {
	//	input.ServicePlanId = val.(string)
	//}
	//if output, err := btpServiceManagementV1Client.UpdateServiceInstance(ctx, input); err != nil {
	//	if output != nil && output.ErrorMessage != "" {
	//		return diag.FromErr(
	//			errors.Errorf("BTP Sub Account ServiceManagement Instances can't be created; %s", output.ErrorMessage))
	//	}
	//} else {
	//	d.Set("service_plan_id", output.ServicePlanId)
	//	d.Set("ready", output.Ready)
	//	d.Set("platform_id", output.PlatformId)
	//	d.Set("dashboard_url", output.DashboardUrl)
	//	d.Set("context", output.Context)
	//	d.Set("maintenance_info", output.MaintenanceInfo)
	//	d.Set("usable", output.Usable)
	//}
	return nil
}

func resourceSapBtpSubAccountServiceManagementInstancesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.DeleteServiceInstanceInput{
		ServiceInstanceID: d.Id(),
		Async:             false,
	}
	if output, err := btpServiceManagementV1Client.DeleteServiceInstance(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Instances can't be created; %s", output.ErrorMessage))
		}
	}
	return nil
}
