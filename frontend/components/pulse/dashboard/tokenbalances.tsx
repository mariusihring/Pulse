import { cn } from "@/lib/utils"
import { ArrowUpRight, ArrowDownRight } from "lucide-react"
import Image from "next/image"

interface Token {
    id: string
    name: string
    symbol: string
    amount: string
    currentPrice: string
    pnl: string
    currentValue: string
    change24h: string
    image: string
}

interface TokenBalancesProps {
    tokens: Token[]
    className?: string
}

export default function TokenBalances({ tokens, className }: TokenBalancesProps) {
    return (
        <div className={cn("w-full overflow-x-auto", className)}>
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-zinc-800 dark:text-gray-400">
                <tr>
                    <th scope="col" className="px-6 py-3">
                        Token
                    </th>
                    <th scope="col" className="px-6 py-3">
                        Amount
                    </th>
                    <th scope="col" className="px-6 py-3">
                        Price
                    </th>
                    <th scope="col" className="px-6 py-3">
                        24h
                    </th>
                    <th scope="col" className="px-6 py-3">
                        PNL
                    </th>
                    <th scope="col" className="px-6 py-3">
                        Value
                    </th>
                </tr>
                </thead>
                <tbody>
                {tokens.map((token) => (
                    <tr key={token.id} className="bg-white border-b dark:bg-zinc-900 dark:border-gray-700">
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            <div className="flex items-center">
                                <Image
                                    src={token.image || "/placeholder.svg"}
                                    alt={token.name}
                                    width={24}
                                    height={24}
                                    className="mr-2"
                                    unoptimized={true}
                                />
                                {token.name} ({token.symbol})
                            </div>
                        </th>
                        <td className="px-6 py-4">{token.amount}</td>
                        <td className="px-6 py-4">${token.currentPrice}</td>
                        <td className="px-6 py-4">
                <span
                    className={cn(
                        "flex items-center",
                        Number.parseFloat(token.change24h) >= 0 ? "text-green-600" : "text-red-600",
                    )}
                >
                  {Number.parseFloat(token.change24h) >= 0 ? (
                      <ArrowUpRight className="w-4 h-4 mr-1" />
                  ) : (
                      <ArrowDownRight className="w-4 h-4 mr-1" />
                  )}
                    {token.change24h}
                </span>
                        </td>
                        <td className="px-6 py-4">
                <span className={cn(Number.parseFloat(token.pnl) >= 0 ? "text-green-600" : "text-red-600")}>
                  {token.pnl}
                </span>
                        </td>
                        <td className="px-6 py-4">${token.currentValue}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    )
}

