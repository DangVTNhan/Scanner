import { Suspense } from "react";
import ComparisonClient from "./comparison-client";

export default function ComparePage() {
  return (
    <Suspense
      fallback={
        <div className="text-center py-8">Loading comparison data...</div>
      }
    >
      <ComparisonClient />
    </Suspense>
  );
}
