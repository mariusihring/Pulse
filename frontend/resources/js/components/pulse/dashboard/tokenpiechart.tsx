import { Pie, PieChart, LabelList } from "recharts";
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart";

export default function TokenPieChart({ data }) {
    // Early return if no data
    if (!data || !data.token_holdings || data.token_holdings.length === 0) {
        return (
            <div className="space-y-6 p-4 h-full">
                <h3 className="text-sm font-medium text-zinc-400">Token Distribution</h3>
                <p className="text-zinc-400 text-center">No tokens found in wallets</p>
            </div>
        );
    }

    // Process token holdings to create chart data
    const tokenMap = data.token_holdings.reduce((acc, holding) => {
        const symbol = holding.token?.symbol || "Unknown";
        const value = parseFloat(holding.value);

        if (!acc[symbol]) {
            acc[symbol] = {
                symbol,
                value: value,
                tokenName: holding.token?.name || "Unknown Token",
            };
        } else {
            acc[symbol].value += value;
        }

        return acc;
    }, {});

    // Convert to array and filter out tokens with 0 value Missouri
    let chartData = Object.values(tokenMap).filter(item => item.value > 0); // Filter out 0-value tokens

    // Calculate total value of non-zero tokens
    const totalValue = chartData.reduce((sum, item) => sum + item.value, 0);

    // Add percentage and other properties, filter out negligible percentages
    chartData = chartData
        .map(item => ({
            ...item,
            percentage: (item.value / totalValue) * 100,
            name: item.symbol,
            id: item.symbol.toLowerCase(),
        }))
        .filter(item => item.percentage >= 0.01); // Filter out tokens < 0.01%

    // Sort by value (largest first)
    chartData.sort((a, b) => b.value - a.value);

    // Define colors for tokens
    const colors = [
        "var(--chart-1)",
        "var(--chart-2)",
        "var(--chart-3)",
        "var(--chart-4)",
        "var(--chart-5)",
    ];

    // Add fill colors to chart data
    chartData.forEach((item, index) => {
        item.fill = colors[index % colors.length];
    });

    // Create dynamic chart config
    const chartConfig = chartData.reduce((config, item, index) => {
        config[item.symbol.toLowerCase()] = {
            label: item.tokenName,
            color: `hsl(var(--chart-${(index % 5) + 1}))`,
        };
        return config;
    }, {
        value: {
            label: "Value (USD)",
        },
    });

    const isSingleToken = chartData.length === 1;

    return (
        <div className="space-y-6 p-4 h-full">
            <h3 className="text-sm font-medium text-zinc-400">Token Distribution</h3>
            {chartData.length === 0 ? (
                <p className="text-zinc-400 text-center">No significant token holdings</p>
            ) : (
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
                            outerRadius={75}
                            label={(props) => {
                                const value = isSingleToken
                                    ? "100%"
                                    : `${(props.percent * 100).toFixed(0)}%`;
                                return value;
                            }}
                            labelLine={{ stroke: "var(--foreground)", strokeWidth: 1 }}
                        >
                            <LabelList
                                dataKey="symbol"
                                position="inside"
                                fill="#fff"
                                stroke="none"
                                fontSize={11}
                            />
                        </Pie>
                    </PieChart>
                </ChartContainer>
            )}
        </div>
    );
}
