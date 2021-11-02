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

type AppCount struct {
	AppName string `json:"appname"`
	Count   int64  `json:"Count"`
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

func dealAppCountMap(appCountMap map[string]int64) AppCounts {

	var appCountList AppCounts
	appCountList = make([]AppCount, 0)
	for appName, count := range appCountMap {
		var appCount AppCount
		appCount.AppName = appName
		appCount.Count = count
		appCountList = append(appCountList, appCount)
	}
	sort.Sort(appCountList)
	return appCountList
}

// @router /get [post]
func (o *DisplayController) GetDisplay() {

	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	var AllAttackCount int64

	appCountMap := make(map[string]int64)

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

			appCountMap[appinfo.Name] = appCount
			AllAttackCount += appCount

			// 获取近期攻击流量趋势
			bucketDateList, err2 := logs.SearchAttackAggrByIndexAndDate("corerasp-attack-alarm-"+appinfo.Id, currentDate)
			if err2 == nil {
				for _, bucket := range bucketDateList {

					keyString := *bucket.KeyAsString

					// 攻击类型排名赋值
					attackDateCount, ok := attackDateCoutMap[keyString]
					if ok {
						attackDateCoutMap[keyString] = attackDateCount + bucket.DocCount
					} else {
						attackDateCoutMap[keyString] = bucket.DocCount
					}
				}
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

	// 被攻击应用排名
	result["app_attack_count"] = dealAppCountMap(appCountMap)

	// 防御攻击次数累计
	result["all_attack_count"] = AllAttackCount

	// 当前数据时间戳
	result["current_time"] = currentTime.Unix()

	o.Serve(result)
}
