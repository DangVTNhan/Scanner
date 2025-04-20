/**
 * API client for weather reports
 */

import { API_BASE_URL, defaultOptions } from '../constants';
import {
  ComparisonRequest,
  ComparisonResult,
  PaginatedReportsRequest,
  PaginatedReportsResponse,
  ReportRequest,
  WeatherReport,
} from "../types";
import { processApiResponse } from "../utils";

/**
 * Generate a new weather report
 * @param request Optional timestamp for the report
 * @returns Promise with the generated weather report
 */
export async function generateReport(
  request: ReportRequest = {}
): Promise<WeatherReport> {
  const response = await fetch(`${API_BASE_URL}/reports`, {
    ...defaultOptions,
    method: "POST",
    body: JSON.stringify(request),
  });

  return processApiResponse<WeatherReport>(
    response,
    "Failed to generate report"
  );
}

/**
 * Get all weather reports (legacy method)
 * @returns Promise with an array of weather reports
 */
export async function getAllReports(): Promise<WeatherReport[]> {
  const response = await fetch(`${API_BASE_URL}/reports`, defaultOptions);

  return processApiResponse<WeatherReport[]>(
    response,
    "Failed to fetch reports"
  );
}

/**
 * Get paginated reports with optional filtering
 * @param params Pagination and filtering parameters
 * @returns Promise with paginated reports response
 */
export async function getPaginatedReports(
  params: PaginatedReportsRequest = {}
): Promise<PaginatedReportsResponse> {
  // Build query string
  const queryParams = new URLSearchParams();

  if (params.limit) {
    queryParams.append("limit", params.limit.toString());
  }

  if (params.lastId) {
    queryParams.append("lastId", params.lastId);
  }

  if (params.fromTime) {
    queryParams.append("fromTime", params.fromTime);
  }

  if (params.toTime) {
    queryParams.append("toTime", params.toTime);
  }

  const url = `${API_BASE_URL}/reports/paginated${
    queryParams.toString() ? `?${queryParams.toString()}` : ""
  }`;

  const response = await fetch(url, defaultOptions);

  return processApiResponse<PaginatedReportsResponse>(
    response,
    "Failed to fetch paginated reports"
  );
}

/**
 * Get a specific weather report by ID
 * @param id The ID of the report to retrieve
 * @returns Promise with the weather report
 */
export async function getReportById(id: string): Promise<WeatherReport> {
  const response = await fetch(`${API_BASE_URL}/reports/${id}`, defaultOptions);

  return processApiResponse<WeatherReport>(response, "Failed to fetch report");
}

/**
 * Compare two weather reports
 * @param request The IDs of the reports to compare
 * @returns Promise with the comparison result
 */
export async function compareReports(
  request: ComparisonRequest
): Promise<ComparisonResult> {
  const response = await fetch(`${API_BASE_URL}/reports/compare`, {
    ...defaultOptions,
    method: "POST",
    body: JSON.stringify(request),
  });

  return processApiResponse<ComparisonResult>(
    response,
    "Failed to compare reports"
  );
}
