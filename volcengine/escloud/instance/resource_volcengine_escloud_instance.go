package instance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ve "github.com/volcengine/terraform-provider-volcengine/common"
)

/*

Import
ESCloud Instance can be imported using the id, e.g.
```
$ terraform import volcengine_escloud_instance.default n769ewmjjqyqh5dv
```

*/

func ResourceVolcengineESCloudInstance() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceVolcengineESCloudInstanceCreate,
		Read:   resourceVolcengineESCloudInstanceRead,
		Update: resourceVolcengineESCloudInstanceUpdate,
		Delete: resourceVolcengineESCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,

				Description: "The configuration of ESCloud instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The version of ESCloud instance.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The region ID of ESCloud instance.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The available zone ID of ESCloud instance.",
						},
						"zone_number": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The zone count of the ESCloud instance used.",
						},
						"enable_https": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Whether Https access is enabled.",
						},
						"admin_user_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of administrator account(should be admin).",
						},
						"admin_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The password of administrator account.",
						},
						"charge_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PostPaid",
								"PrePaid",
							}, false),
							Description: "The charge type of ESCloud instance.",
						},
						"configuration_code": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Configuration code used for billing.",
						},
						"enable_pure_master": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Whether the Master node is independent.",
						},
						"node_specs_assigns": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "The number and configuration of various ESCloud instance node.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The type of node.",
									},
									"number": {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "The number of node.",
									},
									"resource_spec_name": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The name of compute resource spec.",
									},
									"storage_spec_name": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The name of storage spec.",
									},
									"storage_size": {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "The size of storage.",
									},
								},
							},
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of ESCloud instance.",
						},
						"vpc": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Information about the VPC where the instance is located.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of vpc.",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of vpc.",
									},
								},
							},
						},
						"subnet": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The ID of subnet, the subnet must belong to the AZ selected.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of subnet.",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of subnet.",
									},
								},
							},
						},
						"project_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The project name  to which the ESCloud instance belongs.",
						},
						"maintenance_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The maintainable time period for the instance.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// 创建时不存在这个参数，修改时存在这个参数
								return d.Id() == ""
							},
						},
						"maintenance_day": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The maintainable date for the instance.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// 创建时不存在这个参数，修改时存在这个参数
								return d.Id() == ""
							},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}

	return resource
}

func resourceVolcengineESCloudInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewESCloudInstanceService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceVolcengineESCloudInstance())
	if err != nil {
		return fmt.Errorf("Error on creating ESCloud instance %q,%s", d.Id(), err)
	}
	return resourceVolcengineESCloudInstanceRead(d, meta)
}

func resourceVolcengineESCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewESCloudInstanceService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceVolcengineESCloudInstance())
	if err != nil {
		return fmt.Errorf("error on updating ESCloud instance  %q, %s", d.Id(), err)
	}
	return resourceVolcengineESCloudInstanceRead(d, meta)
}

func resourceVolcengineESCloudInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewESCloudInstanceService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceVolcengineESCloudInstance())
	if err != nil {
		return fmt.Errorf("error on deleting ecs instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceVolcengineESCloudInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewESCloudInstanceService(meta.(*ve.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceVolcengineESCloudInstance())
	if err != nil {
		return fmt.Errorf("Error on reading ESCloud instance %q,%s", d.Id(), err)
	}
	return err
}