package tools

import (
	"context"
	"fmt"
	"log"
	"rasp-cloud/config"
	"rasp-cloud/ipsearch"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/olivere/elastic"
)

func TestLong2String(t *testing.T) {

	appkey, _ := long2String(10000001)
	appid, _ := string2Long(appkey)
	checkapp := checkAppKey(appkey)
	log.Printf("TestLong2String.appid[%v] appkey[%v], checkapp[%v]", appid, appkey, checkapp)

	appkey1, appsecret1, _ := createAppkeyAndSecret(10000001, "terry", "WLAN")

	checkflag := CheckAppkeyAndSecret(appkey1, appsecret1)

	log.Printf("TestLong2String.appkey1[%v] appsecret1[%v], checkflag[%v]", appkey1, appsecret1, checkflag)
}

func TestGetNetworkName(t *testing.T) {
	testSerial := getCurrentSerial("WLAN")
	log.Printf("TestGetNetworkName.getCurrentSerial[%v]", testSerial)
}

func TestHTTPDns(t *testing.T) {
	// response, err := http.Get("http://203.107.1.33/174597/d?host=core.xinghuoyouxi.com")
	// if err != nil {
	// 	log.Printf("TestHTTPDns.err[%v]", err)
	// 	return
	// }

	// body, err1 := ioutil.ReadAll(response.Body)
	// if err1 != nil {
	// 	log.Printf("TestHTTPDns.err1[%v]", err1)
	// 	return
	// }

	// log.Printf("TestHTTPDns.body[%v]", string(body))
	ip := GetDnsIP("core.xinghuoyouxi.com")
	log.Printf("TestHTTPDns.ip[%v]", ip)

}

func TestSendMessage(t *testing.T) {

	if _, err := toml.DecodeFile("../conf/local.toml", &config.TOMLConfig); err != nil {
		log.Printf("err: Config: %v\n", err)
		return
	}

	log.Printf("DEBUG: Config: %#v\n", config.TOMLConfig)

	testMap := make(map[string]interface{})

	testMap["test"] = "test"

	SendMessage("key", testMap)
}

func TestTimeSub(t *testing.T) {
	t1 := time.Now()
	t2 := time.Unix(1635322292, 0)
	left := t1.Sub(t2)
	log.Printf("left [%v]", int(left.Hours()/24)+1)

	timestr := "2021-10-08T17:30:29+0800"
	timeaa, _ := time.Parse("2006-01-02T15:04:05+0800", timestr)
	log.Printf("%d", timeaa.Unix())

}

