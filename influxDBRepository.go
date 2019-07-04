package main

import (
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxDBRespository struct {
	client  client.Client
	address string
}

func NewInfluxDBRepository(address string) (InfluxDBRespository, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: address,
	})
	if err != nil {
		log.Fatal(err)
		return InfluxDBRespository{}, err
	}
	return InfluxDBRespository{
		address: address,
		client:  client,
	}, nil
}

func (i *InfluxDBRespository) Close() error {
	if err := i.client.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *InfluxDBRespository) WriteTemperatureState(waterTemperature, ambientTemperature float64, device string) error {
	tags := map[string]string{"device": device}

	temperatureFields := map[string]interface{}{
		"waterTemperature":   waterTemperature,
		"ambientTemperature": ambientTemperature,
	}
	tempPoint, err := client.NewPoint("temperature", tags, temperatureFields, time.Now())
	if err != nil {
		return err
	}
	batchPoints, err := i.CreateBatchPoint(tempPoint)
	if err != nil {
		log.Println("couldnt create batch point for electrical conductivity state")
		return err
	}
	err = i.WriteBatchPoint(batchPoints)
	i.Close()
	return nil
}

func (i *InfluxDBRespository) WriteHumidityState(ambientHumidity float64, device string) error {
	tags := map[string]string{"device": device}

	humidityFields := map[string]interface{}{
		"ambientHumidity": ambientHumidity,
	}
	humidityPoint, err := client.NewPoint("humidity", tags, humidityFields, time.Now())
	if err != nil {
		log.Print("error creating humidity point")
		log.Print(err)
		return err
	}
	batchPoints, err := i.CreateBatchPoint(humidityPoint)
	if err != nil {
		log.Println("couldnt create batch point for electrical conductivity state")
		return err
	}
	err = i.WriteBatchPoint(batchPoints)
	i.Close()

	return err
}

func (i *InfluxDBRespository) WriteElectricalConductivityState(electricalConductivity float64, device string) error {
	tags := map[string]string{"device": device}

	electricalConductivityFields := map[string]interface{}{
		"electricalConductivity": electricalConductivity,
	}
	electricalConductivityPt, err := client.NewPoint("electricalConductivity", tags, electricalConductivityFields, time.Now())
	if err != nil {
		log.Print(err)
		return err
	}

	batchPoints, err := i.CreateBatchPoint(electricalConductivityPt)
	if err != nil {
		log.Println("couldnt create batch point for electrical conductivity state")
		return err
	}
	err = i.WriteBatchPoint(batchPoints)
	i.Close()

	return err
}

func (i InfluxDBRespository) CreateBatchPoint(point *client.Point) (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  SefhydroStateDB,
		Precision: "s",
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	bp.AddPoint(point)

	return bp, nil
}

func (i InfluxDBRespository) WriteBatchPoint(batchPoints client.BatchPoints) error {
	if err := i.client.Write(batchPoints); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
