package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
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

	input := &btpaccounts.GetSubAccountInput{
		SubAccountGuid: d.Get("sub_account_id").(string),
	}
	if val, ok := d.GetOk("derived_authorizations"); ok {
		input.DerivedAuthorizations = val.(string)
	}

	if output, err := btpAccountsClient.GetSubAccount(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account can't be read;  %v", err))
	} else {
		d.SetId(output.Guid)
		d.Set("global_account_id", output.GlobalAccountGuid)
		d.Set("beta_enabled", output.BetaEnabled)
		d.Set("created_by", output.CreatedBy)
		d.Set("created_date", output.CreatedDate.Format(time.RFC3339))
		d.Set("description", output.Description)
		d.Set("display_name", output.DisplayName)
		d.Set("modified_date", output.ModifiedDate.Format(time.RFC3339))
		d.Set("parent_features", output.ParentFeatures)
		d.Set("parent_id", output.ParentGuid)
		d.Set("region", output.Region)
		d.Set("state", output.State)
		d.Set("state_message", output.StateMessage)
		d.Set("subdomain", output.Subdomain)
		d.Set("used_for_production", output.UsedForProduction)
		d.Set("zone_id", output.ZoneId)

		cpm := make([]map[string]interface{}, 0)
		for _, saCP := range output.CustomProperties {
			dm := make(map[string]interface{})
			dm["key"] = saCP.Key
			dm["value"] = saCP.Value
			dm["account_id"] = saCP.AccountGuid

			cpm = append(cpm, dm)
		}
		d.Set("custom_properties", cpm)

	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
