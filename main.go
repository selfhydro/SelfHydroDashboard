package main

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
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

func main() {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://influxdb:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  SefhydroStateDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	localMQTT := mqtt.NewLocalMQTT("influxdb-ingester", "mosquitto:1883")
	localMQTT.ConnectDevice()

	waterTemp := &sensors.WaterTemperature{}
	waterTemp.Subscribe(localMQTT)
	ambientTemperature := &sensors.AmbientTemperature{}
	ambientTemperature.Subscribe(localMQTT)
	ambientHumidity := &sensors.AmbientHumidity{}
	ambientHumidity.Subscribe(localMQTT)
	electricalConductivity := &sensors.WaterEC{}
	electricalConductivity.Subscribe(localMQTT)
	time.Sleep(time.Second * 10)

	for {
		waterTemperature := waterTemp.GetLatestData()
		ambientTemp := ambientTemperature.GetLatestData()
		currentAmbientHumidity := ambientHumidity.GetLatestData()
		currentElectricalConductivity := electricalConductivity.GetLatestData()
		fmt.Println("water temp: ", waterTemperature)
		fmt.Println("ec level: ", currentElectricalConductivity)

		tags := map[string]string{"device": "selfhydro"}
		temperatureFields := map[string]interface{}{
			"waterTemperature":   waterTemperature,
			"ambientTemperature": ambientTemp,
		}
		humidityFields := map[string]interface{}{
			"ambientHumidity": currentAmbientHumidity,
		}
		electricalConductivityFields := map[string]interface{}{
			"electricalConductivity": currentElectricalConductivity,
		}

		tempPt, err := client.NewPoint("temperature", tags, temperatureFields, time.Now())
		humidityPt, err := client.NewPoint("humidity", tags, humidityFields, time.Now())
		electricalConductivityPt, err := client.NewPoint("electricalConductivity", tags, electricalConductivityFields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(tempPt)
		bp.AddPoint(humidityPt)
		bp.AddPoint(electricalConductivityPt)

		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Minute * 5)
	}
}
