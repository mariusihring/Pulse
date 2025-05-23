import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { User } from '@/lib/types/crypto/dashboard/user';
import { DataTable } from '@/components/pulse/crypo/token-table/data-table';
import { columns } from '@/components/pulse/crypo/token-table/columns';

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

export default function Tokens({ user }: { user: User }) {
    console.log(user)
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />
            <div className="container mx-auto p-10">
                <DataTable columns={columns} data={user.tokens} />
            </div>
        </AppLayout>
    );
}
