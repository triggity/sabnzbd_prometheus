package sabnzbd_prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SabnzbdStatsResponse struct {
	Day     int64                                 `json:"day"`
	Week    int64                                 `json:"week"`
	Month   int64                                 `json:"month"`
	Total   int64                                 `json:"total"`
	Servers map[string]SabnzbdStatsResponseServer `json:"servers"`
}
type SabnzbdStatsResponseServer struct {
	Day   int64            `json:"day"`
	Week  int64            `json:"week"`
	Month int64            `json:"month"`
	Total int64            `json:"total"`
	Daily map[string]int64 `json:"daily"`
}

// Queue response output. only putting in fields that are used
type SabnzbdQueueResponse struct {
	Queue SabnzbdQueueResponseQueue `json:"queue"`
}
type SabnzbdQueueResponseQueue struct {
	NoOfSlotsTotal int64  `json:"noofslots_total"`
	KbPerSec       string `json:"kbpersec"`
	MbLeft         string `json:"mbleft"`
	Mb             string `json:"mb"`
	TimeLeft       string `json:"timeleft"`
	SpeedLimit     string `json:"speedlimit"`
	SpeedLimitAbs  string `json:"speedlimit_abs"`
}

type SabNzbdClient interface {
	GetServerStats() (SabnzbdStatsResponse, error)
	GetQueue() (SabnzbdQueueResponse, error)
}

func NewSabNzbdClient(baseUri string, apiKey string) SabNzbdClient {
	return &sabNzbdClient{baseUri, apiKey}
}

type sabNzbdClient struct {
	baseUri string
	apiKey  string
}

func (s *sabNzbdClient) createUri(mode string) string {
	return fmt.Sprintf("%s/api?output=json&apikey=%s&mode=%s", s.baseUri, s.apiKey, mode)
}

func (s *sabNzbdClient) GetServerStats() (SabnzbdStatsResponse, error) {
	fullUri := s.createUri("server_stats")

	r, err := http.Get(fullUri)
	if err != nil {
		return SabnzbdStatsResponse{}, err
	}
	var response SabnzbdStatsResponse
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return SabnzbdStatsResponse{}, err
	}
	return response, nil

}

func (s *sabNzbdClient) GetQueue() (SabnzbdQueueResponse, error) {
	fullUri := s.createUri("queue")
	r, err := http.Get(fullUri)
	if err != nil {
		return SabnzbdQueueResponse{}, err
	}
	var response SabnzbdQueueResponse

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return SabnzbdQueueResponse{}, err
	}
	return response, nil
}
