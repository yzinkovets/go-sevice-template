package main

import (
	"go-service-template/config"
	"go-service-template/utils"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type MqttMessageReceiver struct {
	client     mqtt.Client
	chMessages chan<- *Message
	topic      string
	qos        byte
}

func NewMqttMessageReceiver(cfg config.MqttConfig, client mqtt.Client, chMessages chan<- *Message) (*MqttMessageReceiver, error) {
	m := &MqttMessageReceiver{
		client:     client,
		chMessages: chMessages,
		topic:      cfg.TopicReceive,
		qos:        cfg.Qos,
	}

	token := client.Subscribe(cfg.TopicReceive, cfg.Qos, m.messageHandler)
	token.Wait()
	logrus.Infof("Subscribed to topic %s with QoS %d", cfg.TopicReceive, cfg.Qos)

	return m, nil
}

func (r *MqttMessageReceiver) messageHandler(client mqtt.Client, msg mqtt.Message) {
	logrus.Tracef("Received MQTT message from %v with payload length: %v", msg.Topic(), len(msg.Payload()))

	m := &Message{
		GwUuid:    utils.CutGw(GetUsernameFromTopic(msg.Topic())),
		Data:      msg.Payload(),
		Timestamp: time.Now(),
	}

	r.chMessages <- m

	// Acking is automatic (managed by mqtt.ClientOptions.AutoAckDisabled = false)
	// msg.Ack()
}
