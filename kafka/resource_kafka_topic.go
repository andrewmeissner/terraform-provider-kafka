package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/segmentio/kafka-go"
)

func resourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Description: "Represents simplistic information about a kafka topic.  Altering any of these will force a new resource to be created.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the topic",
				ForceNew:    true,
			},
			"partitions": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "number of partitions",
				ForceNew:    true,
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "replication factor",
				ForceNew:    true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		CreateContext: resourceKafkaTopicCreate,
		ReadContext:   resourceKafkaTopicRead,
		UpdateContext: resourceKafkaTopicUpdate,
		DeleteContext: resourceKafkaTopicDelete,
	}
}

func resourceKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cli := m.(*kafka.Client)

	resp, err := cli.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             d.Get("name").(string),
				NumPartitions:     d.Get("partitions").(int),
				ReplicationFactor: d.Get("replication_factor").(int),
			},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil && len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			diags = append(diags, diag.FromErr(err)...)
		}
		return diags
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceKafkaTopicRead(ctx, d, m)
}

func resourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cli := m.(*kafka.Client)

	conn, err := kafka.Dial(cli.Addr.Network(), cli.Addr.String())
	if err != nil {
		return diag.FromErr(err)
	}
	defer conn.Close()

	name := d.Get("name").(string)

	partitions, err := conn.ReadPartitions(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", name)
	d.Set("partitions", len(partitions))
	if len(partitions) > 0 {
		d.Set("replication_factor", len(partitions[0].Replicas))
	} else {
		d.Set("replication_factor", -1)
	}

	return diags
}

func resourceKafkaTopicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d := resourceKafkaTopicDelete(ctx, d, m); len(d) > 0 {
		return d
	}

	if c := resourceKafkaTopicCreate(ctx, d, m); len(c) > 0 {
		return c
	}
	d.Set("last_updated", time.Now().String())

	return resourceKafkaTopicRead(ctx, d, m)
}

func resourceKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cli := m.(*kafka.Client)

	name := d.Get("name").(string)

	resp, err := cli.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{name},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil && len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "encountered an error",
					Detail:   err.Error(),
				})
			}
		}
	}

	d.SetId("")

	return diags
}
