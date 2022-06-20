package scalyr

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ParserFormatRewriteSchema = map[string]*schema.Schema{
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

func resourceParserFormatRewrite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParserFormatRewriteCreate,
		ReadContext:   resourceParserFormatRewriteRead,
		UpdateContext: resourceParserFormatRewriteUpdate,
		DeleteContext: resourceParserFormatRewriteDelete,

		Schema: ParserFormatRewriteSchema,
	}
}

func resourceParserFormatRewriteCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatRewriteRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatRewriteUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceParserFormatRewriteDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
