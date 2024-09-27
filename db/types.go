package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/guregu/null/v5"
)

/*
CREATE TABLE iot.gateway_updates_wifi_mesh_dec (

	gw_uuid uuid NOT NULL,
	"timestamp" timestamp without time zone NOT NULL,
	mac macaddr NOT NULL,
	rate integer,
	rssi integer,
	channel integer,
	bandwidth text,
	stats jsonb,
	ssid text,
	master_channels text,
	parent macaddr,
	ip inet,
	client_conn_type iot.client_conn_type,
	host text,
	mode text,
	model text,
	snr smallint,
	qoe_score smallint,
	is_active boolean,
	pkt_fail_rate_rx smallint,
	pkt_fail_rate_tx smallint,
	ch_util smallint,
	eff_phy_rx smallint,
	eff_phy_tx smallint,
	data_consume_rx bigint,
	data_consume_tx bigint,
	pred_thrput_udp_full bigint,
	pred_thrput_udp_available bigint,
	pred_thrput_tcp_full bigint,
	pred_thrput_tcp_available bigint,
	rate_tx integer

);
*/

// {"rx": "796446707,2663363,0,0", "tx": "310331047,3940247,104,0", "iface": "ath0"}
type Stat struct {
	Rx    string `json:"rx"`
	Tx    string `json:"tx"`
	Iface string `json:"iface"`
}

type GatewayUpdateWifiMeshDec struct {
	GwUuid                 string          `db:"gw_uuid"`
	Timestamp              time.Time       `db:"timestamp"` // "2024-06-05 16:59:47"
	Mac                    string          `db:"mac"`
	IP                     null.String     `db:"ip"`
	ClientConnType         null.String     `db:"client_conn_type"`
	Host                   null.String     `db:"host"`
	Rate                   null.Int        `db:"rate"`
	RateTx                 null.Int        `db:"rate_tx"`
	Rssi                   null.Int        `db:"rssi"`
	Snr                    null.Int16      `db:"snr"`
	IsActive               null.Bool       `db:"is_active"`
	DataConsumeRx          null.Int64      `db:"data_consume_rx"`
	DataConsumeTx          null.Int64      `db:"data_consume_tx"`
	PredThrputTcpFull      null.Int64      `db:"pred_thrput_tcp_full"`
	PredThrputTcpAvailable null.Int64      `db:"pred_thrput_tcp_available"`
	PredThrputUdpFull      null.Int64      `db:"pred_thrput_udp_full"`
	PredThrputUdpAvailable null.Int64      `db:"pred_thrput_udp_available"`
	PktFailRateRx          null.Int16      `db:"pkt_fail_rate_rx"`
	PktFailRateTx          null.Int16      `db:"pkt_fail_rate_tx"`
	ChUtil                 null.Int16      `db:"ch_util"`
	EffPhyRx               null.Int16      `db:"eff_phy_rx"`
	EffPhyTx               null.Int16      `db:"eff_phy_tx"`
	Channel                null.Int        `db:"channel"`
	Bandwidth              null.String     `db:"bandwidth"`
	Stats                  json.RawMessage `db:"stats"` // []Stat
	Ssid                   null.String     `db:"ssid"`
	MasterChannels         null.String     `db:"master_channels"`
	Parent                 null.String     `db:"parent"`
	Mode                   null.String     `db:"mode"`
	Model                  null.String     `db:"model"`
	// QoeScore               null.Int16  `db:"qoe_score"`
}

type InfoNetworkDevicesIns struct {
	GwUuid      string
	Mac         string
	Host        string
	Fingerprint string
}

func (v InfoNetworkDevicesIns) Value() (driver.Value, error) {
	s := fmt.Sprintf("(%q,%q,%q,%q)",
		v.GwUuid,
		v.Mac,
		v.Host,
		v.Fingerprint,
	)
	return []byte(s), nil
}

type InfoNetworkDevicesFcdIns struct {
	GwUuid     string
	Mac        string
	LastOnline time.Time
}

func (v InfoNetworkDevicesFcdIns) Value() (driver.Value, error) {
	s := fmt.Sprintf("(%q,%q,%q)",
		v.GwUuid,
		v.Mac,
		v.LastOnline.Format(time.RFC3339),
	)
	return []byte(s), nil
}
