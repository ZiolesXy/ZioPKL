
import PutCardTableAxios from "@/components/put/CardTableAxios"
import PutCardTableBarang from "@/components/put/CardTableBarang"
import PutCardTableOwnership from "@/components/put/CardTableOwnership"
import PutCardTableRole from "@/components/put/CardTableRole"
import { Suspense } from "react"
export default function UsersPage() {

  return (
    <>
      <div className="w-full p-5">
        <div className="grid min-w-0 grid-cols-1 gap-5 md:grid-cols-2">
          <Suspense fallback={<p>Loading table...</p>}>
          <div className="min-w-0 md:col-span-2">
            <PutCardTableAxios />
          </div>
          <div className="min-w-0 md:col-span-2">
            <div className="grid grid-cols-1 gap-5 md:grid-cols-3">
              <PutCardTableRole />
              <PutCardTableBarang />
              <PutCardTableOwnership />
            </div>
          </div>
          </Suspense>
        </div>
      </div>
    </>
  )
}
