package app

type Telemetry struct {
	NetworkTopology NetworkTopology `json:"network_topology"`
}

type NetworkTopology struct {
	Re        []Re   `json:"re"`
	Cap       Cap    `json:"cap"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Message   string `json:"msg,omitempty"`
}

type Re struct {
	IP       string   `json:"ip,omitempty"`
	Mac      string   `json:"mac,omitempty"`
	Rate     string   `json:"rate,omitempty"`
	Rssi     string   `json:"rssi,omitempty"`
	Ssid     string   `json:"ssid,omitempty"`
	Model    string   `json:"model,omitempty"`
	Stats    []Stats  `json:"stats,omitempty"`
	Parent   string   `json:"parent,omitempty"`
	Channel  string   `json:"channel,omitempty"`
	Clients  Clients  `json:"clients,omitempty"`
	Backhaul []Client `json:"backhaul,omitempty"`
}

type Cap struct {
	IP      string  `json:"ip,omitempty"`
	Mac     string  `json:"mac,omitempty"`
	Ssid    string  `json:"ssid,omitempty"`
	Model   string  `json:"model,omitempty"`
	Stats   []Stats `json:"stats,omitempty"`
	Parent  string  `json:"parent,omitempty"`
	ChUtil  string  `json:"ch_util,omitempty"`
	Channel string  `json:"channel,omitempty"`
	Clients Clients `json:"clients,omitempty"`
}

type Stats struct {
	Rx    string `json:"rx,omitempty"`
	Tx    string `json:"tx,omitempty"`
	Iface string `json:"iface,omitempty"`
}

type Clients struct {
	Wifi     []Client `json:"wifi"`
	Ethernet []Client `json:"ethernet"`
}

type Client struct {
	IP                     string `json:"ip,omitempty"`
	Mac                    string `json:"mac,omitempty"`
	Host                   string `json:"host,omitempty"`
	Rate                   string `json:"rate,omitempty"`
	Rssi                   string `json:"rssi,omitempty"`
	Channel                string `json:"channel,omitempty"`
	Bandwidth              string `json:"bandwidth,omitempty"`
	DhcpFingerprint        string `json:"dhcp_fingerprint,omitempty"`
	Snr                    string `json:"snr,omitempty"`
	Mode                   string `json:"mode,omitempty"`
	TxRate                 string `json:"tx_rate,omitempty"`
	IsActive               string `json:"is_active,omitempty"`
	EffPhyRx               string `json:"eff_phy_rx,omitempty"`
	EffPhyTx               string `json:"eff_phy_tx,omitempty"`
	DataConsumeRx          string `json:"data_consume_rx,omitempty"`
	DataConsumeTx          string `json:"data_consume_tx,omitempty"`
	ChUtil                 string `json:"ch_util,omitempty"`
	PktFailRateRx          string `json:"pkt_fail_rate_rx,omitempty"`
	PktFailRateTx          string `json:"pkt_fail_rate_tx,omitempty"`
	LastConnectTime        string `json:"last_connect_time,omitempty"`
	PredThrputTcpFull      string `json:"pred_thrput_tcp_full,omitempty"`
	PredThrputUdpFull      string `json:"pred_thrput_udp_full,omitempty"`
	PredThrputTcpAvailable string `json:"pred_thrput_tcp_available,omitempty"`
	PredThrputUdpAvailable string `json:"pred_thrput_udp_available,omitempty"`
	Model                  string `json:"model,omitempty"`
}
