package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/services"
	"github.com/gorilla/mux"
)

// ReportHandler handles HTTP requests related to weather reports
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler creates a new instance of ReportHandler
func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GenerateReport handles requests to generate a new weather report
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req models.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GenerateReport(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetAllReports handles requests to retrieve all weather reports (legacy endpoint)
func (h *ReportHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportService.GetAllReports(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// GetPaginatedReports handles requests to retrieve paginated weather reports with optional filtering
func (h *ReportHandler) GetPaginatedReports(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Create request object
	req := &models.PaginatedReportsRequest{
		LastID: query.Get("lastId"),
		IsFiltered: false,
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		req.Limit = limit
	}

	// Parse from time
	if fromTimeStr := query.Get("fromTime"); fromTimeStr != "" {
		fromTime, err := time.Parse(time.RFC3339, fromTimeStr)
		if err != nil {
			http.Error(w, "Invalid fromTime parameter", http.StatusBadRequest)
			return
		}
		req.FromTime = fromTime
		req.IsFiltered = true
	}

	// Parse to time
	if toTimeStr := query.Get("toTime"); toTimeStr != "" {
		toTime, err := time.Parse(time.RFC3339, toTimeStr)
		if err != nil {
			http.Error(w, "Invalid toTime parameter", http.StatusBadRequest)
			return
		}
		req.ToTime = toTime
		req.IsFiltered = true
	}

	// Get paginated reports
	response, err := h.reportService.GetPaginatedReports(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetReportByID handles requests to retrieve a specific weather report
func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	report, err := h.reportService.GetReportByID(r.Context(), id)
	if err != nil {
		if err.Error() == "report not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// CompareReports handles requests to compare two weather reports
func (h *ReportHandler) CompareReports(w http.ResponseWriter, r *http.Request) {
	var req models.ComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.reportService.CompareReports(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
