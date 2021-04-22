package sap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"time"
)

func resourceSapBtpDirectoryFeatures() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpDirectoryFeaturesCreate,
		ReadContext:   resourceSapBtpDirectoryFeaturesRead,
		UpdateContext: resourceSapBtpDirectoryFeaturesUpdate,
		DeleteContext: resourceSapBtpDirectoryFeaturesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"features": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"admins": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpDirectoryFeaturesCreate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}
	directoryId := d.Get("directory_id").(string)
	input := &btpaccounts.UpdateDirectoryFeaturesInput{
		DirectoryGuid: directoryId,
	}
	if val, ok := d.GetOk("admins"); ok {
		input.DirectoryAdmins = expandStringSet(val.(*schema.Set))
	}
	if val, ok := d.GetOk("features"); ok {
		input.DirectoryFeatures = expandStringSet(val.(*schema.Set))
	}
	return updateDirectoryFeatures(ctx, "created", input, meta)
}

func resourceSapBtpDirectoryFeaturesRead(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceSapBtpDirectoryFeaturesUpdate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//directoryId := d.Get("directory_id").(string)
	//input := &btpaccounts.UpdateDirectoryFeaturesInput{
	//	DirectoryGuid: directoryId,
	//}
	//if val, ok := d.GetOk("admins"); ok {
	//	input.DirectoryAdmins = expandStringSet(val.(*schema.Set))
	//}
	//if val, ok := d.GetOk("features"); ok {
	//	input.DirectoryFeatures = expandStringSet(val.(*schema.Set))
	//}
	//return updateDirectoryFeatures(ctx, "updated", input, meta)
	return nil
}

func resourceSapBtpDirectoryFeaturesDelete(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//directoryId := d.Get("directory_id").(string)
	//input := &btpaccounts.UpdateDirectoryFeaturesInput{
	//	DirectoryGuid:     directoryId,
	//	DirectoryAdmins:   []string{},
	//	DirectoryFeatures: []string{},
	//}
	//return updateDirectoryFeatures(ctx, "deleted", input, meta)
	return nil
}

func updateDirectoryFeatures(ctx context.Context, operation string,
	input *btpaccounts.UpdateDirectoryFeaturesInput, meta interface{}) diag.Diagnostics {
	btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	if output, err := btpAccountsV1Client.UpdateDirectoryFeatures(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Directory Features can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Directory Features can't be %s;  %v", operation, err)
		}
	} else if output.StatusCode != 200 {
		if output != nil && output.Error != nil {
			return diag.Errorf("BTP Directory Features can't be %s; Operation code %v; %s",
				operation, output.StatusCode, sap.StringValue(output.Error.Message))
		} else {
			return diag.Errorf("BTP Directory Features can't be %s; Operation code %v", operation, output.StatusCode)
		}
	}

	return nil
}
