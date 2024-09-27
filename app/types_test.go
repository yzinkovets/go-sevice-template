package app

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNetworkTopology(t *testing.T) {
	// load Telemetry from telemetryFile (minified version)
	telemetryFile := "../data/telemetry_data_re.json"
	telemetryData, err := os.ReadFile(telemetryFile)
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	var tel Telemetry
	if err := json.Unmarshal(telemetryData, &tel); err != nil {
		t.Fatal("Error unmarshalling JSON:", err)
	}

	nt := tel.NetworkTopology

	// check the number of Re elements
	if len(nt.Re) != 2 {
		t.Fatal("Expected 2 Re elements, got ", len(nt.Re))
	}

	if len(nt.Re[0].Stats) != 2 {
		t.Fatal("Expected 2 Stats elements, got ", len(nt.Re[0].Stats))
	}

	// marshal the struct back to JSON
	data, err := json.Marshal(nt)
	if err != nil {
		t.Fatal("Error marshalling JSON:", err)
	}

	networkTopologyFile := "../data/network_topology_min.json"
	originalNetworkTopology, err := os.ReadFile(networkTopologyFile)
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	if string(data) != string(originalNetworkTopology) {
		t.Errorf("JSONs are not equal")
		t.Log("Original:", string(originalNetworkTopology))
		t.Log("Got     :", string(data))
	}
}
