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
import { generateReport, WeatherReport } from "@/lib/api";
import { handleApiError } from "@/lib/api/utils";
import {
  ArrowRight,
  CalendarClock,
  Cloud,
  CloudFog,
  CloudSun,
  Droplets,
  Gauge,
  Plane,
  Sun,
} from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

export default function Home() {
  const [report, setReport] = useState<WeatherReport | null>(null);
  const [loading, setLoading] = useState(false);

  const form = useForm({
    defaultValues: {
      timestamp: "",
    },
  });

  const onSubmit = async (data: { timestamp: string }) => {
    try {
      setLoading(true);
      const timestamp = data.timestamp
        ? new Date(data.timestamp).toISOString()
        : undefined;
      const newReport = await generateReport({ timestamp });
      setReport(newReport);
      toast.success("Weather report generated successfully");
    } catch (error) {
      handleApiError(error, "Failed to generate weather report");
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  // Function to get weather icon based on cloud cover
  const getWeatherIcon = (cloudCover: number) => {
    if (cloudCover < 20) return <Sun className="size-16 text-yellow-500" />;
    if (cloudCover < 50) return <CloudSun className="size-16 text-blue-400" />;
    if (cloudCover < 80) return <Cloud className="size-16 text-gray-400" />;
    return <CloudFog className="size-16 text-gray-500" />;
  };

  return (
    <div className="w-full mx-auto">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-950 dark:to-blue-900 mb-8 p-8 md:p-12">
        <div className="absolute inset-0 bg-grid-slate-200 [mask-image:linear-gradient(0deg,#fff,rgba(255,255,255,0.6))] dark:bg-grid-slate-700/25 dark:[mask-image:linear-gradient(0deg,rgba(255,255,255,0.1),rgba(255,255,255,0.5))]" />
        <div className="absolute -top-24 -right-20 opacity-20">
          <CloudSun className="size-64 text-blue-500" />
        </div>
        <div className="relative flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="space-y-4 max-w-2xl">
            <div className="inline-flex items-center rounded-lg bg-blue-50/50 dark:bg-blue-900/30 px-3 py-1 text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
              <Plane className="mr-1 size-4" /> Changi Airport Weather System
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-slate-900 dark:text-slate-100">
              Real-time Weather Reports
            </h1>
            <p className="text-slate-700 dark:text-slate-300 text-lg">
              Generate accurate weather reports for Changi Airport. View current
              conditions or check historical data with just a few clicks.
            </p>
          </div>
          <div className="flex-shrink-0 hidden md:block">
            <div className="relative bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm rounded-xl p-4 shadow-lg">
              <div className="flex items-center justify-center">
                <CloudSun className="size-24 text-blue-500" />
              </div>
              <div className="text-center mt-2">
                <div className="text-sm font-medium text-slate-500 dark:text-slate-400">
                  Changi Airport
                </div>
                <div className="text-2xl font-bold">Singapore</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid gap-8 lg:grid-cols-2">
        {/* Form Card */}
        <Card className="overflow-hidden border-blue-100 dark:border-blue-900">
          <div className="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-blue-100 to-blue-50 dark:from-blue-900 dark:to-blue-950 rounded-bl-full opacity-50" />
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CalendarClock className="size-5 text-blue-500" />
              Generate New Report
            </CardTitle>
            <CardDescription>
              Generate a weather report for Changi Airport for the current time
              or a specific date.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="timestamp"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Date & Time (Optional)</FormLabel>
                      <FormControl>
                        <Input
                          type="datetime-local"
                          placeholder="Select date and time"
                          className="focus-visible:ring-blue-500/20"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Leave empty to use the current time.
                      </FormDescription>
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 transition-all duration-200"
                >
                  {loading ? (
                    <span className="flex items-center gap-2">
                      <div className="size-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      Generating...
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      Generate Report
                      <ArrowRight className="size-4" />
                    </span>
                  )}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>

        {/* Weather Report Card */}
        {report ? (
          <Card className="overflow-hidden border-blue-100 dark:border-blue-900 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-blue-100 to-blue-50 dark:from-blue-900 dark:to-blue-950 rounded-bl-full opacity-50" />
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CloudSun className="size-5 text-blue-500" />
                Weather Report
              </CardTitle>
              <CardDescription>
                Weather data for Changi Airport at{" "}
                {formatDate(report.timestamp)}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col md:flex-row gap-6 items-center">
                <div className="flex-shrink-0 flex flex-col items-center">
                  {getWeatherIcon(report.cloudCover)}
                  <div className="text-3xl font-bold mt-2">
                    {report.temperature}°C
                  </div>
                  <div className="text-sm text-muted-foreground">
                    Changi Airport
                  </div>
                </div>
                <div className="flex-grow w-full">
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <div className="bg-slate-50 dark:bg-slate-900 p-4 rounded-lg flex flex-col items-center">
                      <Gauge className="size-8 text-blue-500 mb-2" />
                      <div className="text-sm font-medium text-muted-foreground">
                        Pressure
                      </div>
                      <div className="text-xl font-semibold">
                        {report.pressure} hPa
                      </div>
                    </div>
                    <div className="bg-slate-50 dark:bg-slate-900 p-4 rounded-lg flex flex-col items-center">
                      <Droplets className="size-8 text-blue-500 mb-2" />
                      <div className="text-sm font-medium text-muted-foreground">
                        Humidity
                      </div>
                      <div className="text-xl font-semibold">
                        {report.humidity}%
                      </div>
                    </div>
                    <div className="bg-slate-50 dark:bg-slate-900 p-4 rounded-lg flex flex-col items-center">
                      <Cloud className="size-8 text-blue-500 mb-2" />
                      <div className="text-sm font-medium text-muted-foreground">
                        Cloud Cover
                      </div>
                      <div className="text-xl font-semibold">
                        {report.cloudCover}%
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
            <CardFooter className="border-t">
              <div className="text-xs text-muted-foreground">
                Report generated at {formatDate(report.createdAt)}
              </div>
            </CardFooter>
          </Card>
        ) : (
          <Card className="border-dashed border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
            <div className="flex flex-col items-center justify-center h-full py-12">
              <CloudSun className="size-16 text-slate-300 dark:text-slate-700 mb-4" />
              <h3 className="text-xl font-medium text-slate-700 dark:text-slate-300">
                No Report Generated
              </h3>
              <p className="text-sm text-slate-500 dark:text-slate-400 text-center max-w-xs mt-2">
                Generate a weather report using the form to see the current
                weather conditions at Changi Airport.
              </p>
            </div>
          </Card>
        )}
      </div>

      {/* Additional Information Section */}
      <div className="mt-12 bg-slate-50 dark:bg-slate-900/50 rounded-xl p-6 border border-slate-200 dark:border-slate-800">
        <h2 className="text-xl font-semibold mb-4">
          About Changi Airport Weather System
        </h2>
        <p className="text-slate-700 dark:text-slate-300 mb-4">
          This system provides real-time and historical weather data for Changi
          Airport (Lat: 1.3586° N, Long: 103.9899° E) using the OpenWeather API.
          You can generate reports, view historical data, and compare weather
          conditions over time.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
          <div className="flex items-start gap-3">
            <div className="bg-blue-100 dark:bg-blue-900/30 p-2 rounded-lg">
              <CloudSun className="size-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h3 className="font-medium">Real-time Data</h3>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Get current weather conditions with a single click
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="bg-blue-100 dark:bg-blue-900/30 p-2 rounded-lg">
              <CalendarClock className="size-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h3 className="font-medium">Historical Reports</h3>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Access weather data from specific dates and times
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="bg-blue-100 dark:bg-blue-900/30 p-2 rounded-lg">
              <ArrowRight className="size-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h3 className="font-medium">View History</h3>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Browse all previously generated reports
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
