
import {
  BookOpen,
  Bot,
  Settings2,
  SquareTerminal,
  Bitcoin
} from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { NavMain } from "./nav_main";
import { NavUser } from "./nav_user";

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Dashboard",
      url: "/dashboard",
      icon: SquareTerminal,
    },
    {
      title: "Wallets",
      url: "/wallets",
      icon: Bot,
      items: [
        {
          title: "Performance",
          url: "/wallets/performance",
        },
        {
          title: "Subwallets",
          url: "wallets/subwallets",
        },
        
      ],
    },
    {
      title: "Chains",
      url: "/chains",
      icon: BookOpen,
      
    },
    {
      title: "Settings",
      url: "/settings",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "/settings/general",
        },
      ],
    },
  ],
 
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="inset" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <div className="flex">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <Bitcoin className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">Pulse</span>
                  <span className="truncate text-xs">Enterprise</span>
                </div>
                </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        {/* @ts-expect-error: this is supabase type. have to fix TODO: fix this */}
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}
