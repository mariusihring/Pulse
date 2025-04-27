import { Badge } from '@/components/ui/badge';
import {Progress} from "@/components/ui/progress"
import { TrendingUpIcon, WalletIcon, ArrowUpIcon, ArrowDownIcon } from 'lucide-react';
import { useEffect, useState } from 'react';

interface Wallet {
    value: number;
    snapshots: Array<{
        created_at: string;
        value: string;
    }>;
}

interface OverviewProps {
    data: {
        wallets: Wallet[];
    };
}

function summing (data: number[]) {
    let val = 0
    for (const item in data) {
        val += data[item];
    }
    return val
}

export default function DashboardOverview ({data}: OverviewProps) {
    if (data.wallets.length === 0) {
        //TODO: placeholder
        return null
    }
    const sum = summing(data.wallets.flatMap(wallet => wallet.value)) || 0
    let formattedBalance = `$ ${Number(sum).toFixed(2)}` || 0
    useEffect(() => {
        formattedBalance = `$ ${Number(sum).toFixed(2)}` || 0
    }, [sum]);

    const [hideBalance, setHideBalance] = useState(false)
    const [dailyChange, setDailyChange] = useState<number>(0);
    const [monthlyChange, setMonthlyChange] = useState<number>(0);
    useEffect(() => {
        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1);
        const yesterdayDate = yesterday.toISOString().split('T')[0];

        const oneMonthAgo = new Date();
        oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);
        const monthAgoDate = oneMonthAgo.toISOString().split('T')[0];

        const currentTotalValue = data.wallets.reduce((sum: number, wallet: Wallet) => sum + wallet.value, 0);

        const yesterdayTotalValue = data.wallets.reduce((sum: number, wallet: Wallet) => {
            const snapshot = wallet.snapshots.find(snapshot => 
                snapshot.created_at.startsWith(yesterdayDate)
            );
            return sum + (snapshot ? parseFloat(snapshot.value) : 0);
        }, 0);

        const monthAgoTotalValue = data.wallets.reduce((sum: number, wallet: Wallet) => {
            const snapshot = wallet.snapshots.find(snapshot => 
                snapshot.created_at.startsWith(monthAgoDate)
            );
            return sum + (snapshot ? parseFloat(snapshot.value) : 0);
        }, 0);

        if (yesterdayTotalValue > 0) {
            setDailyChange((currentTotalValue / yesterdayTotalValue) * 100 - 100);
        }

        if (monthAgoTotalValue > 0) {
            setMonthlyChange((currentTotalValue / monthAgoTotalValue) * 100 - 100);
        }
    }, []);

    const userData = {
        dailyChange: dailyChange,
        wallets: data.wallets,
        monthlyChange: monthlyChange,
        current_value: sum,
        goal: 100000
    }
    

    const goalPercentage = Math.min((userData.current_value / userData.goal) * 100, 100)
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
                            className={`flex items-center gap-1 ${Number(Math.abs(Number(userData.dailyChange)).toFixed(4)) === 0 ? null : userData.dailyChange > 0 ? "text-green-400 border-green-800" : "text-red-400 border-red-800"}`}
                        >
                            {Number(Math.abs(Number(userData.dailyChange)).toFixed(4)) === 0 ? null : userData.dailyChange > 0 ? <ArrowUpIcon size={12} /> : <ArrowDownIcon size={12}/>}
                            {Math.abs(Number(userData.dailyChange.toFixed(0)))}% today
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
                        <p className="text-2xl font-semibold">{userData.monthlyChange.toFixed(2)}%</p>
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
                <Progress value={goalPercentage} className="h-2" />
            </div>
        </div>
    )
}
