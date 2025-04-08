import { TrendingUp } from "lucide-react"
import { Pie, PieChart, LabelList } from "recharts"
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart"

export default function TokenPieChart({ data }) {
    // Early return if no data
    if (!data || !data.token_holdings || data.token_holdings.length === 0) {
        return (
            <Card className="flex flex-col h-full">
                <CardHeader className="items-center pb-0">
                    <CardTitle>Token Distribution</CardTitle>
                    <CardDescription>No token data available</CardDescription>
                </CardHeader>
                <CardContent className="flex items-center justify-center flex-1">
                    <p className="text-zinc-400">No tokens found in wallets</p>
                </CardContent>
            </Card>
        )
    }

    // Process token holdings to create chart data
    const tokenMap = data.token_holdings.reduce((acc, holding) => {
        const symbol = holding.token?.symbol || "Unknown"

        if (!acc[symbol]) {
            acc[symbol] = {
                symbol,
                value: parseFloat(holding.value),
                tokenName: holding.token?.name || "Unknown Token"
            }
        } else {
            acc[symbol].value += parseFloat(holding.value)
        }

        return acc
    }, {})

    // Convert to array for the chart
    const chartData = Object.values(tokenMap)

    // Calculate total value
    const totalValue = chartData.reduce((sum, item) => sum + item.value, 0)

    // Add percentage to each item and prepare for chart
    chartData.forEach(item => {
        item.percentage = (item.value / totalValue) * 100
        // Add required properties for the chart
        item.name = item.symbol
        item.id = item.symbol.toLowerCase()
    })

    // Sort by value (largest first)
    chartData.sort((a, b) => b.value - a.value)

    // Define colors for tokens
    const colors = [
        "var(--chart-1)",
        "var(--chart-2)",
        "var(--chart-3)",
        "var(--chart-4)",
        "var(--chart-5)"
    ]

    // Add fill colors to chart data
    chartData.forEach((item, index) => {
        item.fill = colors[index % colors.length]
    })

    // Create dynamic chart config
    const chartConfig = chartData.reduce((config, item, index) => {
        config[item.symbol.toLowerCase()] = {
            label: item.tokenName,
            color: `hsl(var(--chart-${(index % 5) + 1}))`,
        }
        return config
    }, {
        value: {
            label: "Value (USD)",
        }
    })

    // Calculate percentage change if available (using snapshots)
    let percentageChange = null
    if (data.wallets && data.wallets[0]?.snapshots?.length > 1) {
        const snapshots = [...data.wallets[0].snapshots].sort((a, b) =>
            new Date(a.created_at) - new Date(b.created_at)
        )
        const oldest = parseFloat(snapshots[0].value)
        const newest = parseFloat(snapshots[snapshots.length - 1].value)
        percentageChange = ((newest - oldest) / oldest) * 100
    }

    // Get current date for description
    const currentDate = new Date().toLocaleDateString('en-US', {
        month: 'long',
        year: 'numeric'
    })
    const isSingleToken = chartData.length === 1;
    return (
        <div className="space-y-6 p-4 h-full">
                <h3  className="text-sm font-medium text-zinc-400">Token Distribution</h3>
            <ChartContainer
                config={chartConfig}
                className="mx-auto aspect-square max-h-[250px] pb-0 [&_.recharts-pie-label-text]:fill-foreground"
            >
                <PieChart margin={{ top: 20, right: 30, bottom: 20, left: 30 }}>
                    <ChartTooltip
                        content={
                            <ChartTooltipContent
                                valueKey="value"
                                nameKey="symbol"
                                formatter={(value) => `$${parseFloat(value).toFixed(2)}`}
                            />
                        }
                    />
                    <Pie
                        data={chartData}
                        dataKey="value"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        outerRadius={75} // Reduced outer radius to leave more space for labels
                        label={(props) => {
                            const value = isSingleToken ? "100%" : `${(props.percent * 100).toFixed(0)}%`;
                            return value;
                        }}
                        labelLine={{ stroke: "var(--foreground)", strokeWidth: 1 }}
                    >
                        <LabelList
                            dataKey="symbol"
                            position="inside"
                            fill="#fff"
                            stroke="none"
                            fontSize={11} // Slightly smaller font for inside labels
                        />
                    </Pie>
                </PieChart>
            </ChartContainer>

        </div>
    )
}
