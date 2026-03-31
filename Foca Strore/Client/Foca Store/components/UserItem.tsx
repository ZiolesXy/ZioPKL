import React from 'react'
import { getCurrentUser } from "@/lib/api/core/current-user"

async function UserItem() {
  const { data: user } = await getCurrentUser()

  const initials = user.name
    ?.split(" ")
    .map((n: string) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2)

  return (
    <div className='flex items-center gap-3 border rounded-[10px] p-2 w-full overflow-hidden'>
      <div className="avatar rounded-full h-10 w-10 min-h-10 min-w-10 bg-amber-400 text-white 
        flex items-center justify-center text-sm font-bold">
        {initials}
      </div>

      <div className='grow min-w-0'>
        <p className='text-sm font-bold text-slate-900 truncate leading-none mb-1'>
          {user.name}
        </p>
        <p className='text-sm text-neutral-500 truncate leading-none'>
          {user.email}
        </p>
      </div>
    </div>
  )
}

export default UserItem
        