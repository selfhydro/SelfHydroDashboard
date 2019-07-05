package main

import (
	mqtt "github.com/selfhydro/selfhydro/mqtt"
	"github.com/selfhydro/selfhydro/sensors"
)

type SelfHydroState struct {
	WaterTemperatureSensor       sensors.MQTTTopic
	AmbientTemperatureSensor     sensors.MQTTTopic
	AmbientHumiditySensor        sensors.MQTTTopic
	ElectricalConductivitySensor sensors.MQTTTopic
}

type state struct {
	WaterTemperature       float64
	AmbientTemperature     float64
	AmbientHumidity        float64
	ElectricalConductivity float64
	Time                   string
	DeviceId               string
}

func NewSelfHydroState(mqtt mqtt.MQTTComms) *SelfHydroState {
	waterTemp := SetupSensor(&sensors.WaterTemperature{}, mqtt)
	ambientTemperature := SetupSensor(&sensors.AmbientTemperature{}, mqtt)
	ambientHumidity := SetupSensor(&sensors.AmbientHumidity{}, mqtt)
	electricalConductivity := SetupSensor(&sensors.WaterElectricalConductivity{}, mqtt)
	selfHydroState := &SelfHydroState{
		WaterTemperatureSensor:       waterTemp,
		AmbientTemperatureSensor:     ambientTemperature,
		AmbientHumiditySensor:        ambientHumidity,
		ElectricalConductivitySensor: electricalConductivity,
	}
	return selfHydroState
}

func SetupSensor(sensor sensors.MQTTTopic, mqtt mqtt.MQTTComms) sensors.MQTTTopic {
	sensor.Subscribe(mqtt)
	return sensor
}

func (s SelfHydroState) GetLatestDeviceState() state {
	return state{
		WaterTemperature:       s.WaterTemperatureSensor.GetLatestData(),
		AmbientTemperature:     s.AmbientTemperatureSensor.GetLatestData(),
		AmbientHumidity:        s.AmbientHumiditySensor.GetLatestData(),
		ElectricalConductivity: s.ElectricalConductivitySensor.GetLatestData(),
		DeviceId:               "selfhydro",
	}
}
