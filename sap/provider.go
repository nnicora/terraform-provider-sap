package sap

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/sap/session"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"log"
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

			"service_endpoint": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},

						"oauth2": {
							Type:     schema.TypeList,
							Optional: true,
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
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sap_btp_global_account": dataSourceSapBtpGlobalAccount(),
			//"sap_btp_global_account_assignments": dataSourceSapBtpGlobalAccountAssignments(),
			"sap_btp_directory":                          dataSourceSapBtpDirectory(),
			"sap_btp_sub_account":                        dataSourceSapBtpSubAccount(),
			"sap_btp_sub_account_custom_properties":      dataSourceSapBtpSubAccountCustomProperties(),
			"sap_btp_sub_account_environments_instances": dataSourceSapBtpSubAccountEnvironmentsInstances(),

			"sap_btp_directory_custom_properties": dataSourceSapBtpDirectoryCustomProperties(),

			"sap_btp_provisioning_available_environments": dataSourceSapBtpProvisioningAvailableEnvironments(),
			"sap_btp_application_registration":            dataSourceSapBtpApplicationRegistration(),
			"sap_btp_application_subscriptions":           dataSourceSapBtpApplicationSubscriptions(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"sap_btp_sub_account":                              resourceSapBtpSubAccount(),
			"sap_btp_sub_account_service_management":           resourceSapBtpSubAccountServiceManagement(),
			"sap_btp_sub_account_service_management_platforms": resourceSapBtpSubAccountServiceManagementPlatforms(),
			"sap_btp_sub_account_service_management_instances": resourceSapBtpSubAccountServiceManagementInstances(),
			"sap_btp_sub_account_service_management_bindings":  resourceSapBtpSubAccountServiceManagementBindings(),

			"sap_btp_directory":                        resourceSapBtpDirectory(),
			"sap_btp_directory_features":               resourceSapBtpDirectoryFeatures(),
			"sap_btp_directory_entitlements":           resourceSapBtpDirectoryEntitlements(),
			"sap_btp_directory_saas_entitlements":      resourceSapBtpDirectoryDynamicEntitlements("saas"),
			"sap_btp_directory_elastic_entitlements":   resourceSapBtpDirectoryDynamicEntitlements("elastic"),
			"sap_btp_directory_unlimited_entitlements": resourceSapBtpDirectoryDynamicEntitlements("unlimited"),

			"sap_btp_entitlements":           resourceSapBtpEntitlements(),
			"sap_btp_saas_entitlements":      resourceSapBtpDynamicEntitlements("saas"),
			"sap_btp_elastic_entitlements":   resourceSapBtpDynamicEntitlements("elastic"),
			"sap_btp_unlimited_entitlements": resourceSapBtpDynamicEntitlements("unlimited"),

			"sap_btp_provisioning_environments": resourceSapBtpProvisioningEnvironments(),

			"sap_btp_tenant_application_subscriptions": resourceSapBtpTenantApplicationSubscriptions(),
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
	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] Default OAuth2 configuration: %v", oauth2Map)

	defaultOAuth2 := oauth2ConfigFrom(oauth2Map)

	rawEndpoints := listFrom(d.Get("service_endpoint"))
	log.Printf("[DEBUG] Service Endpoints: %v", rawEndpoints)

	endpointsCfg := make(map[string]*sap.EndpointConfig)
	for _, rawEndpoint := range rawEndpoints {
		log.Printf("[DEBUG] Processing Service Endpoints: %v", rawEndpoint)

		endpoint := mapFrom(rawEndpoint)

		serviceId := endpoint["id"].(string)
		serviceHost := endpoint["host"].(string)

		serviceOAuth2 := defaultOAuth2
		oauth2Map := mapFrom(endpoint["oauth2"])
		if len(oauth2Map) != 0 {
			serviceOAuth2 = oauth2ConfigFrom(oauth2Map)
		}
		endpointsCfg[serviceId] = &sap.EndpointConfig{
			Host:   serviceHost,
			OAuth2: serviceOAuth2,
		}
	}

	cfg := &sap.Config{
		Endpoints:     endpointsCfg,
		DefaultOAuth2: defaultOAuth2,
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
		//btpProvisioningV1Client:     btpprovisioning.New(sess),
		//btpSaasManagerV1Client: btpsaasmanager.New(sess),
	}, nil
}

func mapFrom(block interface{}) map[string]interface{} {
	log.Printf("[DEBUG] RAW Block configuration: %v", block)

	if block == nil {
		return nil
	}

	if l, ok := block.([]interface{}); ok && len(l) > 0 && l[0] != nil {
		return l[0].(map[string]interface{})
	}
	if m, ok := block.(map[string]interface{}); ok {
		return m
	}

	return nil
}

func listFrom(block interface{}) []interface{} {
	log.Printf("[DEBUG] RAW Block configuration: %v", block)

	if block == nil {
		return nil
	}

	if l, ok := block.([]interface{}); ok {
		return l
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
