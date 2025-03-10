import {cn} from "@/lib/utils"
import {ArrowRight, SendHorizonal} from "lucide-react"
import {Transaction} from "@/lib/gql/graphql";
import {parseTransaction} from "@/lib/solana/parse_transaction";
import Link from "next/link";

interface TransactionListProps {
    transactions: Transaction[]
    className?: string
}


const categoryStyles = {
    transfer: "bg-blue-100 dark:bg-blue-800/30 text-blue-600 dark:text-blue-300",
    swap: "bg-purple-100 dark:bg-purple-800/30 text-purple-600 dark:text-purple-300",
    unknown: "bg-green-100 dark:bg-green-800/30 text-green-600 dark:text-green-300",
}

export default function TransactionList({transactions, className}: TransactionListProps) {

    const visibleTransactions = transactions.slice(0, 5);
    return (
        <div className={cn("space-y-1", className)}>
            {visibleTransactions.map((transaction) => {
                const tx = parseTransaction( transaction, "4g7SgYkTTnxhq1tPE1A4kR2UkUZGYLqKt7B12SKomxw3")
                return (
                    <div
                        key={tx.originalTransaction.result?.block_time}
                        className={cn(
                            "group flex items-center gap-3",
                            "p-2 rounded-lg",
                            "hover:bg-zinc-100 dark:hover:bg-zinc-800/50",
                            "transition-all duration-200",
                        )}
                    >
                        <div
                            className={cn(
                                "p-2 rounded-lg",
                                 categoryStyles[tx.info.type as keyof typeof categoryStyles] ||
                                "bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100",
                            )}
                        >
                            <SendHorizonal className="w-4 h-4" />
                        </div>

                        <div className="flex-1 flex items-center justify-between min-w-0">
                            <div className="space-y-0.5">
                                {/*<h3 className="text-xs font-medium text-zinc-900 dark:text-zinc-100">{transaction.result?.meta.status}</h3>*/}
                                <p className="text-[11px] text-zinc-600 dark:text-zinc-400">{new Date(tx.originalTransaction.result?.block_time * 1000).toLocaleString()}</p>
                                {/*{transaction.result?.transaction.signatures[0] && (*/}
                                {/*  <p className="text-[11px] text-zinc-500 dark:text-zinc-500">From: {transaction.result?.transaction.signatures[0]}</p>*/}
                                {/*)}*/}
                                {/*{transaction?.result?.transaction.signatures[1] && (*/}
                                {/*  <p className="text-[11px] text-zinc-500 dark:text-zinc-500">To: {transaction.result?.transaction.signatures[1]}</p>*/}
                                {/*)}*/}
                                {transaction.result?.transaction.signatures[0] && (
                                  <p className="text-[11px] text-zinc-500 dark:text-zinc-500 truncate max-w-[150px]">
                                    From: {transaction.result?.transaction.signatures[0]}
                                  </p>
                                )}
                                {transaction.result?.transaction.signatures[1] && (
                                  <p className="text-[11px] text-zinc-500 dark:text-zinc-500 truncate max-w-[150px]">
                                    To: {transaction.result?.transaction.signatures[1]}
                                  </p>
                                )}
                            </div>

                            <div className="flex flex-col items-end gap-0.5 pl-3">
                                <span
                                  className={cn(
                                    "text-xs font-medium",
                                     tx.info.overallTransferDirection === "received"
                                       ? "text-emerald-600 dark:text-emerald-400"
                                       : tx.info.overallTransferDirection === "sent"
                                         ? "text-red-600 dark:text-red-400"
                                         :
                                    "text-purple-600 dark:text-purple-400",
                                  )}
                                >
                                  {tx.info.overallTransferDirection === "received" ? "+" : tx.info.overallTransferDirection === "sent" ? "-" : "↔"}
                                  {tx.info.fee}
                                </span>
                                {/*<span className="text-[11px] text-zinc-500 dark:text-zinc-400">{transaction.fiatAmount}</span>*/}
                                <span
                                  className={cn(
                                    "text-[10px] px-1.5 py-0.5 rounded-full",
                                     tx.info.transactionStatus === "completed"
                                       ? "bg-green-100 text-green-800 dark:bg-green-800/30 dark:text-green-300"
                                       : tx.info.transactionStatus === "pending"
                                         ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-800/30 dark:text-yellow-300"
                                         :
                                           "bg-red-100 text-red-800 dark:bg-red-800/30 dark:text-red-300",
                                  )}
                                >
                                  {tx.info.transactionStatus}
                                </span>
                            </div>
                        </div>
                    </div>
                )
            })}
            {/*  <div*/}
            {/*  className={cn(*/}
            {/*    "group flex items-center gap-3",*/}
            {/*    "p-2 rounded-lg",*/}
            {/*    "transition-all duration-200",*/}
            {/*  )}*/}
            {/*>*/}
            {/*  /!* Icon container skeleton *!/*/}
            {/*  <div className="p-2 rounded-lg bg-zinc-100 dark:bg-zinc-800">*/}
            {/*    <Skeleton className="w-4 h-4" />*/}
            {/*  </div>*/}

            {/*  <div className="flex-1 flex items-center justify-between min-w-0">*/}
            {/*    /!* Left side content skeleton *!/*/}
            {/*    <div className="space-y-0.5">*/}
            {/*      <Skeleton className="h-3 w-24" /> /!* Title *!/*/}
            {/*      <Skeleton className="h-2.5 w-20" /> /!* Timestamp *!/*/}
            {/*      <Skeleton className="h-2.5 w-32" /> /!* From wallet *!/*/}
            {/*      <Skeleton className="h-2.5 w-28" /> /!* To wallet *!/*/}
            {/*    </div>*/}

            {/*    /!* Right side content skeleton *!/*/}
            {/*    <div className="flex flex-col items-end gap-0.5 pl-3">*/}
            {/*      <Skeleton className="h-3 w-12" /> /!* Amount *!/*/}
            {/*      <Skeleton className="h-2.5 w-10" /> /!* Fiat amount *!/*/}
            {/*      <Skeleton className="h-3 w-16 rounded-full" /> /!* Status *!/*/}
            {/*    </div>*/}
            {/*  </div>*/}
            {/*</div>*/}
            <Link href="/transactions" className="cursor-pointer">
            <button
                type="button"
                className={cn(
                    "w-full flex items-center justify-center gap-2 cursor-pointer",
                    "py-2 px-3 mt-2 rounded-lg",
                    "text-xs font-medium",
                    "bg-zinc-900 dark:bg-zinc-50",
                    "text-zinc-50 dark:text-zinc-900",
                    "hover:bg-zinc-800 dark:hover:bg-zinc-200",
                    "transition-all duration-200",
                )}
            >
               <span>View All Transactions</span>
                <ArrowRight className="w-3.5 h-3.5"/>
            </button>
            </Link>
        </div>
    )
}

