package main

import (
	"encoding/json"
	"go-service-template/config"
	"go-service-template/utils"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// {"changed":"devices-monitor", "device_macs": ["aa:bb:cc:...", "bb:dd:00:..."]}
type RouterNotificationContent struct {
	Changed    string   `json:"changed"`
	DeviceMacs []string `json:"device_macs"`
}
type RouterNotification struct {
	GwUuid  string                    `json:"gw_uuid"`
	Content RouterNotificationContent `json:"content"`
}

type MqttMessageNotifier struct {
	client               mqtt.Client
	chRouterNotification <-chan *RouterNotification
	topic                string
}

func NewMqttMessageNotifier(cfg config.MqttConfig, client mqtt.Client, chRouterNotification <-chan *RouterNotification) (*MqttMessageNotifier, error) {

	m := &MqttMessageNotifier{
		client:               client,
		chRouterNotification: chRouterNotification,
		topic:                cfg.TopicSend,
	}

	go m.Run()

	return m, nil
}

func (m *MqttMessageNotifier) Run() {
	for notification := range m.chRouterNotification {
		logrus.Debugf("Sending notification: %v", notification)

		topic := strings.Replace(m.topic, "{gw_uuid}", utils.SetGw(notification.GwUuid), 1)
		logrus.Debugf("Send to topic: [%s]", topic)

		payload, err := json.Marshal(notification.Content)
		if err != nil {
			logrus.Error("Can't marshal notification. Error:", err)
			continue
		}
		logrus.Debugf("Payload: %s", payload)

		token := m.client.Publish(topic, 0, false, payload)
		if token.Wait() && token.Error() != nil {
			logrus.Error("Can't notify router. Error:", token.Error())
		}
	}
}
