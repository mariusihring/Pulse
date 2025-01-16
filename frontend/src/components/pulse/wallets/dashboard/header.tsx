import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { RefreshCw, Plus, Wallet } from "lucide-react";

export default function DashboardHeader() {
  return (
    <div className="border-b">
      <div className="container flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-4">
          <h1 className="text-2xl font-semibold">Wallets</h1>
          <div className="flex-1 md:max-w-sm">
            <Input
              placeholder="Search wallets or chains..."
              className="bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-backround/60"
            />
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Select defaultValue="market">
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="market">Highest market value</SelectItem>
              <SelectItem value="name">Name</SelectItem>
              <SelectItem value="recent">Recently updated</SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" size="icon">
            <RefreshCw className="h-4 w-4" />
          </Button>
          {/* <Button variant="outline">
            <Plus className="h-4 w-4 mr-2" />
            Add Transaction
          </Button> */}
          <Button>
            <Wallet className="h-4 w-4 mr-2" />
            Add Wallet
          </Button>
        </div>
      </div>
    </div>
  );
}
