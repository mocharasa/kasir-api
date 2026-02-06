package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	var (
		result interface{}
		err    error
	)

	if startDate == "" && endDate == "" {
		// /api/report/hari-ini
		result, err = h.service.GetTodaySummary()
	} else {
		start, err1 := time.Parse("2006-01-02", startDate)
		end, err2 := time.Parse("2006-01-02", endDate)
		if err1 != nil || err2 != nil {
			http.Error(w, "invalid date format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		end = end.Add(24 * time.Hour)
		result, err = h.service.GetSummaryByDateRange(start, end)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}
