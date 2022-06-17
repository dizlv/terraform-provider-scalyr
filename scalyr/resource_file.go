package scalyr

import (
	"fmt"

	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		Read:   resourceFileRead,
		Delete: resourceFileDelete,
		Update: resourceFileUpdate,
		Create: resourceFileCreate,
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

func resourceFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	_, err := client.PutFile(path, content)
	if err != nil {
		return err
	}
	return resourceFileRead(d, meta)
}
func resourceFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	_, err := client.PutFile(path, content)
	if err != nil {
		return err
	}
	return resourceFileRead(d, meta)
}
func resourceFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	err := client.DeleteFile(path)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	path := d.Get("path").(string)
	res, err := client.GetFile(path)
	if err != nil {
		return err
	}
	d.Set("content", res.Content)
	d.Set("version", res.Version)
	d.Set("create_date", res.CreateDate.String())
	d.Set("mod_date", res.CreateDate.String())

	d.SetId(fmt.Sprintf("%v", res.Version))
	return nil
}
