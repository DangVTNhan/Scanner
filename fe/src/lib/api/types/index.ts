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
