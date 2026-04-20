"use client"

import * as React from "react"
import {
    BellIcon,
    CalculatorIcon,
    CalendarIcon,
    ClipboardPasteIcon,
    CodeIcon,
    CopyIcon,
    CreditCardIcon,
    FileTextIcon,
    FolderIcon,
    FolderPlusIcon,
    HelpCircleIcon,
    HomeIcon,
    ImageIcon,
    InboxIcon,
    LayoutGridIcon,
    ListIcon,
    PlusIcon,
    ScissorsIcon,
    SettingsIcon,
    TrashIcon,
    UserIcon,
    ZoomInIcon,
    Search,
    ZoomOutIcon,
} from "lucide-react"

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
  CommandShortcut,
} from "@/components/ui/command"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"

export function SearchBar() {
  const [open, setOpen] = React.useState(false)

  return (
    <div className="flex flex-col gap-4">
      <InputGroup  onClick={() => setOpen(true)} className="max-w-xs">
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
            <CommandGroup heading="Sales">
              <CommandItem>
                <HomeIcon />
                <span>Control</span>
                <CommandShortcut>⌘H</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <InboxIcon />
                <span>Add Products</span>
                <CommandShortcut>⌘I</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <FileTextIcon />
                <span>Add Coupon</span>
                <CommandShortcut>⌘D</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <FolderIcon />
                <span>Add Category</span>
                <CommandShortcut>⌘F</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <FolderIcon />
                <span>Check Transaction</span>
                <CommandShortcut>⌘F</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <FolderIcon />
                <span>Chat</span>
                <CommandShortcut>⌘F</CommandShortcut>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Navigation">
              <CommandItem>
                <PlusIcon />
                <span>User Page</span>
                <CommandShortcut>⌘N</CommandShortcut>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="System">
              <CommandItem>
                <LayoutGridIcon />
                <span>System</span>
              </CommandItem>
              <CommandItem>
                <ListIcon />
                <span>List View</span>
              </CommandItem>
              <CommandItem>
                <ZoomInIcon />
                <span>Zoom In</span>
                <CommandShortcut>⌘+</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <ZoomOutIcon />
                <span>Zoom Out</span>
                <CommandShortcut>⌘-</CommandShortcut>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Account">
              <CommandItem>
                <UserIcon />
                <span>Edit Profile</span>
                <CommandShortcut>⌘P</CommandShortcut>
              </CommandItem>
            </CommandGroup>
          </CommandList>
        </Command>
      </CommandDialog>
    </div>
  )
}
