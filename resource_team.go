package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Schema: map[string]*schema.Schema{
			"email_address": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTeamCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTeamRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
