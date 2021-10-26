package tools

import (
	"encoding/json"
	"rasp-cloud/config"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego"
)

var kafkaconfig *sarama.Config

func init() {
	kafkaconfig = sarama.NewConfig()
	kafkaconfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaconfig.Producer.Return.Successes = true
}

func SendMessage(key string, val map[string]interface{}) error {
	val["serverId"] = config.TOMLConfig.AppId
	kafka, ok := config.TOMLConfig.GetKafkaServer("online")
	if ok && kafka.Enable && kafka.Topic != "" && kafka.Addr != "" {
		// addrs := strings.Split(kafka.Addr, ",")
		if kafka.User != "" || kafka.Password != "" {
			kafkaconfig.Net.SASL.Enable = true
			kafkaconfig.Net.SASL.User = kafka.User
			kafkaconfig.Net.SASL.Password = kafka.Password
		} else {
			kafkaconfig.Net.SASL.Enable = false
		}
		addrIp := GetDnsIP(kafka.Addr)
		addrIps := strings.Split(addrIp, ",")
		producer, err := sarama.NewSyncProducer(addrIps, kafkaconfig)
		if err != nil {
			beego.Error(err)
			return err
		}
		defer producer.Close()

		var content []byte
		content, err = json.Marshal(val)
		if err != nil {
			return err
		}
		sContent := string(content)
		msg := &sarama.ProducerMessage{
			Partition: int32(1),
			Key:       sarama.StringEncoder(key),
			Value:     sarama.ByteEncoder(sContent),
			Topic:     kafka.Topic,
		}
		_, _, err = producer.SendMessage(msg)
		if err != nil {
			beego.Error("Send message Fail")
			return err
		}
	}
	return nil
}

func SendMessages(appId string, key string, valMaps []interface{}) error {
	var msgs []*sarama.ProducerMessage
	kafka, ok := config.TOMLConfig.GetKafkaServer("online")
	if ok && kafka.Enable && kafka.Topic != "" && kafka.Addr != "" {
		// addr := strings.Split(kafka.KafkaAddr, ",")
		if kafka.User != "" || kafka.Password != "" {
			kafkaconfig.Net.SASL.Enable = true
			kafkaconfig.Net.SASL.User = kafka.User
			kafkaconfig.Net.SASL.Password = kafka.Password
		}
		addrIp := GetDnsIP(kafka.Addr) + ":" + kafka.Port
		addrIps := strings.Split(addrIp, ",")
		producer, err := sarama.NewSyncProducer(addrIps, kafkaconfig)
		if err != nil {
			beego.Error(err)
			return err
		}
		defer producer.Close()

		var content []byte
		for _, val := range valMaps {
			content, err = json.Marshal(val)
			if err != nil {
				return err
			}
			sContent := string(content)
			msg := &sarama.ProducerMessage{
				Partition: int32(1),
				Key:       sarama.StringEncoder(key),
				Value:     sarama.ByteEncoder(sContent),
				Topic:     kafka.Topic,
			}
			msgs = append(msgs, msg)
		}
		err = producer.SendMessages(msgs)

		if err != nil {
			beego.Error("Send message Fail")
			return err
		}
	}
	return nil
}

// func GetKafkaConfig(name string) (kafka *config.KafkaServer, err error) {
// 	config.TOMLConfig
// 	err = mongo.FindId(kafkaAddrCollectionName, appId, &kafka)
// 	if err != nil {
// 		kafka = &Kafka{
// 			KafkaAddr:   conf.AppConfig.KafkaAddr,
// 			KafkaUser:   conf.AppConfig.KafkaUser,
// 			KafkaPwd:    conf.AppConfig.KafkaPwd,
// 			KafkaTopic:  conf.AppConfig.KafkaTopic,
// 			KafkaEnable: conf.AppConfig.KafkaEnable,
// 		}
// 	}
// 	return kafka, err
// }

// func PutKafkaConfig(appId string, kafka *Kafka) error {
// 	err := mongo.UpsertId(kafkaAddrCollectionName, appId, &kafka)
// 	return err
// }
