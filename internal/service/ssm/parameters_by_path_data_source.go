// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssm

import (
	"context"
	"fmt"
	"math/rand"
	"rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_ssm_parameters_by_path", name="Parameters By Path")
func dataSourceParametersByPath() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceParametersReadByPath,

		Schema: map[string]*schema.Schema{
			names.AttrParameters: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"example": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									names.AttrARN: {
										Type:     schema.TypeString,
										Computed: true,
									},
									names.AttrName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									names.AttrType: {
										Type:     schema.TypeString,
										Computed: true,
									},
									names.AttrValue: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			names.AttrPath: {
				Type:     schema.TypeString,
				Required: true,
			},
			"recursive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"with_decryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func randomId() int {
	return rand.Int()
}

func dataSourceParametersReadByPath(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMClient(ctx)

	path := d.Get(names.AttrPath).(string)
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(d.Get("recursive").(bool)),
		WithDecryption: aws.Bool(d.Get("with_decryption").(bool)),
	}
	var output []awstypes.Parameter

	pages := ssm.NewGetParametersByPathPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading SSM Parameters by path (%s): %s", path, err)
		}

		output = append(output, page.Parameters...)
	}

	d.SetId(path)
	fmt.Println("d = ")
	fmt.Println(d)
	myMap := make(map[string]any)
	for _, parameter := range output {
		myMap[aws.ToString(parameter.Name)] = map[string]any{
			names.AttrARN:   aws.ToString(parameter.ARN),
			names.AttrType:  parameter.Type,
			names.AttrValue: aws.ToString(parameter.Value),
		}
	}
	fmt.Println("myMap = ")
	fmt.Println(myMap)
	d.Set(names.AttrParameters, schema.Set{m: myMap})

	return diags
}
