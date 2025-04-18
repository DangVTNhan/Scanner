'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { toast } from 'sonner';
import { getAllReports, WeatherReport } from '@/lib/api';

export default function HistoryPage() {
  const [reports, setReports] = useState<WeatherReport[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedReports, setSelectedReports] = useState<string[]>([]);
  const router = useRouter();

  useEffect(() => {
    const fetchReports = async () => {
      try {
        setLoading(true);
        const data = await getAllReports();
        setReports(data);
      } catch (error) {
        console.error('Failed to fetch reports:', error);
        toast.error('Failed to load weather reports');
      } finally {
        setLoading(false);
      }
    };

    fetchReports();
  }, []);

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
          toast.error('You can only select two reports for comparison');
          return prev;
        }
      }
    });
  };

  const handleCompare = () => {
    if (selectedReports.length !== 2) {
      toast.error('Please select exactly two reports to compare');
      return;
    }

    router.push(`/compare?report1=${selectedReports[0]}&report2=${selectedReports[1]}`);
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Weather Report History</h1>
        <Button 
          onClick={handleCompare} 
          disabled={selectedReports.length !== 2}
        >
          Compare Selected Reports
        </Button>
      </div>

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
          ) : reports.length === 0 ? (
            <div className="text-center py-4">No reports found. Generate a report first.</div>
          ) : (
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
                {reports.map((report) => (
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
          )}
        </CardContent>
      </Card>
    </div>
  );
}
