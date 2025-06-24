package httputil

import (
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type MetricsFromExpr []MetricsInfo // 一个表达式得到的数据

type MetricsInfo struct { // 具体某个时序数据，比如GPU0
	Metric map[string]string `json:"metric"`
	Values []MetricsValues   `json:"values"` // 时序数值
}

type MetricsValues struct {
	Value                 float64 `json:"value"`
	Timestamp             int64   `json:"timestamp"`
	Color                 string  `json:"color"`            // 这是panel上设置的
	BandwidthPercent      float64 `json:"bandwidthPercent"` // 所占带宽百分比
	BandwidthPercentColor string  `json:"bandwidthColor"`   // 带宽百分比颜色
}

func PromCommonModelValue(promValues []model.Value) ([]MetricsFromExpr, error) {
	var ret []MetricsFromExpr

	for _, result := range promValues {
		mData, err := ParseModelValue2MetricsData(result)
		if err != nil {
			return nil, err
		}

		if len(mData) > 0 {
			ret = append(ret, mData)
		}
	}

	return ret, nil
}

// 一个表达式得到的数据
func ParseModelValue2MetricsData(commonModelValue model.Value) (MetricsFromExpr, error) {
	var ret MetricsFromExpr
	switch commonModelValue.Type() {
	case model.ValScalar:
		logrus.Warnf("need to parse 'Scalar' type value")
	case model.ValVector:
		vector := commonModelValue.(model.Vector)
		for _, sample := range vector {
			var mtc = make(map[string]string)
			for k, v := range sample.Metric {
				mtc[string(k)] = string(v)
			}

			ret = append(ret, MetricsInfo{
				Metric: mtc,
				Values: []MetricsValues{
					{
						Value:     float64(sample.Value),
						Timestamp: sample.Timestamp.Unix(),
					},
				},
			})
		}

		logrus.Warnf("need to parse 'Vector' type value")
	case model.ValMatrix:
		matrix := commonModelValue.(model.Matrix)
		for _, sample := range matrix {
			var values []MetricsValues
			for _, value := range sample.Values {
				values = append(values, MetricsValues{
					Value:     float64(value.Value),
					Timestamp: value.Timestamp.Unix(),
				})
			}

			var m = make(map[string]string)
			for k, v := range sample.Metric {
				m[string(k)] = string(v)
			}

			ret = append(ret, MetricsInfo{
				Metric: m,
				Values: values,
			})
		}

	case model.ValString:
		logrus.Warnf("need to parse 'String' type value")
	default:
		return nil, fmt.Errorf("unknown metric type: %s", commonModelValue.Type())
	}

	return ret, nil
}
