import { getUser } from "@/lib/api/user"
import ProfileForm from "@/components/auth/update-form"
import { redirect } from "next/navigation"

export default async function UpdateProfilePage() {
  const user = await getUser()

  if (!user) {
    redirect("/login")
  }

  // Tentukan URL kembali berdasarkan role  
  const backUrl = user.role === "Admin" ? "/overview" : "/"

  return (
    <main className="container py-10 px-4 md:px-0">
      <ProfileForm user={user} backUrl={backUrl} />
    </main>
  )
}
