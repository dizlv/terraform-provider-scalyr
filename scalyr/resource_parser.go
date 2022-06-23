package scalyr

import (
	"context"
	"encoding/json"
	"fmt"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Hardcoded terraform input argument names
const (
	NameArg                   = "name"
	TimeZoneArg               = "timezone"
	IntermittentTimestampsArg = "intermittent_timestamps"
	AliasToArg                = "alias_to"
	FormatsArg                = "formats"
	AttributesArg             = "attributes"
	LineGroupersArg           = "line_groupers"
	PatternsArg               = "patterns"
)

var ParserSchema = map[string]*schema.Schema{
	NameArg: {
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	},

	TimeZoneArg: {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	IntermittentTimestampsArg: {
		Type:        schema.TypeBool,
		Description: "",
		Optional:    true,
	},

	AliasToArg: {
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
	},

	FormatsArg: {
		Type:     schema.TypeList,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},

	AttributesArg: {
		Type:        schema.TypeMap,
		Description: "",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},

	LineGroupersArg: {
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},

	PatternsArg: {
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
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

func formatAttributes(attributes map[string]interface{}) scalyr.Attributes {
	formatted := make(scalyr.Attributes, len(attributes))

	for key, value := range attributes {
		strValue := fmt.Sprintf("%v", value)

		formatted[key] = strValue
	}

	return formatted
}

func resourceParserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	name := data.Get(NameArg).(string)
	formatsRaw := data.Get(FormatsArg).([]interface{})

	formats := make(scalyr.Formats, len(formatsRaw))

	// todo: refactor this bit.
	for _, f := range formatsRaw {
		var format scalyr.Format

		d := []byte(fmt.Sprint(f))

		if err := json.Unmarshal(d, &format); err != nil {
			return diag.FromErr(err)
		}

		formats = append(formats, format)
	}

	input := &scalyr.CreateParserInput{
		Name:    name,
		Formats: formats,
	}

	if v, ok := data.GetOk(TimeZoneArg); ok {
		input.TimeZone = v.(string)
	}

	if v, ok := data.GetOk(IntermittentTimestampsArg); ok {
		input.IntermittentTimestamps = v.(bool)
	}

	if v, ok := data.GetOk(AliasToArg); ok {
		input.AliasTo = v.(string)
	}

	if v, ok := data.GetOk(AttributesArg); ok {
		input.Attributes = formatAttributes(v.(map[string]interface{}))
	}

	if v, ok := data.GetOk(LineGroupersArg); ok {
		input.LineGroupers = v.(scalyr.LineGroupers)
	}

	if v, ok := data.GetOk(PatternsArg); ok {
		input.Patterns = v.(scalyr.Patterns)
	}

	if output, err := client.CreateParser(ctx, input); err != nil {
		return diag.FromErr(err)
	} else {
		data.SetId(output.Name)
	}

	return nil
}

func resourceParserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	input := &scalyr.ReadParserInput{}

	if _, err := client.ReadParser(ctx, input); err != nil {
		return diag.FromErr(err)
	} else {

	}

	return nil
}

func resourceParserUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// find this and re-create whole file with newly generated data
	client := meta.(*scalyr.ScalyrConfig)

	input := &scalyr.UpdateParserInput{}

	input.Formats = data.Get(FormatsArg).(scalyr.Formats)

	if v, ok := data.GetOk(TimeZoneArg); ok {
		input.TimeZone = v.(string)
	}

	if v, ok := data.GetOk(IntermittentTimestampsArg); ok {
		input.IntermittentTimestamps = v.(bool)
	}

	if v, ok := data.GetOk(AliasToArg); ok {
		input.AliasTo = v.(string)
	}

	if v, ok := data.GetOk(AttributesArg); ok {
		input.Attributes = v.(scalyr.Attributes)
	}

	if v, ok := data.GetOk(LineGroupersArg); ok {
		input.LineGroupers = v.(scalyr.LineGroupers)
	}

	if v, ok := data.GetOk(PatternsArg); ok {
		input.Patterns = v.(scalyr.Patterns)
	}

	readParserInput := &scalyr.ReadParserInput{}

	if _, err := client.ReadParser(ctx, readParserInput); err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.UpdateParser(ctx, input); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceParserDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	input := &scalyr.DeleteParserInput{
		Name: data.Id(),
	}

	if _, err := client.DeleteParser(ctx, input); err != nil {
		return diag.FromErr(err)
	} else {
		data.SetId("")
	}

	return nil
}
