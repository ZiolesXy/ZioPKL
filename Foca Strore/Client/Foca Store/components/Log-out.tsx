'use client'

import { useRouter } from 'next/navigation'
import { Button } from './ui/button'

export default function LogoutButton() {
  const router = useRouter()

  const handleLogout = async () => {
    try {
      // Memanggil API route yang sudah Anda buat
      const response = await fetch('/api/auth/logout', {
        method: 'POST',
      })

      if (response.ok) {
        router.push('/')
      }
    } catch (error) {
      console.error('Gagal logout:', error)
    }
  }

  return (
    <Button
      onClick={handleLogout}
      className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 transition"
    >
      Keluar
    </Button>
  )
}