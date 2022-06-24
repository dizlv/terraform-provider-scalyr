package scalyr

import (
	"context"
	"encoding/json"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

var ParserPatternSchema = map[string]*schema.Schema{
	"json": {
		Type:     schema.TypeString,
		Computed: true,
	},

	"name": {
		Type:     schema.TypeString,
		Required: true,
	},

	"value": {
		Type:     schema.TypeString,
		Required: true,
	},
}

func dataSourceParserPattern() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceParserPatternRead,

		Schema: ParserPatternSchema,
	}
}

func resourceParserPatternRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	patterns := make(scalyr.Patterns)

	name := data.Get("name").(string)
	value := data.Get("value").(string)

	patterns[name] = value

	jsonDocument, err := json.MarshalIndent(patterns, "", " ")
	if err != nil {
		return diag.FromErr(err)
	}

	jsonString := string(jsonDocument)

	data.Set("json", jsonString)
	data.SetId(strconv.Itoa(StringHashCode(jsonString)))

	return nil
}
