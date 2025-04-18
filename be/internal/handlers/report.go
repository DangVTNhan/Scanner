package handlers

import (
	"encoding/json"
	"net/http"

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

// GetAllReports handles requests to retrieve all weather reports
func (h *ReportHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportService.GetAllReports(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
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
