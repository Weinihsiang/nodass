package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const externalAPI = "https://nodass.namr.gov.tw/noapi/namr/v1/obs/stations/WRA"

// Station represents the structure of a single station object from the API.
type Station struct {
	ID               string  `json:"Id"`
	StationID        string  `json:"StationID"`
	StationName      string  `json:"StationName"`
	StationNameLocal string  `json:"StationNameLocal"`
	CenterLongitude  float64 `json:"CenterLongitude"`
	CenterLatitude   float64 `json:"CenterLatitude"`
	StationChargeID  string  `json:"StationChargeID"`
	StationTypeID    string  `json:"StationTypeID"`
	EarliestDate     string  `json:"EarliestDate"`
	LatestDate       string  `json:"LatestDate"`
	AccessURL        string  `json:"AccessURL"`
}

func main() {
	http.HandleFunc("/api/stations", getStations)
	log.Println("Go backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getStations(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allows all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	resp, err := http.Get(externalAPI)
	if err != nil {
		http.Error(w, "Failed to fetch data from external API", http.StatusInternalServerError)
		log.Printf("Error fetching from external API: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		log.Printf("Error reading response body: %v", err)
		return
	}

	var stations []Station
	err = json.Unmarshal(body, &stations)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}
