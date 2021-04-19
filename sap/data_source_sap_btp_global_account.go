package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
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

	derivedAuthorizations := d.Get("derived_authorizations").(string)
	expand := d.Get("expand").(bool)

	params := &btpaccounts.GetGlobalAccountInput{
		DerivedAuthorizations: derivedAuthorizations,
		Expand:                expand,
	}
	ga, err := btpAccountsClient.GetGlobalAccount(ctx, params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting API BAT Global Account Resources: %w", err))
	}

	d.SetId(ga.Guid)
	d.Set("display_name", ga.DisplayName)
	d.Set("description", ga.Description)
	d.Set("const_center", ga.CostCenter)

	d.Set("commercial_model", ga.CommercialModel)
	d.Set("consumption_based", ga.ConsumptionBased)
	d.Set("contract_status", ga.ContractStatus)
	d.Set("created_date", ga.CreatedDate.Format(time.RFC3339))
	d.Set("modified_date", ga.ModifiedDate.Format(time.RFC3339))
	d.Set("crm_customer_id", ga.CrmCustomerId)
	d.Set("crm_tenant_id", ga.CrmTenantId)
	d.Set("entity_state", ga.EntityState)
	d.Set("expiry_date", ga.ExpiryDate.Format(time.RFC3339))
	d.Set("geo_access", ga.GeoAccess)
	//d.Set("legal_links", ga.LegalLinks.Privacy)
	d.Set("license_type", ga.LicenseType)
	d.Set("origin", ga.Origin)
	d.Set("parent_id", ga.ParentGuid)
	d.Set("parent_type", ga.ParentType)
	d.Set("renewal_date", ga.RenewalDate.Format(time.RFC3339))
	d.Set("service_id", ga.ServiceId)
	d.Set("state_message", ga.StateMessage)
	d.Set("subdomain", ga.Subdomain)
	d.Set("termination_notification_status", ga.TerminationNotificationStatus)
	d.Set("use_for", ga.UseFor)

	cp := make([]map[string]interface{}, 0)
	{
		for _, gaCP := range ga.CustomProperties {
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
		for _, sa := range ga.Subaccounts {
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
