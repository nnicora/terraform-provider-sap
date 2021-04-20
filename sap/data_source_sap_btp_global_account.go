package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"time"
)

func dataSourceSapBtpGlobalAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpGlobalAccountRead,
		Schema: map[string]*schema.Schema{
			"derived_authorizations": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"expand": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"const_center": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"commercial_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"consumption_based": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"contract_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crm_customer_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crm_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"geo_access": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"renewal_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state_message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"termination_notification_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"use_for": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"sub_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"global_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"beta_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"modified_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parent_features": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
							Set:      schema.HashString,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"state_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subdomain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"used_for_production": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						//"legal_links": {
						//	Type:     schema.TypeString,
						//	Optional: true,
						//},
					},
				},
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpGlobalAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	input := &btpaccounts.GetGlobalAccountInput{}
	if val, ok := d.GetOk("derived_authorizations"); ok {
		input.DerivedAuthorizations = val.(string)
	}
	if val, ok := d.GetOk("expand"); ok {
		input.Expand = val.(bool)
	}
	output, err := btpAccountsClient.GetGlobalAccount(ctx, input)
	if err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Global Account can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(fmt.Errorf("BTP Global Account can't be read; %w", err))
	}

	d.SetId(output.Guid)
	d.Set("display_name", output.DisplayName)
	d.Set("description", output.Description)
	d.Set("const_center", output.CostCenter)

	d.Set("commercial_model", output.CommercialModel)
	d.Set("consumption_based", output.ConsumptionBased)
	d.Set("contract_status", output.ContractStatus)
	d.Set("created_date", output.CreatedDate.Format(time.RFC3339))
	d.Set("modified_date", output.ModifiedDate.Format(time.RFC3339))
	d.Set("crm_customer_id", output.CrmCustomerId)
	d.Set("crm_tenant_id", output.CrmTenantId)
	d.Set("entity_state", output.EntityState)
	d.Set("expiry_date", output.ExpiryDate.Format(time.RFC3339))
	d.Set("geo_access", output.GeoAccess)
	//d.Set("legal_links", output.LegalLinks.Privacy)
	d.Set("license_type", output.LicenseType)
	d.Set("origin", output.Origin)
	d.Set("parent_id", output.ParentGuid)
	d.Set("parent_type", output.ParentType)
	d.Set("renewal_date", output.RenewalDate.Format(time.RFC3339))
	d.Set("service_id", output.ServiceId)
	d.Set("state_message", output.StateMessage)
	d.Set("subdomain", output.Subdomain)
	d.Set("termination_notification_status", output.TerminationNotificationStatus)
	d.Set("use_for", output.UseFor)

	cp := make([]map[string]interface{}, 0)
	{
		for _, gaCP := range output.CustomProperties {
			m := make(map[string]interface{})
			m["key"] = gaCP.Key
			m["value"] = gaCP.Value
			m["account_id"] = gaCP.AccountGuid

			cp = append(cp, m)
		}
	}
	d.Set("custom_properties", cp)

	subAccounts := make([]map[string]interface{}, 0)
	{
		for _, sa := range output.Subaccounts {
			m := make(map[string]interface{})
			m["id"] = sa.Guid
			m["global_account_id"] = sa.GlobalAccountGuid
			m["beta_enabled"] = sa.BetaEnabled
			m["created_by"] = sa.CreatedBy
			m["created_date"] = sa.CreatedDate.Format(time.RFC3339)
			m["description"] = sa.Description
			m["display_name"] = sa.DisplayName
			m["modified_date"] = sa.ModifiedDate.Format(time.RFC3339)
			m["parent_features"] = sa.ParentFeatures
			m["parent_id"] = sa.ParentGuid
			m["region"] = sa.Region
			m["state"] = sa.State
			m["state_message"] = sa.StateMessage
			m["subdomain"] = sa.Subdomain
			m["used_for_production"] = sa.UsedForProduction
			m["zone_id"] = sa.ZoneId

			cpm := make([]map[string]interface{}, 0)
			for _, saCP := range sa.CustomProperties {
				dm := make(map[string]interface{})
				dm["key"] = saCP.Key
				dm["value"] = saCP.Value
				dm["account_id"] = saCP.AccountGuid

				cpm = append(cpm, dm)
			}
			m["custom_properties"] = cpm

			subAccounts = append(subAccounts, m)
		}
	}
	d.Set("sub_accounts", subAccounts)

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
