import DashboardHeader from "@/components/pulse/wallets/dashboard/header";
import StatsCards from "@/components/pulse/wallets/dashboard/stat_cards";
import { createFileRoute } from "@tanstack/react-router";
import WalletList from "@/components/pulse/wallets/dashboard/list";
import { graphql } from "@/graphql";
import { useQuery } from "@tanstack/react-query";
import { execute } from "@/execute";

/*
TODO:
- loader to display skeletons while loading
- load in Route
- detail pages

*/



const WALLET_DASHBOARD_QUERY = graphql(`
query Wallets {
 wallets {
    id
    name
    createdAt
    updatedAt
    subwallets {
      id
      
    }
}
}`)

export const Route = createFileRoute("/_authenticated/wallets/")({
  component: RouteComponent,
});

function RouteComponent() {
  const {data } = useQuery({
    queryKey: ["Wallet_Dashboard_all_wallets"],
    queryFn: () => execute(WALLET_DASHBOARD_QUERY)
  })
  console.log(data)
  return (
    <div>
      <DashboardHeader />
      <div className="container py-6 space-y-8">
        <StatsCards wallets={[]} />
        <div className="rounded-lg border bg-card">
          <WalletList wallets={[]} />
        </div>
      </div>
    </div>
  );
}
