import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useInitials } from "@/hooks/use-initials";
import { Token } from "@/lib/types/crypto/dashboard/user";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<Token>[] = [
    {
        accessorKey: "logo",
        cell: ({ row }) => {

            const getInitials = useInitials();
            return (

                <Avatar className="h-8 w-8 overflow-hidden rounded-full">
                    <AvatarImage src={row.getValue("logo")} alt={"logo"} />
                    <AvatarFallback className="rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                        {getInitials(row.getValue("name"))}
                    </AvatarFallback>
                </Avatar>)
        }


    },
    {
        accessorKey: "name",
        header: "Name"
    },
    {
        accessorKey: "mint"
    },
    {

        accessorKey: "pnl.totalPnL",
        header: "Total Pnl"
    },
    {

        accessorKey: "pnl.realizedPnL",
        header: "Realized Pnl"
    },
    {

        accessorKey: "pnl.unrealizedPnL",
        header: "Unrealized Pnl"
    },
    {
        accessorKey: "wallet"
    }
]
