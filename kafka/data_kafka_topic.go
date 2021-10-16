package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/segmentio/kafka-go"
)

func dataKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Description: "Simplistic information about a kafka topic.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the topic",
			},
			"partitions": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "number of parititions on the topic",
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "replication factor on the topic",
			},
		},
		ReadContext: dataKafkaTopicRead,
	}
}

func dataKafkaTopicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cli := m.(*kafka.Client)

	conn, err := kafka.Dial(cli.Addr.Network(), cli.Addr.String())
	if err != nil {
		return diag.FromErr(err)
	}
	defer conn.Close()

	name := d.Get("name").(string)

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", name)
	d.Set("partitions", len(partitions))
	if len(partitions) > 0 {
		d.Set("replication_factor", len(partitions[0].Replicas))
	} else {
		d.Set("replication_factor", 0)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
