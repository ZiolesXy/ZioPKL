import CardTableAxios from "@/components/delete/CardTableAxios"
import CardTableBarang from "@/components/delete/CardTableBarang"
import CardTableOwnership from "@/components/delete/CardTableOwnership"
import CardTableRole from "@/components/delete/CardTableRole"
import { Suspense } from "react"
export default function UsersPage() {

  return (
    <>
      <div className="w-full p-5">
        <div className="grid min-w-0 grid-cols-1 gap-5 md:grid-cols-2">
          <Suspense fallback={<p>Loading table...</p>}>
          <div className="min-w-0 md:col-span-2">
            <CardTableAxios />
          </div>
          <div className="min-w-0 md:col-span-2">
            <div className="grid grid-cols-1 gap-5 md:grid-cols-3 sticky">
              <CardTableRole />
              <CardTableBarang />
              <CardTableOwnership />
            </div>
          </div>
          </Suspense>
        </div>
      </div>
    </>
  )
}
