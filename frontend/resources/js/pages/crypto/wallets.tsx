import AppLayout from '@/layouts/app-layout';
import { type BreadcrumbItem } from '@/types';
import { Head } from '@inertiajs/react';
import { useEffect, useMemo, useState } from 'react';
import { ArrowDownIcon, ArrowUpIcon, Coins, Plus, RefreshCw, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { WalletCard } from '@/components/pulse/crypo/wallet-card';
import WalletCardSkeleton from '@/components/pulse/crypo/wallet-card-skeleton';
import { User, Wallet as WalletType } from '@/lib/types/crypto/dashboard/user';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import apiClient from '@/lib/apiClient';
import { toast } from 'sonner';

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Crypto', href: '/crypto' },
    { title: 'Wallets', href: '/crypto/wallets' },
];

interface Chain {
  id: string;
  name: string;
}

export default function Dashboard({ user: initialUser }: { user: User }) {
  const queryClient = useQueryClient();
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newWallet, setNewWallet] = useState({ name: '', address: '', chainId: '' });

  // Fetch user wallets
  const { data: user = initialUser, isLoading, isFetching, error } = useQuery<User>({
    queryKey: ['userWallets'],
    queryFn: async () => {
      const { data } = await apiClient.get('/crypto/user/wallets/all');
      return data.user;
    },
    initialData: initialUser,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });

  // Fetch available chains
  const { data: chains = [], isLoading: isChainsLoading } = useQuery<Chain[]>({
    queryKey: ['chains'],
    queryFn: async () => {
      const { data } = await apiClient.get('/crypto/chains');
      return data.chains;
    },
    staleTime: 1000 * 60 * 60, // 1 hour
  });

  // Mutation for updating wallet
  const updateWalletMutation = useMutation({
    mutationFn: async ({ walletId, updates }: { walletId: string; updates: Partial<WalletType> }) => {
      const { data } = await apiClient.patch(`/crypto/user/wallets/${walletId}`, updates);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['userWallets'] });
      toast.success('Wallet updated successfully!');
    },
    onError: (error: any) => {
      toast.error(`Failed to update wallet: ${error.message}`);
    },
  });

  // Mutation for creating wallet
  const createWalletMutation = useMutation({
    mutationFn: async (walletData: { name: string; address: string; chainId: string }) => {
      const { data } = await apiClient.post('/crypto/user/wallets', walletData);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['userWallets'] });
      toast.success('Wallet created successfully!');
      setIsModalOpen(false);
      setNewWallet({ name: '', address: '', chainId: '' });
    },
    onError: (error: any) => {
      toast.error(`Failed to create wallet: ${error.message}`);
    },
  });

  const refreshData = () => {
    queryClient.invalidateQueries({ queryKey: ['userWallets'] });
  };

  // Compute daily and monthly changes
  const [dailyChange, setDailyChange] = useState<number>(0);
  const [monthlyChange, setMonthlyChange] = useState<number>(0);
  useEffect(() => {
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);
    const yesterdayDate = yesterday.toISOString().split('T')[0];

    const oneMonthAgo = new Date();
    oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);
    const monthAgoDate = oneMonthAgo.toISOString().split('T')[0];

    const currentTotalValue = user.wallets.reduce((sum: number, wallet: WalletType) => sum + wallet.value, 0);

    const yesterdayTotalValue = user.wallets.reduce((sum: number, wallet: WalletType) => {
      const snapshot = wallet.snapshots.find(snapshot =>
        snapshot.created_at.startsWith(yesterdayDate)
      );
      return sum + (snapshot ? parseFloat(snapshot.value) : 0);
    }, 0);

    const monthAgoTotalValue = user.wallets.reduce((sum: number, wallet: WalletType) => {
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
  }, [user]);

  const { totalValue, totalTokens, filteredWallets } = useMemo(() => {
    const totalValue = user.wallets.reduce((sum, wallet) => sum + wallet.value, 0);
    const totalTokens = user.wallets.reduce((sum, wallet) => sum + wallet.token_holdings.length, 0);

    const filteredWallets = user.wallets.filter((wallet) => {
      if (
        searchQuery &&
        !wallet.name.toLowerCase().includes(searchQuery.toLowerCase()) &&
        !wallet.address.toLowerCase().includes(searchQuery.toLowerCase())
      ) {
        return false;
      }

      if (activeTab === 'favorites' && !wallet.favorite) {
        return false;
      }

      return true;
    });

    return { totalValue, totalTokens, filteredWallets };
  }, [user, searchQuery, activeTab]);

  const handleAddWallet = () => {
    if (!newWallet.name || !newWallet.address || !newWallet.chainId) {
      toast.error('Please fill in all fields');
      return;
    }
    createWalletMutation.mutate(newWallet);
  };

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title="Wallets" />
      <div className="container max-w-6xl p-8 mx-auto">
        <header className="flex flex-col gap-6 mb-8">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-semibold tracking-tight">Crypto Wallets</h1>
            <Button
              size="sm"
              className="gap-1"
              onClick={() => setIsModalOpen(true)}
              disabled={isLoading || isFetching}
            >
              <Plus className="w-4 h-4" />
              Add Wallet
            </Button>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {isLoading || isFetching ? (
              <>
                <Card>
                  <CardContent className="flex items-center justify-between p-6">
                    <div className="space-y-1">
                      <Skeleton className="h-4 w-20" />
                      <Skeleton className="h-6 w-32" />
                    </div>
                    <Skeleton className="h-10 w-10 rounded-full" />
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="flex items-center justify-between p-6">
                    <div className="space-y-1">
                      <Skeleton className="h-4 w-20" />
                      <Skeleton className="h-6 w-32" />
                    </div>
                    <Skeleton className="h-10 w-10 rounded-full" />
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="flex items-center justify-between p-6">
                    <div className="space-y-1">
                      <Skeleton className="h-4 w-20" />
                      <Skeleton className="h-6 w-32" />
                    </div>
                    <Skeleton className="h-10 w-10 rounded-full" />
                  </CardContent>
                </Card>
              </>
            ) : (
              <>
                <Card>
                  <CardContent className="flex items-center justify-between p-6">
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Total Value</p>
                      <p className="text-2xl font-semibold">${totalValue.toLocaleString()}</p>
                    </div>
                    <div className="p-2 bg-green-100 rounded-full dark:bg-green-900">
                      <Coins className="w-5 h-5 text-green-600 dark:text-green-400" />
                    </div>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="flex items-center justify-between p-6">
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Daily Change:</p>
                      <div
                        className={`flex text-2xl font-semibold items-center gap-1 ${
                          Number(Math.abs(Number(dailyChange)).toFixed(4)) === 0
                            ? ''
                            : dailyChange > 0
                            ? 'text-green-400 border-green-800'
                            : 'text-red-400 border-red-800'
                        }`}
                      >
                        {Number(Math.abs(Number(dailyChange)).toFixed(4)) === 0 ? null : dailyChange > 0 ? (
                          <ArrowUpIcon size={12} />
                        ) : (
                          <ArrowDownIcon size={12} />
                        )}
                        {Math.abs(Number(dailyChange.toFixed(0)))}% today
                      </div>
                    </div>
                    <div className="p-2 bg-blue-100 rounded-full dark:bg-blue-900">
                      <Coins className="w-5 h-5 text-blue-600 dark:text-blue-400" />
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
              </>
            )}
          </div>
        </header>

        <div className="flex flex-col gap-6">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <Tabs
              defaultValue="all"
              className="w-full sm:w-auto"
              onValueChange={setActiveTab}
              disabled={isLoading || isFetching}
            >
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
                disabled={isLoading || isFetching}
              />
            </div>
          </div>

          <div className="flex items-center justify-between">
            {isLoading || isFetching ? (
              <Skeleton className="h-4 w-32" />
            ) : (
              <p className="text-sm text-muted-foreground">
                Showing {filteredWallets.length} of {user.wallets.length} wallets
              </p>
            )}

            <Button
              variant="ghost"
              size="sm"
              className="gap-1 text-xs"
              onClick={refreshData}
              disabled={isLoading || isFetching}
            >
              <RefreshCw className={`w-3 h-3 ${isFetching ? 'animate-spin' : ''}`} />
              {isFetching ? 'Refreshing...' : 'Refresh All'}
            </Button>
          </div>

          {error && (
            <p className="text-sm text-red-500">Failed to load wallet data: {error.message}</p>
          )}

          <div className="grid gap-4 md:grid-cols-2">
            {isLoading || isFetching ? (
              Array.from({ length: 4 }).map((_, index) => (
                <WalletCardSkeleton key={index} />
              ))
            ) : (
              filteredWallets.map((wallet) => (
                <WalletCard
                  key={wallet.id}
                  wallet={wallet}
                  onUpdate={(updates) => updateWalletMutation.mutate({ walletId: wallet.id, updates })}
                />
              ))
            )}
          </div>
        </div>

        {/* Add Wallet Modal */}
        <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add New Wallet</DialogTitle>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Wallet Name</Label>
                <Input
                  id="name"
                  value={newWallet.name}
                  onChange={(e) => setNewWallet({ ...newWallet, name: e.target.value })}
                  placeholder="Enter wallet name"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="address">Wallet Address</Label>
                <Input
                  id="address"
                  value={newWallet.address}
                  onChange={(e) => setNewWallet({ ...newWallet, address: e.target.value })}
                  placeholder="Enter wallet address"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="chain">Blockchain</Label>
                <Select
                  value={newWallet.chainId}
                  onValueChange={(value) => setNewWallet({ ...newWallet, chainId: value })}
                  disabled={isChainsLoading}
                >
                  <SelectTrigger id="chain">
                    <SelectValue placeholder="Select blockchain" />
                  </SelectTrigger>
                  <SelectContent>
                    {chains.map((chain) => (
                      <SelectItem key={chain.id} value={chain.id}>
                        {chain.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setIsModalOpen(false);
                  setNewWallet({ name: '', address: '', chainId: '' });
                }}
              >
                Cancel
              </Button>
              <Button
                onClick={handleAddWallet}
                disabled={createWalletMutation.isPending || isChainsLoading}
              >
                {createWalletMutation.isPending ? 'Adding...' : 'Add Wallet'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </AppLayout>
  );
}
