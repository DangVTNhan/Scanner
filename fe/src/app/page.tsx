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
      console.error("Failed to generate report:", error);
      toast.error("Failed to generate weather report");
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Generate Weather Report</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Generate New Report</CardTitle>
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
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Leave empty to use the current time.
                      </FormDescription>
                    </FormItem>
                  )}
                />
                <Button type="submit" disabled={loading}>
                  {loading ? "Generating..." : "Generate Report"}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>

        {report && (
          <Card>
            <CardHeader>
              <CardTitle>Weather Report</CardTitle>
              <CardDescription>
                Weather data for Changi Airport at{" "}
                {formatDate(report.timestamp)}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-2">
                  <div className="text-sm font-medium">Temperature:</div>
                  <div>{report.temperature} °C</div>

                  <div className="text-sm font-medium">Pressure:</div>
                  <div>{report.pressure} hPa</div>

                  <div className="text-sm font-medium">Humidity:</div>
                  <div>{report.humidity} %</div>

                  <div className="text-sm font-medium">Cloud Cover:</div>
                  <div>{report.cloudCover} %</div>
                </div>
              </div>
            </CardContent>
            <CardFooter>
              <div className="text-xs text-muted-foreground">
                Report generated at {formatDate(report.createdAt)}
              </div>
            </CardFooter>
          </Card>
        )}
      </div>
    </div>
  );
}
