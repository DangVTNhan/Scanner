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
import { SortOrder } from "@/lib/api/types";
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
  const [currentOffset, setCurrentOffset] = useState(0);
  const [sortBy, setSortBy] = useState<string>("timestamp");
  const [sortOrder, setSortOrder] = useState<SortOrder>("desc");

  const router = useRouter();
  const searchParams = useSearchParams();

  // Get URL parameters
  const offsetStr = searchParams.get("offset") || "0";
  const offset = parseInt(offsetStr, 10);
  const fromTime = searchParams.get("fromTime") || "";
  const toTime = searchParams.get("toTime") || "";
  const sortByParam = searchParams.get("sortBy") || "timestamp";
  const sortOrderParam = (searchParams.get("sortOrder") as SortOrder) || "desc";

  // Form for filters
  const form = useForm<FilterForm>({
    defaultValues: {
      fromTime,
      toTime,
    },
  });

  // Set form values and sort state from URL parameters on initial load
  useEffect(() => {
    form.reset({
      fromTime,
      toTime,
    });
    setSortBy(sortByParam);
    setSortOrder(sortOrderParam);
  }, [form, fromTime, toTime, sortByParam, sortOrderParam]);

  const limit = 10; // Number of items per page

  // Update URL with current filter, sorting, and pagination state
  const updateUrl = useCallback(
    (params: {
      offset?: number;
      fromTime?: string;
      toTime?: string;
      sortBy?: string;
      sortOrder?: SortOrder;
    }) => {
      const newParams = new URLSearchParams();

      if (params.offset && params.offset > 0) {
        newParams.set("offset", params.offset.toString());
      }

      if (params.fromTime) {
        newParams.set("fromTime", params.fromTime);
      }

      if (params.toTime) {
        newParams.set("toTime", params.toTime);
      }

      if (params.sortBy) {
        newParams.set("sortBy", params.sortBy);
      }

      if (params.sortOrder) {
        newParams.set("sortOrder", params.sortOrder);
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
      const params: any = {
        limit,
        sortBy,
        sortOrder,
      };

      // Set offset for pagination
      if (offset > 0) {
        params.offset = offset;
      }

      // Update current offset state
      setCurrentOffset(offset);

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
  }, [offset, fromTime, toTime, limit, sortBy, sortOrder]);

  // Initial data fetch
  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  // Handle filter submission
  const onSubmitFilter = (data: FilterForm) => {
    // Update URL with filter parameters and reset pagination
    updateUrl({
      offset: 0, // Reset to first page when applying filters
      fromTime: data.fromTime,
      toTime: data.toTime,
      sortBy,
      sortOrder,
    });
  };

  // Handle filter reset
  const resetFilters = () => {
    form.reset({
      fromTime: "",
      toTime: "",
    });

    // Clear URL parameters and reset pagination state
    router.push("/history", { scroll: false });
  };

  // Handle next page
  const handleNextPage = () => {
    const info = getPaginationInfo();
    if (info && info.currentPage < info.totalPages) {
      // Calculate the new offset for the next page
      const newOffset = currentOffset + limit;

      // Update URL with new offset
      updateUrl({
        offset: newOffset,
        fromTime: form.getValues("fromTime"),
        toTime: form.getValues("toTime"),
        sortBy,
        sortOrder,
      });
    }
  };

  // Handle previous page
  const handlePreviousPage = () => {
    if (currentOffset > 0) {
      // Calculate the new offset for the previous page
      const newOffset = Math.max(0, currentOffset - limit);

      // Update URL with new offset
      updateUrl({
        offset: newOffset,
        fromTime: form.getValues("fromTime"),
        toTime: form.getValues("toTime"),
        sortBy,
        sortOrder,
      });
    } else {
      // If we're already at the first page, just refresh
      updateUrl({
        fromTime: form.getValues("fromTime"),
        toTime: form.getValues("toTime"),
        sortBy,
        sortOrder,
      });
    }
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

  // Calculate pagination information
  const getPaginationInfo = () => {
    if (!paginatedData) return null;

    const totalItems = paginatedData.totalCount;
    const totalPages = Math.ceil(totalItems / limit);
    const currentPage = Math.floor(currentOffset / limit) + 1;

    return {
      totalItems,
      totalPages,
      currentPage,
      currentOffset,
      itemsPerPage: limit,
      displayedItems: paginatedData.reports.length,
    };
  };

  // Generate pagination text
  const getPaginationText = () => {
    const info = getPaginationInfo();
    if (!info) return "";

    const start = currentOffset + 1;
    const end = currentOffset + info.displayedItems;

    return `${start}-${end} of ${info.totalItems} records`;
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
          ) : !paginatedData ||
            !paginatedData.reports ||
            paginatedData.reports.length === 0 ? (
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
                      <TableHead>
                        <div
                          className="flex items-center gap-1 cursor-pointer"
                          onClick={() => {
                            const newSortOrder =
                              sortBy === "timestamp" && sortOrder === "desc"
                                ? "asc"
                                : "desc";
                            setSortBy("timestamp");
                            setSortOrder(newSortOrder as SortOrder);
                            updateUrl({
                              offset: 0,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy: "timestamp",
                              sortOrder: newSortOrder as SortOrder,
                            });
                          }}
                        >
                          Timestamp
                          {sortBy === "timestamp" && (
                            <span className="ml-1">
                              {sortOrder === "desc" ? "↓" : "↑"}
                            </span>
                          )}
                        </div>
                      </TableHead>
                      <TableHead>
                        <div
                          className="flex items-center gap-1 cursor-pointer"
                          onClick={() => {
                            const newSortOrder =
                              sortBy === "temperature" && sortOrder === "desc"
                                ? "asc"
                                : "desc";
                            setSortBy("temperature");
                            setSortOrder(newSortOrder as SortOrder);
                            updateUrl({
                              offset: 0,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy: "temperature",
                              sortOrder: newSortOrder as SortOrder,
                            });
                          }}
                        >
                          Temperature (°C)
                          {sortBy === "temperature" && (
                            <span className="ml-1">
                              {sortOrder === "desc" ? "↓" : "↑"}
                            </span>
                          )}
                        </div>
                      </TableHead>
                      <TableHead>
                        <div
                          className="flex items-center gap-1 cursor-pointer"
                          onClick={() => {
                            const newSortOrder =
                              sortBy === "pressure" && sortOrder === "desc"
                                ? "asc"
                                : "desc";
                            setSortBy("pressure");
                            setSortOrder(newSortOrder as SortOrder);
                            updateUrl({
                              offset: 0,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy: "pressure",
                              sortOrder: newSortOrder as SortOrder,
                            });
                          }}
                        >
                          Pressure (hPa)
                          {sortBy === "pressure" && (
                            <span className="ml-1">
                              {sortOrder === "desc" ? "↓" : "↑"}
                            </span>
                          )}
                        </div>
                      </TableHead>
                      <TableHead>
                        <div
                          className="flex items-center gap-1 cursor-pointer"
                          onClick={() => {
                            const newSortOrder =
                              sortBy === "humidity" && sortOrder === "desc"
                                ? "asc"
                                : "desc";
                            setSortBy("humidity");
                            setSortOrder(newSortOrder as SortOrder);
                            updateUrl({
                              offset: 0,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy: "humidity",
                              sortOrder: newSortOrder as SortOrder,
                            });
                          }}
                        >
                          Humidity (%)
                          {sortBy === "humidity" && (
                            <span className="ml-1">
                              {sortOrder === "desc" ? "↓" : "↑"}
                            </span>
                          )}
                        </div>
                      </TableHead>
                      <TableHead>
                        <div
                          className="flex items-center gap-1 cursor-pointer"
                          onClick={() => {
                            const newSortOrder =
                              sortBy === "cloudCover" && sortOrder === "desc"
                                ? "asc"
                                : "desc";
                            setSortBy("cloudCover");
                            setSortOrder(newSortOrder as SortOrder);
                            updateUrl({
                              offset: 0,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy: "cloudCover",
                              sortOrder: newSortOrder as SortOrder,
                            });
                          }}
                        >
                          Cloud Cover (%)
                          {sortBy === "cloudCover" && (
                            <span className="ml-1">
                              {sortOrder === "desc" ? "↓" : "↑"}
                            </span>
                          )}
                        </div>
                      </TableHead>
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
        <CardFooter className="flex flex-col md:flex-row justify-between items-center border-t py-4 gap-4">
          <div className="text-sm text-muted-foreground flex flex-col md:flex-row md:items-center gap-2">
            <div>{paginatedData && getPaginationText()}</div>
          </div>

          {/* Pagination Controls */}
          <div className="flex flex-wrap items-center justify-center gap-2">
            {/* Previous Button */}
            <Button
              variant="outline"
              size="sm"
              onClick={handlePreviousPage}
              disabled={loading || currentOffset <= 0}
              className="flex items-center gap-1"
            >
              <ChevronLeft className="size-3" />
              Prev
            </Button>

            {/* Page Numbers */}
            {paginatedData && getPaginationInfo() && (
              <div className="flex items-center gap-1">
                {(() => {
                  const info = getPaginationInfo()!;
                  const { currentPage, totalPages } = info;
                  const pageButtons = [];

                  // Function to create a page button
                  const createPageButton = (
                    pageNum: number,
                    active: boolean = false
                  ) => {
                    return (
                      <Button
                        key={`page-${pageNum}`}
                        variant={active ? "default" : "outline"}
                        size="sm"
                        className={`size-8 p-0 ${
                          active ? "bg-blue-500 hover:bg-blue-600" : ""
                        }`}
                        onClick={() => {
                          if (pageNum !== currentPage) {
                            const newOffset = (pageNum - 1) * limit;
                            updateUrl({
                              offset: newOffset,
                              fromTime: form.getValues("fromTime"),
                              toTime: form.getValues("toTime"),
                              sortBy,
                              sortOrder,
                            });
                          }
                        }}
                      >
                        {pageNum}
                      </Button>
                    );
                  };

                  // Show first page if we're not on it
                  if (currentPage > 1) {
                    pageButtons.push(createPageButton(1));

                    // Add ellipsis if there's a gap
                    if (currentPage > 2) {
                      pageButtons.push(
                        <div key="ellipsis-start" className="px-2">
                          ...
                        </div>
                      );
                    }
                  }

                  // Always show current page
                  pageButtons.push(createPageButton(currentPage, true));

                  // Show next 2 pages if they exist
                  for (let i = 1; i <= 2; i++) {
                    const pageNum = currentPage + i;
                    if (pageNum <= totalPages) {
                      pageButtons.push(createPageButton(pageNum));
                    }
                  }

                  // If there are more pages after the 2 we've shown, add ellipsis and last page
                  if (currentPage + 3 <= totalPages) {
                    pageButtons.push(
                      <div key="ellipsis" className="px-2">
                        ...
                      </div>
                    );
                    pageButtons.push(createPageButton(totalPages));
                  }

                  return pageButtons;
                })()}
              </div>
            )}

            {/* Next Button */}
            <Button
              variant="outline"
              size="sm"
              onClick={handleNextPage}
              disabled={
                loading ||
                !paginatedData ||
                getPaginationInfo()?.currentPage ===
                  getPaginationInfo()?.totalPages
              }
              className="flex items-center gap-1"
            >
              Next
              <ChevronRight className="size-3" />
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
