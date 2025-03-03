import { CreditCard, Wallet, SendHorizontal, ArrowDownLeft, ArrowLeftRight, ShoppingCart, Coins } from "lucide-react"
import WalletList from "./walletlist"
 import TransactionList from "./transactionlist"
// import TokenBalances from "./token-balances"
// import PieChart from "./pie-chart"

export default function () {
  const tokens = [
    {
      id: "1",
      name: "Ethereum",
      symbol: "ETH",
      amount: "5.5",
      currentPrice: "1850.00",
      pnl: "+12.5%",
      currentValue: "10175.00",
      change24h: "+2.3%",
      image: "https://cryptologos.cc/logos/ethereum-eth-logo.png",
    },
    {
      id: "2",
      name: "Bitcoin",
      symbol: "BTC",
      amount: "0.5",
      currentPrice: "29000.00",
      pnl: "-5.2%",
      currentValue: "14500.00",
      change24h: "-1.8%",
      image: "https://cryptologos.cc/logos/bitcoin-btc-logo.png",
    },
    {
      id: "3",
      name: "Solana",
      symbol: "SOL",
      amount: "100",
      currentPrice: "35.00",
      pnl: "+25.0%",
      currentValue: "3500.00",
      change24h: "+5.6%",
      image: "https://cryptologos.cc/logos/solana-sol-logo.png",
    },
    {
      id: "4",
      name: "Polygon",
      symbol: "MATIC",
      amount: "2000",
      currentPrice: "0.80",
      pnl: "-8.3%",
      currentValue: "1600.00",
      change24h: "-3.1%",
      image: "https://cryptologos.cc/logos/polygon-matic-logo.png",
    },
    {
      id: "5",
      name: "Avalanche",
      symbol: "AVAX",
      amount: "50",
      currentPrice: "16.00",
      pnl: "+15.7%",
      currentValue: "800.00",
      change24h: "+4.2%",
      image: "https://cryptologos.cc/logos/avalanche-avax-logo.png",
    },
  ]

  const totalValue = tokens.reduce((sum, token) => sum + Number.parseFloat(token.currentValue), 0)

  // Mock data for charts
  const walletBalanceData = [
    { date: "2023-01-01", balance: 20000 },
    { date: "2023-02-01", balance: 22000 },
    { date: "2023-03-01", balance: 21000 },
    { date: "2023-04-01", balance: 23000 },
    { date: "2023-05-01", balance: 25000 },
    { date: "2023-06-01", balance: 26540 },
  ]

  const inflowData = [
    { date: "2023-01", amount: 5000 },
    { date: "2023-02", amount: 3000 },
    { date: "2023-03", amount: 2000 },
    { date: "2023-04", amount: 4000 },
    { date: "2023-05", amount: 3500 },
    { date: "2023-06", amount: 4500 },
  ]

  const outflowData = [
    { date: "2023-01", amount: 3000 },
    { date: "2023-02", amount: 2000 },
    { date: "2023-03", amount: 3000 },
    { date: "2023-04", amount: 2500 },
    { date: "2023-05", amount: 1500 },
    { date: "2023-06", amount: 3000 },
  ]

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
            <h3 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">$26,540.25</h3>
            <WalletList
              accounts={[
                {
                  id: "1",
                  title: "Ethereum",
                  description: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
                  balance: "$12,459.45",
                  type: "ethereum",
                },
                {
                  id: "2",
                  title: "Bitcoin",
                  description: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
                  balance: "$8,850.00",
                  type: "bitcoin",
                },
                {
                  id: "3",
                  title: "Solana",
                  description: "5YNmS1R9nNSCDzb5a7mMJ1dwK9uHeAAQmx5c8DQQvXbP",
                  balance: "$3,230.80",
                  type: "solana",
                },
                {
                  id: "4",
                  title: "Polygon",
                  description: "0x8ba1f109551bD432803012645Ac136ddd64DBA72",
                  balance: "$1,200.00",
                  type: "polygon",
                },
                {
                  id: "5",
                  title: "Avalanche",
                  description: "0x2E7D2C03a9507ae265ecFe10859D88aCEA45fd2f",
                  balance: "$800.00",
                  type: "avalanche",
                },
              ]}
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
              transactions={[
                {
                  id: "1",
                  title: "Sent ETH",
                  amount: "0.5 ETH",
                  fiatAmount: "$925.00",
                  type: "outgoing",
                  category: "transfer",
                  icon: SendHorizontal,
                  timestamp: "Today, 2:45 PM",
                  status: "completed",
                  fromWallet: "Ethereum",
                  toAddress: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
                },
                {
                  id: "2",
                  title: "Received BTC",
                  amount: "0.03 BTC",
                  fiatAmount: "$870.00",
                  type: "incoming",
                  category: "transfer",
                  icon: ArrowDownLeft,
                  timestamp: "Today, 9:00 AM",
                  status: "completed",
                  toWallet: "Bitcoin",
                  fromAddress: "3FZbgi29cpjq2GjdwV8eyHuJJnkLtktZc5",
                },
                {
                  id: "3",
                  title: "Swapped SOL to USDC",
                  amount: "50 SOL",
                  fiatAmount: "$1,750.00",
                  type: "swap",
                  category: "exchange",
                  icon: ArrowLeftRight,
                  timestamp: "Yesterday, 3:30 PM",
                  status: "completed",
                  fromWallet: "Solana",
                  toWallet: "Solana",
                },
                {
                  id: "4",
                  title: "Bought MATIC",
                  amount: "500 MATIC",
                  fiatAmount: "$400.00",
                  type: "incoming",
                  category: "purchase",
                  icon: ShoppingCart,
                  timestamp: "2 days ago",
                  status: "completed",
                  toWallet: "Polygon",
                },
                {
                  id: "5",
                  title: "Sent AVAX",
                  amount: "10 AVAX",
                  fiatAmount: "$230.00",
                  type: "outgoing",
                  category: "transfer",
                  icon: SendHorizontal,
                  timestamp: "3 days ago",
                  status: "pending",
                  fromWallet: "Avalanche",
                  toAddress: "0x1234...5678",
                },
              ]}
              className=""
            />
          </div> 
        </div>
      </div>

      <div className="bg-white dark:bg-[#0F0F12] rounded-xl p-6 flex flex-col items-start justify-start border border-gray-200 dark:border-[#1F1F23]">
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-4 text-left flex items-center gap-2 w-full">
          <Coins className="w-3.5 h-3.5 text-zinc-900 dark:text-zinc-50" />
          Token Balances and Analytics
        </h2>
        <div className="w-full flex flex-col lg:flex-row gap-6">
          <div className="w-full lg:w-1/2">
            {/* <TokenBalances tokens={tokens} />
            <div className="mt-4">
              <h3 className="text-md font-semibold text-gray-900 dark:text-white mb-2">Token Distribution</h3>
              <PieChart
                data={tokens.map((token) => ({
                  name: token.symbol,
                  value: Number.parseFloat(token.currentValue),
                  color: token.color,
                }))}
                totalValue={totalValue}
              />
            </div> */}
          </div>
        </div>
      </div>
    </div>
  )
}

