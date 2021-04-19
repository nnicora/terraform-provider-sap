package sap

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/sap/oauth2"
	"github.com/nnicora/sap-sdk-go/sap/session"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"log"
	"time"
)

var endpointServiceNames []string

func init() {
	endpointServiceNames = []string{}
}

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
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

			"endpoints": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"btp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accounts": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP BTP Accounts REST API Host",
									},
									"entitlements": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "SAP BTP Entitlements REST API Host",
									},
								},
							},
						},
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sap_btp_global_account":             dataSourceSapBtpGlobalAccount(),
			"sap_btp_global_account_assignments": dataSourceSapBtpGlobalAccountAssignments(),

			"sap_btp_sub_account":                   dataSourceSapBtpSubAccount(),
			"sap_btp_sub_account_custom_properties": dataSourceSapBtpSubAccountCustomProperties(),
			"sap_btp_account_directory":             dataSourceSapBtpAccountDirectory(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"sap_btp_sub_account":                            resourceSapBtpSubAccount(),
			"sap_btp_sub_account_service_management_binding": resourceSapBtpSubAccountServiceManagementBinding(),
			"sap_btp_account_directory":                      resourceSapBtpAccountDirectory(),

			"sap_btp_entitlement":           resourceSapBtpEntitlementFixedAssignments(),
			"sap_btp_saas_entitlement":      resourceSapBtpDynamicEntitlement("saas"),
			"sap_btp_elastic_entitlement":   resourceSapBtpDynamicEntitlement("elastic"),
			"sap_btp_unlimited_entitlement": resourceSapBtpDynamicEntitlement("unlimited"),
		},
	}

	provider.ConfigureContextFunc = func(context context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(context, d, terraformVersion)
	}

	return provider
}

type KubeProvider struct {
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	oauth2s := mapFromFromSet(d.Get("oauth2"))
	log.Printf("[DEBUG] SAP OAuth2 configuration: %v", oauth2s)

	endpoints := mapFromFromSet(d.Get("endpoints"))
	log.Printf("[DEBUG] SAP Endpoints: %v", endpoints)

	btpEndpoints := mapFromFromSet(endpoints["btp"])
	log.Printf("[DEBUG] SAP BTP Endpoints: %v", btpEndpoints)

	endpointsCfg := make(map[string]*sap.EndpointConfig, 0)
	endpointsCfg["accounts"] = &sap.EndpointConfig{
		Host: getOr(btpEndpoints, "accounts", "").(string),
	}
	endpointsCfg["entitlements"] = &sap.EndpointConfig{
		Host: getOr(btpEndpoints, "entitlements", "").(string),
	}

	cfg := &sap.Config{
		Endpoints: endpointsCfg,

		DefaultOAuth2: &oauth2.Config{
			GrantType:    getOr(oauth2s, "grant_type", "").(string),
			ClientID:     getOr(oauth2s, "client_id", "").(string),
			ClientSecret: getOr(oauth2s, "client_secret", "").(string),
			TokenURL:     getOr(oauth2s, "token_url", "").(string),
			AuthURL:      getOr(oauth2s, "authorization_url", "").(string),
			RedirectURL:  getOr(oauth2s, "redirect_url", "").(string),
			Username:     getOr(oauth2s, "username", "").(string),
			Password:     getOr(oauth2s, "password", "").(string),
			Timeout:      time.Duration(getOr(oauth2s, "timeout_seconds", 60).(int)) * time.Second,
		},
	}

	cfgBytes, _ := json.Marshal(cfg)
	log.Printf("[DEBUG] SAP BTP OAuth2 config: %s", string(cfgBytes))

	sess, err := session.BuildFromConfig(cfg)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &SAPClient{
		session:                 sess,
		btpAccountsV1Client:     btpaccounts.New(sess),
		btpEntitlementsV1Client: btpentitlements.New(sess),
	}, nil
}

func mapFromFromSet(block interface{}) map[string]interface{} {
	log.Printf("[DEBUG] RAW Block configuration: %v", block)

	if block == nil {
		return nil
	}

	if l, ok := block.([]interface{}); ok && len(l) > 0 && l[0] != nil {
		return l[0].(map[string]interface{})
	}
	if m, ok := block.(map[string]interface{}); ok {
		return m[""].(map[string]interface{})
	}

	return nil
}

func getOr(m map[string]interface{}, key string, def interface{}) interface{} {
	if m == nil {
		return def
	}
	v := m[key]
	if v == nil {
		return def
	}
	return v
}
