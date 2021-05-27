package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/service/btpmanagment"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func resourceSapBtpSubAccountServiceManagementBindings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountServiceManagementBindingsCreate,
		ReadContext:   resourceSapBtpSubAccountServiceManagementBindingsRead,
		//UpdateContext: resourceSapBtpSubAccountServiceManagementBindingsUpdate,
		UpdateContext: resourceSapBtpSubAccountServiceManagementBindingsRead,
		DeleteContext: resourceSapBtpSubAccountServiceManagementBindingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_management": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
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
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resources": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeList},
			},

			"ready": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"context": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpSubAccountServiceManagementBindingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session
	serviceList := d.Get("service_management").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("Service management is required")
	}

	err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Service Management OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.CreateServiceBindingInput{
		Async:             false,
		Name:              d.Get("name").(string),
		ServiceInstanceId: d.Get("service_instance_id").(string),
		Parameters:        expandMapString(d.Get("parameters")),
		BindResource:      expandMapString(d.Get("resources")),
		Labels:            expandMapListString(d.Get("labels")),
	}

	if output, err := btpServiceManagementV1Client.CreateServiceBinding(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Bindings can't be created; %s", output.ErrorMessage))
		}
	} else {
		d.SetId(output.Id)
		d.Set("ready", output.Ready)
		d.Set("context", output.Context)

		data := make(map[string]string)
		flatMap("", output.Credentials, data)
		d.Set("credentials", data)
	}

	return nil
}

func resourceSapBtpSubAccountServiceManagementBindingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session
	serviceList := d.Get("service_management").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("Service management service is required")
	}

	err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Service Management OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.GetServiceBindingInput{
		ServiceBindingID: d.Id(),
	}
	if output, err := btpServiceManagementV1Client.GetServiceBinding(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Bindings can't be created; %s", output.ErrorMessage))
		}
	} else {
		d.Set("ready", output.Ready)
		d.Set("context", output.Context)

		data := make(map[string]string)
		flatMap("", output.Credentials, data)
		d.Set("credentials", data)
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementBindingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSapBtpSubAccountServiceManagementBindingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session
	serviceList := d.Get("service_management").([]interface{})
	if len(serviceList) < 1 {
		return diag.Errorf("Service management service is required")
	}

	err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, extractEndpointConfig(serviceList))
	if err != nil {
		return diag.FromErr(errors.Errorf("BTP Service Management OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.DeleteServiceBindingInput{
		ServiceBindingID: d.Id(),
		Async:            false,
	}
	if output, err := btpServiceManagementV1Client.DeleteServiceBinding(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Bindings can't be created; %s", output.ErrorMessage))
		}
	}
	return nil
}

func flatMap(prefix string, src map[string]interface{}, dst map[string]string) {
	for k, v := range src {
		if nestedMap, ok := v.(map[string]interface{}); ok {
			flatMap(k, nestedMap, dst)
		} else {
			if prefix != "" && !strings.HasSuffix(prefix, ".") {
				prefix = prefix + "."
			}
			key := prefix + k
			dst[key] = v.(string)
		}
	}
}
