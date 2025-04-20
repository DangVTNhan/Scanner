import { ComparisonResult } from "@/lib/api";
import {
  mockFetchError,
  mockFetchResponse,
  render,
  screen,
  waitFor,
} from "@/lib/test-utils";
import "@testing-library/jest-dom";
import { useSearchParams } from "next/navigation";
import { toast } from "sonner";
import ComparisonClient from "../comparison-client";

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

// Sample comparison data for tests
const mockComparisonData: ComparisonResult = {
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

describe("ComparisonClient", () => {
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
  });

  test("renders loading state initially", () => {
    // Mock search params with valid report IDs
    (useSearchParams as jest.Mock).mockReturnValue({
      get: (param: string) => {
        if (param === "report1") return "report1";
        if (param === "report2") return "report2";
        return null;
      },
    });

    // Mock API call but don't resolve it yet
    global.fetch = jest.fn(() => new Promise(() => {}));

    render(<ComparisonClient />);

    // Check if loading state is displayed
    expect(screen.getByText("Loading Comparison Data")).toBeInTheDocument();
    expect(
      screen.getByText("Please wait while we analyze the weather reports...")
    ).toBeInTheDocument();
  });

  test("renders error state when report IDs are missing", async () => {
    // Mock search params with missing report IDs
    (useSearchParams as jest.Mock).mockReturnValue({
      get: () => null,
    });

    render(<ComparisonClient />);

    // Wait for the component to update
    await waitFor(() => {
      expect(screen.getByText("Error Loading Comparison")).toBeInTheDocument();
      expect(
        screen.getByText("Two report IDs are required for comparison")
      ).toBeInTheDocument();
    });
  });

  test("renders error state when API call fails", async () => {
    // Mock search params with valid report IDs
    (useSearchParams as jest.Mock).mockReturnValue({
      get: (param: string) => {
        if (param === "report1") return "report1";
        if (param === "report2") return "report2";
        return null;
      },
    });

    // Mock API call failure
    mockFetchError("Failed to compare reports");

    render(<ComparisonClient />);

    // Wait for the component to update
    await waitFor(() => {
      expect(screen.getByText("Error Loading Comparison")).toBeInTheDocument();
      expect(
        screen.getByText("Failed to compare the selected reports")
      ).toBeInTheDocument();
    });

    // Verify error handling
    expect(toast.error).toHaveBeenCalled();
    expect(console.error).toHaveBeenCalled();
  });

  test("renders comparison data successfully", async () => {
    // Mock search params with valid report IDs
    (useSearchParams as jest.Mock).mockReturnValue({
      get: (param: string) => {
        if (param === "report1") return "report1";
        if (param === "report2") return "report2";
        return null;
      },
    });

    // Mock successful API response
    mockFetchResponse(mockComparisonData);

    render(<ComparisonClient />);

    // Wait for the component to update and display comparison data
    await waitFor(() => {
      expect(screen.getByText("Weather Report Comparison")).toBeInTheDocument();
      expect(screen.getByText("Comparison Summary")).toBeInTheDocument();
    });

    // Check if report data is displayed using a more specific approach
    // Use getAllByText and check the first occurrence which should be the card heading
    const report1Elements = screen.getAllByText("Report 1");
    const report2Elements = screen.getAllByText("Report 2");
    expect(report1Elements.length).toBeGreaterThan(0);
    expect(report2Elements.length).toBeGreaterThan(0);

    // Check if temperature values are displayed
    // Use a more specific approach to find temperature values
    const temperatureValues = screen.getAllByText(/\d+\.\d+ °C/);
    expect(temperatureValues.length).toBeGreaterThanOrEqual(2);
    expect(
      temperatureValues.some((el) => el.textContent?.includes("30.5"))
    ).toBe(true);
    expect(
      temperatureValues.some((el) => el.textContent?.includes("31.2"))
    ).toBe(true);

    // Check if the table contains the parameter names
    expect(screen.getByText("Temperature (°C)")).toBeInTheDocument();
    expect(screen.getByText("Pressure (hPa)")).toBeInTheDocument();
    expect(screen.getByText("Humidity (%)")).toBeInTheDocument();
    expect(screen.getByText("Cloud Cover (%)")).toBeInTheDocument();

    // Check if the back button is present
    expect(screen.getByText("Back to History")).toBeInTheDocument();
  });
});
