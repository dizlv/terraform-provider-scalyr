package scalyr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParserPatternSchema = map[string]*schema.Schema{
	"input": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	"output": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	"match": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	"replace": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	"replace_all": {
		Type:        schema.TypeBool,
		Description: "",
		Optional:    true,
	},
}

func resourceParserPattern() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceParserPatternRead,

		Schema: ParserPatternSchema,
	}
}

func resourceParserPatternRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}