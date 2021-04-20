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

func dataSourceSapBtpAccountDirectory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpAccountDirectoryRead,
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"derived_authorizations": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			// Computed
			"contract_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeBool,
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
			"features": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"modified_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state_message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sub_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"global_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
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

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpAccountDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	input := &btpaccounts.GetDirectoryInput{
		DirectoryGuid: d.Get("directory_id").(string),
	}
	if val, ok := d.GetOk("derived_authorizations"); ok {
		input.DerivedAuthorizations = val.(string)
	}

	if output, err := btpAccountsClient.GetDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory can't be read; %s", sap.StringValue(output.Error.Message)))
		}

		return diag.FromErr(errors.Errorf("BTP Directory can't be read;  %v", err))
	} else {
		d.SetId(output.Guid)
		d.Set("contract_status", output.ContractStatus)
		d.Set("created_by", output.CreatedBy)
		d.Set("created_date", output.CreatedDate.Format(time.RFC3339))
		d.Set("description", output.Description)
		d.Set("display_name", output.DisplayName)
		d.Set("modified_date", output.ModifiedDate.Format(time.RFC3339))
		d.Set("features", output.DirectoryFeatures)
		d.Set("parent_id", output.ParentGuid)
		d.Set("state_message", output.StateMessage)
		d.Set("subdomain", output.Subdomain)
		d.Set("children", output.Children)
		d.Set("entity_state", output.EntityState)

		cpm := make([]map[string]interface{}, len(output.CustomProperties))
		for idx, saCP := range output.CustomProperties {
			dm := make(map[string]interface{})
			dm["key"] = saCP.Key
			dm["value"] = saCP.Value
			dm["account_id"] = saCP.AccountGuid

			cpm[idx] = dm
		}
		d.Set("custom_properties", cpm)

		subAccs := make([]map[string]string, len(output.SubAccounts))
		for idx, subAcc := range output.SubAccounts {
			dm := make(map[string]string)
			dm["account_id"] = subAcc.Guid
			dm["global_account_id"] = subAcc.GlobalAccountGuid

			subAccs[idx] = dm
		}
		d.Set("sub_accounts", cpm)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
