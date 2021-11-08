package api

import (
	"fmt"
	"log"
	"rasp-cloud/config"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
	"rasp-cloud/models/logs"
	"sort"
	"strings"
	"time"
)

// Operations about display
type DisplayController struct {
	controllers.BaseController
}

// 攻击类型
type AttackTypeInfo struct {
	Count   int64   `json:"count"`
	Percent float32 `json:"percent"`
}

func getUseDate() int {
	t1 := time.Now()
	t2 := time.Unix(models.ServerstartTime, 0)
	left := t1.Sub(t2)
	return int(left.Hours()/24) + 1
}

func getOrderListBySize(list logs.AttackInfoLists, size int) logs.AttackInfoLists {

	sort.Sort(list)
	return list
}

func dealAttackTypeCountMap(attackTypeMap map[string]int64, allCount int64) map[string]AttackTypeInfo {

	returnMap := make(map[string]AttackTypeInfo)
	for keyName, attackTypeCount := range attackTypeMap {

		var attackTypeInfo AttackTypeInfo
		attackTypeInfo.Count = attackTypeCount
		attackTypeInfo.Percent = float32(attackTypeCount) / float32(allCount)
		returnMap[keyName] = attackTypeInfo

	}
	return returnMap
}

// 城市攻击排名
type AttackCountryInfo struct {
	CountyNameZH string `json:"country_name_zh_cn"`
	CountyNameEN string `json:"country_name_en"`
	Count        int64  `json:"count"`
}

type AttackCountryInfos []AttackCountryInfo

func (m AttackCountryInfos) Len() int {
	return len(m)
}

func (m AttackCountryInfos) Less(i, j int) bool {
	return m[i].Count > m[j].Count
}

func (m AttackCountryInfos) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type AppCount struct {
	AppName      string `json:"appname"`
	AppHeadImage string `json:"head_image_base64"`
	Count        int64  `json:"Count"`
}

type AppCounts []AppCount

func (m AppCounts) Len() int {
	return len(m)
}

func (m AppCounts) Less(i, j int) bool {
	return m[i].Count > m[j].Count
}

func (m AppCounts) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func dealAppCountMap(appCountMap map[string]AppCount) AppCounts {

	var appCountList AppCounts
	appCountList = make([]AppCount, 0)
	for appName, count := range appCountMap {
		var appCount AppCount
		appCount.AppName = appName
		appCount.Count = count.Count
		appCount.AppHeadImage = count.AppHeadImage
		appCountList = append(appCountList, appCount)
	}
	sort.Sort(appCountList)
	return appCountList
}

func dealAttackCountryMap(attackCountryMap map[string]AttackCountryInfo) AttackCountryInfos {

	var attackCountryInfos AttackCountryInfos
	attackCountryInfos = make([]AttackCountryInfo, 0)
	for _, info := range attackCountryMap {
		var attackCountryInfo AttackCountryInfo
		attackCountryInfo.Count = info.Count
		attackCountryInfo.CountyNameEN = info.CountyNameEN
		attackCountryInfo.CountyNameZH = info.CountyNameZH
		attackCountryInfos = append(attackCountryInfos, attackCountryInfo)
	}
	sort.Sort(attackCountryInfos)
	return attackCountryInfos
}

// @router /getconfig [get]
func (o *DisplayController) GetConfig() {
	var result = make(map[string]interface{})
	result["mobile"] = config.TOMLConfig.AppMobile
	result["email"] = config.TOMLConfig.AppEmail
	result["version"] = config.TOMLConfig.AppVersion
	o.Serve(result)

}

