"use client"

import { useTheme } from "next-themes"
import { useEffect, useState } from "react"
import { Moon, Sun } from "lucide-react"

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  // Avoid hydration mismatch — only render after mount
  useEffect(() => setMounted(true), [])

  if (!mounted) {
    return (
      <div className="p-2 rounded-md">
        <div className="size-4" />
      </div>
    )
  }

  const isDark = theme === "dark"

  return (
    <button
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="hover:bg-accent p-2 rounded-md cursor-pointer transition-colors"
      aria-label={`Switch to ${isDark ? "light" : "dark"} theme`}
      title={`Switch to ${isDark ? "light" : "dark"} theme`}
    >
      {isDark ? (
        <Sun className="size-4 text-foreground transition-transform duration-300 rotate-0 scale-100" />
      ) : (
        <Moon className="size-4 text-foreground transition-transform duration-300 rotate-0 scale-100" />
      )}
    </button>
  )
}
