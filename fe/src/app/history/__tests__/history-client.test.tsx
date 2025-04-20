import { PaginatedReportsResponse, WeatherReport } from "@/lib/api";
import {
  fireEvent,
  mockFetchResponse,
  render,
  screen,
  waitFor,
} from "@/lib/test-utils";
import "@testing-library/jest-dom";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";
import HistoryClient from "../history-client";

// Mock the next/navigation hooks
jest.mock("next/navigation", () => ({
  useRouter: jest.fn(() => ({
    push: jest.fn(),
  })),
  useSearchParams: jest.fn(),
}));

// Mock console.error to prevent test output pollution
const originalConsoleError = console.error;
beforeEach(() => {
  console.error = jest.fn();
});
afterEach(() => {
  console.error = originalConsoleError;
});

// Sample reports data for tests
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

describe("HistoryClient", () => {
  const mockPush = jest.fn();

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();

    // Setup router mock
    (useRouter as jest.Mock).mockReturnValue({
      push: mockPush,
    });

    // Default search params
    (useSearchParams as jest.Mock).mockReturnValue({
      get: (param: string) => {
        switch (param) {
          case "offset":
            return "0";
          case "fromTime":
            return "";
          case "toTime":
            return "";
          case "sortBy":
            return "timestamp";
          case "sortOrder":
            return "desc";
          default:
            return null;
        }
      },
    });
  });

  test("renders loading state initially", () => {
    // Mock API call but don't resolve it yet
    global.fetch = jest.fn(() => new Promise(() => {}));

    render(<HistoryClient />);

    // Check if loading state is displayed
    expect(screen.getByText("Loading reports...")).toBeInTheDocument();
  });

  test("renders empty state when no reports are found", async () => {
    // Mock empty response
    mockFetchResponse({ reports: [], totalCount: 0 });

    render(<HistoryClient />);

    // Wait for the component to update
    await waitFor(() => {
      expect(screen.getByText("No Reports Found")).toBeInTheDocument();
    });
  });

  test("renders reports data successfully", async () => {
    // Mock successful API response
    mockFetchResponse(mockPaginatedResponse);

    render(<HistoryClient />);

    // Wait for the component to update and display reports with a longer timeout
    await waitFor(
      () => {
        expect(screen.getByText("Weather Report History")).toBeInTheDocument();
        expect(screen.getByText("Historical Reports")).toBeInTheDocument();

        // Wait for the table to be rendered
        expect(screen.getByRole("table")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Check if table headers are displayed using a more flexible approach
    // Look for elements containing the text within the table headers
    const tableHeaders = screen.getAllByRole("columnheader");
    expect(tableHeaders.length).toBeGreaterThan(0);

    // Check for the presence of specific text in the rendered HTML
    const html = document.body.innerHTML;
    expect(html).toContain("Temperature");
    expect(html).toContain("Pressure");
    expect(html).toContain("Humidity");
    expect(html).toContain("Cloud Cover");

    // Check if report data is displayed
    expect(screen.getByText("30.5")).toBeInTheDocument();
    expect(screen.getByText("31.2")).toBeInTheDocument();
    expect(screen.getByText("1013")).toBeInTheDocument();
    expect(screen.getByText("1010")).toBeInTheDocument();
  });

  test("handles report selection correctly", async () => {
    // Mock successful API response
    mockFetchResponse(mockPaginatedResponse);

    render(<HistoryClient />);

    // Wait for the component to update and display reports
    await waitFor(
      () => {
        expect(screen.getByText("Historical Reports")).toBeInTheDocument();
        // Wait for the table to be rendered
        expect(screen.getByRole("table")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Find and click the select buttons for both reports
    const selectButtons = screen.getAllByRole("button", { name: "" });
    fireEvent.click(selectButtons[0]); // Select first report
    fireEvent.click(selectButtons[1]); // Select second report

    // Get all compare buttons (there are two - one for desktop and one for mobile)
    const compareButtons = screen.getAllByText(/Compare Selected Reports/);
    expect(compareButtons.length).toBeGreaterThan(0);

    // Get the first button and check if it's enabled
    const compareButton = compareButtons[0];
    expect(compareButton).not.toBeDisabled();

    // Find the actual button element (the span is inside a button)
    const buttonElement = compareButton.closest("button");
    expect(buttonElement).not.toBeNull();

    // Click the compare button
    if (buttonElement) {
      fireEvent.click(buttonElement);
    }

    // Verify that router.push was called with the correct URL
    expect(mockPush).toHaveBeenCalledWith(
      "/compare?report1=report1&report2=report2"
    );
  });

  test("prevents selecting more than two reports", async () => {
    // Create a mock with 3 reports
    const threeReports = [
      ...mockReports,
      {
        id: "report3",
        timestamp: "2023-06-17T10:00:00Z",
        temperature: 32.0,
        pressure: 1008,
        humidity: 85,
        cloudCover: 40,
        createdAt: "2023-06-17T10:05:00Z",
      },
    ];

    mockFetchResponse({
      reports: threeReports,
      totalCount: 3,
    });

    render(<HistoryClient />);

    // Wait for the component to update and display reports
    await waitFor(
      () => {
        expect(screen.getByText("Historical Reports")).toBeInTheDocument();
        // Wait for the table to be rendered
        expect(screen.getByRole("table")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Find and click the select buttons for all three reports
    const selectButtons = screen.getAllByRole("button", { name: "" });
    fireEvent.click(selectButtons[0]); // Select first report
    fireEvent.click(selectButtons[1]); // Select second report
    fireEvent.click(selectButtons[2]); // Try to select third report

    // Verify that toast.error was called
    expect(toast.error).toHaveBeenCalledWith(
      "You can only select two reports for comparison"
    );
  });

  test("handles filter submission correctly", async () => {
    // Mock successful API response
    mockFetchResponse(mockPaginatedResponse);

    render(<HistoryClient />);

    // Wait for the component to update
    await waitFor(
      () => {
        expect(screen.getByText("Filter Reports")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Fill in the filter form
    const fromTimeInput = screen.getByLabelText(/From Date & Time/);
    const toTimeInput = screen.getByLabelText(/To Date & Time/);

    fireEvent.change(fromTimeInput, { target: { value: "2023-06-15T00:00" } });
    fireEvent.change(toTimeInput, { target: { value: "2023-06-16T23:59" } });

    // Submit the form
    const applyButton = screen.getByText("Apply Filters");
    fireEvent.click(applyButton);

    // Wait for the router.push to be called
    await waitFor(
      () => {
        expect(mockPush).toHaveBeenCalled();
      },
      { timeout: 1000 }
    );

    // Verify that router.push was called with the correct URL parameters
    // The URL is encoded, so we need to check for the encoded values
    const pushCall = mockPush.mock.calls[0][0];
    expect(pushCall).toContain("fromTime=");
    expect(pushCall).toContain("toTime=");
    expect(pushCall).toContain(encodeURIComponent("2023-06-15T00:00"));
    expect(pushCall).toContain(encodeURIComponent("2023-06-16T23:59"));
  });

  test("handles pagination correctly", async () => {
    // Mock successful API response
    mockFetchResponse(mockPaginatedResponse);

    render(<HistoryClient />);

    // Wait for the component to update
    await waitFor(
      () => {
        expect(screen.getByText("Historical Reports")).toBeInTheDocument();
        // Wait for the table to be rendered
        expect(screen.getByRole("table")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Click the next page button
    const nextButton = screen.getByText("Next");
    fireEvent.click(nextButton);

    // Wait for the router.push to be called
    await waitFor(
      () => {
        expect(mockPush).toHaveBeenCalled();
      },
      { timeout: 1000 }
    );

    // Verify that router.push was called with the correct offset
    const pushCall = mockPush.mock.calls[0][0];
    expect(pushCall).toContain("offset=10");
  });

  test("handles sorting correctly", async () => {
    // Mock successful API response
    mockFetchResponse(mockPaginatedResponse);

    render(<HistoryClient />);

    // Wait for the component to update
    await waitFor(
      () => {
        expect(screen.getByText("Historical Reports")).toBeInTheDocument();
        // Wait for the table to be rendered
        expect(screen.getByRole("table")).toBeInTheDocument();
      },
      { timeout: 3000 }
    );

    // Find the temperature column header using a more flexible approach
    // Look for the div with the onClick handler that contains "Temperature"
    const temperatureHeaderDiv = screen
      .getByText(/Temperature/i)
      .closest("div");

    // Make sure we found the temperature header div
    expect(temperatureHeaderDiv).toBeTruthy();

    // Click on the temperature header div to trigger the sort
    if (temperatureHeaderDiv) {
      fireEvent.click(temperatureHeaderDiv);
    }

    // Wait for the router.push to be called
    await waitFor(
      () => {
        expect(mockPush).toHaveBeenCalled();
      },
      { timeout: 1000 }
    );

    // Verify that router.push was called with the correct sort parameters
    const pushCall = mockPush.mock.calls[0][0];
    expect(pushCall).toContain("sortBy=temperature");
    expect(pushCall).toContain("sortOrder=desc");
  });
});
