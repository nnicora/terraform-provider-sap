package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"log"
	"time"
)

//STARTED, CREATING, UPDATING, MOVING, PROCESSING, DELETING, OK, PENDING_REVIEW, CANCELED,
//CREATION_FAILED, UPDATE_FAILED, UPDATE_ACCOUNT_TYPE_FAILED, UPDATE_DIRECTORY_TYPE_FAILED, PROCESSING_FAILED,
//DELETION_FAILED, MOVE_FAILED, MIGRATING, MIGRATION_FAILED, ROLLBACK_MIGRATION_PROCESSING, MIGRATED

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

	if gaId, ok := d.GetOk("global_account_id"); ok {
		if region, ok := d.GetOk("region"); ok {

			betaEnabled, _ := d.GetOk("beta_enabled")
			usedForProduction, _ := d.GetOk("used_for_production")
			displayName, _ := d.GetOk("display_name")
			description, _ := d.GetOk("description")
			subdomain, _ := d.GetOk("subdomain")
			origin, _ := d.GetOk("origin")

			customProperties := make([]btpaccounts.KeyValue, 0)
			if v, ok := d.GetOk("custom_properties"); ok {
				customProperties = expandSapBtpAccountCustomPropertiesParameters(v.([]interface{}))
			}

			subaccountAdmins := make([]string, 0)
			if v, ok := d.GetOk("sub_account_admins"); ok {
				subaccountAdmins = expandStringList(v.([]interface{}))
			}

			req := &btpaccounts.CreateSubAccountInput{
				ParentGuid:        gaId.(string),
				Region:            region.(string),
				BetaEnabled:       betaEnabled.(bool),
				DisplayName:       displayName.(string),
				Description:       description.(string),
				UsedForProduction: usedForProduction.(string),
				Subdomain:         subdomain.(string),
				CustomProperties:  customProperties,
				SubaccountAdmins:  subaccountAdmins,
				Origin:            origin.(string),
			}

			if respCreateSubAccount, err := btpAccountsClient.CreateSubAccount(ctx, req); err != nil {
				return diag.FromErr(errors.Errorf("BTP Sub Account can't be created:  %v", err))
			} else {
				aId := respCreateSubAccount.Guid
				retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
					respSubAccount, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
						SubAccountGuid: aId,
					})
					if gAcErr != nil {
						return resource.RetryableError(gAcErr)
					} else {
						if respSubAccount.State == "CREATING" {
							return resource.RetryableError(fmt.Errorf("account creation in progress"))
						}
						if respSubAccount.State == "OK" {
							return nil
						} else {
							return resource.RetryableError(fmt.Errorf("btp subaccount not yet started"))
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
					return diag.FromErr(errors.New("btp sub-account not found"))
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
		} else {
			return diag.FromErr(errors.New("region must be set when want to create an sub-account"))
		}
	} else {
		return diag.FromErr(errors.New("global_account_id must be set when want to create an sub-account"))
	}
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
		delete := param["delete"].(bool)

		prop := btpaccounts.UpdateSubAccountProperties{
			KeyValue: btpaccounts.KeyValue{
				Key:   key,
				Value: value,
			},
			Delete: delete,
		}
		parameters = append(parameters, prop)
	}

	return parameters
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func resourceSapBtpSubAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if resp, err := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
		SubAccountGuid: d.Id(),
	}); err != nil {
		return diag.FromErr(fmt.Errorf("BTP sub account can't be read: %#v", err))
	} else {
		d.SetId(resp.Guid)
		d.Set("beta_enabled", resp.BetaEnabled)
		d.Set("created_by", resp.CreatedBy)
		d.Set("created_date", resp.CreatedDate.Format(time.RFC3339))
		d.Set("modified_date", resp.ModifiedDate.Format(time.RFC3339))
		d.Set("parent_features", resp.ParentFeatures)
		d.Set("parent_id", resp.ParentGuid)
		d.Set("state", resp.State)
		d.Set("state_message", resp.StateMessage)
		d.Set("description", resp.Description)
		d.Set("display_name", resp.DisplayName)
		d.Set("global_account_id", resp.GlobalAccountGuid)
		d.Set("region", resp.Region)
		d.Set("zone_id", resp.ZoneId)
		d.Set("used_for_production", resp.UsedForProduction)
		d.Set("subdomain", resp.Subdomain)

		cp := make([]map[string]interface{}, 0)
		{
			for _, gaCP := range resp.CustomProperties {
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

	betaEnabled, _ := d.GetOk("beta_enabled")
	usedForProduction, _ := d.GetOk("used_for_production")
	displayName, _ := d.GetOk("display_name")
	description, _ := d.GetOk("description")

	customProperties := []btpaccounts.UpdateSubAccountProperties{}
	if v, ok := d.GetOk("custom_properties"); ok {
		customProperties = expandSapBtpAccountUpdateCustomPropertiesParameters(v.([]interface{}))
	}

	reqInput := &btpaccounts.UpdateSubAccountInput{
		SubAccountGuid:    d.Id(),
		BetaEnabled:       betaEnabled.(bool),
		CustomProperties:  customProperties,
		Description:       description.(string),
		DisplayName:       displayName.(string),
		UsedForProduction: usedForProduction.(string),
	}
	if resp, err := btpAccountsClient.UpdateSubAccount(ctx, reqInput); err != nil {
		return diag.FromErr(fmt.Errorf("BTP sub account can't be read: %#v", err))
	} else {
		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			if c, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
				SubAccountGuid: resp.Guid,
			}); gAcErr != nil {
				return resource.NonRetryableError(gAcErr)
			} else {
				if c.State == "UPDATING" {
					return resource.RetryableError(fmt.Errorf("account updating in progress"))
				}
				if c.State == "OK" {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("btp subaccount not yet started"))
				}
			}

			return nil
		})

		if retryErr != nil && isResourceTimeoutError(retryErr) {
			return diag.FromErr(retryErr)
		}

		saInput := &btpaccounts.GetSubAccountInput{
			SubAccountGuid: resp.Guid,
		}
		if saOutput, err := btpAccountsClient.GetSubAccount(ctx, saInput); err != nil {
			return diag.FromErr(retryErr)
		} else {

			d.SetId(saOutput.Guid)
			d.Set("beta_enabled", saOutput.BetaEnabled)
			d.Set("created_by", saOutput.CreatedBy)
			d.Set("created_date", saOutput.CreatedDate.Format(time.RFC3339))
			d.Set("modified_date", saOutput.ModifiedDate.Format(time.RFC3339))
			d.Set("parent_features", saOutput.ParentFeatures)
			d.Set("parent_id", saOutput.ParentGuid)
			d.Set("state", saOutput.State)
			d.Set("state_message", saOutput.StateMessage)
			d.Set("description", saOutput.Description)
			d.Set("display_name", saOutput.DisplayName)
			d.Set("global_account_id", saOutput.GlobalAccountGuid)
			d.Set("region", saOutput.Region)
			d.Set("zone_id", saOutput.ZoneId)
			d.Set("used_for_production", saOutput.UsedForProduction)
			d.Set("subdomain", saOutput.Subdomain)

			cp := make([]map[string]interface{}, 0)
			{
				for _, gaCP := range saOutput.CustomProperties {
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

	log.Printf("[INFO] Deleting Btp Subaccount: %s", aId)

	_, err := btpAccountsClient.DeleteSubAccount(ctx, &btpaccounts.DeleteSubAccountInput{
		SubAccountGuid: aId,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting BTO Sub-Account(%s): %w", d.Id(), err))
	}

	retryErr := resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		if acc, gAcErr := btpAccountsClient.GetSubAccount(ctx, &btpaccounts.GetSubAccountInput{
			SubAccountGuid: aId,
		}); gAcErr == nil {
			if acc.State == "DELETING" {
				return resource.RetryableError(fmt.Errorf("account still exist, having deletion in progress"))
			}
			return resource.RetryableError(fmt.Errorf("account still exist, having deletion in progress: %s", acc.State))
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
