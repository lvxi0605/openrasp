//Copyright 2021-2021 corecna Inc.

package models

import (
	"context"
	"rasp-cloud/es"
	"time"

	"github.com/olivere/elastic"
)

type ReportData struct {
	RaspId     string `json:"rasp_id"`
	Time       int64  `json:"time"`
	RequestSum int64  `json:"request_sum"`
	InsertTime int64  `json:"@timestamp"`
}

var (
	ReportIndexName      = "corerasp-report-data"
	AliasReportIndexName = "real-corerasp-report-data"
	reportType           = "report-data"
)

func init() {
	es.RegisterTTL(24*100*time.Hour, AliasReportIndexName+"-*")
}

func CreateReportDataEsIndex(appId string) error {
	return es.CreateEsIndex(ReportIndexName+"-"+appId,
		AliasReportIndexName+"-"+appId, reportType+"-template")
}

func CreateDependencyEsIndex(appId string) error {
	return es.CreateEsIndex(DependencyIndexName+"-"+appId,
		AliasDependencyIndexName+"-"+appId, "dependency-data-template")
}

func AddReportData(reportData *ReportData, appId string) error {
	reportData.InsertTime = time.Now().Unix() * 1000
	return es.Insert(AliasReportIndexName+"-"+appId, reportType, reportData)
}

func GetHistoryRequestSum(startTime int64, endTime int64, interval string, timeZone string,
	appId string) (error, []map[string]interface{}) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	timeAggrName := "aggr_time"
	sumAggrName := "request_sum"
	timeAggr := elastic.NewDateHistogramAggregation().Field("time").TimeZone(timeZone).
		Interval(interval).ExtendedBounds(startTime, endTime)
	requestSumAggr := elastic.NewSumAggregation().Field("request_sum")
	timeAggr.SubAggregation(sumAggrName, requestSumAggr)
	timeQuery := elastic.NewRangeQuery("time").Gte(startTime).Lte(endTime)
	aggrResult, err := es.ElasticClient.Search(AliasReportIndexName+"-"+appId).
		Query(timeQuery).
		Aggregation(timeAggrName, timeAggr).
		Size(0).
		Do(ctx)
	if err != nil {
		return err, nil
	}
	result := make([]map[string]interface{}, 0)
	if aggrResult != nil && aggrResult.Aggregations != nil {
		if terms, ok := aggrResult.Aggregations.Terms(timeAggrName); ok && terms.Buckets != nil {
			result = make([]map[string]interface{}, len(terms.Buckets))
			for index, item := range terms.Buckets {
				result[index] = make(map[string]interface{})
				result[index]["start_time"] = item.Key
				if sumItem, ok := item.Sum(sumAggrName); ok {
					result[index]["request_sum"] = sumItem.Value
				} else {
					result[index]["request_sum"] = 0
				}
			}
		}
	}
	return nil, result
}
