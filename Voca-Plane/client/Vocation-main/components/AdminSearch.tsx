"use client"

import * as React from "react"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command"
import {
  CalculatorIcon,
  CalendarDays,
  CopyIcon,
  CreditCardIcon,
  FolderPlusIcon,
  LayoutGridIcon,
  ListIcon,
  PlusIcon,
  SettingsIcon,
  UserIcon,
} from "lucide-react"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"
import { Search } from "lucide-react"
import { useRouter } from "next/navigation"

export function AdminSearch() {
  const router = useRouter()
  const [open, setOpen] = React.useState(false)

  const navigate = (path: string) => {
    setOpen(false)
    router.push(path)
  }
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen((open) => !open)
      }
    }

    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [])


  return (
    <div className="flex flex-col gap-4">
      <InputGroup onClick={() => setOpen(true)} className="w-fit">
        <InputGroupInput placeholder="Search..." />
        <InputGroupAddon>
          <Search />
        </InputGroupAddon>
        <InputGroupAddon align="inline-end">12 results</InputGroupAddon>
      </InputGroup>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <Command>
          <CommandInput placeholder="Type a command or search..." />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup heading="Main">
              <CommandItem onSelect={() => navigate("/dashboard")}>
                <LayoutGridIcon />
                <span>Dashboard</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("/flights-schedule")}>
                <CalendarDays />
                <span>Flight Schedule</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Flight Operations">
              <CommandItem onSelect={() => navigate("flights")}>
                <PlusIcon />
                <span>Flight Master</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("airlines")}>
                <FolderPlusIcon />
                <span>Airline</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("airports")}>
                <CopyIcon />
                <span>Airport</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Sales Marketing">
              <CommandItem onSelect={() => navigate("transactions")}>
                <LayoutGridIcon />
                <span>Transaction</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("promo-codes")}>
                <ListIcon />
                <span>Promo Code</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="User Management">
              <CommandItem onSelect={() => navigate("users-monitoring")}>
                <UserIcon />
                <span>User Access</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("admin-tools")}>
                <CreditCardIcon />
                <span>Admin Tools</span>
              </CommandItem>

              <CommandItem onSelect={() => navigate("profile")}>
                <SettingsIcon />
                <span>Profile</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="System">
              <CommandItem onSelect={() => navigate("settings")}>
                <CalculatorIcon />
                <span>Settings</span>
              </CommandItem>
            </CommandGroup>
          </CommandList>
        </Command>
      </CommandDialog>
    </div >
  )
}
