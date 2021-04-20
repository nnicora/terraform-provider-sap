package sap

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"github.com/pkg/errors"
)

func dataSourceSapBtpGlobalAccountAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpGlobalAccountAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"include_auto_managed_plans": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sub_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"entitled_services": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"business_category": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"owner_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_plans": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"unlimited": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"unique_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"provisioning_method": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"amount": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"remaining_amount": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"provided_by": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"beta": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"available_for_internal": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"internal_quota_limit": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"auto_assign": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"auto_distribute_amount": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"max_allowed_sub_account_quota": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"category": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_entitlements": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"data_centers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"resources": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceSapBtpGlobalAccountAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	btpEntitlementsV1Client := meta.(*SAPClient).btpEntitlementsV1Client

	input := &btpentitlements.GlobalAccountAssignmentsInput{}
	if val, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = val.(string)
	}
	if val, ok := d.GetOk("include_auto_managed_plans"); ok {
		input.IncludeAutoManagedPlans = val.(bool)
	}
	if val, ok := d.GetOk("sub_account_id"); ok {
		input.SubAccountGuid = val.(string)
	}

	if output, err := btpEntitlementsV1Client.GetGlobalAccountAssignments(ctx, input); err != nil {
		if output != nil && output.Error != nil {
			return diag.FromErr(
				errors.Errorf("BTP Global Account assignments can't be read; %s", sap.StringValue(output.Error.Message)))
		}
		return diag.FromErr(fmt.Errorf("BTP Global Account assignments can't be read; %w", err))
	} else {
		d.Set("entitled_services", output.EntitledServices)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	if uuidString, err := uuid.GenerateUUID(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(uuidString)
	}

	return nil
}
