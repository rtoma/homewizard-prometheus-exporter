package scraper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	watermeterActiveLiterLPM = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_watermeter_active_liter_lpm",
			Help: "Liters per minute",
		},
		[]string{"name"},
	)
	watermeterTotalLiterM3 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "homewizard_watermeter_total_liter_m3",
			Help: "Total M3 of water measured",
		},
		[]string{"name"},
	)
)

type WatermeterData struct {
	GenericData
	ActiveLiterLPM float64 `json:"active_liter_lpm"`
	TotalLiterM3   float64 `json:"total_liter_m3"`
	// TotalLiterOffsetM3 float64 `json:"total_liter_offset_m3"`
}

type WatermeterScraper struct {
	deviceBase
}

func (s *WatermeterScraper) Scrape() {
	// s.l.Print("Scrape")
	data := &WatermeterData{}
	if err := s.getData(data); err != nil {
		s.l.Printf("failed to fetch measurement data: %v", err)
		return
	}
	if s.debugEnabled {
		s.l.Printf("DEBUG: Scraped data: %+v", data)
	}

	s.collectDeviceInfo(&data.GenericData)

	watermeterActiveLiterLPM.WithLabelValues(s.Name).Set(data.ActiveLiterLPM)
	watermeterTotalLiterM3.WithLabelValues(s.Name).Set(data.TotalLiterM3)
}
