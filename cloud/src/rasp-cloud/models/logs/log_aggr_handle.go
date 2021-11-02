package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"rasp-cloud/es"
	"time"

	"github.com/olivere/elastic"
)

// 攻击列表
type AttackInfoList struct {
	AttackSource string `json:"attack_source"`
	AttackType   string `json:"attack_type"`
	EventTime    string `json:"event_time"`
}

type AttackInfoLists []AttackInfoList

func (m AttackInfoLists) Len() int {
	return len(m)
}

func (m AttackInfoLists) Less(i, j int) bool {
	// timestr := "2021-10-08T17:30:29+0800"
	timei, _ := time.Parse("2006-01-02T15:04:05+0800", m[i].EventTime)
	timej, _ := time.Parse("2006-01-02T15:04:05+0800", m[j].EventTime)
	return timei.Unix() > timej.Unix()
}

func (m AttackInfoLists) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// 根据攻击时间分表
func SearchAttackAggrByIndexAndDate(index string, startTime string) ([]*elastic.AggregationBucketKeyItem, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	aggs := elastic.NewDateHistogramAggregation().
		Field("event_time"). // 根据date字段值，对数据进行分组
		//  分组间隔：month代表每月、支持minute（每分钟）、hour（每小时）、day（每天）、week（每周）、year（每年)
		Interval("day").
		// 设置返回结果中桶key的时间格式
		Format("yyyy-MM-dd")

	searchResult, err := es.ElasticClient.Search().
		Index(index). // 设置索引名
		// Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
		Query(elastic.NewRangeQuery("event_time").Gte(startTime)). // 设置查询条件
		Aggregation("attack_event_time", aggs).                    // 设置聚合条件，并为聚合条件设置一个名字
		Size(0).                                                   // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
		Do(ctx)                                                    // 执行请求

	if err != nil {
		return nil, err
	}

	// 使用Terms函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.Terms("attack_event_time")
	if !found {
		// log.Fatal("没有找到聚合数据")
		return nil, fmt.Errorf("没有找到聚合数据")
	}

	return agg.Buckets, nil

}

// 根据攻击类型分表
func SearchAttackAggrByIndex(index string) ([]*elastic.AggregationBucketKeyItem, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	aggs := elastic.NewTermsAggregation().
		Field("attack_type") // 根据attack_type字段值，对数据进行分组

	searchResult, err := es.ElasticClient.Search().
		Index(index).                      // 设置索引名
		Query(elastic.NewMatchAllQuery()). // 设置查询条件
		Aggregation("attack_type", aggs).  // 设置聚合条件，并为聚合条件设置一个名字
		Size(0).                           // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
		Do(ctx)                            // 执行请求

	if err != nil {
		return nil, err
	}

	// 使用Terms函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.Terms("attack_type")
	if !found {
		// log.Fatal("没有找到聚合数据")
		return nil, fmt.Errorf("没有找到聚合数据")
	}

	return agg.Buckets, nil
}

// 获取最新列表
func SearchAttackList(index string, sizeList int) (AttackInfoLists, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	if sizeList <= 0 {
		sizeList = 5
	}

	resultaaa := make([]AttackInfoList, 0)

	// 获取最新数据
	queryResult, err := es.ElasticClient.Search().
		Index(index).                      // 设置索引名
		Query(elastic.NewMatchAllQuery()). // 设置查询条件
		Sort("event_time", false).         // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0).                           // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(sizeList).                    // 设置分页参数 - 每页大小
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if queryResult != nil && queryResult.Hits != nil && queryResult.Hits.Hits != nil {
		hits := queryResult.Hits.Hits
		// 	// total = queryResult.Hits.TotalHits
		result := make([]map[string]interface{}, len(hits))
		for index, item := range hits {
			result[index] = make(map[string]interface{})
			// tempResult := make(map[string]interface{})
			// var filterId string
			err := json.Unmarshal(*item.Source, &result[index])
			if err != nil {
				continue
			}

			var attackInfoList AttackInfoList
			attackSource, _ := result[index]["attack_source"]
			attackInfoList.AttackSource = attackSource.(string)
			attackType, _ := result[index]["attack_type"]
			attackInfoList.AttackType = attackType.(string)
			eventTime, _ := result[index]["event_time"]
			attackInfoList.EventTime = eventTime.(string)

			resultaaa = append(resultaaa, attackInfoList)
		}
	}

	return resultaaa, nil

}