// @router /get [post]
func (o *DisplayController) GetDisplay() {

	currentTime := time.Now()
	currentDate := currentTime.AddDate(0, 0, -6).Format("2006-01-02")

	var AllAttackCount int64

	appCountMap := make(map[string]AppCount)
	attackCountryMap := make(map[string]AttackCountryInfo)

	// var attackCountryList AttackCountryInfos
	// attackCountryList = make([]AttackCountryInfo, 0)

	attackTypeCountMap := make(map[string]int64)

	attackLevelCountMap := make(map[string]int64)
	attackLevelCountMap["高危"] = 0
	attackLevelCountMap["中危"] = 0
	attackLevelCountMap["低危"] = 0

	attackDateCoutMap := make(map[string]int64)
	for i := 0; i < 7; i++ {
		dateStr := currentTime.AddDate(0, 0, 0-i).Format("2006-01-02")
		attackDateCoutMap[dateStr] = 0
	}

	attackListArray := make([]logs.AttackInfoList, 0)

	countApp, appList, err := models.GetAllApp(1, 1000, false)
	if err == nil {
		for _, appinfo := range appList {
			//获取攻击类型排序
			bucketList, err2 := logs.SearchAttackAggrByIndex("corerasp-attack-alarm-" + appinfo.Id)
			var appCount int64
			if err2 == nil {
				for _, bucket := range bucketList {
					// 应用排名赋值
					appCount += bucket.DocCount
					keyString := fmt.Sprintf("%q", bucket.Key)

					// 攻击类型排名赋值
					keyString = strings.ReplaceAll(keyString, "\"", "")
					attackTypeCount, ok := attackTypeCountMap[keyString]
					if ok {
						attackTypeCountMap[keyString] = attackTypeCount + bucket.DocCount
					} else {
						attackTypeCountMap[keyString] = bucket.DocCount
					}

					attackInfo, ok1 := config.TOMLConfig.AttackTypes[keyString]
					if ok1 {
						attackLevelCount, ok := attackLevelCountMap[attackInfo.Level]
						if ok {
							attackLevelCountMap[attackInfo.Level] = attackLevelCount + bucket.DocCount
						} else {
							attackLevelCountMap[attackInfo.Level] = bucket.DocCount
						}
					}
				}
			}

			var appInfo AppCount
			appInfo.Count = appCount
			appInfo.AppHeadImage = appinfo.HeadImageBase64
			appInfo.AppName = appinfo.Name
			appCountMap[appinfo.Name] = appInfo
			AllAttackCount += appCount

			// 获取城市排名
			bucketCountryList, err3 := logs.SearchAttackAggrByCountry("corerasp-attack-alarm-" + appinfo.Id)
			if err3 == nil {
				for _, bucket := range bucketCountryList {
					childagg, isfind2 := bucket.Aggregations.Terms("location_zh_cn")
					if isfind2 {
						for _, bucket1 := range childagg.Buckets {
							// var attackCountryInfo AttackCountryInfo
							// attackCountryInfo.Count = bucket1.DocCount
							countyNameEN := fmt.Sprintf("%q", bucket.Key)
							countyNameEN = strings.ReplaceAll(countyNameEN, "\"", "")
							countyNameCH := fmt.Sprintf("%q", bucket1.Key)
							countyNameCH = strings.ReplaceAll(countyNameCH, "\"", "")
							attackCountryInfo, ok := attackCountryMap[countyNameEN+":"+countyNameCH]
							if ok {
								attackCountryInfo.Count = attackCountryInfo.Count + bucket1.DocCount
								attackCountryMap[countyNameEN+":"+countyNameCH] = attackCountryInfo
							} else {
								var tempattackCountryInfo AttackCountryInfo
								tempattackCountryInfo.Count = bucket1.DocCount
								tempattackCountryInfo.CountyNameEN = countyNameEN
								tempattackCountryInfo.CountyNameZH = countyNameCH
								attackCountryMap[countyNameEN+":"+countyNameCH] = tempattackCountryInfo
							}
							// fmt.Printf("bucket = %q  bucket1 = %q 文档总数 = %d\n", bucket.Key, bucket1.Key, bucket1.DocCount)
						}
					}
				}
			}

			// 获取近期攻击流量趋势
			bucketDateList, err2 := logs.SearchAttackAggrByIndexAndDate("corerasp-attack-alarm-"+appinfo.Id, currentDate)
			if err2 == nil {
				for _, bucket := range bucketDateList {

					keyString := *bucket.KeyAsString

					fmt.Printf("bucket = %v 文档总数 = %d\n", keyString, bucket.DocCount)
					// 攻击类型排名赋值
					attackDateCount, ok := attackDateCoutMap[keyString]
					if ok {
						attackDateCoutMap[keyString] = attackDateCount + bucket.DocCount
					} else {
						attackDateCoutMap[keyString] = bucket.DocCount
					}
				}
			} else {
				fmt.Printf("err2 = %v\n", err2)
			}

			attackList, err3 := logs.SearchAttackList("corerasp-attack-alarm-"+appinfo.Id, 5)

			if err3 == nil {
				attackListArray = append(attackListArray, attackList...)
			} else {
				log.Printf("SearchAttackList error[%v]", err3)
			}
			log.Printf("attackList len[%d], name[%s]", len(attackList), appinfo.Name)
		}

	}

	var result = make(map[string]interface{})
	// 应用数
	result["app_count"] = countApp
	// 已保护时间
	result["server_start_day"] = getUseDate()
	// 已收录攻击类型
	result["attack_types"] = len(config.TOMLConfig.AttackTypes)

	// 最新攻击
	result["attack_list_array"] = getOrderListBySize(attackListArray, 5)

	// 攻击类型占比和数量
	result["attack_type_count"] = dealAttackTypeCountMap(attackTypeCountMap, AllAttackCount)

	// 近期攻击流量趋势
	result["attack_date_cout"] = attackDateCoutMap

	// 攻击等级排名
	result["attack_level_count"] = attackLevelCountMap

	// 攻击来源排名
	result["attack_country_count"] = dealAttackCountryMap(attackCountryMap)

	// 被攻击应用排名
	result["app_attack_count"] = dealAppCountMap(appCountMap)

	// 防御攻击次数累计
	result["all_attack_count"] = AllAttackCount

	// 当前数据时间戳
	result["current_time"] = currentTime.Unix()

	o.Serve(result)
}
