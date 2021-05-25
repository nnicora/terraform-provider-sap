package sap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
	"github.com/nnicora/sap-sdk-go/sap/oauth2"
	"time"
)

func expandStringSlice(s []interface{}) []string {
	result := make([]string, len(s), len(s))
	for k, v := range s {
		// Handle the Terraform parser bug which turns empty strings in lists to nil.
		if v == nil {
			result[k] = ""
		} else {
			result[k] = v.(string)
		}
	}
	return result
}

// Takes the result of schema.Set of strings and returns a []*string
func expandStringPointerSet(configured *schema.Set) []*string {
	return expandStringPointerList(configured.List())
}

// Takes the result of schema.Set of strings and returns a []*string
func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringPointerList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, sap.String(v.(string)))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandMapString(block interface{}) map[string]string {
	result := make(map[string]string)
	if m, ok := block.(map[string]interface{}); ok {
		for k, v := range m {
			result[k] = v.(string)
		}
	}
	return result
}

func expandMapListString(block interface{}) map[string][]string {
	result := make(map[string][]string)
	if m, ok := block.(map[string]interface{}); ok {
		for k, v := range m {
			list := make([]string, 0)
			if l, ok := v.([]interface{}); ok {
				for _, item := range l {
					list = append(list, item.(string))
				}
			}
			result[k] = list
		}
	}
	return result
}

func pointersMapToStringList(pointers map[string]*string) map[string]interface{} {
	list := make(map[string]interface{}, len(pointers))
	for i, v := range pointers {
		list[i] = *v
	}
	return list
}

func stringMapToPointers(m map[string]interface{}) map[string]*string {
	list := make(map[string]*string, len(m))
	for i, v := range m {
		list[i] = sap.String(v.(string))
	}
	return list
}

func oauth2ConfigFrom(oauth2Map map[string]interface{}) *oauth2.Config {
	return &oauth2.Config{
		GrantType:    getOr(oauth2Map, "grant_type", "").(string),
		ClientID:     getOr(oauth2Map, "client_id", "").(string),
		ClientSecret: getOr(oauth2Map, "client_secret", "").(string),
		TokenURL:     getOr(oauth2Map, "token_url", "").(string),
		AuthURL:      getOr(oauth2Map, "authorization_url", "").(string),
		RedirectURL:  getOr(oauth2Map, "redirect_url", "").(string),
		Username:     getOr(oauth2Map, "username", "").(string),
		Password:     getOr(oauth2Map, "password", "").(string),
		Timeout:      time.Duration(getOr(oauth2Map, "timeout_seconds", 60).(int)) * time.Second,
	}
}
