import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ArrowUpRight,
  ArrowDownRight,
  ArrowRightLeft,
  WalletIcon,
} from "lucide-react";
import type { Subwallet, Wallet } from "@/graphql/graphql";

export default function StatsCards({ wallets }: { wallets: Wallet[] }) {
  const totalBalance = 1000000;
  //wallets.reduce(
  //  (sum, wallet) => sum + (wallet.usdBalance ?? 0),
  //  0,
  //);
  const totalTransactions = 69420
  // wallets.reduce(
  //   (sum, wallet) => sum + wallet.subwallets.length,
  //   0,
  // );
  const averageChange = 12;
  // wallets.reduce((sum, wallet) => sum + (wallet.change24h ?? 0), 0) /
  //wallets.length;

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Balance</CardTitle>
          <WalletIcon className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            ${totalBalance.toLocaleString()}
          </div>
          <p className="text-xs text-muted-foreground">
            Across {wallets.length} wallets
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Portfolio Change
          </CardTitle>
          {averageChange >= 0 ? (
            <ArrowUpRight className="h-4 w-4 text-green-500" />
          ) : (
            <ArrowDownRight className="h-4 w-4 text-red-500" />
          )}
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${averageChange >= 0 ? "text-green-500" : "text-red-500"}`}
          >
            {averageChange >= 0 ? "+" : ""}
            {averageChange.toFixed(2)}%
          </div>
          <p className="text-xs text-muted-foreground">24h average change</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            Total Transactions
          </CardTitle>
          <ArrowRightLeft className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{totalTransactions}</div>
          <p className="text-xs text-muted-foreground">Across all wallets</p>
        </CardContent>
      </Card>
    </div>
  );
}
