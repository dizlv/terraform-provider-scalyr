package scalyr

import (
	scalyr "ansoni/terraform-provider-scalyr/scalyr-go"
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
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
		ReadContext: resourceParserLineGrouperRead,

		Schema: ParserLineGrouperSchema,
	}
}

func resourceParserLineGrouperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	document := &scalyr.LineGrouper{}

	if v, ok := data.GetOk("start"); ok {
		document.Start = v.(string)
	}

	if v, ok := data.GetOk("continue_through"); ok {
		document.ContinueThrough = v.(string)
	}

	if v, ok := data.GetOk("continue_past"); ok {
		document.ContinuePast = v.(string)
	}

	if v, ok := data.GetOk("halt_before"); ok {
		document.HaltBefore = v.(string)
	}

	if v, ok := data.GetOk("halt_with"); ok {
		document.HaltWith = v.(string)
	}

	if v, ok := data.GetOk("max_chars"); ok {
		document.MaxChars = v.(int)
	}

	if v, ok := data.GetOk("max_lines"); ok {
		document.MaxLines = v.(int)
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
