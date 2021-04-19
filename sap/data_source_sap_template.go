package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSapTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpTemplateRead,
		Schema: map[string]*schema.Schema{
			"field1": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpAccountsV1Client := meta.(*SAPClient).btpAccountsV1Client

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
