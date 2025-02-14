import { graphql } from '@/graphql';
import { queryOptions, useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import { execute } from '@/execute';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import WalletHeader from '@/components/pulse/wallets/details/header';
import type { Wallet } from '@/graphql/graphql';


const WALLET_DETAIL_QUERY = graphql(`
    query WalletDetailQuery($id: UUID!) {
     
        wallet(id: $id) {
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
          totalBalance
      }
      }
    
  `)



export const Route = createFileRoute(
  '/_authenticated/wallets/wallet/$walletId',
)({
  component: RouteComponent,
  loader: async ({ context: { queryClient }, params }) =>
		queryClient.ensureQueryData(
			queryOptions({
				queryKey: [`wallet_detail_${params.walletId}`],
				queryFn: () => execute(WALLET_DETAIL_QUERY, {id: params.walletId}),
			}),
		),
    errorComponent: () => <div className='text-red-700 flex w-full h-full items-center justify-center text-center text-cl font-bold'>Oooops smth went wrong !</div>
});


function RouteComponent() {
  
  const { walletId } = Route.useParams()
  const {
		data: { wallet },
		error,
	} = useQuery({
		queryKey: [`wallet_detail_${walletId}`],
		queryFn: () => execute(WALLET_DETAIL_QUERY, {id: walletId}),
		initialData: Route.useLoaderData(),
	});
  return (
    <div className='flex min-h-screen flex-col bg-background'>
      <div className='container py-8'>
        <WalletHeader wallet={wallet as Wallet} />
        <Tabs defaultValue="overview" className='space-y-8'>
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="transactions">Transactions</TabsTrigger>
          </TabsList>
          <TabsContent value="overview" className='space-y-8'>
            {/* Some type of graph */}
            {/* <SubwalletsGrid wallet={wallet} /> */}
            </TabsContent>
            <TabsContent value="transactions" className='space-y-8'>
            {/* <TransactionHistory wallet={wallet} /> */}
            </TabsContent>
        </Tabs>

      </div>
    </div>
  )
}
