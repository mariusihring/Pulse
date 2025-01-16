import DashboardHeader from "@/components/pulse/wallets/dashboard/header";
import StatsCards from "@/components/pulse/wallets/dashboard/stat_cards";
import { createFileRoute } from "@tanstack/react-router";
import WalletList from "@/components/pulse/wallets/dashboard/list";
import { graphql } from "@/graphql";
import { useQuery } from "@tanstack/react-query";
import { execute } from "@/execute";
import type { Wallet } from "@/graphql/graphql";

/*
TODO:
- loader to display skeletons while loading
- load in Route
- detail pages

*/

const WALLET_DASHBOARD_QUERY = graphql(`
  query Wallets {
      wallets {
          subwallets {
              id
              createdAt
              updatedAt
              name
              
              tokens {
                  amount
                  valueUsd
                  totalPnl
              }
              snapshots {
                  snapshotDate
                  totalPnl
                  totalValue
                  id
                  createdAt
              }
          }
          createdAt
          id
          updatedAt
          name
      }
}`);

export const Route = createFileRoute("/_authenticated/wallets/")({
	component: RouteComponent,
});

function RouteComponent() {
	const { data: { wallets = [] } = {} } = useQuery({
    queryKey: ["Wallet_Dashboard_all_wallets"],
    queryFn: () => execute(WALLET_DASHBOARD_QUERY),
  });
	return (
		<div>
			<DashboardHeader />
			<div className="container py-6 space-y-8">
				<StatsCards wallets={wallets as Wallet[]} />
				<div className="rounded-lg border bg-card">
					<WalletList wallets={wallets as Wallet[]} />
				</div>
			</div>
		</div>
	);
}
