import { Suspense } from 'react';
import HistoryClient from './history-client';

export default function HistoryPage() {
  return (
    <Suspense fallback={<div className="text-center py-8">Loading history page...</div>}>
      <HistoryClient />
    </Suspense>
  );
}
