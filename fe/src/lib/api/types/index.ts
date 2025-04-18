/**
 * API types for the application
 */

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
 * Request parameters for paginated reports
 */
export interface PaginatedReportsRequest {
  limit?: number;      // Number of reports to return per page
  lastId?: string;     // ID of the last report from the previous page
  fromTime?: string;   // Filter reports from this time (ISO string)
  toTime?: string;     // Filter reports until this time (ISO string)
}

/**
 * Response for paginated reports
 */
export interface PaginatedReportsResponse {
  reports: WeatherReport[];  // List of reports for the current page
  totalCount?: number;       // Total number of reports (only when not filtered)
  hasMore: boolean;          // Whether there are more reports to fetch
  currentPage: number;       // Current page number
  fromNumber: number;        // Starting index of the current page
  toNumber: number;          // Ending index of the current page
}
