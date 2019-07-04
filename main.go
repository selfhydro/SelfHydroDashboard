package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/selfhydro/selfhydro/mqtt"
	"github.com/selfhydro/selfhydro/sensors"
)

const (
	SefhydroStateDB = "selfhydro"
)

type SeflhydroState struct {
	WaterTemperature float64 `json:"waterTemperature"`
	Time             string  `json:"time"`
	deviceId         string
}

var influxDBAddress = "http://influxdb:8086"

func main() {
	InfluxDB, err := NewInfluxDBRepository(influxDBAddress)
	CheckError(err)
	localMQTT := mqtt.NewLocalMQTT("influxdb-ingester", "mosquitto:1883")
	localMQTT.ConnectDevice()

	waterTemp := SetupSensor(&sensors.WaterTemperature{}, localMQTT)
	ambientTemperature := SetupSensor(&sensors.AmbientTemperature{}, localMQTT)
	ambientHumidity := SetupSensor(&sensors.AmbientHumidity{}, localMQTT)
	electricalConductivity := SetupSensor(&sensors.WaterElectricalConductivity{}, localMQTT)
	time.Sleep(time.Second * 10)

	for {
		waterTemperature := waterTemp.GetLatestData()
		ambientTemperature := ambientTemperature.GetLatestData()
		currentAmbientHumidity := ambientHumidity.GetLatestData()
		currentElectricalConductivity := electricalConductivity.GetLatestData()
		fmt.Println("water temp: ", waterTemperature)
		fmt.Println("ec level: ", currentElectricalConductivity)

		err := InfluxDB.WriteElectricalConductivityState(currentElectricalConductivity, "selfhydro")
		CheckError(err)
		err = InfluxDB.WriteTemperatureState(waterTemperature, ambientTemperature, "selfhydro")
		CheckError(err)
		err = InfluxDB.WriteHumidityState(currentAmbientHumidity, "selfhydro")
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

func SetupSensor(sensor sensors.MQTTTopic, mqtt mqtt.MQTTComms) sensors.MQTTTopic {
	sensor.Subscribe(mqtt)
	return sensor
}
