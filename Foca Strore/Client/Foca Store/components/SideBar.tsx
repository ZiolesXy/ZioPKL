import { CircleFadingArrowUp, HatGlasses, Logs, Settings, Trash, Upload, User } from "lucide-react"
import { Command, CommandGroup,  CommandItem, CommandList, CommandSeparator } from "./ui/command"
import UserItem from "./UserItem"
import Link from "next/link"
import LogoutButton from "./Log-out"
function Sidebar() {
    const menuList = [
        {
            group: 'General',
            items: [
                {
                    link: '/dashboard',
                    icon: <User />,
                    text: 'Dashboard'
                },
                {
                    link: '/delete',
                    icon: <Trash />,
                    text: 'Delete'
                },
                {
                    link: '/post',
                    icon: <Upload />,
                    text: 'Post'
                },
                {
                    link: '/put',
                    icon: <CircleFadingArrowUp />,
                    text: 'Put'
                }
            ]
        },
        {
            group: 'Setting',
            items: [
                {
                    link: '/setting',
                    icon: <Settings />,
                    text: 'General Setting'
                },
                {
                    link: '/privacy',
                    icon: <HatGlasses />,
                    text: 'Privacy'
                },
                {
                    link: '/logs',
                    icon: <Logs />,
                    text: 'logs'
                }
            ]
        }
    ]
    return (

        <div className="sticky top-0 flex flex-col w-75 min-w-75 p-4 h-screen">
            <UserItem />
            <div className="grow">
                <Command style={{ overflow: "visible" }}>
                    <CommandList style={{ overflow: "visible" }}>
                        {menuList.map((menu: any, key: number) => (
                            <div key={key}>
                                {menu.group === 'Setting' && <CommandSeparator className="my-2" />}
                                <CommandGroup heading={menu.group}>
                                    {menu.items.map((option: any, optionKey: number) => (
                                        <CommandItem key={optionKey} className="flex gap-2 cursor-pointer p-0">
                                            <Link href={option.link} className="flex gap-2 items-center w-full p-2">
                                                {option.icon}
                                                {option.text}
                                            </Link>
                                        </CommandItem>
                                    ))}
                                </CommandGroup>
                            </div>
                        ))}
                    </CommandList>
                </Command>
            </div>
            <div className="grid grid-cols-1 gap-5">
                Settings
                <LogoutButton />
            </div>
           
        </div>
    )
}

export default Sidebar
