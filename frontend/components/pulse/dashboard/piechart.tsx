'use client'
import { PieChart as RechartsPieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from "recharts"

interface PieChartProps {
    data: {
        name: string
        value: number
        color: string
    }[]
    totalValue: number
}

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#8884D8"]

export default function PieChart({ data, totalValue }: PieChartProps) {
    return (
        <div className="w-full h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
                <RechartsPieChart>
                    <Pie data={data} cx="50%" cy="50%" labelLine={false} outerRadius={80} fill="#8884d8" dataKey="value">
                        {data.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                    </Pie>
                    <Tooltip
                        formatter={(value: number) => `$${value.toFixed(2)} (${((value / totalValue) * 100).toFixed(2)}%)`}
                    />
                    <Legend />
                </RechartsPieChart>
            </ResponsiveContainer>
        </div>
    )
}

