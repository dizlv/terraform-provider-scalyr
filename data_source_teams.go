package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
)

func datasourceTeams() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTeamRead,
		Schema: map[string]*schema.Schema{
			"teams": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	teams, err := client.ListTeams()
	if err != nil {
		return fmt.Errorf("Error retrieving teams: %s", err)
	}
	if err := d.Set("teams", teams); err != nil {
		return fmt.Errorf("Error setting teams: %s", err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
