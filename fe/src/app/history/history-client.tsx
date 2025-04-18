/*eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getPaginatedReports, PaginatedReportsResponse } from "@/lib/api";
import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

// Define the filter form type
interface FilterForm {
  fromTime: string;
  toTime: string;
}

export default function HistoryClient() {
  const [paginatedData, setPaginatedData] =
    useState<PaginatedReportsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedReports, setSelectedReports] = useState<string[]>([]);

  const router = useRouter();
  const searchParams = useSearchParams();

  // Get URL parameters
  const lastId = searchParams.get("lastId") || undefined;
  const fromTime = searchParams.get("fromTime") || "";
  const toTime = searchParams.get("toTime") || "";
  const currentPage = parseInt(searchParams.get("page") || "1", 10);

  // Form for filters
  const form = useForm<FilterForm>({
    defaultValues: {
      fromTime,
      toTime,
    },
  });

  // Set form values from URL parameters on initial load
  useEffect(() => {
    form.reset({
      fromTime,
      toTime,
    });
  }, [form, fromTime, toTime]);

  const limit = 10; // Number of items per page

  // Update URL with current filter and pagination state
  const updateUrl = useCallback(
    (params: {
      lastId?: string;
      fromTime?: string;
      toTime?: string;
      page?: number;
    }) => {
      const newParams = new URLSearchParams();

      if (params.lastId) {
        newParams.set("lastId", params.lastId);
      }

      if (params.fromTime) {
        newParams.set("fromTime", params.fromTime);
      }

      if (params.toTime) {
        newParams.set("toTime", params.toTime);
      }

      if (params.page && params.page > 1) {
        newParams.set("page", params.page.toString());
      }

      const newUrl = newParams.toString() ? `?${newParams.toString()}` : "";

      router.push(`/history${newUrl}`, { scroll: false });
    },
    [router]
  );

  // Fetch reports with pagination and optional filtering
  const fetchReports = useCallback(async () => {
    try {
      setLoading(true);

      // Prepare request parameters
      const params: any = { limit };

      if (lastId) {
        params.lastId = lastId;
      }

      if (fromTime) {
        try {
          params.fromTime = new Date(fromTime).toISOString();
        } catch (e) {
          console.error("Invalid fromTime format:", e);
        }
      }

      if (toTime) {
        try {
          params.toTime = new Date(toTime).toISOString();
        } catch (e) {
          console.error("Invalid toTime format:", e);
        }
      }

      // Fetch paginated reports
      const data = await getPaginatedReports(params);
      setPaginatedData(data);
    } catch (error) {
      console.error("Failed to fetch reports:", error);
      toast.error("Failed to load weather reports");
    } finally {
      setLoading(false);
    }
  }, [lastId, fromTime, toTime, limit]);

  // Initial data fetch
  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  // Handle filter submission
  const onSubmitFilter = (data: FilterForm) => {
    // Update URL with filter parameters and reset pagination
    updateUrl({
      fromTime: data.fromTime,
      toTime: data.toTime,
      page: 1,
    });
  };

  // Handle filter reset
  const resetFilters = () => {
    form.reset({
      fromTime: "",
      toTime: "",
    });

    // Clear URL parameters
    router.push("/history", { scroll: false });
  };

  // Handle next page
  const handleNextPage = () => {
    if (
      paginatedData &&
      paginatedData.hasMore &&
      paginatedData.reports.length > 0
    ) {
      const newLastId =
        paginatedData.reports[paginatedData.reports.length - 1].id;

      // Update URL with new lastId and increment page
      updateUrl({
        lastId: newLastId,
        fromTime: form.getValues("fromTime"),
        toTime: form.getValues("toTime"),
        page: currentPage + 1,
      });
    }
  };

  // Handle previous page
  const handlePreviousPage = () => {
    // For cursor-based pagination, going back is tricky
    // For simplicity, we'll just go back to the first page
    updateUrl({
      fromTime: form.getValues("fromTime"),
      toTime: form.getValues("toTime"),
      page: 1,
    });
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const handleReportSelect = (id: string) => {
    setSelectedReports((prev) => {
      if (prev.includes(id)) {
        return prev.filter((reportId) => reportId !== id);
      } else {
        if (prev.length < 2) {
          return [...prev, id];
        } else {
          toast.error("You can only select two reports for comparison");
          return prev;
        }
      }
    });
  };

  const handleCompare = () => {
    if (selectedReports.length !== 2) {
      toast.error("Please select exactly two reports to compare");
      return;
    }

    router.push(
      `/compare?report1=${selectedReports[0]}&report2=${selectedReports[1]}`
    );
  };

  // Generate pagination text
  const getPaginationText = () => {
    if (!paginatedData) return "";

    const { fromNumber, toNumber, totalCount } = paginatedData;

    // Check if filtering is applied
    const isFiltered = !!(fromTime || toTime);

    if (isFiltered) {
      return `${fromNumber}-${toNumber} of many recorded history`;
    } else {
      return `${fromNumber}-${toNumber} of ${totalCount} total records`;
    }
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Weather Report History</h1>
        <Button onClick={handleCompare} disabled={selectedReports.length !== 2}>
          Compare Selected Reports
        </Button>
      </div>

      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Filter Reports</CardTitle>
          <CardDescription>Filter reports by time range.</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmitFilter)}
              className="space-y-4"
            >
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="fromTime"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>From Date & Time</FormLabel>
                      <FormControl>
                        <Input
                          type="datetime-local"
                          placeholder="From date"
                          {...field}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="toTime"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>To Date & Time</FormLabel>
                      <FormControl>
                        <Input
                          type="datetime-local"
                          placeholder="To date"
                          {...field}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
              <div className="flex justify-end space-x-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={resetFilters}
                  disabled={loading}
                >
                  Reset
                </Button>
                <Button type="submit" disabled={loading}>
                  Apply Filters
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Historical Reports</CardTitle>
          <CardDescription>
            Select two reports to compare their weather data.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-4">Loading reports...</div>
          ) : !paginatedData || paginatedData.reports.length === 0 ? (
            <div className="text-center py-4">
              No reports found. Generate a report first.
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12">Select</TableHead>
                    <TableHead>Timestamp</TableHead>
                    <TableHead>Temperature (°C)</TableHead>
                    <TableHead>Pressure (hPa)</TableHead>
                    <TableHead>Humidity (%)</TableHead>
                    <TableHead>Cloud Cover (%)</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {paginatedData.reports.map((report) => (
                    <TableRow key={report.id}>
                      <TableCell>
                        <input
                          type="checkbox"
                          checked={selectedReports.includes(report.id)}
                          onChange={() => handleReportSelect(report.id)}
                          className="h-4 w-4 rounded border-gray-300"
                        />
                      </TableCell>
                      <TableCell>{formatDate(report.timestamp)}</TableCell>
                      <TableCell>{report.temperature}</TableCell>
                      <TableCell>{report.pressure}</TableCell>
                      <TableCell>{report.humidity}</TableCell>
                      <TableCell>{report.cloudCover}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </>
          )}
        </CardContent>
        <CardFooter className="flex justify-between">
          <div className="text-sm text-muted-foreground">
            {paginatedData && getPaginationText()}
          </div>
          <div className="flex space-x-2">
            <Button
              variant="outline"
              onClick={handlePreviousPage}
              disabled={loading || currentPage <= 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              onClick={handleNextPage}
              disabled={loading || !paginatedData?.hasMore}
            >
              Next
            </Button>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
