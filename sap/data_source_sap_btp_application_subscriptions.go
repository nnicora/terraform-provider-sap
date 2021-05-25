package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpsaasprovisioning"
)

func dataSourceSapBtpApplicationSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSapBtpApplicationSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},

			"oauth2": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"grant_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "client_credentials",
							Description: "SAP OAuth2 Grant Type.",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SAP OAuth2 Client Id.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SAP OAuth2 Client Secret.",
						},
						"token_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SAP OAuth2 Token Url.",
						},
						"authorization_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "SAP OAuth2 Authorization Url.",
						},
						"redirect_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "SAP OAuth2 Redirect Url.",
						},

						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "SAP OAuth2 Username. Used in case if 'grant_type=password'.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "SAP OAuth2 Password. Used in case if 'grant_type=password'.",
						},

						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     60,
							Description: "SAP OAuth2 HTTP Client timeout.",
						},
					},
				},
			},

			"global_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subscriptions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"changed_on": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"consumer_tenant_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"error": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"global_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_consumer_tenant_active": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"license_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"service_instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"subdomain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dependencies": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"error": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"xsappname": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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

func dataSourceSapBtpApplicationSubscriptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//btpSaaSProvisioningV1Client := meta.(*SAPClient).btpSaaSProvisioningV1Client
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	provisioning := &sap.EndpointConfig{
		Host:   d.Get("host").(string),
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	session.AddEndpointWithReplace("saas-manager", provisioning)
	btpSaaSProvisioningV1Client := btpsaasprovisioning.New(session)

	input := &btpsaasprovisioning.GetApplicationSubscriptionsInput{}
	if val, ok := d.GetOk("global_account_id"); ok {
		strVal := val.(string)
		if len(strVal) > 0 {
			input.GlobalAccountId = strVal
		}
	}
	if val, ok := d.GetOk("sub_account_id"); ok {
		strVal := val.(string)
		if len(strVal) > 0 {
			input.SubAccountId = strVal
		}
	}
	if val, ok := d.GetOk("tenant_id"); ok {
		strVal := val.(string)
		if len(strVal) > 0 {
			input.TenantId = strVal
		}
	}
	if val, ok := d.GetOk("state"); ok {
		strVal := val.(string)
		if len(strVal) > 0 {
			input.State = strVal
		}
	}
	if output, err := btpSaaSProvisioningV1Client.GetApplicationSubscriptions(ctx, input); err != nil {
		return diag.Errorf("BTP SaaS Subscription to an application can't be done;  %v", err)
	} else {
		subs := make([]map[string]interface{}, 0, len(output.Values))
		for _, sub := range output.Values {
			subMap := make(map[string]interface{})
			subMap["amount"] = sub.Amount
			subMap["app_name"] = sub.AppName
			subMap["changed_on"] = sub.ChangedOn
			subMap["code"] = sub.Code
			subMap["consumer_tenant_id"] = sub.ConsumerTenantId
			subMap["created_on"] = sub.CreatedOn
			subMap["error"] = sub.Error
			subMap["global_account_id"] = sub.GlobalAccountId
			subMap["is_consumer_tenant_active"] = sub.IsConsumerTenantActive
			subMap["license_type"] = sub.LicenseType
			subMap["service_instance_id"] = sub.ServiceInstanceId
			subMap["state"] = sub.State
			subMap["sub_account_id"] = sub.SubAccountId
			subMap["subdomain"] = sub.Subdomain
			subMap["url"] = sub.Url

			deps := make(map[string]interface{})
			for _, dep := range sub.Dependencies {
				deps["app_name"] = dep.AppName
				deps["error"] = dep.Error
				deps["xsappname"] = dep.XSAppName
			}
			subMap["dependencies"] = deps

			subs = append(subs, subMap)
		}

		d.Set("subs", subs)
	}

	tags := make(map[string]interface{})
	{
		// TODO
	}
	d.Set("tags", tags)

	return nil
}
