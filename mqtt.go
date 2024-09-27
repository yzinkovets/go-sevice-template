package main

import (
	"go-service-template/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func GetMqttClient(cfg config.MqttConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Host)
	opts.SetClientID(cfg.ClientId)
	opts.SetUsername(cfg.Username)
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetKeepAlive(time.Duration(cfg.KeepAliveTimeoutSec) * time.Second)
	opts.SetPingTimeout(time.Duration(cfg.PingTimeoutSec) * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}
