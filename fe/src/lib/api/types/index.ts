/**
 * API types for the application
 */

/**
 * Standardized API response structure
 */
export interface ApiResponse<T> {
  message: string;    // Human-readable message about the response
  status: string;     // Status of the response (success, error)
  errorCode: string;  // Error code (empty if status is success)
  data: T;            // Response data
}

/**
 * Weather report data structure
 */
export interface WeatherReport {
  id: string;
  timestamp: string;
  temperature: number;
  pressure: number;
  humidity: number;
  cloudCover: number;
  createdAt: string;
}

/**
 * Request parameters for generating a weather report
 */
export interface ReportRequest {
  timestamp?: string;
}

/**
 * Request parameters for comparing two reports
 */
export interface ComparisonRequest {
  reportId1: string;
  reportId2: string;
}

/**
 * Result of comparing two weather reports
 */
export interface ComparisonResult {
  report1: WeatherReport;
  report2: WeatherReport;
  deviation: {
    temperature: number;
    pressure: number;
    humidity: number;
    cloudCover: number;
  };
}

/**
 * Sort order for paginated reports
 */
export type SortOrder = 'asc' | 'desc';

/**
 * Request parameters for paginated reports
 */
export interface PaginatedReportsRequest {
  limit?: number; // Number of reports to return per page
  offset?: number; // Number of reports to skip (for pagination)
  fromTime?: string; // Filter reports from this time (ISO string)
  toTime?: string; // Filter reports until this time (ISO string)
  sortBy?: string; // Field to sort by (default: "timestamp")
  sortOrder?: SortOrder; // Sort order (default: "desc")
}

/**
 * Response for paginated reports
 */
export interface PaginatedReportsResponse {
  reports: WeatherReport[]; // List of reports for the current page
  totalCount: number; // Total number of reports (for calculating total pages)
}
