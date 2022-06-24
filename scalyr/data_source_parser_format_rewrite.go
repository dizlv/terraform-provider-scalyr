package scalyr

import (
	"context"
	"encoding/json"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

var ParserFormatRewriteSchema = map[string]*schema.Schema{
	"json": {
		Type:     schema.TypeString,
		Computed: true,
	},

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

func dataSourceParserFormatRewrite() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceParserFormatRewriteRead,

		Schema: ParserFormatRewriteSchema,
	}
}

func resourceParserFormatRewriteRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := data.Get("input").(string)
	output := data.Get("output").(string)
	match := data.Get("match").(string)
	replace := data.Get("replace").(string)

	document := &scalyr.Rewrite{
		Input:   input,
		Output:  output,
		Match:   match,
		Replace: replace,
	}

	if v, ok := data.GetOk("replace_all"); ok {
		document.ReplaceAll = v.(bool)
	}

	jsonDocument, err := json.MarshalIndent(document, "", " ")
	if err != nil {
		return diag.FromErr(err)
	}

	jsonString := string(jsonDocument)

	data.Set("json", jsonString)
	data.SetId(strconv.Itoa(StringHashCode(jsonString)))

	return nil
}
