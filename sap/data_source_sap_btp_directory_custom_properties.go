package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
)

func dataSourceSapBtpDirectoryCustomProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpDirectoryCustomPropertiesRead,
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpDirectoryCustomPropertiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	directoryId := d.Get("directory_id")
	input := &btpaccounts.GetDirectoryCustomPropertiesInput{
		DirectoryGuid: directoryId.(string),
	}
	if output, err := btpAccountsClient.GetDirectorCustomProperties(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Directory Custom Properties can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(errors.Errorf("BTP Directory Custom Properties can't be read;  %v", err))
	} else {
		d.SetId(directoryId.(string))

		cp := make([]map[string]interface{}, 0)
		for _, cpValue := range output.Value {
			m := make(map[string]interface{})
			m["key"] = cpValue.Key
			m["value"] = cpValue.Value
			m["account_id"] = cpValue.AccountGuid

			cp = append(cp, m)
		}
		d.Set("custom_properties", cp)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
