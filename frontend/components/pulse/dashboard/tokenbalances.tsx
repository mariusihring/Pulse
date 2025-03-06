import Image from "next/image"
import {useWalletStore} from "@/lib/providers/wallet_provider";
import type {Token} from "@/lib/gql/graphql";
import {
    type ColumnDef,
    flexRender,
    getCoreRowModel,
    getPaginationRowModel,
    useReactTable,
} from "@tanstack/react-table"
import {useState} from "react";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
    Pagination,
    PaginationContent,
    PaginationItem,
    PaginationNext,
    PaginationPrevious,
} from "@/components/ui/pagination"

interface TokenBalancesProps {
    className?: string
}

const columns: ColumnDef<Token>[] = [
    {
        accessoryKey: "name",
        header: "Token",
        cell: ({row}) => {
            return (
           <div className="flex items-center">
               <Image src={row.original.image} alt={row.original.name} width={24} height={24} className="mr-2 rounded-full" unoptimized={true} />
               <span className="font-medium text-gray-900 dark:text-white">{row.original.name}</span>
           </div>

            )
        }

    },
    {
        accessoryKey: "amount",
        header: "Amount",
        cell: ({row}) => {
            return <span>{row.original.amount.toFixed(2)}</span>
        }
    },
    {
        accessoryKey: "price",
        header: "Price",
        cell: ({ row }) => {
            return <span>${row.original.price.toFixed(4)}</span>
        },
    },
    {
        accessorKey: "change24h",
        header: "24h",
        cell: () => {
            return <span>24h change</span>
        },
    },
    {
        accessorKey: "pnl",
        header: "PNL",
        cell: ({ row }) => {
            const pnl = row.original.pnl
            const isPositive = pnl >= 0
            return <span className={isPositive ? "text-green-600" : "text-red-600"}>{pnl.toFixed(4)}</span>
        },
    },
    {
        accessorKey: "value",
        header: "Value",
        cell: ({ row }) => {
            return <span>${row.original.value.toFixed(4)}</span>
        },
    }
]
export default function TokenBalances({  className }: TokenBalancesProps) {
    const {wallets} = useWalletStore(state => state)
    const tokens = wallets.flatMap(wallet => wallet.Wallet.tokens)
    const [pagination, setPagination] = useState({
        pageIndex: 0,
        pageSize: 5,
    })

    const table = useReactTable(({
        data: tokens,
        columns,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        onPaginationChange: setPagination,
        state: {
            pagination
        }
    }))

    return (
        <div>
            <div className="rounded-md border border-zinc-200 dark:border-zinc-800">
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow key={headerGroup.id}>
                                {headerGroup.headers.map((header) => (
                                    <TableHead key={header.id}>
                                        {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                                    </TableHead>
                                ))}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-24 text-center">
                                    No tokens found.
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className="mt-2">
                <Pagination>
                    <PaginationContent>
                        <PaginationItem>
                            <PaginationPrevious
                                onClick={() => table.previousPage()}
                                disabled={!table.getCanPreviousPage()}
                                className={!table.getCanPreviousPage() ? "opacity-50 cursor-not-allowed" : ""}
                            />
                        </PaginationItem>
                        {Array.from({ length: table.getPageCount() }).map((_, index) => (
                            <PaginationItem key={index}>
                                    {index + 1}
                            </PaginationItem>
                        ))}
                        <PaginationItem>
                            <PaginationNext
                                onClick={() => table.nextPage()}
                                disabled={!table.getCanNextPage()}
                                className={!table.getCanNextPage() ? "opacity-50 cursor-not-allowed" : ""}
                            />
                        </PaginationItem>
                    </PaginationContent>
                </Pagination>
            </div>
        </div>
    )
}

