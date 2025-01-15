import DashboardHeader from "@/components/pulse/wallets/dashboard/header";
import StatsCards from "@/components/pulse/wallets/dashboard/stat_cards";
import { createFileRoute } from "@tanstack/react-router";
import type { Wallet } from "@/graphql/graphql";
import WalletList from "@/components/pulse/wallets/dashboard/list";

export const Route = createFileRoute("/_authenticated/wallets/")({
  component: RouteComponent,
});

function RouteComponent() {
  const wallets: Wallet[] = [
    {
      createdAt: undefined,
      id: "someid",
      name: "my first wallet",
      subwallets: [],
      updatedAt: undefined,
    },
    {
      createdAt: undefined,
      id: "someid2",
      name: "my first wallet 2",
      subwallets: [],
      updatedAt: undefined,
    },
  ];
  return (
    <div>
      <DashboardHeader />
      <div className="container py-6 space-y-8">
        <StatsCards wallets={wallets} />
        <div className="rounded-lg border bg-card">
          <WalletList wallets={wallets} />
        </div>
      </div>
    </div>
  );
}
