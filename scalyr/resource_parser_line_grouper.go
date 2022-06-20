package scalyr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParserLineGrouperSchema = map[string]*schema.Schema{
	"start": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"continue_through": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"continue_past": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"halt_before": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"halt_with": {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	"max_chars": {
		Type:        schema.TypeInt,
		Description: "",
		Optional:    true,
	},

	"max_lines": {
		Type:        schema.TypeInt,
		Description: "",
		Optional:    true,
	},
}

func resourceParserLineGrouper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParserLineGrouperCreate,
		ReadContext:   resourceParserLineGrouperRead,
		UpdateContext: resourceParserLineGrouperUpdate,
		DeleteContext: resourceParserLineGrouperDelete,

		Schema: ParserLineGrouperSchema,
	}
}

func resourceParserLineGrouperCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserLineGrouperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserLineGrouperUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserLineGrouperDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
