import PostCardTableAxios from "@/components/post/CardTableAxios"
import PostCardTableBarang from "@/components/post/CardTableBarang"
import PostCardTableOwnership from "@/components/post/CardTableOwnership"
import PostCardTableRole from "@/components/post/CardTableRole"
import { Suspense } from "react"
export default function UsersPage() {

  return (
    <>
      <div className="w-full p-5">
        <div className="grid min-w-0 grid-cols-1 gap-5 md:grid-cols-2">
          <Suspense fallback={<p>Loading table...</p>}>
          <div className="min-w-0 md:col-span-2">
            <PostCardTableAxios />
          </div>
          <div className="min-w-0 md:col-span-2">
            <div className="grid grid-cols-1 gap-5 md:grid-cols-3">
              <PostCardTableRole />
              <PostCardTableBarang />
              <PostCardTableOwnership />
            </div>                                      
          </div>
          </Suspense>
        </div>
      </div>
    </>
  )
}
