import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { MoreHorizontal } from "lucide-react";
import MiniGraph from "./mini_graph";
import type { Wallet } from "@/graphql/graphql";

export default function WalletList({ wallets }: { wallets: Wallet[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Holdings</TableHead>
          <TableHead>Balance</TableHead>
          <TableHead>Transactions</TableHead>
          <TableHead>Performance</TableHead>
          <TableHead />
        </TableRow>
      </TableHeader>
      <TableBody>
        {wallets.map((wallet) => {
          const graphColor = wallet.change24h >= 0 ? "#22c55e" : "#ef4444";
          return (
            <TableRow key={wallet.id}>
              <TableCell>
                <div className="flex items-center gap-2">
                  {/* <wallet.icon className="h-5 w-5" /> */}
                  <span className="font-medium">{wallet.name}</span>
                </div>
              </TableCell>
              <TableCell>
                {wallet.balance} {wallet.currency}
              </TableCell>
              <TableCell className="font-medium">0</TableCell>
              <TableCell>{wallet.subwallets.length}</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <div className="w-24 h-8">
                    <MiniGraph
                      data={wallet.historicalData}
                      color={graphColor}
                    />
                  </div>
                  <span
                    className={
                      wallet.change24h >= 0 ? "text-green-500" : "text-red-500"
                    }
                  >
                    {wallet.change24h >= 0 ? "+" : ""}
                    {wallet.change24h}%
                  </span>
                </div>
              </TableCell>
              <TableCell>
                <Button variant="ghost" size="icon">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
