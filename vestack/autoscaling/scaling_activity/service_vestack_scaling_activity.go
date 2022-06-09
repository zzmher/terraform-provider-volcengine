package scaling_activity

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ve "github.com/volcengine/terraform-provider-vestack/common"
	"github.com/volcengine/terraform-provider-vestack/logger"
)

type VestackScalingActivityService struct {
	Client     *ve.SdkClient
	Dispatcher *ve.Dispatcher
}

func NewScalingActivityService(c *ve.SdkClient) *VestackScalingActivityService {
	return &VestackScalingActivityService{
		Client:     c,
		Dispatcher: &ve.Dispatcher{},
	}
}

func (s *VestackScalingActivityService) GetClient() *ve.SdkClient {
	return s.Client
}

func (s *VestackScalingActivityService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return ve.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		asClient := s.Client.AutoScalingClient
		action := "DescribeScalingActivities"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = asClient.DescribeScalingActivitiesCommon(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = asClient.DescribeScalingActivitiesCommon(&condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, resp)
		results, err = ve.ObtainSdkValue("Result.ScalingActivities", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingActivities is not Slice")
		}
		return data, err
	})
}

func (s *VestackScalingActivityService) ReadResource(resourceData *schema.ResourceData, activityId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if activityId == "" {
		activityId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ScalingActivityIds.1": activityId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Scaling Activity %s not exist ", activityId)
	}
	return data, err
}

func (s *VestackScalingActivityService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *VestackScalingActivityService) WithResourceResponseHandlers(data map[string]interface{}) []ve.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]ve.ResponseConvert, error) {
		return data, nil, nil
	}
	return []ve.ResourceResponseHandler{handler}
}

func (s *VestackScalingActivityService) CreateResource(*schema.ResourceData, *schema.Resource) []ve.Callback {
	return nil
}

func (s *VestackScalingActivityService) ModifyResource(*schema.ResourceData, *schema.Resource) []ve.Callback {
	return nil
}

func (s *VestackScalingActivityService) RemoveResource(*schema.ResourceData, *schema.Resource) []ve.Callback {
	return nil
}

func (s *VestackScalingActivityService) DatasourceResources(*schema.ResourceData, *schema.Resource) ve.DataSourceInfo {
	return ve.DataSourceInfo{
		RequestConverts: map[string]ve.RequestConvert{
			"ids": {
				TargetField: "ScalingActivityIds",
				ConvertType: ve.ConvertWithN,
			},
		},
		IdField:      "ScalingActivityId",
		CollectField: "activities",
		ResponseConverts: map[string]ve.ResponseConvert{
			"ScalingActivityId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *VestackScalingActivityService) ReadResourceId(id string) string {
	return id
}
