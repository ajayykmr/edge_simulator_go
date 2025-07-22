package data_generator

import (
	"time"

	utils "github.com/ajayykmr/edge_simulator_go/utils"
)

type CNCData struct {
	MachineID        string  `json:"machine_id"`
	SpindleRPM       int     `json:"spindle_rpm"`
	SpindleTemp      float64 `json:"spindle_temp"`
	MotorLoad        float64 `json:"motor_load"`        // %
	VibrationX       float64 `json:"vibration_x"`       // mm/s
	VibrationY       float64 `json:"vibration_y"`       // mm/s
	PowerConsumption float64 `json:"power_consumption"` // kW
	CoolantFlow      float64 `json:"coolant_flow"`      // L/min
	CycleTimeSec     int     `json:"cycle_time_sec"`    // sec
	ToolID           string  `json:"tool_id"`
	Timestamp        string  `json:"timestamp"`
}

func GenerateCNCData(machineID string) CNCData {
	return CNCData{
		MachineID:        machineID,
		SpindleRPM:       utils.RandInt(800, 12000),
		SpindleTemp:      utils.RandFloat(40.0, 75.0),
		MotorLoad:        utils.RandFloat(30.0, 80.0),
		VibrationX:       utils.RandFloat(0.01, 0.2),
		VibrationY:       utils.RandFloat(0.01, 0.2),
		PowerConsumption: utils.RandFloat(2.5, 8.0),
		CoolantFlow:      utils.RandFloat(1.5, 4.0),
		CycleTimeSec:     utils.RandInt(15, 300),
		ToolID:           "TOOL-" + string(rune(utils.RandInt(0, 5)+'A')),
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
	}
}
