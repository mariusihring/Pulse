import { useState } from "react"
import { formatDistanceToNow } from "date-fns"
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
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { cn } from "@/lib/utils"
import { Wallet } from '@/lib/types/crypto/dashboard/user';
import Solana from "@/components/icons/solana"
import {toast} from "sonner"


export function WalletCard({ wallet }: {wallet: Wallet}) {

  const [isFavorite, setIsFavorite] = useState(wallet.favorite)

  const truncateAddress = (address: string) => {
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`
  }

  const getChainIcon = (chain: string) => {
    switch (chain.toLowerCase()) {
      case "ethereum":
        return <Ethereum className="w-4 h-4 text-[#627EEA]" />
      case "bitcoin":
        return <Bitcoin className="w-4 h-4 text-[#F7931A]" />
      case "solana":
          return <Solana className="w-6 h-6 text-[#F7931A]" />
      default:
        return <Coins className="w-4 h-4" />
    }
  }

  const toggleFavorite = () => {
    setIsFavorite(!isFavorite)
  }

    const copyAddress = async (text: string) => {
        try {
            await navigator.clipboard.writeText(text);
            toast.success("Copied!")
        } catch (err) {

            try {
                const textArea = document.createElement('textarea');
                textArea.value = text;
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                toast.success("Copied!")
            } catch (fallbackErr) {
                toast.error('Failed to copy');

            }
        }
    }
  return (
    <Card className="overflow-hidden transition-all hover:shadow-md">
      <CardHeader className="flex flex-row items-center justify-between p-4 space-y-0 border-b">
        <div className="flex items-center gap-2">
          {getChainIcon(wallet.chain.name)}
          <h3 className="font-medium">{wallet.name}</h3>
        </div>
        <div className="flex items-center gap-1">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="w-8 h-8" onClick={toggleFavorite}>
                  <Star
                    className={cn("w-4 h-4", isFavorite ? "fill-yellow-400 text-yellow-400" : "text-muted-foreground")}
                  />
                  <span className="sr-only">Toggle favorite</span>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>{isFavorite ? "Remove from favorites" : "Add to favorites"}</p>
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
              <DropdownMenuItem>
                <RefreshCw className="w-4 h-4 mr-2" />
                Refresh
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => copyAddress(wallet.address)}>
                <Copy className="w-4 h-4 mr-2" />
                Copy Address
              </DropdownMenuItem>
              <DropdownMenuItem >
                  <a href={ `https://solscan.io/account/${wallet.address} `} target="_blank" className="flex justify-between w-full">
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
  )
}
