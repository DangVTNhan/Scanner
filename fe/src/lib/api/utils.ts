/**
 * API utility functions
 */

import { toast } from "sonner";
import { ApiResponse } from "./types";

/**
 * Handle API errors in a consistent way
 * @param error The error object
 * @param fallbackMessage A fallback message to display if the error doesn't have a message
 */
export function handleApiError(error: unknown, fallbackMessage: string): void {
  console.error(fallbackMessage, error);
  
  if (error instanceof Error) {
    toast.error(error.message || fallbackMessage);
  } else if (typeof error === 'object' && error !== null && 'message' in error) {
    toast.error((error as { message: string }).message || fallbackMessage);
  } else {
    toast.error(fallbackMessage);
  }
}

/**
 * Process API response and handle errors
 * @param response The fetch response
 * @param fallbackErrorMessage A fallback error message
 * @returns The processed API response data
 */
export async function processApiResponse<T>(
  response: Response,
  fallbackErrorMessage: string
): Promise<T> {
  try {
    const responseData: ApiResponse<T> = await response.json();
    
    if (!response.ok || responseData.status === 'error') {
      // If there's an error code, include it in the error message
      const errorMessage = responseData.errorCode 
        ? `${responseData.message} (${responseData.errorCode})`
        : responseData.message;
        
      throw new Error(errorMessage || fallbackErrorMessage);
    }
    
    return responseData.data;
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error(fallbackErrorMessage);
  }
}
