package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceSapBtpTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpTemplateCreate,
		ReadContext:   resourceSapBtpTemplateRead,
		UpdateContext: resourceSapBtpTemplateUpdate,
		DeleteContext: resourceSapBtpTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"field1": {
				Type:     schema.TypeString,
				Required: true,
			},

			"computed_field": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	return nil
}

func resourceSapBtpTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	return nil
}

func resourceSapBtpTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	return nil
}

func resourceSapBtpTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	return nil
}
