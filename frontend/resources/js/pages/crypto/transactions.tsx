import { PlaceholderPattern } from '@/components/ui/placeholder-pattern';
import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { Suspense } from 'react';
import { Skeleton } from '@/components/ui/skeleton';
import DashboardOverview from '@/components/pulse/dashboard/overview';
import TokenPieChart from '@/components/pulse/dashboard/tokenpiechart';
import { SwapTable } from '@/components/pulse/dashboard/Tokenswaptable';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Crypto',
        href: '/crypto',
    },
    {
        title: 'Transactions',
        href: '/crypto/transactions',
    },
];

export default function Dashboard({ user }) {

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />
            <div>Transactions</div>
        </AppLayout>
    );
}
