package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
	"time"
)

func resourceSapBtpSubAccountServiceManagementBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountServiceManagementBindingCreate,
		ReadContext:   resourceSapBtpSubAccountServiceManagementBindingRead,
		UpdateContext: resourceSapBtpSubAccountServiceManagementBindingRead,
		DeleteContext: resourceSapBtpSubAccountServiceManagementBindingDelete,
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

func resourceSapBtpSubAccountServiceManagementBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if subAccountId, ok := d.GetOk("sub_account_id"); ok {

		csmbInput := &btpaccounts.CreateServiceManagementBindingInput{
			SubAccountGuid: subAccountId.(string),
		}

		if csmbOutput, err := btpAccountsClient.CreateSubAccountServiceManagementBinding(ctx, csmbInput); err != nil {
			return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagementBinding can't be created:  %v", err))
		} else {
			d.Set("client_id", csmbOutput.ClientId)
			d.Set("client_secret", csmbOutput.ClientSecret)
			d.Set("service_management_url", csmbOutput.SMUrl)
			d.Set("authentication_server_url", csmbOutput.Url)
			d.Set("application_name", csmbOutput.XsAppName)
		}
	} else {
		return diag.FromErr(errors.New("sub_account_id must be set when want to create an sub-account bindings"))
	}

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}

	return nil
}

func resourceSapBtpSubAccountServiceManagementBindingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if subAccountId, ok := d.GetOk("sub_account_id"); ok {

		csmbInput := &btpaccounts.GetServiceManagementBindingInput{
			SubAccountGuid: subAccountId.(string),
		}

		if csmbInputOutput, err := btpAccountsClient.GetSubAccountServiceManagementBinding(ctx, csmbInput); err != nil {
			return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagementBinding can't be read:  %v", err))
		} else {
			d.SetId(d.Id())
			d.Set("client_id", csmbInputOutput.ClientId)
			d.Set("client_secret", csmbInputOutput.ClientSecret)
			d.Set("service_management_url", csmbInputOutput.SMUrl)
			d.Set("authentication_server_url", csmbInputOutput.Url)
			d.Set("application_name", csmbInputOutput.XsAppName)
		}
	} else {
		return diag.FromErr(errors.New("sub_account_id must be set when want to read an sub-account bindings"))
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementBindingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if subAccountId, ok := d.GetOk("sub_account_id"); ok {
		csmbInput := &btpaccounts.DeleteServiceManagementBindingInput{
			SubAccountGuid: subAccountId.(string),
		}

		if _, err := btpAccountsClient.DeleteSubAccountServiceManagementBinding(ctx, csmbInput); err != nil {
			return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagementBinding can't be deleted:  %v", err))
		}
	} else {
		return diag.FromErr(errors.New("sub_account_id must be set when want to delete an sub-account bindings"))
	}
	return nil
}
