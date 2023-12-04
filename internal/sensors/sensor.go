package sensors

import (
	"math"
	"math/rand"
	"time"
)

type sensorType struct {
	name            string
	minValue        float64
	maxValue        float64
	maxStep         float64
	maxTargetSpace  float64
	nightValueDelta float64
	unit            string
}

var WindSensor sensorType = sensorType{"Wind", 0.0, 90, 4, 10, 0, "kt"}
var HeatSensor sensorType = sensorType{"Heat", -10.0, 25, 0.6, 8, 0, "°C"}
var PressureSensor sensorType = sensorType{"Pressure", 900.0, 1030, 1.5, 50, 0, "hPa"}

type Sensor struct {
	id                 int
	airportId          string
	sensorType         *sensorType
	lastGeneratedValue float64
	currentTargetValue float64
	currentMaxStep     float64
}

func NewSensor(id int, airportId string, sensorType *sensorType) *Sensor {
	s := new(Sensor)
	s.id = id
	s.airportId = airportId
	s.sensorType = sensorType
	s.currentTargetValue = s.sensorType.minValue + ((s.sensorType.maxValue - s.sensorType.minValue) / 2)
	s.lastGeneratedValue, _ = s.randomSource()
	s.currentMaxStep = s.sensorType.maxStep
	return s
}

func (s *Sensor) randomSource() (float64, float64) {
	rand.Seed(time.Now().UnixMicro())
	r := rand.Float64()
	return (s.sensorType.minValue + r*(s.sensorType.maxValue+math.Abs(s.sensorType.minValue))), r
}

func (s *Sensor) defineNewTarget() {

	_, rand := s.randomSource()
	rand = rand*2.0 - 1.0 // Pour avoir des valeurs négative
	temp := s.currentTargetValue
	for temp = s.currentTargetValue + s.sensorType.maxTargetSpace*rand; temp < s.sensorType.minValue || temp > s.sensorType.maxValue; {
		_, rand := s.randomSource()
		rand = rand*2.0 - 1.0
		temp = s.currentTargetValue + s.sensorType.maxTargetSpace*rand
	}
	s.currentTargetValue = temp
}

func (s *Sensor) GenerateNextData() float64 {
	if math.Abs(s.lastGeneratedValue-s.currentTargetValue) < s.sensorType.maxStep {
		s.defineNewTarget()
	}

	if math.Abs(s.lastGeneratedValue-s.currentTargetValue) < s.sensorType.maxStep*5 {
		s.currentMaxStep = s.sensorType.maxStep * (math.Abs(s.lastGeneratedValue-s.currentTargetValue) / (s.sensorType.maxStep * 5))
	} else {
		if s.currentMaxStep < s.sensorType.maxStep {
			s.currentMaxStep = s.currentMaxStep + 1/(math.Abs(s.lastGeneratedValue-s.currentTargetValue)/(s.sensorType.maxStep*5))
		} else {
			s.currentMaxStep = s.sensorType.maxStep
		}
	}
	_, rand := s.randomSource()

	nextData := s.currentMaxStep * rand
	_, variation := s.randomSource()
	if s.lastGeneratedValue > s.currentTargetValue {
		nextData *= -1
		if variation < 0.25 {
			nextData *= -1
		}
		s.lastGeneratedValue += nextData
		return s.lastGeneratedValue
	} else {
		if variation < 0.25 {
			nextData *= -1
		}
		s.lastGeneratedValue += nextData
		return s.lastGeneratedValue
	}
}
