import { useState } from "react";
import { formatDistanceToNow } from "date-fns";
import {
  Bitcoin,
  Clock,
  Coins,
  Copy,
  EclipseIcon as Ethereum,
  ExternalLink,
  MoreHorizontal,
  RefreshCw,
  Star,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipProvider, TooltipTrigger, TooltipContent } from "@/components/ui/tooltip";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { Wallet } from '@/lib/types/crypto/dashboard/user';
import Solana from "@/components/icons/solana";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/apiClient";
import WalletCardSkeleton from '@/components/pulse/crypo/wallet-card-skeleton';

interface WalletCardProps {
    wallet: Wallet;
    onUpdate: (updates: Partial<Wallet>) => void;
}

export function WalletCard({ wallet, onUpdate }: WalletCardProps) {
    const queryClient = useQueryClient();
    const [isEditing, setIsEditing] = useState(false);
    const [editName, setEditName] = useState(wallet.name);
    const [isRefreshing, setIsRefreshing] = useState(false);

    const reloadWalletMutation = useMutation({
        mutationFn: async (address: string) => {
            const { data } = await apiClient.post(`/crypto/user/wallet/refresh?address=${address}`);
            return data
        },
        onMutate: () => {
            setIsRefreshing(true);
        },
        onSuccess: (updatedUser) => {
            // Update the specific wallet in the userWallets query cache
            queryClient.setQueryData(['userWallets'], updatedUser.user);


            toast.success("Wallet reloaded");
        },
        onError: (e: any) => {
            console.error('Reload Error:', e); // Debug log
            toast.error("Error reloading the wallet");
        },
        onSettled: () => {
            setIsRefreshing(false);
        },
    });

    const truncateAddress = (address: string) => {
        return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
    };

    const getChainIcon = (chain: string) => {
        switch (chain.toLowerCase()) {
            case "ethereum":
                return <Ethereum className="w-4 h-4 text-[#627EEA]" />;
            case "bitcoin":
                return <Bitcoin className="w-4 h-4 text-[#F7931A]" />;
            case "solana":
                return <Solana className="w-6 h-6 text-[#F7931A]" />;
            default:
                return <Coins className="w-4 h-4" />;
        }
    };

    const handleToggleFavorite = () => {
        onUpdate({ favorite: !wallet.favorite });
    };

    const handleEditName = () => {
        setIsEditing(true);
    };

    const handleSaveName = () => {
        if (editName.trim() && editName !== wallet.name) {
            onUpdate({ name: editName.trim() });
        }
        setIsEditing(false);
    };

    const handleCancelEdit = () => {
        setEditName(wallet.name);
        setIsEditing(false);
    };

    const copyAddress = async (text: string) => {
        try {
            await navigator.clipboard.writeText(text);
            toast.success("Address copied!");
        } catch (err) {
            toast.error('Failed to copy address');
        }
    };

    if (isRefreshing) {
        return <WalletCardSkeleton />;
    }

    return (
        <Card className="overflow-hidden transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between p-4 space-y-0 border-b">
                <div className="flex items-center gap-2">
                    {getChainIcon(wallet.chain.name)}
                    {isEditing ? (
                        <div className="flex items-center gap-2">
                            <Input
                                value={editName}
                                onChange={(e) => setEditName(e.target.value)}
                                className="w-48"
                            />
                            <Button size="sm" onClick={handleSaveName} disabled={!editName.trim()}>
                                Save
                            </Button>
                            <Button size="sm" variant="outline" onClick={handleCancelEdit}>
                                Cancel
                            </Button>
                        </div>
                    ) : (
                        <h3 className="font-medium">
                            {wallet.name}
                            <Button
                                variant="ghost"
                                size="sm"
                                onClick={handleEditName}
                                className="ml-2 text-xs"
                            >
                                Edit
                            </Button>
                        </h3>
                    )}
                </div>
                <div className="flex items-center gap-1">
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button variant="ghost" size="icon" className="w-8 h-8" onClick={handleToggleFavorite}>
                                    <Star
                                        className={cn(
                                            "w-4 h-4",
                                            wallet.favorite ? "fill-yellow-400 text-yellow-400" : "text-muted-foreground"
                                        )}
                                    />
                                    <span className="sr-only">Toggle favorite</span>
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>{wallet.favorite ? "Remove from favorites" : "Add to favorites"}</p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>

                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="w-8 h-8">
                                <MoreHorizontal className="w-4 h-4" />
                                <span className="sr-only">More options</span>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => reloadWalletMutation.mutate(wallet.address)}>
                                <RefreshCw className="w-4 h-4 mr-2" />
                                Refresh
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => copyAddress(wallet.address)}>
                                <Copy className="w-4 h-4 mr-2" />
                                Copy Address
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                                <a
                                    href={`https://solscan.io/account/${wallet.address}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex justify-between w-full"
                                >
                                    <ExternalLink className="w-4 h-4 mr-4" />
                                    View on Explorer
                                </a>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </CardHeader>

            <CardContent className="p-4">
                <div className="grid gap-4">
                    <div className="flex items-center justify-between">
                        <p className="text-sm font-medium text-muted-foreground">Address</p>
                        <p className="text-sm font-mono">{truncateAddress(wallet.address)}</p>
                    </div>

                    <div className="flex items-center justify-between">
                        <p className="text-sm font-medium text-muted-foreground">Value</p>
                        <p className="text-lg font-semibold">${wallet.value.toLocaleString()}</p>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="flex flex-col">
                            <p className="text-xs font-medium text-muted-foreground">Tokens</p>
                            <div className="flex items-center gap-1">
                                <Coins className="w-3 h-3 text-muted-foreground" />
                                <p className="text-sm">{wallet.token_holdings.length}</p>
                            </div>
                        </div>

                        <div className="flex flex-col">
                            <p className="text-xs font-medium text-muted-foreground">Chain</p>
                            <p className="text-sm capitalize">{wallet.chain.name}</p>
                        </div>
                    </div>
                </div>
            </CardContent>

            <CardFooter className="flex items-center justify-between p-4 text-xs text-muted-foreground border-t">
                <div className="flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    <span>Added {formatDistanceToNow(wallet.created_at, { addSuffix: true })}</span>
                </div>

                <div className="flex items-center gap-1">
                    <RefreshCw className="w-3 h-3" />
                    <span>Updated {formatDistanceToNow(wallet.updated_at, { addSuffix: true })}</span>
                </div>
            </CardFooter>
        </Card>
    );
}
