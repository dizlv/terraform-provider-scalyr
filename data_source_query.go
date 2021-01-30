package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
)

func fixTypes(in *[]map[string]interface{}) *[]map[string]interface{} {
	out := make([]map[string]interface{}, len(*in))
	for i, valueMap := range *in {
		out[i] = make(map[string]interface{})
		for k, v := range valueMap {
			if k == "timestamp" {
				secs, nsecs := int64(v.(float64))/1000000000, int64(v.(float64))%1000000000
				out[i][k] = time.Unix(secs, nsecs).Format(time.RFC3339)
				continue
			}
			switch v.(type) {
			case float64:
				out[i][k] = fmt.Sprintf("%v", v)
			default:
				out[i][k] = v
			}
		}
	}
	return &out
}

func datasourceQuery() *schema.Resource {
	return &schema.Resource{
		Read: datasourceQueryRead,
		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expected_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"retry_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"retry_wait": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pq",
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10 mins",
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"max_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func datasourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scalyr.ScalyrConfig)
	queryType := d.Get("query_type").(string)
	query := strings.TrimSuffix(d.Get("query").(string), "\n")
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	expectedCount := d.Get("expected_count").(int)
	retryWait := d.Get("retry_wait").(int)
	retryCount := d.Get("retry_count").(int)

	var request scalyr.Query
	if queryType == "pq" {
		request = client.NewPowerQuery(query)
	} /* else {
		request = client.NewLogQuery(query)
		request.MaxCount(max_count)
	} */
	request.Range(startTime, endTime)

	for i := 0; i < retryCount; i++ {
		res, err := request.Fetch()
		if err != nil {
			log.Fatalf("Error Executing Query: %v - %v", query, err)
		}

		log.Printf("Res: %v", &res)
		if err = d.Set("results", *fixTypes(&res)); err != nil {
			return fmt.Errorf("Error setting results - %v", err)
		}

		if expectedCount < 0 || expectedCount == int(request.Size()) {
			d.SetId(time.Now().UTC().String())
			return nil
		} else {
			time.Sleep(time.Duration(retryWait) * time.Second)
		}
	}
	return fmt.Errorf("Error Executing Query: %v", query)
}
