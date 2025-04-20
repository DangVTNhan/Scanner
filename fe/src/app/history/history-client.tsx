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
  FormDescription,
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
import { handleApiError } from "@/lib/api/utils";
import {
  ArrowRight,
  CalendarClock,
  CalendarDays,
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  CloudSun,
  Filter,
  History,
  Plane,
  RefreshCw,
  Search,
} from "lucide-react";
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
      handleApiError(error, "Failed to load weather reports");
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
    <div className="w-full max-w-6xl mx-auto px-4">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-950 dark:to-blue-900 mb-8 p-8 md:p-12">
        <div className="absolute inset-0 bg-grid-slate-200 [mask-image:linear-gradient(0deg,#fff,rgba(255,255,255,0.6))] dark:bg-grid-slate-700/25 dark:[mask-image:linear-gradient(0deg,rgba(255,255,255,0.1),rgba(255,255,255,0.5))]" />
        <div className="absolute -top-24 -right-20 opacity-20">
          <History className="size-64 text-blue-500" />
        </div>
        <div className="relative flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="space-y-4 max-w-2xl">
            <div className="inline-flex items-center rounded-lg bg-blue-50/50 dark:bg-blue-900/30 px-3 py-1 text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
              <Plane className="mr-1 size-4" /> Changi Airport Weather System
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
              Weather Report History
            </h1>
          </div>
        </div>
      </div>

      {/* Filter Card */}
      <Card className="mb-8 overflow-hidden border-blue-100 dark:border-blue-900">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="size-5 text-blue-500" />
            Filter Reports
          </CardTitle>
          <CardDescription>
            Filter reports by time range to find specific weather data.
          </CardDescription>
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
                      <FormLabel className="flex items-center gap-1">
                        <CalendarClock className="size-4 text-blue-500" />
                        From Date & Time
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="datetime-local"
                          placeholder="From date"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Start date for filtering reports
                      </FormDescription>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="toTime"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="flex items-center gap-1">
                        <CalendarClock className="size-4 text-blue-500" />
                        To Date & Time
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="datetime-local"
                          placeholder="To date"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        End date for filtering reports
                      </FormDescription>
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
                  className="flex items-center gap-1"
                >
                  <RefreshCw className="size-4" />
                  Reset
                </Button>
                <Button
                  type="submit"
                  disabled={loading}
                  className="bg-blue-500 hover:bg-blue-600 transition-all duration-200 flex items-center gap-1"
                >
                  <Search className="size-4" />
                  Apply Filters
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>

      {/* Reports Table Card */}
      <Card className="overflow-hidden border-blue-100 dark:border-blue-900">
        <CardHeader className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <CardTitle className="flex items-center gap-2">
              <CloudSun className="size-5 text-blue-500" />
              Historical Reports
            </CardTitle>
            <CardDescription>
              Select two reports to compare their weather data.
            </CardDescription>
          </div>
          <div className="hidden md:block">
            <Button
              onClick={handleCompare}
              disabled={selectedReports.length !== 2}
              className="bg-blue-500 hover:bg-blue-600 transition-all duration-200"
            >
              <span className="flex items-center gap-2">
                Compare Selected Reports ({selectedReports.length}/2)
                <ArrowRight className="size-4" />
              </span>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-12">
              <div className="size-8 border-4 border-blue-200 border-t-blue-500 rounded-full animate-spin mx-auto mb-4"></div>
              <p className="text-slate-600 dark:text-slate-400">
                Loading reports...
              </p>
            </div>
          ) : !paginatedData || paginatedData.reports.length === 0 ? (
            <div className="text-center py-12 bg-slate-50/50 dark:bg-slate-900/50 rounded-lg border border-dashed border-slate-200 dark:border-slate-800">
              <CalendarDays className="size-12 text-slate-300 dark:text-slate-700 mx-auto mb-4" />
              <h3 className="text-xl font-medium text-slate-700 dark:text-slate-300 mb-2">
                No Reports Found
              </h3>
              <p className="text-sm text-slate-500 dark:text-slate-400 max-w-md mx-auto">
                No weather reports match your criteria. Try adjusting your
                filters or generate a new report first.
              </p>
            </div>
          ) : (
            <>
              <div className="rounded-lg border overflow-hidden">
                <Table>
                  <TableHeader className="bg-slate-50 dark:bg-slate-900">
                    <TableRow>
                      <TableHead className="w-12 text-center">Select</TableHead>
                      <TableHead>Timestamp</TableHead>
                      <TableHead>Temperature (°C)</TableHead>
                      <TableHead>Pressure (hPa)</TableHead>
                      <TableHead>Humidity (%)</TableHead>
                      <TableHead>Cloud Cover (%)</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {paginatedData.reports.map((report) => (
                      <TableRow
                        key={report.id}
                        className={
                          selectedReports.includes(report.id)
                            ? "bg-blue-50/50 dark:bg-blue-900/20"
                            : ""
                        }
                      >
                        <TableCell className="text-center">
                          <button
                            type="button"
                            onClick={() => handleReportSelect(report.id)}
                            className={`flex items-center justify-center size-6 rounded-full mx-auto ${
                              selectedReports.includes(report.id)
                                ? "bg-blue-500 text-white"
                                : "border border-slate-300 dark:border-slate-700 text-transparent hover:border-blue-500 dark:hover:border-blue-500"
                            }`}
                          >
                            <CheckCircle2 className="size-4" />
                          </button>
                        </TableCell>
                        <TableCell className="font-medium">
                          {formatDate(report.timestamp)}
                        </TableCell>
                        <TableCell>{report.temperature}</TableCell>
                        <TableCell>{report.pressure}</TableCell>
                        <TableCell>{report.humidity}</TableCell>
                        <TableCell>{report.cloudCover}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </>
          )}
        </CardContent>
        <CardFooter className="flex justify-between border-t py-4">
          <div className="text-sm text-muted-foreground">
            {paginatedData && getPaginationText()}
          </div>
          <div className="flex space-x-2">
            <Button
              variant="outline"
              onClick={handlePreviousPage}
              disabled={loading || currentPage <= 1}
              className="flex items-center gap-1"
            >
              <ChevronLeft className="size-4" />
              Previous
            </Button>
            <Button
              variant="outline"
              onClick={handleNextPage}
              disabled={loading || !paginatedData?.hasMore}
              className="flex items-center gap-1"
            >
              Next
              <ChevronRight className="size-4" />
            </Button>
          </div>
        </CardFooter>
      </Card>

      {/* Mobile Compare Button */}
      <div className="md:hidden mt-6 mb-8">
        <Button
          onClick={handleCompare}
          disabled={selectedReports.length !== 2}
          className="w-full bg-blue-500 hover:bg-blue-600 transition-all duration-200"
        >
          <span className="flex items-center gap-2 justify-center">
            Compare Selected Reports ({selectedReports.length}/2)
            <ArrowRight className="size-4" />
          </span>
        </Button>
      </div>
    </div>
  );
}
