package sap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/service/btpmanagment"
	"github.com/pkg/errors"
	"log"
	"time"
)

func resourceSapBtpSubAccountServiceManagementPlatforms() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSapBtpSubAccountServiceManagementPlatformsCreate,
		ReadContext:   resourceSapBtpSubAccountServiceManagementPlatformsRead,
		UpdateContext: resourceSapBtpSubAccountServiceManagementPlatformsUpdate,
		DeleteContext: resourceSapBtpSubAccountServiceManagementPlatformsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_management_url": {
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceSapBtpSubAccountServiceManagementPlatformsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.CreatePlatformInput{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Labels:      expandMapListString(d.Get("labels")),
	}
	if output, err := btpServiceManagementV1Client.CreatePlatform(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Platforms can't be created; %s", output.ErrorMessage))
		}
	} else {
		d.SetId(output.Id)
		d.Set("ready", output.Ready)

		basic := map[string]string{
			"username": output.Credentials.Basic.Username,
			"password": output.Credentials.Basic.Password,
		}

		credentials := make([]map[string]string, 1)
		credentials[0] = basic
		d.Set("credentials", credentials)
	}

	return nil
}

func resourceSapBtpSubAccountServiceManagementPlatformsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.GetPlatformInput{
		PlatformID: d.Id(),
	}
	if output, err := btpServiceManagementV1Client.GetPlatform(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Platform can't be retrived; %s", output.ErrorMessage))
		}
	} else {
		d.SetId(output.Id)
		d.Set("ready", output.Ready)
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementPlatformsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.UpdatePlatformInput{
		PlatformID:  d.Id(),
		Id:          d.Id(),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
	}
	if output, err := btpServiceManagementV1Client.UpdatePlatform(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Platform can't be updated; %s", output.ErrorMessage))
		}
	}
	return nil
}

func resourceSapBtpSubAccountServiceManagementPlatformsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	session := meta.(*SAPClient).session

	oauth2Map := mapFrom(d.Get("oauth2"))
	log.Printf("[DEBUG] OAuth2 configuration: %v", oauth2Map)

	serviceManagementUrl := d.Get("service_management_url").(string)

	endpointConfig := &sap.EndpointConfig{
		Host:   serviceManagementUrl,
		OAuth2: oauth2ConfigFrom(oauth2Map),
	}
	if err := session.AddEndpointWithReplace(btpmanagment.EndpointsID, endpointConfig); err != nil {
		return diag.FromErr(errors.Errorf("BTP Sub Account ServiceManagement OAuth2;  %v", err))
	}

	btpServiceManagementV1Client := btpmanagment.New(session)

	input := &btpmanagment.DeletePlatformInput{
		PlatformID: d.Id(),
		Cascade:    true,
	}
	if output, err := btpServiceManagementV1Client.DeletePlatform(ctx, input); err != nil {
		if output != nil && output.ErrorMessage != "" {
			return diag.FromErr(
				errors.Errorf("BTP Sub Account ServiceManagement Platform can't be retrived; %s", output.ErrorMessage))
		}
	}
	return nil
}
