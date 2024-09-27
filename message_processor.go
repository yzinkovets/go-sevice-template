package main

import (
	"encoding/json"
	"go-service-template/app"
	"go-service-template/db"
	"go-service-template/utils"
	"time"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Data      []byte
	GwUuid    string
	Timestamp time.Time
}

type MessageProcessor struct {
	db                   *db.DbConnection
	chMessages           <-chan *Message
	chRouterNotification chan<- *RouterNotification
}

func NewMessageProcessor(db *db.DbConnection, chMessages <-chan *Message, chRouterNotification chan<- *RouterNotification) *MessageProcessor {

	p := MessageProcessor{
		db:                   db,
		chMessages:           chMessages,
		chRouterNotification: chRouterNotification,
	}
	go p.Run()

	return &p
}

func (p *MessageProcessor) Run() {
	for m := range p.chMessages {
		p.ProcessMessage(m)
	}
}

func (p *MessageProcessor) ProcessMessage(message *Message) error {
	if message == nil {
		logrus.Error("Message is nil")
		return nil
	}

	if logrus.GetLevel() >= logrus.TraceLevel {
		logrus.Tracef("Incoming message: GwUuid:%s Data:%s\n", message.GwUuid, utils.CompactString(string(message.Data)))
	}

	var telemetry app.Telemetry
	if err := json.Unmarshal(message.Data, &telemetry); err != nil {
		return err
	}

	if telemetry.NetworkTopology.Message == "Device is in RE mode" || telemetry.NetworkTopology.Message == "QCA-MESH-API service is not supported. Please check hyd!" {
		return nil
	}

	if err := p.processNetworkTopology(telemetry.NetworkTopology, message.GwUuid, message.Timestamp); err != nil {
		return err
	}

	return nil
}

func (this *MessageProcessor) processNetworkTopology(t app.NetworkTopology, gwUuid string, timestamp time.Time) error {
	return nil
}
