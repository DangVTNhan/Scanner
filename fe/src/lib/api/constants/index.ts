/**
 * API constants for the application
 */

/**
 * Base URL for the API
 */
export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

/**
 * Default fetch options for API requests
 */
export const defaultOptions: RequestInit = {
  mode: "cors",
  credentials: "same-origin",
  headers: {
    "Content-Type": "application/json",
  },
};
