"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Home,
  Inbox,
  User,
  Settings,
  Store,
  ChevronDown,
  Navigation,
  BookOpenText,
} from "lucide-react";

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { group } from "console";
import { is } from "date-fns/locale";

const menuGroups = [
  {
    group: "Dashboard",
    items: [
      { title: "Overview", url: "/overview", icon: Home },
      { title: "System", url: "/system", icon: Settings },
    ],
  },
  {
    group: "Management",
    icon: Inbox,
    isCollapsible: true,
    items: [
      { title: "Add Products", url: "/add-products" },
      { title: "Add Category", url: "/add-category"},
      { title: "Add Coupon", url: "/add-coupon" },
      { title: "Check Transaction", url: "/check-transaction" },
      { title: "Product Management", url: "/product-management" },
    ],
  },
  {
    group: "Navigation",
    icon: Navigation,
    isCollapsible: true,
    items: [
      { title: "Home", url: "/" },
      { title: "Products", url: "/products" },
      { title: "Category", url: "/category" },
      { title: "Coupon", url: "/coupon" },
    ],
  },
  {
    group: "Rules",
    icon: BookOpenText,
    isCollapsible: true,
    items: [
      { title: "Privacy", url: "/privacy"},
      { title: "Refund", url: "/refund"},
      { title: "Terms", url: "/terms"},
    ]
  },
  {
    group: "User & Tools",
    icon: User,
    isCollapsible: true,
    items: [
      { title: "Admin Tools", url: "/admin-tool"},
      { title: "Chat", url: "/customer-service"},
    ],
  },
];
export function AdminSidebar() {
  const pathname = usePathname();

  return (
    <Sidebar collapsible="icon" className="border-r border-border/50">
      <SidebarHeader className="py-4">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href="/overview">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                  <Store className="size-5" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-bold text-lg">Voca Store</span>
                  <span className="text-xs text-muted-foreground">Admin Panel</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="no-scrollbar">
        {menuGroups.map((group) => {
          const GroupIcon = group.icon;
          return (
            <SidebarGroup key={group.group}>
              <SidebarGroupContent>
                <SidebarMenu>
                  {group.isCollapsible ? (
                    <Collapsible className="group/collapsible">
                      <SidebarMenuItem>
                        <CollapsibleTrigger asChild>
                          <SidebarMenuButton tooltip={group.group}>
                            {GroupIcon && <GroupIcon className="size-5" />}

                            <span>{group.group}</span>
                            <ChevronDown
                              className="ml-auto h-4 w-4 shrink-0 transition-transform duration-300 ease-in-out group-data-[state=open]/collapsible:rotate-180"
                            />
                          </SidebarMenuButton>
                        </CollapsibleTrigger>

                        <CollapsibleContent className="overflow-hidden data-[state=open]:animate-collapsible-down data-[state=closed]:animate-collapsible-up">
                          <SidebarMenuSub>
                            {group.items.map((item) => (
                              <SidebarMenuSubItem key={item.title}>
                                <SidebarMenuSubButton asChild isActive={pathname === item.url}>
                                  <Link href={item.url}>{item.title}</Link>
                                </SidebarMenuSubButton>
                              </SidebarMenuSubItem>
                            ))}
                          </SidebarMenuSub>
                        </CollapsibleContent>
                      </SidebarMenuItem>
                    </Collapsible>
                  ) : (
                    group.items.map((item) => (
                      <SidebarMenuItem key={item.title}>
                        <SidebarMenuButton asChild isActive={pathname === item.url}>
                          <Link href={item.url}>
                            <item.icon />
                            <span>{item.title}</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    ))
                  )}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          );
        })}
      </SidebarContent>

      <SidebarFooter className="p-4 border-t border-border">
        <SidebarMenu>
          <SidebarMenuItem>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}