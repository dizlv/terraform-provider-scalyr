package scalyr

import (
	"context"
	"encoding/json"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"hash/crc32"
	"strconv"
)

func StringHashCode(input string) int {
	value := int(crc32.ChecksumIEEE([]byte(input)))
	if value >= 0 {
		return value
	}

	if -value >= 0 {
		return -value
	}

	return 0
}

var parserFormatSchema = map[string]*schema.Schema{
	"json": {
		Type:     schema.TypeString,
		Computed: true,
	},

	"name": {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

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
		Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceParserFormat() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceParserFormatRead,

		Schema: parserFormatSchema,
	}
}

func dataSourceParserFormatRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := data.Get("name").(string)
	format := data.Get("format").(string)

	mergedDocument := &scalyr.Format{
		Name:   name,
		Format: format,
	}

	if v, ok := data.GetOk("halt"); ok {
		mergedDocument.Halt = v.(bool)
	}

	if v, ok := data.GetOk("repeat"); ok {
		mergedDocument.Repeat = v.(bool)
	}

	if v, ok := data.GetOk("discard"); ok {
		mergedDocument.Discard = v.(bool)
	}

	// load rewrites from string
	if v, ok := data.GetOk("rewrites"); ok {
		rewrites, err := TransformType[[]any, scalyr.Rewrite](v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		mergedDocument.Rewrites = rewrites
	}

	jsonDocument, err := json.MarshalIndent(mergedDocument, "", " ")
	if err != nil {
		return diag.FromErr(err)
	}

	jsonString := string(jsonDocument)

	data.Set("json", jsonString)
	data.SetId(strconv.Itoa(StringHashCode(jsonString)))

	return nil
}
