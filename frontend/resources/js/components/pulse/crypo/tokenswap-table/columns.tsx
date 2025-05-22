import { DialogHeader } from "@/components/ui/dialog"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { SubCategory, TokenSwap, TransactionType } from "@/lib/types/crypto/dashboard/user"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Dialog, DialogTrigger, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import {  ArrowUpDown, MoreHorizontal } from "lucide-react"

const subCategoryLabels = {
    [SubCategory.NewPosition]: "New Position",
    [SubCategory.Accumulation]: "Accumulation",
    [SubCategory.SellAll]: "Sell All",
    [SubCategory.PartialSell]: "Partial Sell",
  }
  
  const transactionTypeStyles = {
    [TransactionType.Buy]: "bg-green-100 text-green-800 dark:bg-green-800/30 dark:text-green-300",
    [TransactionType.Sell]: "bg-red-100 text-red-800 dark:bg-red-800/30 dark:text-red-300",
  }
  
  const subCategoryStyles = {
    [SubCategory.NewPosition]: "bg-blue-100 text-blue-800 dark:bg-blue-800/30 dark:text-blue-300",
    [SubCategory.Accumulation]: "bg-purple-100 text-purple-800 dark:bg-purple-800/30 dark:text-purple-300",
    [SubCategory.SellAll]: "bg-orange-100 text-orange-800 dark:bg-orange-800/30 dark:text-orange-300",
    [SubCategory.PartialSell]: "bg-yellow-100 text-yellow-800 dark:bg-yellow-800/30 dark:text-yellow-300",
  }
  
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "short",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    }).format(date)
  }
  
  const formatAddress = (address: string) => {
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`
  }
  
  const formatNumber = (num: string | number, maxDecimals = 6) => {
    return Number.parseFloat(num.toString()).toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: maxDecimals,
    })
  }
  
  export const columns: ColumnDef<TokenSwap>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && "indeterminate")}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    // {
    //   accessorKey: "transaction_type",
    //   header: "Type",
    //   cell: ({ row }) => {
    //     const type = row.getValue("transaction_type") as TransactionType
    //     return (
    //       <Badge variant={type === TransactionType.Buy ? "default" : "destructive"} className="capitalize">
    //         {type}
    //       </Badge>
    //     )
    //   },
    // },
    {
      accessorKey: "sub_category",
      header: "Type",
      cell: ({ row }) => {
        const subCategory = row.getValue("sub_category") as SubCategory
        return (
          <span className={cn("text-xs px-2 py-1 rounded-full", subCategoryStyles[subCategory])}>
            {subCategoryLabels[subCategory]}
          </span>
        )
      },
    },
    {
      accessorKey: "block_timestamp",
      header: ({ column }) => {
        return (
          <Button variant="ghost" onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}>
            Date
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => formatDate(row.getValue("block_timestamp")),
      sortingFn: "datetime",
    },
    {
      accessorKey: "pair_label",
      header: "Pair",
      cell: ({ row }) => <div>{row.getValue("pair_label")}</div>,
    },
    {
      accessorKey: "exchange_name",
      header: "Exchange",
      cell: ({ row }) => {
        const swap = row.original
        return (
          <div className="flex items-center gap-2">
            {swap.exchange_logo && (
              <img
                src={swap.exchange_logo || "/placeholder.svg"}
                alt={swap.exchange_name}
                width={16}
                height={16}
                className="rounded-full"
              />
            )}
            <span>{swap.exchange_name}</span>
          </div>
        )
      },
    },
    {
      id: "tokens",
      header: "Tokens",
      cell: ({ row }) => {
        const swap = row.original
        const isBuy = swap.transaction_type === TransactionType.Buy
  
        return (
          <div className="flex flex-col text-sm">
            <div className="flex items-center gap-1">
              <span className={isBuy ? "text-green-600 dark:text-green-400" : "text-red-600 dark:text-red-400"}>
                {isBuy ? "Bought:" : "Sold:"}
              </span>
              <div className="flex items-center gap-1">
                <span>{formatNumber(isBuy ? swap.bought.amount : swap.sold.amount, 4)}</span>
                <span className="font-medium">{isBuy ? swap.bought.symbol : swap.sold.symbol}</span>
                {(isBuy ? swap.bought.logo : swap.sold.logo) && (
                  <img
                    src={(isBuy ? swap.bought.logo : swap.sold.logo) || "/placeholder.svg"}
                    alt={isBuy ? swap.bought.symbol : swap.sold.symbol}
                    width={12}
                    height={12}
                    className="rounded-full"
                  />
                )}
              </div>
            </div>
            <div className="flex items-center gap-1">
              <span className="text-muted-foreground">For:</span>
              <div className="flex items-center gap-1">
                <span>{formatNumber(isBuy ? swap.sold.amount : swap.bought.amount, 4)}</span>
                <span className="font-medium">{isBuy ? swap.sold.symbol : swap.bought.symbol}</span>
                {(isBuy ? swap.sold.logo : swap.bought.logo) && (
                  <img
                    src={(isBuy ? swap.sold.logo : swap.bought.logo) || "/placeholder.svg"}
                    alt={isBuy ? swap.sold.symbol : swap.bought.symbol}
                    width={12}
                    height={12}
                    className="rounded-full"
                  />
                )}
              </div>
            </div>
          </div>
        )
      },
    },
    {
      accessorKey: "total_value_usd",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            className="whitespace-nowrap"
          >
            USD Value
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        )
      },
      cell: ({ row }) => {
        const amount = Number.parseFloat(row.getValue("total_value_usd"))
        const formatted = new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
        }).format(amount)
        return <div className="font-medium">{formatted}</div>
      },
      sortingFn: "basic",
    },
    {
      accessorKey: "transaction_hash",
      header: "Tx Hash",
      cell: ({ row }) => {
        return <div className="font-mono text-xs">{formatAddress(row.getValue("transaction_hash"))}</div>
      },
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        const swap = row.original
  
        return (
          <Dialog>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DialogTrigger asChild>
                  <DropdownMenuItem>View Details</DropdownMenuItem>
                </DialogTrigger>
                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(swap.transaction_hash)}>
                  Copy Tx Hash
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(swap.wallet_address)}>
                  Copy Wallet Address
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>View on Explorer</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
  
            <DialogContent className="max-w-3xl">
              <DialogHeader>
                <DialogTitle>Transaction Details</DialogTitle>
              </DialogHeader>
              <TransactionDetails swap={swap} />
            </DialogContent>
          </Dialog>
        )
      },
    },
  ]
  
  interface TransactionDetailsProps {
    swap: TokenSwap
  }
  
  function TransactionDetails({ swap }: TransactionDetailsProps) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <Badge
            variant={swap.transaction_type === TransactionType.Buy ? "default" : "destructive"}
            className="capitalize"
          >
            {swap.transaction_type}
          </Badge>
          <span className={cn("text-xs px-2 py-1 rounded-full", subCategoryStyles[swap.sub_category])}>
            {subCategoryLabels[swap.sub_category]}
          </span>
          <span className="text-sm text-muted-foreground">{formatDate(swap.block_timestamp)}</span>
        </div>
  
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2 border rounded-lg p-4">
            <h3 className="font-medium text-sm">
              {swap.transaction_type === TransactionType.Buy ? "Bought" : "Sold"} Token
            </h3>
            <div className="flex items-center gap-2">
              {(swap.transaction_type === TransactionType.Buy ? swap.bought.logo : swap.sold.logo) && (
                <img
                  src={
                    (swap.transaction_type === TransactionType.Buy ? swap.bought.logo : swap.sold.logo) ||
                    "/placeholder.svg"
                  }
                  alt={swap.transaction_type === TransactionType.Buy ? swap.bought.symbol : swap.sold.symbol}
                  width={24}
                  height={24}
                  className="rounded-full"
                />
              )}
              <span className="font-medium">
                {swap.transaction_type === TransactionType.Buy ? swap.bought.name : swap.sold.name} (
                {swap.transaction_type === TransactionType.Buy ? swap.bought.symbol : swap.sold.symbol})
              </span>
            </div>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <p className="text-muted-foreground">Amount:</p>
                <p className="font-medium">
                  {formatNumber(swap.transaction_type === TransactionType.Buy ? swap.bought.amount : swap.sold.amount)}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">USD Price:</p>
                <p className="font-medium">
                  $
                  {formatNumber(
                    swap.transaction_type === TransactionType.Buy ? swap.bought.usdPrice : swap.sold.usdPrice,
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">USD Amount:</p>
                <p className="font-medium">
                  $
                  {formatNumber(
                    swap.transaction_type === TransactionType.Buy ? swap.bought.usdAmount : swap.sold.usdAmount,
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">Token Type:</p>
                <p className="font-medium">
                  {swap.transaction_type === TransactionType.Buy ? swap.bought.tokenType : swap.sold.tokenType}
                </p>
              </div>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Address:</p>
              <p className="font-mono text-xs break-all">
                {swap.transaction_type === TransactionType.Buy ? swap.bought.address : swap.sold.address}
              </p>
            </div>
          </div>
  
          <div className="space-y-2 border rounded-lg p-4">
            <h3 className="font-medium text-sm">
              {swap.transaction_type === TransactionType.Buy ? "Sold" : "Bought"} Token
            </h3>
            <div className="flex items-center gap-2">
              {(swap.transaction_type === TransactionType.Buy ? swap.sold.logo : swap.bought.logo) && (
                <img
                  src={
                    (swap.transaction_type === TransactionType.Buy ? swap.sold.logo : swap.bought.logo) ||
                    "/placeholder.svg"
                  }
                  alt={swap.transaction_type === TransactionType.Buy ? swap.sold.symbol : swap.bought.symbol}
                  width={24}
                  height={24}
                  className="rounded-full"
                />
              )}
              <span className="font-medium">
                {swap.transaction_type === TransactionType.Buy ? swap.sold.name : swap.bought.name} (
                {swap.transaction_type === TransactionType.Buy ? swap.sold.symbol : swap.bought.symbol})
              </span>
            </div>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <p className="text-muted-foreground">Amount:</p>
                <p className="font-medium">
                  {formatNumber(swap.transaction_type === TransactionType.Buy ? swap.sold.amount : swap.bought.amount)}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">USD Price:</p>
                <p className="font-medium">
                  $
                  {formatNumber(
                    swap.transaction_type === TransactionType.Buy ? swap.sold.usdPrice : swap.bought.usdPrice,
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">USD Amount:</p>
                <p className="font-medium">
                  $
                  {formatNumber(
                    swap.transaction_type === TransactionType.Buy ? swap.sold.usdAmount : swap.bought.usdAmount,
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">Token Type:</p>
                <p className="font-medium">
                  {swap.transaction_type === TransactionType.Buy ? swap.sold.tokenType : swap.bought.tokenType}
                </p>
              </div>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Address:</p>
              <p className="font-mono text-xs break-all">
                {swap.transaction_type === TransactionType.Buy ? swap.sold.address : swap.bought.address}
              </p>
            </div>
          </div>
        </div>
  
        <div className="space-y-2 border rounded-lg p-4">
          <h3 className="font-medium text-sm">Transaction Details</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-muted-foreground text-xs">Transaction Hash:</p>
              <p className="font-mono text-xs break-all">{swap.transaction_hash}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Wallet Address:</p>
              <p className="font-mono text-xs break-all">{swap.wallet_address}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Pair Address:</p>
              <p className="font-mono text-xs break-all">{swap.pair_address}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Exchange Address:</p>
              <p className="font-mono text-xs break-all">{swap.exchange_address}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Base/Quote Price:</p>
              <p className="font-medium">{formatNumber(swap.base_quote_price)}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Total Value (USD):</p>
              <p className="font-medium">${formatNumber(swap.total_value_usd)}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Block Number:</p>
              <p className="font-medium">{swap.block_number}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Transaction Index:</p>
              <p className="font-medium">{swap.transaction_index}</p>
            </div>
          </div>
        </div>
      </div>
    )
  }