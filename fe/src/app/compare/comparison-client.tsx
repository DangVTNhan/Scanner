'use client';

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
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

  if (loading) {
    return <div className="text-center py-8">Loading comparison data...</div>;
  }

  if (error || !comparison) {
    return (
      <div className="max-w-4xl mx-auto text-center py-8">
        <h1 className="text-3xl font-bold mb-4">Error</h1>
        <p className="mb-6">{error || "Failed to load comparison data"}</p>
        <Button asChild>
          <Link href="/history">Back to History</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Weather Report Comparison</h1>
        <Button asChild>
          <Link href="/history">Back to History</Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Comparison Results</CardTitle>
          <CardDescription>
            Comparing weather reports from{" "}
            {formatDate(comparison.report1.timestamp)} and{" "}
            {formatDate(comparison.report2.timestamp)}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Parameter</TableHead>
                <TableHead>Report 1</TableHead>
                <TableHead>Report 2</TableHead>
                <TableHead>Deviation</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell className="font-medium">Timestamp</TableCell>
                <TableCell>
                  {formatDate(comparison.report1.timestamp)}
                </TableCell>
                <TableCell>
                  {formatDate(comparison.report2.timestamp)}
                </TableCell>
                <TableCell>-</TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">Temperature (°C)</TableCell>
                <TableCell>{comparison.report1.temperature}</TableCell>
                <TableCell>{comparison.report2.temperature}</TableCell>
                <TableCell>
                  {comparison.deviation.temperature.toFixed(2)}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">Pressure (hPa)</TableCell>
                <TableCell>{comparison.report1.pressure}</TableCell>
                <TableCell>{comparison.report2.pressure}</TableCell>
                <TableCell>
                  {comparison.deviation.pressure.toFixed(2)}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">Humidity (%)</TableCell>
                <TableCell>{comparison.report1.humidity}</TableCell>
                <TableCell>{comparison.report2.humidity}</TableCell>
                <TableCell>
                  {comparison.deviation.humidity.toFixed(2)}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="font-medium">Cloud Cover (%)</TableCell>
                <TableCell>{comparison.report1.cloudCover}</TableCell>
                <TableCell>{comparison.report2.cloudCover}</TableCell>
                <TableCell>
                  {comparison.deviation.cloudCover.toFixed(2)}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
