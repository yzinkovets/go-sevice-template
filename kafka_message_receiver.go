package main

import (
	"fmt"
	"go-service-template/config"
	"go-service-template/utils"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	// time formatted as "YYYY-MM-DD HH:MM:SS.ssssss" (UTC)
	// where ssssss is the microsecond part
	// Example: 2024-05-01 13:38:42.123456
	// This format could be improved by avoiding the use of dashes, colons and spaces but it is less readable
	DateTimeFormat = "2006-01-02 15:04:05.999999"
)

type KafkaMessageReceiver struct {
	kafkaConsumer *kafka.Consumer
	run           bool
}

func NewKafkaMessageReceiver(cfg config.KafkaConfig) (*KafkaMessageReceiver, error) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
		"group.id":          cfg.ClientId,
		"auto.offset.reset": cfg.AutoOffsetReset,
	})

	if err != nil {
		return nil, err
	}

	topics := strings.Split(cfg.Topic, ",")

	if err := kafkaConsumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return &KafkaMessageReceiver{
		kafkaConsumer: kafkaConsumer,
	}, nil
}

func (k *KafkaMessageReceiver) Run(chMessages chan<- *Message) {
	if k.run {
		logrus.Error("already running")
		return
	}

	k.run = true

	for k.run {
		kafkaMessage, err := k.kafkaConsumer.ReadMessage(time.Second)
		if err != nil {
			if err.(kafka.Error).IsTimeout() {
				continue
			}
			// The client will automatically try to recover from all errors.
			logrus.Errorf("Consumer error: %v (%v)", err, kafkaMessage)
			continue
		}

		message, err := k.getMessageFromKafkaMessage(kafkaMessage)
		if err != nil {
			logrus.Errorf("Error parsing message: %v", err)
			continue
		}

		// we don't set gw_uuid here because it IS in the message json
		chMessages <- message
	}
}

func (k *KafkaMessageReceiver) Stop() {
	k.run = false
}

// Message string representation is in the following format:
// <topic>;<qos>;<datetime>;<payload>
// datetime is in "YYYY-MM-DD HH:MM:SS.ssssss" format
// Example: /router/e3e259fe44be4c8d932707e771c6de0c/telemetry;1;2024-05-01 13:38:42.123456;{"temperature": 25.5, "humidity": 50}
func (k *KafkaMessageReceiver) getMessageFromKafkaMessage(kafkaMessage *kafka.Message) (*Message, error) {
	value := string(kafkaMessage.Value)

	parts := strings.Split(value, ";")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid message format: %s", value[:min(len(value), 100)])
	}

	t, err := time.Parse(DateTimeFormat, parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %s", parts[2])
	}

	// get int from string
	qos, err := strconv.Atoi(parts[1])
	if err != nil || qos < 0 || qos > 2 {
		return nil, fmt.Errorf("invalid qos: %s", parts[1])
	}

	message := Message{
		GwUuid:    utils.CutGw(GetUsernameFromTopic(parts[0])),
		Timestamp: t,
		Data:      []byte(parts[3]),
	}

	return &message, nil
}

// Returns username from the topic
// Topic is in /(router|relay)/<username>/... format
func GetUsernameFromTopic(topic string) string {
	idxUsername := 0 // start of username
	s := 0           // number of slashes ('/') in a topic

	for i := 0; i < len(topic); i++ {
		if topic[i] == '/' {
			s++
			if s == 2 {
				idxUsername = i
			} else if s == 3 {
				return topic[idxUsername+1 : i]
			}
		}
	}
	return ""
}
