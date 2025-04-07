import { Badge } from '@/components/ui/badge';
import {Progress} from "@/components/ui/progress"
import { TrendingUpIcon, WalletIcon, ArrowUpIcon, ArrowDownIcon } from 'lucide-react';
import { useState } from 'react';

export default function DashboardOverview ({data}) {
    console.log(data)
    const sum = data.wallets.flatMap(wallet => wallet.value)
    const formattedBalance = `$ ${parseFloat(sum).toFixed(2)}`
    const [hideBalance, setHideBalance] = useState(false)
    const userData = {
        dailyChange: 2,
        wallets: data.wallets,
        monthlyChange: 3.14,
        current_value: sum,
        goal: 1000000
    }

    const goalPercentage = (userData.current_value / userData.goal) * 100
    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(amount)
    }
    return (
        <div className="space-y-6 p-4 h-full">
            <div>
                <div className="flex items-center justify-between mb-2">
                    <h3 className="text-sm font-medium text-zinc-400">Total Balance</h3>
                    <div className="flex items-center">
                        <Badge
                            variant="outline"
                            className={`flex items-center gap-1 ${userData.dailyChange >= 0 ? "text-green-400 border-green-800" : "text-red-400 border-red-800"}`}
                        >
                            {userData.dailyChange >= 0 ? <ArrowUpIcon size={12} /> : <ArrowDownIcon size={12} />}
                            {Math.abs(userData.dailyChange)}% today
                        </Badge>
                    </div>
                </div>
                <p className="text-3xl font-bold tracking-tight">{hideBalance ? "••••••" : formattedBalance}</p>
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div className="bg-gray-100 dark:bg-zinc-900 p-4 rounded-lg">
                    <div className="flex items-center gap-2 mb-2">
                        <WalletIcon size={16} className="text-zinc-400" />
                        <h4 className="text-sm font-medium text-zinc-400">Wallets</h4>
                    </div>
                    <p className="text-2xl font-semibold">{userData.wallets.length}</p>
                </div>

                <div className="bg-gray-100  dark:bg-zinc-900 p-4 rounded-lg">
                    <div className="flex items-center gap-2 mb-2">
                        <TrendingUpIcon size={16} className="text-zinc-400" />
                        <h4 className="text-sm font-medium text-zinc-400">Monthly</h4>
                    </div>
                    <div className="flex items-center gap-2">
                        <p className="text-2xl font-semibold">{userData.monthlyChange}%</p>
                        <span className={userData.monthlyChange >= 0 ? "text-green-400" : "text-red-400"}>
                        {userData.monthlyChange >= 0 ? "↑" : "↓"}
                      </span>
                    </div>
                </div>
            </div>

            <div className="h-full m-auto p-auto">
                <div className="flex items-center justify-between mb-2">
                    <h3 className="text-sm font-medium text-zinc-400">Networth Goal</h3>
                    <span className="text-sm text-zinc-400">
                      {formatCurrency(userData.current_value)} / {formatCurrency(userData.goal)}
                    </span>
                </div>
                <Progress
                    value={goalPercentage}
                    className="h-2 bg-zinc-800"
                />
            </div>
        </div>
    )
}
