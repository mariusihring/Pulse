import React from 'react';
import { Skeleton } from '@/components/ui/skeleton';

const WalletCardSkeleton: React.FC = () => {
  return (
    <div className="p-4 border rounded shadow-sm">
      <div className="flex items-center justify-between mb-4">
        <Skeleton className="h-6 w-1/3" />
        <Skeleton className="h-4 w-4 rounded-full" />
      </div>
      <div className="space-y-2">
        <Skeleton className="h-4 w-2/3" />
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-4 w-1/4" />
      </div>
      <div className="flex justify-end mt-4">
        <Skeleton className="h-8 w-20" />
      </div>
    </div>
  );
};

export default WalletCardSkeleton;
