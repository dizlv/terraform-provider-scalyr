package scalyr

import (
	"context"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func datasourceFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceFileRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mod_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	res, err := client.GetFile(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("content", res.Content)
	d.Set("version", res.Version)
	d.Set("create_date", res.CreateDate.String())
	d.Set("mod_date", res.CreateDate.String())

	d.SetId(time.Now().UTC().String())
	return nil
}
