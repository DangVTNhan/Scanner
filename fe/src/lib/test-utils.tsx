import { render, RenderOptions } from "@testing-library/react";
import { ReactElement } from "react";

// Create a custom render function that includes any global providers
const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">
) => render(ui, { ...options });

// Re-export everything from testing-library
export * from "@testing-library/react";

// Override the render method
export { customRender as render };

// Mock API response helper
export const mockApiResponse = <T,>(
  data: T,
  status = "success",
  message = "Success",
  errorCode = ""
) => {
  return {
    status,
    message,
    errorCode,
    data,
  };
};

// Helper to mock fetch for API calls
export const mockFetchResponse = <T,>(data: T, ok = true, status = 200) => {
  global.fetch = jest.fn().mockResolvedValueOnce({
    ok,
    status,
    json: async () => mockApiResponse(data),
  });
};

// Helper to mock fetch error
export const mockFetchError = (message: string, status = 500) => {
  global.fetch = jest.fn().mockResolvedValueOnce({
    ok: false,
    status,
    json: async () => ({
      status: "error",
      message,
      errorCode: "ERR1000",
      data: null,
    }),
  });
};