func TestLogAggrHandle(t *testing.T) {

	if _, err := toml.DecodeFile("D:/code/ncpost/openrasp/cloud/src/rasp-cloud/conf/local.toml", &config.TOMLConfig); err != nil {
		return
	}

	log.Printf("DEBUG: Config: %#v\n", config.TOMLConfig)

	esAddr := []string{"http://111.74.2.60:9200"}
	client, err := elastic.NewSimpleClient(elastic.SetURL(esAddr...),
		elastic.SetBasicAuth("", ""),
		elastic.SetSnifferTimeoutStartup(5*time.Second),
		elastic.SetSnifferTimeout(5*time.Second),
		elastic.SetSnifferInterval(30*time.Minute))
	if err != nil {
		log.Printf("init ES failed: %v", err)
	}

	// ips, _ := ipsearch.New()

	Version, err := client.ElasticsearchVersion(esAddr[0])
	if err != nil {
		log.Printf("failed to get es version: %v", err)
	}
	log.Printf("ES version: %v", Version)

	ctx := context.Background()

	// // 创建Value Count指标聚合
	// aggs := elastic.NewValueCountAggregation().
	// 	Field("attack_type") // 设置统计字段

	// searchResult, err := client.Search().
	// 	Index("corerasp-attack-alarm-a749803ba7d6653fcf793b4d13569ad58a841e12"). // 设置索引名
	// 	Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
	// 	Aggregation("total", aggs).                                              // 设置聚合条件，并为聚合条件设置一个名字, 支持添加多个聚合条件，命名不一样即可。
	// 	Size(0).                                                                 // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
	// 	Do(ctx)                                                                  // 执行请求

	// if err != nil {
	// 	// Handle error
	// 	// panic(err)
	// 	log.Printf("init ES failed: %v", err)
	// }

	// // 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	// agg, found := searchResult.Aggregations.ValueCount("total")
	// if found {
	// 	// 打印结果，注意：这里使用的是取值运算符
	// 	fmt.Println(*agg.Value)
	// }

	// 创建Terms桶聚合
	// aggs := elastic.NewTermsAggregation().
	// 	Field("attack_type") // 根据attack_type字段值，对数据进行分组

	// searchResult, err := client.Search().
	// 	Index("corerasp-attack-alarm-a749803ba7d6653fcf793b4d13569ad58a841e12"). // 设置索引名
	// 	Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
	// 	Aggregation("attack_type", aggs).                                        // 设置聚合条件，并为聚合条件设置一个名字
	// 	Size(0).                                                                 // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
	// 	Do(ctx)                                                                  // 执行请求

	// if err != nil {
	// 	// Handle error
	// 	// panic(err)
	// }

	// 使用Terms函数和前面定义的聚合条件名称，查询结果
	// agg, found := searchResult.Aggregations.Terms("attack_type")
	// if !found {
	// 	log.Fatal("没有找到聚合数据")
	// }

	// // 遍历桶数据
	// for _, bucket := range agg.Buckets {
	// 	// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
	// 	bucketValue := bucket.Key

	// 	// 打印结果， 默认桶聚合查询，都是统计文档总数
	// 	fmt.Printf("bucket = %q 文档总数 = %d\n", bucketValue, bucket.DocCount)
	// }

	// aggs := elastic.NewDateHistogramAggregation().
	// 	Field("event_time"). // 根据date字段值，对数据进行分组
	// 	//  分组间隔：month代表每月、支持minute（每分钟）、hour（每小时）、day（每天）、week（每周）、year（每年)
	// 	Interval("day").
	// 	// 设置返回结果中桶key的时间格式
	// 	Format("yyyy-MM-dd")

	// // elastic.NewRangeQuery("event_time").Gt("2021-10-21")

	// searchResult, err := client.Search().
	// 	Index("corerasp-attack-alarm-a749803ba7d6653fcf793b4d13569ad58a841e12"). // 设置索引名
	// 	// Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
	// 	Query(elastic.NewRangeQuery("event_time").Gte("2021-10-09")). // 设置查询条件
	// 	Aggregation("sales_over_time", aggs).                         // 设置聚合条件，并为聚合条件设置一个名字
	// 	Size(0).                                                      // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
	// 	Do(ctx)                                                       // 执行请求

	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }

	// // 使用DateHistogram函数和前面定义的聚合条件名称，查询结果
	// agg, found := searchResult.Aggregations.DateHistogram("sales_over_time")
	// if !found {
	// 	log.Fatal("没有找到聚合数据")
	// }

	aggs := elastic.NewTermsAggregation().
		Field("new_location_en.keyword") // 根据attack_type字段值，对数据进行分组

	// 创建Sum指标聚合
	countAggs := elastic.NewTermsAggregation().Field("new_location_zh_cn.keyword")
	aggs.SubAggregation("new_location_zh_cn", countAggs)

	searchResult, err := client.Search().
		Index("corerasp-attack-alarm-97d54c6615b69342ac525896f02d491ee852157e"). // 设置索引名
		Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
		Aggregation("new_location_en", aggs).                                    // 设置聚合条件，并为聚合条件设置一个名字
		Size(0).                                                                 // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
		Do(ctx)                                                                  // 执行请求

	if err != nil {
		// return nil, err
		fmt.Printf("%v", err)
		return
	}

	// 使用Terms函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.Terms("new_location_en")
	if !found {
		// log.Fatal("没有找到聚合数据")
		// return nil, fmt.Errorf("没有找到聚合数据")
	}

	// agg1, found1 := searchResult.Aggregations.Terms("new_location_zh_cn")
	// if !found1 {

	// }
	// for _, bucket := range agg1.Buckets {
	// 	// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
	// 	// bucketValue := *bucket.KeyAsString

	// 	// 打印结果， 默认桶聚合查询，都是统计文档总数
	// 	// keyString := fmt.Sprintf("%q", bucket.Aggregations)
	// 	fmt.Printf("bucket = %q 文档总数 = %d\n", bucket.Key, bucket.DocCount)
	// }

	// 遍历桶数据
	for _, bucket := range agg.Buckets {
		// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
		// bucketValue := *bucket.KeyAsString

		// 打印结果， 默认桶聚合查询，都是统计文档总数
		// keyString := fmt.Sprintf("%q", bucket.Aggregations)
		childagg, isfind2 := bucket.Aggregations.Terms("new_location_zh_cn")
		// aggchild, isfind2 := childagg.Terms("buckets")
		if isfind2 {
			for _, bucket1 := range childagg.Buckets {
				fmt.Printf("bucket = %q  bucket1 = %q 文档总数 = %d\n", bucket.Key, bucket1.Key, bucket1.DocCount)
			}
		}
		fmt.Printf("childagg: %v\n", childagg)
		// for _, bucket1 := range childagg.Meta {

		// }

		// fmt.Printf("bucket = %q 文档总数 = %d, %v\n", bucket.Key, bucket.DocCount, isfind)
	}

	// queryResult, err := client.Search().
	// 	Index("corerasp-attack-alarm-97d54c6615b69342ac525896f02d491ee852157e"). // 设置索引名
	// 	Query(elastic.NewMatchAllQuery()).                                       // 设置查询条件
	// 	Sort("event_time", false).                                               // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
	// 	From(0).                                                                 // 设置分页参数 - 起始偏移量，从第0行记录开始
	// 	Size(5000).                                                              // 设置分页参数 - 每页大小
	// 	Do(ctx)

	// if err != nil {
	// 	fmt.Printf("queryResult err:%d\n", err)
	// }

	// fmt.Printf("queryResult TotalHits:%d\n", len(queryResult.Hits.Hits))

	// if queryResult != nil && queryResult.Hits != nil && queryResult.Hits.Hits != nil {
	// 	hits := queryResult.Hits.Hits
	// 	// 	// total = queryResult.Hits.TotalHits
	// 	result := make([]map[string]interface{}, len(hits))
	// 	// result1 := make([]map[string]interface{}, len(hits))
	// 	for index, item := range hits {
	// 		result[index] = make(map[string]interface{})
	// 		resultTemp := make(map[string]interface{})
	// 		// var filterId string
	// 		err := json.Unmarshal(*item.Source, &result[index])
	// 		if err != nil {
	// 			// return  nil, err
	// 		}

	// 		upsertId, _ := result[index]["upsert_id"]

	// 		attackSource, _ := result[index]["attack_source"]
	// 		resultTemp["attack_source"] = attackSource
	// 		attackType, _ := result[index]["attack_type"]
	// 		resultTemp["attack_type"] = attackType
	// 		eventTime, _ := result[index]["event_time"]
	// 		resultTemp["event_time"] = eventTime

	// 		log.Printf("%v, %v", resultTemp, upsertId)

	// 		log.Printf("%v, %v", result[index], upsertId)

	// 		localtionStr := ips.Get(attackSource.(string))

	// 		localtioninfo := ips.ParseLocaltion(localtionStr)

	// 		result[index]["new_location_zh_cn"] = localtioninfo.Country
	// 		result[index]["new_location_en"] = localtioninfo.CountryEN
	// 		result[index]["new_latitude"] = localtioninfo.Latitude
	// 		result[index]["new_longitude"] = localtioninfo.Longitude

	// 		res, err := client.Update().
	// 			Index("corerasp-attack-alarm-97d54c6615b69342ac525896f02d491ee852157e").
	// 			Type("attack-alarm").
	// 			Id(upsertId.(string)).
	// 			Doc(result[index]).Do(ctx)

	// 		if err != nil {
	// 			fmt.Printf("Update err:%v\n", err)
	// 		}

	// 		fmt.Printf("Update localtion:%v\n", res.Result)

	// 		// log.Printf("%v", resultTemp)
	// 		// 		if typeIndex == "attack" {
	// 		// requestId := result[index]["request_id"].(string)
	// 		// stackMd5 := result[index]["stack_md5"].(string)
	// 		// attackType := result[index]["attack_type"].(string)
	// 		// 			pluginAlgorithm := result[index]["plugin_algorithm"].(string)
	// 		// 			urlString := result[index]["url"].(string)
	// 		// 			if pluginAlgorithm == "response_dataLeak" {
	// 		// 				urlParse, err := url.Parse(urlString)
	// 		// 				if err != nil {
	// 		// 					return 0, nil, err
	// 		// 				}
	// 		// 				filterId = urlParse.Scheme + "://" + urlParse.Host + urlParse.Path
	// 		// 			} else {
	// 		// 				filterId = requestId + stackMd5 + attackType
	// 		// 			}
	// 		// 			result[index]["filter_id"] = filterId
	// 		// 		}
	// 		// 		es.HandleSearchResult(result[index], item.Id)
	// 	}
	// }

}

func TestGet(t *testing.T) {
	fmt.Println("Test Get IP ...")
	p, _ := ipsearch.New()
	// ip := "59.53.213.120"
	ip := "127.0.0.1"
	ipstr := p.Get(ip)
	fmt.Println(ipstr)
	if ipstr != `亚洲|中国|湖北| |潜江|联通|429005|China|CN|112.896866|30.421215` {
		t.Fatal("the IP convert by ipSearch component is not correct!")
	}
}
