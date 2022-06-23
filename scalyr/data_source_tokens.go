package scalyr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"

	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceTokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceTokenRead,
		Schema: map[string]*schema.Schema{
			"tokens": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap, Elem: &schema.Schema{Type: schema.TypeString}},
			},
		},
	}
}

func datasourceTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	tokens, err := client.ListTokens(ctx)
	tfTokens := make([]map[string]string, len(*tokens))
	for i, token := range *tokens {
		tfTokens[i] = make(map[string]string)
		tfTokens[i]["creator"] = token.Creator
		tfTokens[i]["permission"] = token.Permission
		tfTokens[i]["id"] = token.ID
		tfTokens[i]["label"] = token.Label
		tfTokens[i]["create_date"] = token.CreateDate.String()
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving tokens: %s", err))
	}
	if err := d.Set("tokens", tfTokens); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tokens: %s", err))
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
