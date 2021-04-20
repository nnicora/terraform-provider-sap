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

func resourceSapBtpAccountDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpAccountDirectoryCreate,
		ReadContext:   resourceSapBtpAccountDirectoryRead,
		UpdateContext: resourceSapBtpAccountDirectoryUpdate,
		DeleteContext: resourceSapBtpAccountDirectoryDelete,
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
			},
			"derived_authorizations": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Type:     schema.TypeBool,
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

func resourceSapBtpAccountDirectoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	var dirCPM []btpaccounts.CustomProperties

	resCP := d.Get("custom_properties")
	if resCP != nil {
		dirCPM = make([]btpaccounts.CustomProperties, 0)
		cpValues := resCP.([]map[string]string)
		for _, m := range cpValues {
			key := ""
			if val, ok := m["key"]; ok {
				key = val
			}
			value := ""
			if val, ok := m["value"]; ok {
				value = val
			}

			if len(key) <= 0 {
				return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'key' is empty:  %v", m))
			}
			if len(value) <= 0 {
				return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'value' is empty:  %v", m))
			}

			cp := btpaccounts.CustomProperties{
				KeyValue: btpaccounts.KeyValue{
					Key:   key,
					Value: value,
				},
			}
			dirCPM = append(dirCPM, cp)
		}
	}

	input := &btpaccounts.CreateDirectoryInput{
		CustomProperties: dirCPM,
	}
	if val, ok := d.GetOk("subdomain"); ok {
		input.Subdomain = val.(string)
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}
	if val, ok := d.GetOk("admins"); ok {
		input.DirectoryAdmins = val.([]string)
	}
	if val, ok := d.GetOk("features"); ok {
		input.DirectoryFeatures = val.([]string)
	}
	if val, ok := d.GetOk("display_name"); ok {
		input.DisplayName = val.(string)
	}

	if output, err := btpAccountsClient.CreateDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account Directory can't be created; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be created;  %v", err))
	} else {
		d.SetId(output.Guid)
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func resourceSapBtpAccountDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if output, err := btpAccountsClient.GetDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account Directory can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be read;  %v", err))
	} else {
		d.SetId(d.Id())
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func resourceSapBtpAccountDirectoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	var dirCPM []btpaccounts.CustomProperties

	resCP := d.Get("custom_properties")
	if resCP != nil {
		dirCPM = make([]btpaccounts.CustomProperties, 0)
		cpValues := resCP.([]map[string]string)
		for _, m := range cpValues {
			key := ""
			if val, ok := m["key"]; ok {
				key = val
			}
			value := ""
			if val, ok := m["value"]; ok {
				value = val
			}

			if len(key) <= 0 {
				return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'key' is empty:  %v", m))
			}
			if len(value) <= 0 {
				return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'value' is empty:  %v", m))
			}

			cp := btpaccounts.CustomProperties{
				KeyValue: btpaccounts.KeyValue{
					Key:   key,
					Value: value,
				},
			}
			dirCPM = append(dirCPM, cp)
		}
	}

	input := &btpaccounts.UpdateDirectoryInput{
		DirectoryGuid:    d.Id(),
		CustomProperties: dirCPM,
	}
	if val, ok := d.GetOk("display_name"); ok {
		input.DisplayName = val.(string)
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}

	if output, err := btpAccountsClient.UpdateDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account Directory can't be updated; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be updated;  %v", err))
	} else {
		d.SetId(output.Guid)
		readFromDirectoryIntoResourceData(output.Directory, d)
	}

	return nil
}

func readFromDirectoryIntoResourceData(dir btpaccounts.Directory, d *schema.ResourceData) {
	d.Set("contract_status", dir.ContractStatus)
	d.Set("created_by", dir.CreatedBy)
	d.Set("created_date", dir.CreatedDate.Format(time.RFC3339))
	d.Set("description", dir.Description)
	d.Set("display_name", dir.DisplayName)
	d.Set("modified_date", dir.ModifiedDate.Format(time.RFC3339))
	d.Set("features", dir.DirectoryFeatures)
	d.Set("parent_id", dir.ParentGuid)
	d.Set("state_message", dir.StateMessage)
	d.Set("subdomain", dir.Subdomain)
	d.Set("children", dir.Children)
	d.Set("entity_state", dir.EntityState)

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

func resourceSapBtpAccountDirectoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	input := &btpaccounts.DeleteDirectoryInput{
		DirectoryGuid: d.Id(),
	}
	if val, ok := d.GetOk("force_delete"); ok {
		input.ForceDelete = val.(bool)
	}
	if output, err := btpAccountsClient.DeleteDirectory(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account Directory can't be deleted; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be deleted;  %v", err))
	}
	return nil
}
