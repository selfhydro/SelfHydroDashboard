package main

import (
	"log"
	"time"

	mqtt "github.com/selfhydro/selfhydro/mqtt"
)

const (
	SefhydroStateDB = "selfhydro"
)

var influxDBAddress = "http://influxdb:8086"

func main() {
	InfluxDB, err := NewInfluxDBRepository(influxDBAddress)
	CheckError(err)
	localMQTT := mqtt.NewLocalMQTT("influxdb-ingester", "mosquitto:1883")
	localMQTT.ConnectDevice()
	selfhydroState := NewSelfHydroState(localMQTT)
	time.Sleep(time.Second * 10)

	for {
		state := selfhydroState.GetLatestDeviceState()
		err := InfluxDB.WriteElectricalConductivityState(state.ElectricalConductivity, "selfhydro")
		CheckError(err)
		err = InfluxDB.WriteTemperatureState(state.WaterTemperature, state.AmbientTemperature, "selfhydro")
		CheckError(err)
		err = InfluxDB.WriteHumidityState(state.AmbientHumidity, "selfhydro")
		CheckError(err)
		time.Sleep(time.Minute * 15)
	}
}

func CheckError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}
