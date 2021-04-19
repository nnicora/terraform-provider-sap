package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"time"
)

func dataSourceSapBtpSubAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpSubAccountRead,
		Schema: map[string]*schema.Schema{
			"sub_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"derived_authorizations": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			// Computed
			"global_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"beta_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},
			"region": {
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
			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"used_for_production": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
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
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"account_id": {
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

func dataSourceSapBtpSubAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if subAccountId, ok := d.GetOk("sub_account_id"); ok {

		saInput := &btpaccounts.GetSubAccountInput{
			SubAccountGuid:        subAccountId.(string),
			DerivedAuthorizations: d.Get("derived_authorizations").(string),
		}

		if saOutput, err := btpAccountsClient.GetSubAccount(ctx, saInput); err != nil {
			return diag.FromErr(errors.Errorf("BTP Sub Account Custom Properties can't be read:  %v", err))
		} else {
			d.SetId(saOutput.Guid)
			d.Set("global_account_id", saOutput.GlobalAccountGuid)
			d.Set("beta_enabled", saOutput.BetaEnabled)
			d.Set("created_by", saOutput.CreatedBy)
			d.Set("created_date", saOutput.CreatedDate.Format(time.RFC3339))
			d.Set("description", saOutput.Description)
			d.Set("display_name", saOutput.DisplayName)
			d.Set("modified_date", saOutput.ModifiedDate.Format(time.RFC3339))
			d.Set("parent_features", saOutput.ParentFeatures)
			d.Set("parent_id", saOutput.ParentGuid)
			d.Set("region", saOutput.Region)
			d.Set("state", saOutput.State)
			d.Set("state_message", saOutput.StateMessage)
			d.Set("subdomain", saOutput.Subdomain)
			d.Set("used_for_production", saOutput.UsedForProduction)
			d.Set("zone_id", saOutput.ZoneId)

			cpm := make([]map[string]interface{}, 0)
			for _, saCP := range saOutput.CustomProperties {
				dm := make(map[string]interface{})
				dm["key"] = saCP.Key
				dm["value"] = saCP.Value
				dm["account_id"] = saCP.AccountGuid

				cpm = append(cpm, dm)
			}
			d.Set("custom_properties", cpm)

		}
	} else {
		return diag.FromErr(errors.New("sub_account_id must be set when want to read an sub-account custom properties"))
	}
	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
