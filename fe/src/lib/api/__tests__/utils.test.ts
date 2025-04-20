import "@testing-library/jest-dom";
import { toast } from "sonner";
import { handleApiError, processApiResponse } from "../utils";

// Mock toast
jest.mock("sonner", () => ({
  toast: {
    error: jest.fn(),
  },
}));

// Mock console.error
const originalConsoleError = console.error;
beforeEach(() => {
  console.error = jest.fn();
});
afterEach(() => {
  console.error = originalConsoleError;
});

describe("API Utils", () => {
  describe("handleApiError", () => {
    test("handles Error objects correctly", () => {
      const error = new Error("Test error message");
      handleApiError(error, "Fallback message");

      expect(console.error).toHaveBeenCalledWith("Fallback message", error);
      expect(toast.error).toHaveBeenCalledWith("Test error message");
    });

    test("handles objects with message property", () => {
      const error = { message: "Object error message" };
      handleApiError(error, "Fallback message");

      expect(console.error).toHaveBeenCalledWith("Fallback message", error);
      expect(toast.error).toHaveBeenCalledWith("Object error message");
    });

    test("uses fallback message when error has no message", () => {
      const error = new Error();
      handleApiError(error, "Fallback message");

      expect(console.error).toHaveBeenCalledWith("Fallback message", error);
      expect(toast.error).toHaveBeenCalledWith("Fallback message");
    });

    test("uses fallback message for non-error objects", () => {
      handleApiError("string error", "Fallback message");

      expect(console.error).toHaveBeenCalledWith(
        "Fallback message",
        "string error"
      );
      expect(toast.error).toHaveBeenCalledWith("Fallback message");
    });
  });

  describe("processApiResponse", () => {
    test("returns data for successful response", async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue({
          status: "success",
          message: "Success",
          errorCode: "",
          data: { id: "123", name: "Test" },
        }),
      };

      const result = await processApiResponse(
        mockResponse as unknown as Response,
        "Fallback error"
      );

      expect(result).toEqual({ id: "123", name: "Test" });
    });

    test("throws error for non-ok response", async () => {
      const mockResponse = {
        ok: false,
        json: jest.fn().mockResolvedValue({
          status: "error",
          message: "API error",
          errorCode: "ERR1000",
          data: null,
        }),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("API error (ERR1000)");
    });

    test("throws error for response with error status", async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue({
          status: "error",
          message: "Business logic error",
          errorCode: "ERR2000",
          data: null,
        }),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("Business logic error (ERR2000)");
    });

    test("includes error code in message when available", async () => {
      const mockResponse = {
        ok: false,
        json: jest.fn().mockResolvedValue({
          status: "error",
          message: "Not found",
          errorCode: "ERR1003",
          data: null,
        }),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("Not found (ERR1003)");
    });

    test("uses fallback message when response has no message", async () => {
      const mockResponse = {
        ok: false,
        json: jest.fn().mockResolvedValue({
          status: "error",
          message: "",
          errorCode: "",
          data: null,
        }),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("Fallback error");
    });

    test("handles JSON parsing errors with Error objects", async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockRejectedValue(new Error("JSON parse error")),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("JSON parse error");
    });

    test("handles JSON parsing errors with non-Error objects", async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockRejectedValue("Not an Error object"),
      };

      await expect(
        processApiResponse(
          mockResponse as unknown as Response,
          "Fallback error"
        )
      ).rejects.toThrow("Fallback error");
    });
  });
});
