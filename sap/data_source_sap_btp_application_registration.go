package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpsaasprovisioning"
)

func dataSourceSapBtpApplicationRegistration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpApplicationRegistrationRead,
		Schema: map[string]*schema.Schema{
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
	btpSaaSProvisioningV1Client := meta.(*SAPClient).btpSaaSProvisioningV1Client

	input := &btpsaasprovisioning.GetApplicationRegistrationInput{}
	if output, err := btpSaaSProvisioningV1Client.GetApplicationRegistration(ctx, input); err != nil {
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
