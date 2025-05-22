import { PlaceholderPattern } from '@/components/ui/placeholder-pattern';
import AppLayout from '@/layouts/app-layout';
import {  type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { Suspense } from 'react';
import { Skeleton } from '@/components/ui/skeleton';
import DashboardOverview from '@/components/pulse/dashboard/overview';
import TokenPieChart from '@/components/pulse/dashboard/tokenpiechart';
import { SwapTable } from '@/components/pulse/dashboard/Tokenswaptable';
import { User } from '@/lib/types/crypto/dashboard/user';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Crypto Dashboard',
        href: '/crypto/dashboard',
    },
];

export default function Dashboard({ user }: {user: User}) {
    
    const allTokenSwaps = user.wallets.reduce((accumulator, wallet) => {
        //@ts-ignore
        return accumulator.concat(wallet.tokenswaps);
    }, []);
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />
            <div className="flex h-full flex-1 flex-col gap-4 rounded-xl p-4">
                <div className="grid auto-rows-min gap-4 md:grid-cols-2">
                    <div className="border-sidebar-border/70 dark:border-sidebar-border relative overflow-hidden rounded-xl border">
                        <Suspense fallback={
                            <Skeleton className="absolute inset-0 size-full stroke-neutral-900/20 dark:stroke-neutral-100/20 animate-pulse" />
                        }>
                            <DashboardOverview data={user} />
                        </Suspense>
                    </div>
                    <div className="border-sidebar-border/70 dark:border-sidebar-border relative overflow-hidden rounded-xl border">
                        <Suspense fallback={
                            <Skeleton className="absolute inset-0 size-full stroke-neutral-900/20 dark:stroke-neutral-100/20 animate-pulse" />
                        }>
                            <TokenPieChart data={user} />
                        </Suspense>
                    </div>

                </div>
                <div className="border-sidebar-border/70 dark:border-sidebar-border relative min-h-[100vh] flex-1 overflow-hidden rounded-xl border md:min-h-min p-4">
                    <Suspense fallback={
                        <Skeleton className="absolute inset-0 size-full stroke-neutral-900/20 dark:stroke-neutral-100/20 animate-pulse" />
                    }>
                        <SwapTable swaps={allTokenSwaps} />
                    </Suspense>
                </div>
            </div>
        </AppLayout>
    );
}
