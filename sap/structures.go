package sap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nnicora/sap-sdk-go/sap"
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
