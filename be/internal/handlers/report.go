package handlers

import (
	"encoding/json"
	"github.com/DangVTNhan/Scanner/be/internal/interfaces"
	"github.com/DangVTNhan/Scanner/be/internal/models/errors"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// ReportHandler handles HTTP requests related to weather reports
type ReportHandler struct {
	reportService interfaces.IReportService
}

// NewReportHandler creates a new instance of ReportHandler
func NewReportHandler(reportService interfaces.IReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GenerateReport handles requests to generate a new weather report
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req request.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", errors.ErrCodeInvalidRequest, nil, http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GenerateReport(r.Context(), &req)
	if err != nil {
		errorCode := errors.ErrCodeServerError
		statusCode := http.StatusInternalServerError

		// Determine specific error code based on error message
		if strings.Contains(err.Error(), "failed to get weather data") {
			errorCode = errors.ErrCodeWeatherServiceResponse
		} else if strings.Contains(err.Error(), "failed to save report") {
			errorCode = errors.ErrCodeDatabaseInsert
		}

		respondWithError(w, err.Error(), errorCode, nil, statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	responseData := response.NewSuccessResponse("Report generated successfully", report)
	json.NewEncoder(w).Encode(responseData)
}

// GetAllReports handles requests to retrieve all weather reports (legacy endpoint)
func (h *ReportHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportService.GetAllReports(r.Context())
	if err != nil {
		errorCode := errors.ErrCodeServerError
		if strings.Contains(err.Error(), "failed to retrieve reports") {
			errorCode = errors.ErrCodeDatabaseQuery
		}
		respondWithError(w, err.Error(), errorCode, nil, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseData := response.NewSuccessResponse("Reports retrieved successfully", reports)
	json.NewEncoder(w).Encode(responseData)
}

// GetPaginatedReports handles requests to retrieve paginated weather reports with optional filtering
func (h *ReportHandler) GetPaginatedReports(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Create request object
	req := &request.PaginatedReportsRequest{
		IsFiltered: false,
		SortBy:     query.Get("sortBy"),
		SortOrder:  request.SortOrder(query.Get("sortOrder")),
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondWithError(w, "Invalid limit parameter", errors.ErrCodeInvalidParameters, nil, http.StatusBadRequest)
			return
		}
		req.Limit = limit
	}

	// Parse offset
	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondWithError(w, "Invalid offset parameter", errors.ErrCodeInvalidParameters, nil, http.StatusBadRequest)
			return
		}
		req.Offset = offset
	}

	// Parse from time
	if fromTimeStr := query.Get("fromTime"); fromTimeStr != "" {
		fromTime, err := time.Parse(time.RFC3339, fromTimeStr)
		if err != nil {
			respondWithError(w, "Invalid fromTime parameter", errors.ErrCodeInvalidParameters, nil, http.StatusBadRequest)
			return
		}
		req.FromTime = fromTime
		req.IsFiltered = true
	}

	// Parse to time
	if toTimeStr := query.Get("toTime"); toTimeStr != "" {
		toTime, err := time.Parse(time.RFC3339, toTimeStr)
		if err != nil {
			respondWithError(w, "Invalid toTime parameter", errors.ErrCodeInvalidParameters, nil, http.StatusBadRequest)
			return
		}
		req.ToTime = toTime
		req.IsFiltered = true
	}

	// Get paginated reports
	paginatedResponse, err := h.reportService.GetPaginatedReports(r.Context(), req)
	if err != nil {
		errorCode := errors.ErrCodeServerError
		if strings.Contains(err.Error(), "failed to retrieve reports") {
			errorCode = errors.ErrCodeDatabaseQuery
		}
		respondWithError(w, err.Error(), errorCode, nil, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseData := response.NewSuccessResponse("Reports retrieved successfully", paginatedResponse)
	json.NewEncoder(w).Encode(responseData)
}

// GetReportByID handles requests to retrieve a specific weather report
func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	report, err := h.reportService.GetReportByID(r.Context(), id)
	if err != nil {
		if err.Error() == "report not found" {
			respondWithError(w, "Report not found", errors.ErrCodeReportNotFound, nil, http.StatusNotFound)
		} else {
			errorCode := errors.ErrCodeServerError
			if strings.Contains(err.Error(), "failed to retrieve report") {
				errorCode = errors.ErrCodeDatabaseQuery
			}
			respondWithError(w, err.Error(), errorCode, nil, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseData := response.NewSuccessResponse("Report retrieved successfully", report)
	json.NewEncoder(w).Encode(responseData)
}

// CompareReports handles requests to compare two weather reports
func (h *ReportHandler) CompareReports(w http.ResponseWriter, r *http.Request) {
	var req request.ComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", errors.ErrCodeInvalidRequest, nil, http.StatusBadRequest)
		return
	}

	result, err := h.reportService.CompareReports(r.Context(), &req)
	if err != nil {
		errorCode := errors.ErrCodeServerError
		statusCode := http.StatusInternalServerError

		if strings.Contains(err.Error(), "failed to retrieve first report") ||
			strings.Contains(err.Error(), "failed to retrieve second report") {
			if strings.Contains(err.Error(), "report not found") {
				errorCode = errors.ErrCodeReportNotFound
				statusCode = http.StatusNotFound
			} else {
				errorCode = errors.ErrCodeDatabaseQuery
			}
		} else if strings.Contains(err.Error(), "failed to compare reports") {
			errorCode = errors.ErrCodeReportComparison
		}

		respondWithError(w, err.Error(), errorCode, nil, statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseData := response.NewSuccessResponse("Reports compared successfully", result)
	json.NewEncoder(w).Encode(responseData)
}

// respondWithError is a helper function to send standardized error responses
func respondWithError(w http.ResponseWriter, message string, errorCode string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := response.NewErrorResponse(message, errorCode, data)
	json.NewEncoder(w).Encode(response)
}
