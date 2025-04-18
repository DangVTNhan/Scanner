/**
 * API client for weather reports
 */

import { API_BASE_URL, defaultOptions } from '../constants';
import { 
  WeatherReport, 
  ReportRequest, 
  ComparisonRequest, 
  ComparisonResult 
} from '../types';

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

  if (!response.ok) {
    throw new Error(`Failed to generate report: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Get all weather reports
 * @returns Promise with an array of weather reports
 */
export async function getAllReports(): Promise<WeatherReport[]> {
  const response = await fetch(`${API_BASE_URL}/reports`, defaultOptions);

  if (!response.ok) {
    throw new Error(`Failed to fetch reports: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Get a specific weather report by ID
 * @param id The ID of the report to retrieve
 * @returns Promise with the weather report
 */
export async function getReportById(id: string): Promise<WeatherReport> {
  const response = await fetch(`${API_BASE_URL}/reports/${id}`, defaultOptions);

  if (!response.ok) {
    throw new Error(`Failed to fetch report: ${response.statusText}`);
  }

  return response.json();
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

  if (!response.ok) {
    throw new Error(`Failed to compare reports: ${response.statusText}`);
  }

  return response.json();
}
