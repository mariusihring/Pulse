'use client'
import { CreditCard, Wallet, SendHorizontal, ArrowDownLeft, ArrowLeftRight, ShoppingCart, Coins } from "lucide-react"
import WalletList from "./walletlist"
 import TransactionList from "./transactionlist"
import TokenBalances from "@/components/pulse/dashboard/tokenbalances";
import PieChart from "@/components/pulse/dashboard/piechart";
import {useWalletStore} from "@/lib/providers/wallet_provider";
// import TokenBalances from "./token-balances"
// import PieChart from "./pie-chart"

export default function () {
  const {wallets} = useWalletStore(state => state)
  const sum_balance = wallets.reduce((sum, wallet) => {
    return sum + (wallet?.Wallet?.wallet_value || 0);
  }, 0);
  const allTokens = wallets.flatMap(wallet => wallet.Wallet.tokens)
  const allTransactions = wallets.flatMap(wallet => wallet.Wallet.transactions)


  const totalValue = allTokens.reduce((sum, token) => sum + token.value, 0)


  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white dark:bg-[#0F0F12] rounded-xl p-6 flex flex-col border border-gray-200 dark:border-[#1F1F23]">
          <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-4 text-left flex items-center gap-2">
            <Wallet className="w-3.5 h-3.5 text-zinc-900 dark:text-zinc-50" />
            Crypto Wallets
          </h2>
          <div className="flex-1">
            <p className="text-sm text-zinc-600 dark:text-zinc-400 mb-2">Total Balance</p>
            <h3 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">${sum_balance.toFixed(2)}</h3>
            <WalletList
              className="mt-4"
            />
          </div>
        </div>
        <div className="bg-white dark:bg-[#0F0F12] rounded-xl p-6 flex flex-col border border-gray-200 dark:border-[#1F1F23]">
          <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-4 text-left flex items-center gap-2">
            <CreditCard className="w-3.5 h-3.5 text-zinc-900 dark:text-zinc-50" />
            Recent Transactions
          </h2>
           <div className="flex-1">
            <TransactionList
              transactions={allTransactions}
              className=""
            />
          </div> 
        </div>
      </div>

      <div className="bg-white dark:bg-[#0F0F12] rounded-xl p-6 flex flex-col items-start justify-start border border-gray-200 dark:border-[#1F1F23]">
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-4 text-left flex items-center gap-2">
          <Coins className="w-3.5 h-3.5 text-zinc-900 dark:text-zinc-50" />
          Token Balances
        </h2>
        <div className="w-full flex flex-col lg:flex-row gap-6">
          <div className="w-full lg:w-2/3">
            <TokenBalances />
          </div>
          <div className="w-full lg:w-1/2">
            <PieChart
                data={allTokens.map((token) => ({
                  name: token.name,
                  value: token.value,
                  color: token.color,
                }))}
                totalValue={totalValue}
            />
          </div>
        </div>
      </div>
    </div>
  )
}

