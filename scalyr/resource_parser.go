package scalyr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParserSchema = map[string]*schema.Schema{
	"timezone": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"intermittent_timestamps": {
		Type:        schema.TypeBool,
		Description: "",
		Optional:    true,
	},

	"alias_to": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"formats": {
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
	},

	"attributes": {
		Type:        schema.TypeMap,
		Description: "",
		Optional:    true,
	},

	"line_groupers": {
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
	},

	"patterns": {
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
	},
}

func resourceParser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParserCreate,
		ReadContext:   resourceParserRead,
		UpdateContext: resourceParserUpdate,
		DeleteContext: resourceParserDelete,

		Schema: ParserSchema,
	}
}

func resourceParserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
