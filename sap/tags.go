package sap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// tagsSchema returns the schema to use for tags.
//
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func tagsSchemaForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func tagsSchemaConflictsWith(conflictsWith []string) *schema.Schema {
	return &schema.Schema{
		ConflictsWith: conflictsWith,
		Type:          schema.TypeMap,
		Optional:      true,
		Elem:          &schema.Schema{Type: schema.TypeString},
	}
}
