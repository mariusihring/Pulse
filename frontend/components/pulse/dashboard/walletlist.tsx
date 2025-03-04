'use client'
import { cn } from "@/lib/utils"
import { Wallet } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import {useWalletStore} from "@/lib/providers/wallet_provider";


interface WalletListProps {
  className?: string
}

export default function WalletList({  className }: WalletListProps) {
  const {wallets} = useWalletStore(state => state)
  return (
    <div className={cn("space-y-1", className)}>
      {wallets && wallets?.map((wallet) => (
            wallet.Wallet ?
                <div
                    key={wallet.JobID}
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
                          "bg-blue-100 dark:bg-blue-900/30": wallet.Wallet?.network === "ethereum",
                          "bg-orange-100 dark:bg-orange-900/30": wallet.Wallet?.network === "bitcoin",
                          "bg-purple-100 dark:bg-purple-900/30": wallet.Wallet?.network === "solana",
                          "bg-indigo-100 dark:bg-indigo-900/30": wallet.Wallet?.network === "polygon",
                          "bg-red-100 dark:bg-red-900/30": wallet.Wallet?.network === "avalanche",
                        })}
                    >
                      <Wallet
                          className={cn("w-3.5 h-3.5", {
                            "text-blue-600 dark:text-blue-400": wallet.Wallet?.network === "ethereum",
                            "text-orange-600 dark:text-orange-400": wallet.Wallet?.network === "bitcoin",
                            "text-purple-600 dark:text-purple-400": wallet.Wallet?.network === "solana",
                            "text-indigo-600 dark:text-indigo-400": wallet.Wallet?.network === "polygon",
                            "text-red-600 dark:text-red-400": wallet.Wallet?.network === "avalanche",
                          })}
                      />
                    </div>
                    <div>
                      <h3 className="text-xs font-medium text-zinc-900 dark:text-zinc-100">{wallet.Wallet?.name}</h3>
                      {wallet.Wallet?.address && (
                          <p className="text-[11px] text-zinc-600 dark:text-zinc-400 truncate max-w-[125px] md:max-w-[350px]">
                            {wallet.Wallet?.address}
                          </p>
                      )}
                    </div>
                  </div>

                  <div className="text-right">
                    <span
                        className="text-xs font-medium text-zinc-900 dark:text-zinc-100">$ {wallet.Wallet?.wallet_value.toFixed(2)}</span>
                  </div>
                </div> : wallet.Progress !== 100 ?
                <>
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
                        <Skeleton
                            className="text-[11px] text-zinc-600 dark:text-zinc-400 truncate max-w-[180px] w-[120px] h-2.5 mt-1.5">
                        </Skeleton>
                      </div>
                    </div>

                    <div className="text-right">
                      <Skeleton className="text-xs font-medium text-zinc-900 dark:text-zinc-100 h-4 w-12"></Skeleton>
                    </div>
                  </div>
                </> : null
      ))}

    </div>
  )
}

