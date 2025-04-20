import { mockFetchError, mockFetchResponse } from "@/lib/test-utils";
import "@testing-library/jest-dom";
import {
  compareReports,
  generateReport,
  getPaginatedReports,
  getReportById,
} from "../client/reports";
import { API_BASE_URL } from "../constants";
import {
  ComparisonResult,
  PaginatedReportsResponse,
  WeatherReport,
} from "../types";

// Mock global fetch
global.fetch = jest.fn();

describe("API Client - Reports", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("generateReport", () => {
    const mockReport: WeatherReport = {
      id: "report1",
      timestamp: "2023-06-15T10:00:00Z",
      temperature: 30.5,
      pressure: 1013,
      humidity: 75,
      cloudCover: 20,
      createdAt: "2023-06-15T10:05:00Z",
    };

    test("calls the correct endpoint with POST method", async () => {
      // Mock successful response
      mockFetchResponse(mockReport);

      // Call the function
      await generateReport();

      // Verify fetch was called correctly
      expect(fetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/reports`,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({}),
        })
      );
    });

    test("passes timestamp when provided", async () => {
      // Mock successful response
      mockFetchResponse(mockReport);

      // Call the function with timestamp
      const timestamp = "2023-06-15T10:00:00Z";
      await generateReport({ timestamp });

      // Verify fetch was called with the timestamp
      expect(fetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/reports`,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({ timestamp }),
        })
      );
    });

    test("throws error when API call fails", async () => {
      // Mock error response
      mockFetchError("Failed to generate report");

      // Call the function and expect it to throw
      await expect(generateReport()).rejects.toThrow(
        "Failed to generate report"
      );
    });
  });

  describe("getPaginatedReports", () => {
    const mockReports: WeatherReport[] = [
      {
        id: "report1",
        timestamp: "2023-06-15T10:00:00Z",
        temperature: 30.5,
        pressure: 1013,
        humidity: 75,
        cloudCover: 20,
        createdAt: "2023-06-15T10:05:00Z",
      },
      {
        id: "report2",
        timestamp: "2023-06-16T10:00:00Z",
        temperature: 31.2,
        pressure: 1010,
        humidity: 80,
        cloudCover: 30,
        createdAt: "2023-06-16T10:05:00Z",
      },
    ];

    const mockPaginatedResponse: PaginatedReportsResponse = {
      reports: mockReports,
      totalCount: 20,
    };

    test("calls the correct endpoint with GET method", async () => {
      // Mock successful response
      mockFetchResponse(mockPaginatedResponse);

      // Call the function
      await getPaginatedReports();

      // Verify fetch was called correctly
      expect(fetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/reports/paginated`,
        expect.anything()
      );
    });

    test("includes query parameters when provided", async () => {
      // Mock successful response
      mockFetchResponse(mockPaginatedResponse);

      // Call the function with parameters
      await getPaginatedReports({
        limit: 10,
        offset: 20,
        sortBy: "temperature",
        sortOrder: "desc",
        fromTime: "2023-06-15T00:00:00Z",
        toTime: "2023-06-16T23:59:59Z",
      });

      // Verify fetch was called with the correct query parameters
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining(`${API_BASE_URL}/reports/paginated?`),
        expect.anything()
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("limit=10"),
        expect.anything()
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("offset=20"),
        expect.anything()
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("sortBy=temperature"),
        expect.anything()
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("sortOrder=desc"),
        expect.anything()
      );
      // Check for encoded URL parameters
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("fromTime="),
        expect.anything()
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("toTime="),
        expect.anything()
      );

      // Verify the actual URL contains the encoded timestamps
      const actualUrl = (fetch as jest.Mock).mock.calls[0][0];
      expect(actualUrl).toContain(encodeURIComponent("2023-06-15T00:00:00Z"));
      expect(actualUrl).toContain(encodeURIComponent("2023-06-16T23:59:59Z"));
    });

    test("throws error when API call fails", async () => {
      // Mock error response
      mockFetchError("Failed to fetch paginated reports");

      // Call the function and expect it to throw
      await expect(getPaginatedReports()).rejects.toThrow(
        "Failed to fetch paginated reports"
      );
    });
  });

  describe("getReportById", () => {
    const mockReport: WeatherReport = {
      id: "report1",
      timestamp: "2023-06-15T10:00:00Z",
      temperature: 30.5,
      pressure: 1013,
      humidity: 75,
      cloudCover: 20,
      createdAt: "2023-06-15T10:05:00Z",
    };

    test("calls the correct endpoint with GET method", async () => {
      // Mock successful response
      mockFetchResponse(mockReport);

      // Call the function
      await getReportById("report1");

      // Verify fetch was called correctly
      expect(fetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/reports/report1`,
        expect.anything()
      );
    });

    test("throws error when API call fails", async () => {
      // Mock error response
      mockFetchError("Failed to fetch report");

      // Call the function and expect it to throw
      await expect(getReportById("report1")).rejects.toThrow(
        "Failed to fetch report"
      );
    });
  });

  describe("compareReports", () => {
    const mockComparisonResult: ComparisonResult = {
      report1: {
        id: "report1",
        timestamp: "2023-06-15T10:00:00Z",
        temperature: 30.5,
        pressure: 1013,
        humidity: 75,
        cloudCover: 20,
        createdAt: "2023-06-15T10:05:00Z",
      },
      report2: {
        id: "report2",
        timestamp: "2023-06-16T10:00:00Z",
        temperature: 31.2,
        pressure: 1010,
        humidity: 80,
        cloudCover: 30,
        createdAt: "2023-06-16T10:05:00Z",
      },
      deviation: {
        temperature: 0.7,
        pressure: -3,
        humidity: 5,
        cloudCover: 10,
      },
    };

    test("calls the correct endpoint with POST method", async () => {
      // Mock successful response
      mockFetchResponse(mockComparisonResult);

      // Call the function
      await compareReports({
        reportId1: "report1",
        reportId2: "report2",
      });

      // Verify fetch was called correctly
      expect(fetch).toHaveBeenCalledWith(
        `${API_BASE_URL}/reports/compare`,
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify({
            reportId1: "report1",
            reportId2: "report2",
          }),
        })
      );
    });

    test("throws error when API call fails", async () => {
      // Mock error response
      mockFetchError("Failed to compare reports");

      // Call the function and expect it to throw
      await expect(
        compareReports({
          reportId1: "report1",
          reportId2: "report2",
        })
      ).rejects.toThrow("Failed to compare reports");
    });
  });
});
