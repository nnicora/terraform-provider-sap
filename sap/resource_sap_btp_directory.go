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

func resourceSapBtpDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpDirectoryCreate,
		ReadContext:   resourceSapBtpDirectoryRead,
		UpdateContext: resourceSapBtpDirectoryUpdate,
		DeleteContext: resourceSapBtpDirectoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"derived_authorizations": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"admins": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"features": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
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

			"contract_status": {
				Type:     schema.TypeString,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpDirectoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	cp, err := toBtpAccountsCustomProperties(d.Get("custom_properties"))
	if err != nil {
		diag.FromErr(err)
	}
	input := &btpaccounts.CreateDirectoryInput{
		CustomProperties: cp,
	}
	if val, ok := d.GetOk("subdomain"); ok {
		input.Subdomain = val.(string)
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}
	if val, ok := d.GetOk("admins"); ok {
		input.DirectoryAdmins = expandStringSet(val.(*schema.Set))
	}
	if val, ok := d.GetOk("features"); ok {
		input.DirectoryFeatures = expandStringSet(val.(*schema.Set))
	}
	if val, ok := d.GetOk("display_name"); ok {
		input.DisplayName = val.(string)
	}

	logDebug(input, "CreateDirectory Input")
	if output, err := btpAccountsClient.CreateDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory can't be created; Status Code: %d; %s",
					sap.Int32Value(output.Error.Code), sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Directory can't be created;  %v", err))
	} else {
		d.SetId(output.Guid)
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func resourceSapBtpDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	input := &btpaccounts.GetDirectoryInput{
		DirectoryGuid: d.Id(),
	}
	if val, ok := d.GetOk("derived_authorizations"); ok {
		input.DerivedAuthorizations = val.(string)
	}
	if val, ok := d.GetOk("expand"); ok {
		input.Expand = val.(bool)
	}

	logDebug(input, "GetDirectory Input")
	if output, err := btpAccountsClient.GetDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory can't be read; Status Code: %d; %s",
					sap.Int32Value(output.Error.Code), sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Directory can't be read;  %v", err))
	} else {
		d.SetId(d.Id())
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func resourceSapBtpDirectoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	cp, err := toBtpAccountsCustomProperties(d.Get("custom_properties"))
	if err != nil {
		diag.FromErr(err)
	}
	input := &btpaccounts.UpdateDirectoryInput{
		DirectoryGuid:    d.Id(),
		CustomProperties: cp,
	}
	if val, ok := d.GetOk("display_name"); ok {
		input.DisplayName = val.(string)
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}

	logDebug(input, "UpdateDirectory Input")
	if output, err := btpAccountsClient.UpdateDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory can't be updated; Status Code: %d; %s",
					sap.Int32Value(output.Error.Code), sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Directory can't be updated;  %v", err))
	} else {
		d.SetId(output.Guid)
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func readFromDirectoryIntoResourceData(dir btpaccounts.Directory, d *schema.ResourceData) {
	d.Set("contract_status", dir.ContractStatus)
	d.Set("created_date", dir.CreatedDate.Format(time.RFC3339))
	d.Set("display_name", dir.DisplayName)
	d.Set("modified_date", dir.ModifiedDate.Format(time.RFC3339))
	d.Set("features", dir.DirectoryFeatures)
	d.Set("parent_id", dir.ParentGuid)
	d.Set("state_message", dir.StateMessage)
	d.Set("children", dir.Children)
	d.Set("entity_state", dir.EntityState)

	if len(dir.Subdomain) > 0 {
		d.Set("subdomain", dir.Subdomain)
	}
	if len(dir.CreatedBy) > 0 {
		d.Set("created_by", dir.CreatedBy)
	}
	if len(dir.Description) > 0 {
		d.Set("description", dir.Description)
	}

	cpm := make([]map[string]interface{}, len(dir.CustomProperties))
	for idx, saCP := range dir.CustomProperties {
		dm := make(map[string]interface{})
		dm["key"] = saCP.Key
		dm["value"] = saCP.Value
		dm["account_id"] = saCP.AccountGuid

		cpm[idx] = dm
	}
	d.Set("custom_properties", cpm)

	subAccs := make([]map[string]string, len(dir.SubAccounts))
	for idx, subAcc := range dir.SubAccounts {
		dm := make(map[string]string)
		dm["account_id"] = subAcc.Guid
		dm["global_account_id"] = subAcc.GlobalAccountGuid

		subAccs[idx] = dm
	}
	d.Set("sub_accounts", cpm)
}

func resourceSapBtpDirectoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	input := &btpaccounts.DeleteDirectoryInput{
		DirectoryGuid: d.Id(),
	}
	if val, ok := d.GetOk("force_delete"); ok {
		input.ForceDelete = val.(bool)
	}
	logDebug(input, "DeleteDirectory Input")
	if output, err := btpAccountsClient.DeleteDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory can't be deleted; Status Code: %d; %s",
					sap.Int32Value(output.Error.Code), sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Directory can't be deleted;  %v", err))
	}
	return nil
}

func toBtpAccountsCustomProperties(input interface{}) ([]btpaccounts.CustomProperties, error) {
	if input != nil {
		if cpValues, ok := input.([]interface{}); ok {
			result := make([]btpaccounts.CustomProperties, 0, len(cpValues))
			for idx := range cpValues {
				mps, ok := cpValues[idx].(map[string]interface{})
				if !ok {
					continue
				}

				key := ""
				if val, ok := mps["key"]; ok {
					key = val.(string)
				}
				value := ""
				if val, ok := mps["value"]; ok {
					value = val.(string)
				}

				if len(key) <= 0 {
					return nil, errors.Errorf("BTP Directory; Custom Properties 'key' is empty:  %v", mps)
				}
				if len(value) <= 0 {
					return nil, errors.Errorf("BTP Directory; Custom Properties 'value' is empty:  %v", mps)
				}

				cp := btpaccounts.CustomProperties{
					KeyValue: btpaccounts.KeyValue{
						Key:   key,
						Value: value,
					},
				}
				result = append(result, cp)
			}
			return result, nil
		}
	}
	return nil, nil
}
