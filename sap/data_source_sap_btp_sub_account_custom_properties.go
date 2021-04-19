package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/pkg/errors"
)

func dataSourceSapBtpSubAccountCustomProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpSubAccountCustomPropertiesRead,
		Schema: map[string]*schema.Schema{
			"sub_account_id": {
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

func dataSourceSapBtpSubAccountCustomPropertiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpAccountsClient := meta.(*SAPClient).btpAccountsV1Client

	if subAccountId, ok := d.GetOk("sub_account_id"); ok {

		sacpInput := &btpaccounts.GetCustomPropertiesInput{
			SubAccountGuid: subAccountId.(string),
		}

		if csmbInputOutput, err := btpAccountsClient.GetSubAccountCustomProperties(ctx, sacpInput); err != nil {
			return diag.FromErr(errors.Errorf("BTP Sub Account Custom Properties can't be read:  %v", err))
		} else {
			d.SetId(subAccountId.(string))

			cp := make([]map[string]interface{}, 0)
			for _, cpValue := range csmbInputOutput.Value {
				m := make(map[string]interface{})
				m["key"] = cpValue.Key
				m["value"] = cpValue.Value
				m["account_id"] = cpValue.AccountGuid

				cp = append(cp, m)
			}
			d.Set("custom_properties", cp)
		}
	} else {
		return diag.FromErr(errors.New("sub_account_id must be set when want to read an sub-account custom properties"))
	}
	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
