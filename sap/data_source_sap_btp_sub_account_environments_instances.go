package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
)

func dataSourceSapBtpSubAccountEnvironmentsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpSubAccountEnvironmentsInstancesRead,
		Schema: map[string]*schema.Schema{
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
							Type:     schema.TypeString,
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
							Type:     schema.TypeString,
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
	btpProvisioningV1Client := meta.(*SAPClient).btpProvisioningV1Client

	if output, err := btpProvisioningV1Client.GetEnvironmentInstances(ctx); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Sub Account Environments can't be read; Operation code %v; %s",
				output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Sub Account Environments can't be read;  %v", err)
		}
	} else {
		result := make([]map[string]string, 0, len(output.Environments))
		for _, outEnv := range output.Environments {
			m := map[string]string{
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
