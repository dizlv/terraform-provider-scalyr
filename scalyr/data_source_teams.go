package scalyr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"

	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceTeamRead,
		Schema: map[string]*schema.Schema{
			"teams": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	teams, err := client.ListTeams(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving teams: %s", err))
	}
	if err := d.Set("teams", teams); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting teams: %s", err))
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
