import { PlaceholderPattern } from '@/components/ui/placeholder-pattern';
import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { Suspense } from 'react';
import { Skeleton } from '@/components/ui/skeleton';
import DashboardOverview from '@/components/pulse/dashboard/overview';
import TokenPieChart from '@/components/pulse/dashboard/tokenpiechart';
import { SwapTable } from '@/components/pulse/dashboard/Tokenswaptable';
import { User } from '@/lib/types/crypto/dashboard/user';
import { DataTable } from '@/components/pulse/crypo/tokenswap-table/data-table';
import { columns } from '@/components/pulse/crypo/tokenswap-table/columns';

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

export default function Dashboard({ user }: {user: User}) {
    console.log(user)
    const tokenswaps = user.wallets.reduce((accumulator, wallet) => accumulator.concat(wallet.tokenswaps), [])
    
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />
            <div className="container mx-auto p-10">
                <DataTable columns={columns} data={tokenswaps} />
            </div>
        </AppLayout>
    );
}
