package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"time"
)

func resourceSapBtpSubAccountServiceManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountServiceManagementCreate,
		ReadContext:   resourceSapBtpSubAccountServiceManagementRead,
		UpdateContext: resourceSapBtpSubAccountServiceManagementRead,
		DeleteContext: resourceSapBtpSubAccountServiceManagementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sub_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_management_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"authentication_server_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpSubAccountServiceManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}

	subAccountId := d.Get("sub_account_id")
	input := &btpaccounts.CreateServiceManagementBindingInput{
		SubAccountGuid: subAccountId.(string),
	}
	if output, err := btpAccountsClient.CreateSubAccountServiceManagementBinding(ctx, input); err != nil {
		return resourceSapBtpSubAccountServiceManagementRead(ctx, d, meta)
	} else {
		d.Set("client_id", output.ClientId)
		d.Set("client_secret", output.ClientSecret)
		d.Set("service_management_url", output.SMUrl)
		d.Set("authentication_server_url", output.Url)
		d.Set("application_name", output.XsAppName)
	}

	return nil
}

func resourceSapBtpSubAccountServiceManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	subAccountId := d.Get("sub_account_id")
	input := &btpaccounts.GetServiceManagementBindingInput{
		SubAccountGuid: subAccountId.(string),
	}
	if output, err := btpAccountsClient.GetSubAccountServiceManagementBinding(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagementBinding can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagementBinding can't be read;  %v", err))
	} else {
		d.SetId(d.Id())
		d.Set("client_id", output.ClientId)
		d.Set("client_secret", output.ClientSecret)
		d.Set("service_management_url", output.SMUrl)
		d.Set("authentication_server_url", output.Url)
		d.Set("application_name", output.XsAppName)
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	subAccountId := d.Get("sub_account_id")
	input := &btpaccounts.DeleteServiceManagementBindingInput{
		SubAccountGuid: subAccountId.(string),
	}
	if output, err := btpAccountsClient.DeleteSubAccountServiceManagementBinding(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagementBinding can't be deleted; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagementBinding can't be deleted;  %v", err))
	}
	return nil
}
