package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"subdomain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
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

	dirCPM := make([]btpaccounts.CustomProperties, 0)
	cpValues := d.Get("custom_properties").([]map[string]string)
	for _, m := range cpValues {
		key := ""
		if val, ok := m["key"]; ok {
			key = val
		}
		value := ""
		if val, ok := m["value"]; ok {
			value = val
		}
		accountId := ""
		if val, ok := m["account_id"]; ok {
			accountId = val
		}

		if len(key) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'key' is empty:  %v", m))
		}
		if len(value) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'value' is empty:  %v", m))
		}
		if len(accountId) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'account_id' is empty:  %v", m))
		}

		cp := btpaccounts.CustomProperties{
			KeyValue: btpaccounts.KeyValue{
				Key:   key,
				Value: value,
			},
			AccountGuid: accountId,
		}
		dirCPM = append(dirCPM, cp)
	}

	dirInput := &btpaccounts.CreateDirectoryInput{
		CustomProperties:  dirCPM,
		Description:       d.Get("description").(string),
		DirectoryAdmins:   d.Get("admins").([]string),
		DirectoryFeatures: d.Get("features").([]string),
		DisplayName:       d.Get("display_name").(string),
		Subdomain:         d.Get("subdomain").(string),
	}

	if dirOutput, err := btpAccountsClient.CreateDirectory(ctx, dirInput); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be created:  %v", err))
	} else {
		d.SetId(dirOutput.Guid)
		readFromDirectoryIntoResourceData(dirOutput.Directory, d)
	}

	return nil
}

func resourceSapBtpAccountDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	dirInput := &btpaccounts.GetDirectoryInput{
		DirectoryGuid:         d.Id(),
		DerivedAuthorizations: d.Get("derived_authorizations").(string),
		Expand:                d.Get("expand").(bool),
	}
	if dirOutput, err := btpAccountsClient.GetDirectory(ctx, dirInput); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be read:  %v", err))
	} else {
		d.SetId(d.Id())
		readFromDirectoryIntoResourceData(dirOutput.Directory, d)
	}

	return nil
}

func resourceSapBtpAccountDirectoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	dirCPM := make([]btpaccounts.CustomProperties, 0)
	cpValues := d.Get("custom_properties").([]map[string]string)
	for _, m := range cpValues {
		key := ""
		if val, ok := m["key"]; ok {
			key = val
		}
		value := ""
		if val, ok := m["value"]; ok {
			value = val
		}
		accountId := ""
		if val, ok := m["account_id"]; ok {
			accountId = val
		}

		if len(key) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'key' is empty:  %v", m))
		}
		if len(value) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'value' is empty:  %v", m))
		}
		if len(accountId) <= 0 {
			return diag.FromErr(errors.Errorf("BTP Account Directory; Custom Properties 'account_id' is empty:  %v", m))
		}

		cp := btpaccounts.CustomProperties{
			KeyValue: btpaccounts.KeyValue{
				Key:   key,
				Value: value,
			},
			AccountGuid: accountId,
		}
		dirCPM = append(dirCPM, cp)
	}

	dirInput := &btpaccounts.UpdateDirectoryInput{
		DirectoryGuid:    d.Id(),
		CustomProperties: dirCPM,
		DisplayName:      d.Get("display_name").(string),
		Description:      d.Get("description").(string),
	}

	if dirOutput, err := btpAccountsClient.UpdateDirectory(ctx, dirInput); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be created:  %v", err))
	} else {
		d.SetId(dirOutput.Guid)
		readFromDirectoryIntoResourceData(dirOutput.Directory, d)
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

	dirInput := &btpaccounts.DeleteDirectoryInput{
		DirectoryGuid: d.Id(),
		ForceDelete:   d.Get("force_Delete").(bool),
	}
	if _, err := btpAccountsClient.DeleteDirectory(ctx, dirInput); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account Directory can't be read:  %v", err))
	}
	return nil
}
