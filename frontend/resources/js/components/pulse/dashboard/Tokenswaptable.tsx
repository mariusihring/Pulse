
import { useState } from "react"
import { useReactTable, getCoreRowModel, getFilteredRowModel } from "@tanstack/react-table"
import { Input } from "@/components/ui/input"
import { ArrowUpCircle, ArrowDownCircle, Layers, ArrowRight } from 'lucide-react'
import { cn } from "@/lib/utils"

const categoryStyles = {
    accumulation: "bg-blue-100 dark:bg-blue-800/30 text-blue-600 dark:text-blue-300",
    buy: "bg-green-100 dark:bg-green-800/30 text-green-600 dark:text-green-300",
    sell: "bg-red-100 dark:bg-red-800/30 text-red-600 dark:text-red-300",
}

type Swap = {
    id: string
    transaction_type: string
    sub_category: string
    pair_label: string
    block_timestamp: string
    total_value_usd: string
    bought: { symbol: string; amount: string }
    sold: { symbol: string; amount: string }
    wallet_address: string
    transaction_hash: string
}

interface SwapTableProps {
    swaps?: Swap[]
}

export function SwapTable({ swaps = [] }: SwapTableProps) {
    const [globalFilter, setGlobalFilter] = useState("")
    const data = Array.isArray(swaps) ? swaps.slice(0, 9) : []

    const table = useReactTable({
        data,
        columns: [],
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        state: { globalFilter },
        onGlobalFilterChange: setGlobalFilter,
        filterFns: {
            fuzzy: (row, columnId, value) => {
                if (columnId === "transaction_type" || columnId === "pair_label") {
                    const cellValue = row.getValue(columnId) as string
                    return cellValue?.toLowerCase().includes(value.toLowerCase())
                }
                return false
            },
        },
        globalFilterFn: "fuzzy",
    })

    const getTransactionIcon = (type: string) => {
        if (type === "buy") return ArrowUpCircle
        if (type === "sell") return ArrowDownCircle
        return Layers
    }

    const getTransactionType = (type: string) => {
        if (type === "buy") return "incoming"
        if (type === "sell") return "outgoing"
        return "swap"
    }

    return (
        <div className="space-y-4 pt-4">
            <Input
                placeholder="Search by type or pair (e.g., buy, Butthole/SOL)..."
                value={globalFilter}
                onChange={(e) => setGlobalFilter(e.target.value)}
                className="max-w-sm"
            />
            <div className="space-y-1">
                {table.getRowModel().rows.map((row) => {
                    const swap = row.original
                    const Icon = getTransactionIcon(swap.transaction_type)
                    return (
                        <div
                            key={swap.id}
                            className={cn(
                                "group flex items-center gap-3",
                                "p-2 rounded-lg",
                                "hover:bg-zinc-100 dark:hover:bg-zinc-800/50",
                                "transition-all duration-200"
                            )}
                        >
                            <div
                                className={cn(
                                    "p-2 rounded-lg",
                                    categoryStyles[swap.transaction_type as keyof typeof categoryStyles] ||
                                    "bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100"
                                )}
                            >
                                {/*<Icon className="w-4 h-4" />*/}
                                <img src={swap.exchange_logo} className="w-4 h-4" />
                            </div>
                            <div className="flex-1 flex items-center justify-between min-w-0">
                                <div className="space-y-0.5">
                                    <h3 className="text-xs font-medium text-zinc-900 dark:text-zinc-100">Exchange {swap.sold.symbol} for {swap.bought.symbol}</h3>
                                    <p className="text-[11px] text-zinc-600 dark:text-zinc-400">
                                        {new Date(swap.block_timestamp).toLocaleString()}
                                    </p>
                                    <p className="text-[11px] text-zinc-500 dark:text-zinc-500 truncate max-w-[150px]">
                                        From: {swap.wallet_address}
                                    </p>
                                </div>
                                <div className="flex flex-col items-end gap-0.5 pl-3">
                                    <span
                                        className={cn(
                                            "text-xs font-medium",
                                            swap.transaction_type === "buy"
                                                ? "text-emerald-600 dark:text-emerald-400"
                                                : swap.transaction_type === "sell"
                                                    ? "text-red-600 dark:text-red-400"
                                                    : "text-purple-600 dark:text-purple-400"
                                        )}
                                    >
                                        {swap.transaction_type === "buy" ? "+" : swap.transaction_type === "sell" ? "-" : "↔"}
                                        {swap.bought.amount} {swap.bought.symbol}
                                    </span>
                                    <span className="text-[11px] text-zinc-500 dark:text-zinc-400">
                                        ${parseFloat(swap.total_value_usd).toFixed(2)}
                                    </span>
                                    <span
                                        className={cn(
                                            "text-[10px] px-1.5 py-0.5 rounded-full",
                                            "bg-green-100 text-green-800 dark:bg-green-800/30 dark:text-green-300"
                                        )}
                                    >
                                        completed
                                    </span>
                                </div>
                            </div>
                        </div>
                    )
                })}
                <button
                    type="button"
                    className={cn(
                        "w-full flex items-center justify-center gap-2",
                        "py-2 px-3 mt-2 rounded-lg",
                        "text-xs font-medium",
                        "bg-zinc-900 dark:bg-zinc-50",
                        "text-zinc-50 dark:text-zinc-900",
                        "hover:bg-zinc-800 dark:hover:bg-zinc-200",
                        "transition-all duration-200"
                    )}
                >
                    <span>View All Transactions</span>
                    <ArrowRight className="w-3.5 h-3.5" />
                </button>
            </div>
        </div>
    )
}
