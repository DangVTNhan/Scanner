'use client';

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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { compareReports, ComparisonResult } from "@/lib/api";
import { handleApiError } from "@/lib/api/utils";
import {
  AlertCircle,
  ArrowLeft,
  ArrowRight,
  ArrowUpDown,
  Calendar,
  CloudSun,
  Droplets,
  Gauge,
  History,
  Plane,
  ThermometerSun,
  Timer,
  TrendingDown,
  TrendingUp,
} from "lucide-react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function ComparisonClient() {
  const searchParams = useSearchParams();
  const report1Id = searchParams.get("report1");
  const report2Id = searchParams.get("report2");

  const [comparison, setComparison] = useState<ComparisonResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchComparison = async () => {
      if (!report1Id || !report2Id) {
        setError("Two report IDs are required for comparison");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const result = await compareReports({
          reportId1: report1Id,
          reportId2: report2Id,
        });
        setComparison(result);
      } catch (err) {
        handleApiError(err, "Failed to compare reports");
        setError("Failed to compare the selected reports");
      } finally {
        setLoading(false);
      }
    };

    fetchComparison();
  }, [report1Id, report2Id]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  // Helper function to determine trend icon and color based on deviation
  const getTrendIndicator = (value: number) => {
    if (value > 0) {
      return {
        icon: <TrendingUp className="size-4 text-green-500" />,
        color: "text-green-600 dark:text-green-400",
        bgColor: "bg-green-50 dark:bg-green-950/30",
      };
    } else if (value < 0) {
      return {
        icon: <TrendingDown className="size-4 text-red-500" />,
        color: "text-red-600 dark:text-red-400",
        bgColor: "bg-red-50 dark:bg-red-950/30",
      };
    } else {
      return {
        icon: <ArrowRight className="size-4 text-blue-500" />,
        color: "text-blue-600 dark:text-blue-400",
        bgColor: "bg-blue-50 dark:bg-blue-950/30",
      };
    }
  };

  // Helper function to get parameter icon
  const getParameterIcon = (parameter: string) => {
    switch (parameter) {
      case "Temperature":
        return <ThermometerSun className="size-5 text-orange-500" />;
      case "Pressure":
        return <Gauge className="size-5 text-blue-500" />;
      case "Humidity":
        return <Droplets className="size-5 text-blue-500" />;
      case "Cloud Cover":
        return <CloudSun className="size-5 text-gray-500" />;
      case "Timestamp":
        return <Calendar className="size-5 text-blue-500" />;
      default:
        return <ArrowUpDown className="size-5 text-blue-500" />;
    }
  };

  if (loading) {
    return (
      <div className="w-full max-w-6xl mx-auto px-4 py-12">
        <div className="flex flex-col items-center justify-center">
          <div className="size-12 border-4 border-blue-200 border-t-blue-500 rounded-full animate-spin mb-4"></div>
          <h2 className="text-xl font-medium text-slate-700 dark:text-slate-300 mb-2">
            Loading Comparison Data
          </h2>
          <p className="text-slate-500 dark:text-slate-400">
            Please wait while we analyze the weather reports...
          </p>
        </div>
      </div>
    );
  }

  if (error || !comparison) {
    return (
      <div className="w-full max-w-6xl mx-auto px-4 py-12">
        <div className="bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-900 rounded-xl p-8 text-center">
          <AlertCircle className="size-12 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-red-700 dark:text-red-400 mb-4">
            Error Loading Comparison
          </h1>
          <p className="mb-6 text-slate-700 dark:text-slate-300">
            {error || "Failed to load comparison data"}
          </p>
          <Button
            asChild
            className="bg-blue-500 hover:bg-blue-600 transition-all duration-200"
          >
            <Link href="/history" className="flex items-center gap-2">
              <ArrowLeft className="size-4" />
              Back to History
            </Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full max-w-6xl mx-auto px-4">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-950 dark:to-blue-900 mb-8 p-8 md:p-12">
        <div className="absolute inset-0 bg-grid-slate-200 [mask-image:linear-gradient(0deg,#fff,rgba(255,255,255,0.6))] dark:bg-grid-slate-700/25 dark:[mask-image:linear-gradient(0deg,rgba(255,255,255,0.1),rgba(255,255,255,0.5))]" />
        <div className="absolute -top-24 -right-20 opacity-20">
          <ArrowUpDown className="size-64 text-blue-500" />
        </div>
        <div className="relative flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="space-y-4 max-w-2xl">
            <div className="inline-flex items-center rounded-lg bg-blue-50/50 dark:bg-blue-900/30 px-3 py-1 text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
              <Plane className="mr-1 size-4" /> Changi Airport Weather System
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
              Weather Report Comparison
            </h1>
            <p className="text-slate-700 dark:text-slate-300 text-lg">
              Analyze and compare weather conditions between different time
              periods at Changi Airport.
            </p>
          </div>
          <div className="flex-shrink-0">
            <Button
              asChild
              className="bg-blue-500 hover:bg-blue-600 transition-all duration-200"
            >
              <Link href="/history" className="flex items-center gap-2">
                <History className="size-4" />
                Back to History
              </Link>
            </Button>
          </div>
        </div>
      </div>

      {/* Comparison Summary Card */}
      <Card className="mb-8 overflow-hidden border-blue-100 dark:border-blue-900">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ArrowUpDown className="size-5 text-blue-500" />
            Comparison Summary
          </CardTitle>
          <CardDescription>
            Comparing weather reports from{" "}
            <span className="font-medium">
              {formatDate(comparison.report1.timestamp)}
            </span>{" "}
            and{" "}
            <span className="font-medium">
              {formatDate(comparison.report2.timestamp)}
            </span>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div className="bg-slate-50 dark:bg-slate-900/50 rounded-lg p-4 border border-slate-200 dark:border-slate-800">
              <div className="flex items-center gap-2 mb-2">
                <Calendar className="size-5 text-blue-500" />
                <h3 className="font-medium">Report 1</h3>
              </div>
              <div className="text-sm text-slate-500 dark:text-slate-400 mb-3">
                {formatDate(comparison.report1.timestamp)}
              </div>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm">Temperature:</span>
                  <span className="font-medium">
                    {comparison.report1.temperature} °C
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Pressure:</span>
                  <span className="font-medium">
                    {comparison.report1.pressure} hPa
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Humidity:</span>
                  <span className="font-medium">
                    {comparison.report1.humidity} %
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Cloud Cover:</span>
                  <span className="font-medium">
                    {comparison.report1.cloudCover} %
                  </span>
                </div>
              </div>
            </div>

            <div className="bg-slate-50 dark:bg-slate-900/50 rounded-lg p-4 border border-slate-200 dark:border-slate-800">
              <div className="flex items-center gap-2 mb-2">
                <Calendar className="size-5 text-blue-500" />
                <h3 className="font-medium">Report 2</h3>
              </div>
              <div className="text-sm text-slate-500 dark:text-slate-400 mb-3">
                {formatDate(comparison.report2.timestamp)}
              </div>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm">Temperature:</span>
                  <span className="font-medium">
                    {comparison.report2.temperature} °C
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Pressure:</span>
                  <span className="font-medium">
                    {comparison.report2.pressure} hPa
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Humidity:</span>
                  <span className="font-medium">
                    {comparison.report2.humidity} %
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Cloud Cover:</span>
                  <span className="font-medium">
                    {comparison.report2.cloudCover} %
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="rounded-lg border overflow-hidden">
            <Table>
              <TableHeader className="bg-slate-50 dark:bg-slate-900">
                <TableRow>
                  <TableHead>Parameter</TableHead>
                  <TableHead>Report 1</TableHead>
                  <TableHead>Report 2</TableHead>
                  <TableHead>Deviation</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow>
                  <TableCell className="font-medium flex items-center gap-2">
                    {getParameterIcon("Timestamp")}
                    Timestamp
                  </TableCell>
                  <TableCell>
                    {formatDate(comparison.report1.timestamp)}
                  </TableCell>
                  <TableCell>
                    {formatDate(comparison.report2.timestamp)}
                  </TableCell>
                  <TableCell>-</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium flex items-center gap-2">
                    {getParameterIcon("Temperature")}
                    Temperature (°C)
                  </TableCell>
                  <TableCell>{comparison.report1.temperature}</TableCell>
                  <TableCell>{comparison.report2.temperature}</TableCell>
                  <TableCell>
                    {(() => {
                      const trend = getTrendIndicator(
                        comparison.deviation.temperature
                      );
                      return (
                        <div
                          className={`inline-flex items-center gap-1 px-2 py-1 rounded-full ${trend.bgColor} ${trend.color}`}
                        >
                          {trend.icon}
                          {Math.abs(comparison.deviation.temperature).toFixed(
                            2
                          )}
                        </div>
                      );
                    })()}
                  </TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium flex items-center gap-2">
                    {getParameterIcon("Pressure")}
                    Pressure (hPa)
                  </TableCell>
                  <TableCell>{comparison.report1.pressure}</TableCell>
                  <TableCell>{comparison.report2.pressure}</TableCell>
                  <TableCell>
                    {(() => {
                      const trend = getTrendIndicator(
                        comparison.deviation.pressure
                      );
                      return (
                        <div
                          className={`inline-flex items-center gap-1 px-2 py-1 rounded-full ${trend.bgColor} ${trend.color}`}
                        >
                          {trend.icon}
                          {Math.abs(comparison.deviation.pressure).toFixed(2)}
                        </div>
                      );
                    })()}
                  </TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium flex items-center gap-2">
                    {getParameterIcon("Humidity")}
                    Humidity (%)
                  </TableCell>
                  <TableCell>{comparison.report1.humidity}</TableCell>
                  <TableCell>{comparison.report2.humidity}</TableCell>
                  <TableCell>
                    {(() => {
                      const trend = getTrendIndicator(
                        comparison.deviation.humidity
                      );
                      return (
                        <div
                          className={`inline-flex items-center gap-1 px-2 py-1 rounded-full ${trend.bgColor} ${trend.color}`}
                        >
                          {trend.icon}
                          {Math.abs(comparison.deviation.humidity).toFixed(2)}
                        </div>
                      );
                    })()}
                  </TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium flex items-center gap-2">
                    {getParameterIcon("Cloud Cover")}
                    Cloud Cover (%)
                  </TableCell>
                  <TableCell>{comparison.report1.cloudCover}</TableCell>
                  <TableCell>{comparison.report2.cloudCover}</TableCell>
                  <TableCell>
                    {(() => {
                      const trend = getTrendIndicator(
                        comparison.deviation.cloudCover
                      );
                      return (
                        <div
                          className={`inline-flex items-center gap-1 px-2 py-1 rounded-full ${trend.bgColor} ${trend.color}`}
                        >
                          {trend.icon}
                          {Math.abs(comparison.deviation.cloudCover).toFixed(2)}
                        </div>
                      );
                    })()}
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </CardContent>
        <CardFooter className="border-t py-4">
          <div className="text-xs text-muted-foreground">
            <div className="flex items-center gap-1">
              <Timer className="size-3" />
              Report 1 generated at {formatDate(comparison.report1.createdAt)}
            </div>
            <div className="flex items-center gap-1 mt-1">
              <Timer className="size-3" />
              Report 2 generated at {formatDate(comparison.report2.createdAt)}
            </div>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
