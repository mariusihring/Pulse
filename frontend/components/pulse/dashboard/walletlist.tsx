import { cn } from "@/lib/utils"
import { Wallet } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import {WalletUpdate} from "@/lib/gql/graphql";

interface AccountItem {
  id: string
  title: string
  description?: string
  balance: string
  type: "ethereum" | "bitcoin" | "solana" | "polygon" | "avalanche"
}

interface WalletListProps {
  accounts?: WalletUpdate
  className?: string
}

export default function WalletList({ accounts = [], className }: WalletListProps) {
  return (
    <div className={cn("space-y-1", className)}>
      {accounts?.Wallet?.tokens.map((account) => (
        <div
          key={account.id}
          className={cn(
            "group flex items-center justify-between",
            "p-2 rounded-lg",
            "hover:bg-zinc-100 dark:hover:bg-zinc-800/50",
            "transition-all duration-200",
          )}
        >
          <div className="flex items-center gap-2">
            <div
              className={cn("p-1.5 rounded-lg", {
                "bg-blue-100 dark:bg-blue-900/30": account.type === "ethereum",
                "bg-orange-100 dark:bg-orange-900/30": account.type === "bitcoin",
                "bg-purple-100 dark:bg-purple-900/30": account.type === "solana",
                "bg-indigo-100 dark:bg-indigo-900/30": account.type === "polygon",
                "bg-red-100 dark:bg-red-900/30": account.type === "avalanche",
              })}
            >
              <Wallet
                className={cn("w-3.5 h-3.5", {
                  "text-blue-600 dark:text-blue-400": account.type === "ethereum",
                  "text-orange-600 dark:text-orange-400": account.type === "bitcoin",
                  "text-purple-600 dark:text-purple-400": account.type === "solana",
                  "text-indigo-600 dark:text-indigo-400": account.type === "polygon",
                  "text-red-600 dark:text-red-400": account.type === "avalanche",
                })}
              />
            </div>
            <div>
              <h3 className="text-xs font-medium text-zinc-900 dark:text-zinc-100">{account.title}</h3>
              {account.description && (
                <p className="text-[11px] text-zinc-600 dark:text-zinc-400 truncate max-w-[180px]">
                  {account.description}
                </p>
              )}
            </div>
          </div>

          <div className="text-right">
            <span className="text-xs font-medium text-zinc-900 dark:text-zinc-100">{account.balance}</span>
          </div>
        </div>
      ))}
      <div
        className={cn(
          "group flex items-center justify-between",
          "p-2 rounded-lg",
          "hover:bg-zinc-100 dark:hover:bg-zinc-800/50",
          "transition-all duration-200",
        )}
      >
        <div className="flex items-center gap-2">
          <Skeleton
            className={cn("p-1.5 rounded-lg")}>
            <Skeleton
              className={cn("w-3.5 h-3.5")}
            />
          </Skeleton>
          <div>
            <Skeleton className="text-xs font-medium text-zinc-900 dark:text-zinc-100 w-12 h-3"></Skeleton>
            <Skeleton className="text-[11px] text-zinc-600 dark:text-zinc-400 truncate max-w-[180px] w-[120px] h-2.5 mt-1.5">
            </Skeleton>
          </div>
        </div>

        <div className="text-right">
          <Skeleton className="text-xs font-medium text-zinc-900 dark:text-zinc-100 h-4 w-12"></Skeleton>
        </div>
      </div>
    </div>
  )
}

