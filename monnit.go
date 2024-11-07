package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

const DATE_FORMAT = "Mon, 02 Jan 2006 15:04:05"
const CACHE_FILE = "monnit.json"

type Monnit struct {
	sync.Mutex
	sensorId     string
	apiKeyId     string
	apiSecretKey string
	apiUrl       string
	lastData     *SensorDataMessages
}

func NewMonnit(sensorID, apiKeyID, apiSecretKey, url string, interval time.Duration) *Monnit {
	monnit := Monnit{
		sensorId:     sensorID,
		apiKeyId:     apiKeyID,
		apiSecretKey: apiSecretKey,
		apiUrl:       url,
	}

	// Load cached data
	if f, err := os.Open(CACHE_FILE); err == nil {
		if err = json.NewDecoder(f).Decode(&monnit.lastData); err != nil {
			slog.Error("unable to restore cached values", "error", err)
		}
		slog.Info("loaded cached Monnit data", "cache", CACHE_FILE)
		slog.Info("latest reading", "measurement", monnit.LastReading())
	} else {
		slog.Info("cached Monnit data not found", "cache", CACHE_FILE)
		if err = monnit.LoadData(); err != nil {
			slog.Warn("problem loading data on startup", "error", err)
		}
		slog.Info("latest reading", "measurement", monnit.LastReading())
	}

	go monnit.refresh(interval)

	return &monnit
}

// refresh automatically updates sensor data at the interval
func (m *Monnit) refresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := m.LoadData()
		if err != nil {
			slog.Error("failed to load data", "error", err)
		}
		slog.Debug("refreshed data", "interval", interval)
		slog.Info("latest reading", "measurement", m.LastReading())
	}
}

func (m *Monnit) LoadData() error {
	m.Lock()
	defer m.Unlock()

	sdm := SensorDataMessages{
		LastUpdated: time.Now(),
	}

	// Requesting from current time to seven days in the past
	toDate := time.Now().UTC()
	fromDate := toDate.AddDate(0, 0, -7).UTC()
	req, err := http.NewRequest("GET", m.apiUrl, nil)
	if err != nil {
		slog.Error("error creating request", "error", err)
		return err
	}

	// Set API keys in HTTP headers
	req.Header.Set("APIKeyID", m.apiKeyId)
	req.Header.Set("APISecretKey", m.apiSecretKey)

	// Pass sensor ID and date range as query params
	q := req.URL.Query()
	q.Add("sensorID", m.sensorId)
	q.Add("fromDate", fromDate.Format(DATE_FORMAT))
	q.Add("toDate", toDate.Format(DATE_FORMAT))
	req.URL.RawQuery = q.Encode()

	slog.Debug("loading data from Monnit API", "url", req.URL.String(),
		"fromDate", fromDate.Format(DATE_FORMAT),
		"toDate", toDate.Format(DATE_FORMAT),
	)

	// Make request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("error sending request", "error", err)
		return err
	}
	defer res.Body.Close()

	// Write cache while reading response body
	f, _ := os.Create(CACHE_FILE)
	r := io.TeeReader(res.Body, f)

	// Parse response
	err = json.NewDecoder(r).Decode(&sdm)
	if err != nil {
		slog.Error("error decoding JSON response", "error", err)
		return err
	}

	// Update data
	m.lastData = &sdm

	return nil
}

func (m *Monnit) LastReading() *SensorDataMessage {
	m.Lock()
	defer m.Unlock()

	if m.lastData == nil {
		return &SensorDataMessage{}
	}
	return m.lastData.GetLast()
}

// SensorDataMessages represents the structure for sensor data communication.
// It contains the method used and a slice of SensorDataMessage structs.
type SensorDataMessages struct {
	Method      string              `json:"Method"`
	Messages    []SensorDataMessage `json:"Result"`
	LastUpdated time.Time
}

func (sdm *SensorDataMessages) GetLast() *SensorDataMessage {
	if len(sdm.Messages) == 0 {
		return &SensorDataMessage{}
	}
	return &sdm.Messages[0]
}

type SensorDataMessage struct {
	DataMessageGUID             string      `json:"DataMessageGUID"`
	SensorID                    int         `json:"SensorID"`
	MessageDate                 MessageDate `json:"MessageDate"`
	State                       int         `json:"State"`
	SignalStrength              int         `json:"SignalStrength"`
	Voltage                     float64     `json:"Voltage"`
	Battery                     int         `json:"Battery"`
	Data                        string      `json:"Data"`
	DisplayData                 string      `json:"DisplayData"`
	Temperature                 Temperature `json:"PlotValue"`
	MetNotificationRequirements bool        `json:"MetNotificationRequirements"`
	GatewayID                   int         `json:"GatewayID"`
	DataValues                  string      `json:"DataValues"`
	DataTypes                   string      `json:"DataTypes"`
	PlotValues                  string      `json:"PlotValues"`
	PlotLabels                  string      `json:"PlotLabels"`
}

func (m SensorDataMessage) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Time("date", time.Time(m.MessageDate)),
		slog.String("temperature", m.Temperature.String()),
		slog.Int("signal_strength", m.SignalStrength),
	)
}
