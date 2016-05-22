package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAwsVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsVpcRead,

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_options_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"values": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_tenancy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceAwsVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	req := &ec2.DescribeVpcsInput{}

	if id := d.Get("id"); id != "" {
		req.VpcIds = []*string{aws.String(id.(string))}
	}

	req.Filters = buildEC2FilterList(d, map[string]string{
		"cidr_block":      "cidr",
		"dhcp_options_id": "dhcp-options-id",
		"default":         "isDefault",
		"state":           "state",
	})

	log.Printf("[DEBUG] DescribeVpcs %#v\n", req)
	resp, err := conn.DescribeVpcs(req)
	if err != nil {
		return err
	}
	if resp == nil || len(resp.Vpcs) == 0 {
		return fmt.Errorf("no matching VPC found")
	}
	if len(resp.Vpcs) > 1 {
		return fmt.Errorf("multiple VPCs matched; use additional constraints to reduce matches to a single VPC")
	}

	vpc := resp.Vpcs[0]

	d.SetId(*vpc.VpcId)
	d.Set("id", vpc.VpcId)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("dhcp_options_id", vpc.DhcpOptionsId)
	d.Set("instance_tenancy", vpc.InstanceTenancy)
	d.Set("default", vpc.IsDefault)
	d.Set("state", vpc.State)
	d.Set("tags", tagsToMap(vpc.Tags))

	return nil
}
