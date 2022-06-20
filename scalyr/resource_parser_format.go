package scalyr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParserFormatSchema = map[string]*schema.Schema{
	"format": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	"halt": {
		Type:        schema.TypeBool,
		Description: "",
		Optional:    true,
	},

	"rewrites": {
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
	},

	"repeat": {
		Type:     schema.TypeBool,
		Optional: true,
	},

	"discard": {
		Type:        schema.TypeBool,
		Description: "",
		Optional:    true,
	},
}

func resourceParserFormat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParserFormatCreate,
		ReadContext:   resourceParserFormatRead,
		UpdateContext: resourceParserFormatUpdate,
		DeleteContext: resourceParserFormatDelete,

		Schema: ParserFormatSchema,
	}
}

func resourceParserFormatCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
