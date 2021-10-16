package kafka

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/segmentio/kafka-go"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bootstrap_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_BOOTSTRAP_SERVERS", nil),
				Description: "a list of kafka brokers",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kafka_topic": resourceKafkaTopic(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kafka_topic": dataKafkaTopic(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	iaddrs := d.Get("bootstrap_servers").([]interface{})

	addrs := make([]string, len(iaddrs))
	for i := range iaddrs {
		addrs[i] = iaddrs[i].(string)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addrs[0])
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &kafka.Client{Addr: tcpAddr}, diags
}
