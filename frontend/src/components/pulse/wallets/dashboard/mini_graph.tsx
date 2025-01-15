import { Line, LineChart, ResponsiveContainer } from "recharts";

interface MiniGraphProps {
  data: { date: string; value: number }[];
  color: string;
}

export default function MiniGraph({ data, color }: MiniGraphProps) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart data={data}>
        <Line
          type="monotone"
          dataKey="value"
          stroke={color}
          strokeWidth={2}
          dot={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
