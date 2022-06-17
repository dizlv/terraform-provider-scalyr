package scalyr

import (
	"os"
	"strconv"
	"time"

	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var u, _ = uuid.NewUUID()
var session = u.String()

func resourceEvent() *schema.Resource {
	return &schema.Resource{
		Read:   resourceEventRead,
		Create: resourceEventCreate,
		Delete: resourceEventDelete,
		Schema: map[string]*schema.Schema{
			"parser": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"message": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceEventRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceEventDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceEventCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	event := &scalyr.Event{}
	event.Ts = strconv.FormatInt(time.Now().UnixNano(), 10)
	event.Attrs = make(map[string]interface{})
	if message, ok := d.GetOk("message"); ok {
		event.Attrs["message"] = message
	}
	if parser, ok := d.GetOk("parser"); ok {
		event.Attrs["parser"] = parser
	} else {
		event.Attrs["parser"] = "terraform"
	}
	if attrs, ok := d.GetOk("attributes"); ok {
		for k, v := range attrs.(map[string]interface{}) {
			event.Attrs[k] = v
		}
	}
	hostname, _ := os.Hostname()
	client.SendEvent(event, nil, session, &scalyr.SessionInfo{ServerID: hostname})
	d.SetId(time.Now().UTC().String())
	return nil
}
