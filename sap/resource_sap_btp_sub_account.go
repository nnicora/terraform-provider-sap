package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"time"
)

func resourceSapBtpSubAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountCreate,
		ReadContext:   resourceSapBtpSubAccountRead,
		UpdateContext: resourceSapBtpSubAccountUpdate,
		DeleteContext: resourceSapBtpSubAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"used_for_production": {
				Type:     schema.TypeString,
				Required: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Required: true,
			},

			"beta_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"sub_account_admins": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"zone_id": {
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
			"parent_features": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"parent_id": {
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

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpSubAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	customProperties := make([]btpaccounts.KeyValue, 0)
	if v, ok := d.GetOk("custom_properties"); ok {
		customProperties = expandSapBtpAccountCustomPropertiesParameters(v.([]interface{}))
	}

	subAccountAdmins := make([]string, 0)
	if v, ok := d.GetOk("sub_account_admins"); ok {
		subAccountAdmins = expandStringList(v.([]interface{}))
	}

	//"subdomain": {
	//"used_for_production": {
	//"origin": {
	gaId := d.Get("global_account_id")
	region := d.Get("region")

	usedForProduction := d.Get("used_for_production")
	displayName := d.Get("display_name")
	subdomain := d.Get("subdomain")
	origin := d.Get("origin")
	input := &btpaccounts.CreateSubAccountInput{
		ParentGuid:  gaId.(string),
		Region:      region.(string),
		DisplayName: displayName.(string),

		UsedForProduction: usedForProduction.(string),
		Subdomain:         subdomain.(string),
		CustomProperties:  customProperties,
		SubaccountAdmins:  subAccountAdmins,
		Origin:            origin.(string),
	}
	if val, ok := d.GetOk("beta_enabled"); ok {
		input.BetaEnabled = val.(bool)
	}
	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}

	if output, err := btpAccountsClient.CreateSubAccount(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account can't be created; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account can't be created;  %v", err))
	} else {
		aId := output.Guid
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			respSubAccount, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
				SubAccountGuid: aId,
			})
			if gAcErr != nil {
				if respSubAccount != nil && respSubAccount.Error != nil {
					return resource.RetryableError(
						errors.Errorf("BTP Sub Account can't be read; %s", sap.StringValue(output.Error.Message)))
				}
				return resource.RetryableError(gAcErr)
			} else {
				if respSubAccount.State == "CREATING" {
					return resource.RetryableError(fmt.Errorf("BTP Sub Account creation in progress"))
				}
				if respSubAccount.State == "OK" {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("BTP Sub Account not yet started"))
				}
			}

			return nil
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}

		respSubAccount, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
			SubAccountGuid: aId,
		})
		if gAcErr != nil {
			if respSubAccount != nil && respSubAccount.Error != nil {
				return diag.FromErr(
					errors.Errorf("BTP Sub Account can't be read; %s", sap.StringValue(output.Error.Message)))
			}
			return diag.FromErr(errors.New("BTP Sub Account not found"))
		}

		d.SetId(aId)
		d.Set("beta_enabled", respSubAccount.BetaEnabled)
		d.Set("created_by", respSubAccount.CreatedBy)
		d.Set("created_date", respSubAccount.CreatedDate.Format(time.RFC3339))
		d.Set("modified_date", respSubAccount.ModifiedDate.Format(time.RFC3339))
		d.Set("parent_features", respSubAccount.ParentFeatures)
		d.Set("parent_id", respSubAccount.ParentGuid)
		d.Set("state", respSubAccount.State)
		d.Set("state_message", respSubAccount.StateMessage)
		d.Set("description", respSubAccount.Description)
		d.Set("display_name", respSubAccount.DisplayName)
		d.Set("global_account_id", respSubAccount.GlobalAccountGuid)
		d.Set("region", respSubAccount.Region)
		d.Set("zone_id", respSubAccount.ZoneId)
		d.Set("used_for_production", respSubAccount.UsedForProduction)
		d.Set("subdomain", respSubAccount.Subdomain)

		cp := make([]map[string]interface{}, 0)
		{
			for _, gaCP := range respSubAccount.CustomProperties {
				m := make(map[string]interface{})
				m["key"] = gaCP.Key
				m["value"] = gaCP.Value
				cp = append(cp, m)
			}
		}
		d.Set("custom_properties", cp)
	}
	return nil
}

func expandSapBtpAccountCustomPropertiesParameters(config []interface{}) []btpaccounts.KeyValue {
	parameters := make([]btpaccounts.KeyValue, 0)

	for _, c := range config {
		param := c.(map[string]interface{})
		key := param["key"].(string)
		value := param["value"].(string)

		parameters = append(parameters, btpaccounts.KeyValue{
			Key:   key,
			Value: value,
		})
	}

	return parameters
}

func expandSapBtpAccountUpdateCustomPropertiesParameters(config []interface{}) []btpaccounts.UpdateSubAccountProperties {
	parameters := make([]btpaccounts.UpdateSubAccountProperties, 0)

	for _, c := range config {
		param := c.(map[string]interface{})
		key := param["key"].(string)
		value := param["value"].(string)
		deleteFlag := param["delete"].(bool)

		prop := btpaccounts.UpdateSubAccountProperties{
			KeyValue: btpaccounts.KeyValue{
				Key:   key,
				Value: value,
			},
			Delete: deleteFlag,
		}
		parameters = append(parameters, prop)
	}

	return parameters
}

func resourceSapBtpSubAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if output, err := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
		SubAccountGuid: d.Id(),
	}); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(fmt.Errorf("BTP Sub Account can't be read:; %#v", err))
	} else {
		d.SetId(output.Guid)
		d.Set("beta_enabled", output.BetaEnabled)
		d.Set("created_by", output.CreatedBy)
		d.Set("created_date", output.CreatedDate.Format(time.RFC3339))
		d.Set("modified_date", output.ModifiedDate.Format(time.RFC3339))
		d.Set("parent_features", output.ParentFeatures)
		d.Set("parent_id", output.ParentGuid)
		d.Set("state", output.State)
		d.Set("state_message", output.StateMessage)
		d.Set("description", output.Description)
		d.Set("display_name", output.DisplayName)
		d.Set("global_account_id", output.GlobalAccountGuid)
		d.Set("region", output.Region)
		d.Set("zone_id", output.ZoneId)
		d.Set("used_for_production", output.UsedForProduction)
		d.Set("subdomain", output.Subdomain)

		cp := make([]map[string]interface{}, 0)
		{
			for _, gaCP := range output.CustomProperties {
				m := make(map[string]interface{})
				m["key"] = gaCP.Key
				m["value"] = gaCP.Value
				cp = append(cp, m)
			}
		}
		d.Set("custom_properties", cp)
	}
	return nil
}

func resourceSapBtpSubAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	customProperties := make([]btpaccounts.UpdateSubAccountProperties, 0)
	if v, ok := d.GetOk("custom_properties"); ok {
		customProperties = expandSapBtpAccountUpdateCustomPropertiesParameters(v.([]interface{}))
	}

	usedForProduction := d.Get("used_for_production")
	displayName := d.Get("display_name")
	description := d.Get("description")
	input := &btpaccounts.UpdateSubAccountInput{
		SubAccountGuid:    d.Id(),
		CustomProperties:  customProperties,
		Description:       description.(string),
		DisplayName:       displayName.(string),
		UsedForProduction: usedForProduction.(string),
	}
	if val, ok := d.GetOk("beta_enabled"); ok {
		input.BetaEnabled = val.(bool)
	}
	if output, err := btpAccountsClient.UpdateSubAccount(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account can't be updated; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(fmt.Errorf("BTP Sub Account can't be updated; %#v", err))
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			if c, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
				SubAccountGuid: output.Guid,
			}); gAcErr != nil {
				return resource.NonRetryableError(gAcErr)
			} else {
				if c.State == "UPDATING" {
					return resource.RetryableError(fmt.Errorf("BTP Sub Account updating in progress"))
				}
				if c.State == "OK" {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("state not identified"))
				}
			}
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}

		input := &btpaccounts.GetSubAccountInput{
			SubAccountGuid: output.Guid,
		}
		if output, err := btpAccountsClient.GetSubAccount(ctx, input); err != nil {
			if output != nil && output.Error != nil {
				return diag.FromErr(
					errors.Errorf("BTP Sub Account can't be updated; %s", sap.StringValue(output.Error.Message)))
			}
			return diag.FromErr(retryErr)
		} else {

			d.SetId(output.Guid)
			d.Set("beta_enabled", output.BetaEnabled)
			d.Set("created_by", output.CreatedBy)
			d.Set("created_date", output.CreatedDate.Format(time.RFC3339))
			d.Set("modified_date", output.ModifiedDate.Format(time.RFC3339))
			d.Set("parent_features", output.ParentFeatures)
			d.Set("parent_id", output.ParentGuid)
			d.Set("state", output.State)
			d.Set("state_message", output.StateMessage)
			d.Set("description", output.Description)
			d.Set("display_name", output.DisplayName)
			d.Set("global_account_id", output.GlobalAccountGuid)
			d.Set("region", output.Region)
			d.Set("zone_id", output.ZoneId)
			d.Set("used_for_production", output.UsedForProduction)
			d.Set("subdomain", output.Subdomain)

			cp := make([]map[string]interface{}, 0)
			{
				for _, gaCP := range output.CustomProperties {
					m := make(map[string]interface{})
					m["key"] = gaCP.Key
					m["value"] = gaCP.Value
					cp = append(cp, m)
				}
			}
			d.Set("custom_properties", cp)
		}
	}
	return nil
}

func resourceSapBtpSubAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client
	aId := d.Id()

	output, err := btpAccountsClient.DeleteSubAccount(ctx, &btpaccounts.DeleteSubAccountInput{
		SubAccountGuid: aId,
	})
	if err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account can't be deleted; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(fmt.Errorf("BTP Sub Account can't be deleted; %#v", err))
	}

	retryErr := resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		if acc, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
			SubAccountGuid: aId,
		}); gAcErr == nil {
			if acc.State == "DELETING" {
				return resource.RetryableError(fmt.Errorf("BTP Sub Account still exist, having deletion in progress"))
			}
			return resource.RetryableError(
				fmt.Errorf("BTP Sub Account still exist, having deletion in progress, having status %s", acc.State))
		} else {
			return resource.NonRetryableError(gAcErr)
		}
		return nil
	})
	if retryErr != nil && isResourceTimeoutError(retryErr) {
		return diag.FromErr(retryErr)
	}

	return nil
}
