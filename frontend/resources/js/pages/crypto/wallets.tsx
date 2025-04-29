
import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { useState } from "react"
import { Coins, Plus, RefreshCw, Search, Wallet } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { WalletCard } from '@/components/pulse/crypo/wallet-card';
import { User } from '@/lib/types/crypto/dashboard/user';
const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Crypto',
        href: '/crypto',
    },
    {
        title: 'Wallets',
        href: '/crypto/wallets',
    },
];

export default function Dashboard({ user }: {user: User}) {
    const [searchQuery, setSearchQuery] = useState("")
    const [activeTab, setActiveTab] = useState("all")

    const totalValue = user.wallets.reduce((sum, wallet) => sum + wallet.value, 0)
    const totalWallets = user.wallets.length
    const totalTokens = user.wallets.reduce((sum, wallet) => sum + wallet.token_holdings.length, 0)

    const filteredWallets = user.wallets.filter((wallet) => {
      // Filter by search query
      if (
        searchQuery &&
        !wallet.name.toLowerCase().includes(searchQuery.toLowerCase()) &&
        !wallet.address.toLowerCase().includes(searchQuery.toLowerCase())
      ) {
        return false
      }

      // Filter by tab
      if (activeTab === "favorites" && !wallet.isFavorite) {
        return false
      }

      return true
    })
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Wallets" />
            <div className="container max-w-6xl p-8 mx-auto">
                <header className="flex flex-col gap-6 mb-8">
                  <div className="flex items-center justify-between">
                    <h1 className="text-2xl font-semibold tracking-tight">Crypto Wallets</h1>
                    <Button size="sm" className="gap-1">
                      <Plus className="w-4 h-4" />
                      Add Wallet
                    </Button>
                  </div>

                  <div className="grid gap-6 md:grid-cols-3">
                    <Card>
                      <CardContent className="flex items-center justify-between p-6">
                        <div className="space-y-1">
                          <p className="text-sm text-muted-foreground">Total Value</p>
                          <p className="text-2xl font-semibold">${totalValue.toLocaleString()}</p>
                        </div>
                        <div className="p-2 bg-green-100 rounded-full dark:bg-green-900">
                          <Wallet className="w-5 h-5 text-green-600 dark:text-green-400" />
                        </div>
                      </CardContent>
                    </Card>

                    <Card>
                      <CardContent className="flex items-center justify-between p-6">
                        <div className="space-y-1">
                          <p className="text-sm text-muted-foreground">Total Wallets</p>
                          <p className="text-2xl font-semibold">{totalWallets}</p>
                        </div>
                        <div className="p-2 bg-blue-100 rounded-full dark:bg-blue-900">
                          <Wallet className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                        </div>
                      </CardContent>
                    </Card>

                    <Card>
                      <CardContent className="flex items-center justify-between p-6">
                        <div className="space-y-1">
                          <p className="text-sm text-muted-foreground">Total Tokens</p>
                          <p className="text-2xl font-semibold">{totalTokens}</p>
                        </div>
                        <div className="p-2 bg-purple-100 rounded-full dark:bg-purple-900">
                          <Coins className="w-5 h-5 text-purple-600 dark:text-purple-400" />
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </header>

                <div className="flex flex-col gap-6">
                  <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                    <Tabs defaultValue="all" className="w-full sm:w-auto" onValueChange={setActiveTab}>
                      <TabsList>
                        <TabsTrigger value="all">All Wallets</TabsTrigger>
                        <TabsTrigger value="favorites">Favorites</TabsTrigger>
                      </TabsList>
                    </Tabs>

                    <div className="relative w-full sm:w-64">
                      <Search className="absolute w-4 h-4 text-muted-foreground left-3 top-3" />
                      <Input
                        placeholder="Search wallets..."
                        className="pl-9"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                      />
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <p className="text-sm text-muted-foreground">
                      Showing {filteredWallets.length} of {user.wallets.length} wallets
                    </p>

                    <Button variant="ghost" size="sm" className="gap-1 text-xs">
                      <RefreshCw className="w-3 h-3" />
                      Refresh All
                    </Button>
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    {filteredWallets.map((wallet) => (
                      <WalletCard key={wallet.id} wallet={wallet} />
                    ))}
                  </div>
                </div>
              </div>
        </AppLayout>
    );
}
