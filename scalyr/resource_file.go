package scalyr

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceFileRead,
		DeleteContext: resourceFileDelete,
		UpdateContext: resourceFileUpdate,
		CreateContext: resourceFileCreate,
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
				Required: true,
			},
		},
	}
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	_, err := client.PutFile(ctx, path, content)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceFileRead(ctx, d, meta)
}
func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	_, err := client.PutFile(ctx, path, content)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceFileRead(ctx, d, meta)
}
func resourceFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	err := client.DeleteFile(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId(fmt.Sprintf("%v", res.Version))
	return nil
}
