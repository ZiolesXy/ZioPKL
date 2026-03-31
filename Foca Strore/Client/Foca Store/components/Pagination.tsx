import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import Link from "next/link"

type SearchParams = Record<string, string | undefined>

function toSearchParamsString(params: SearchParams) {
  const sp = new URLSearchParams()
  Object.entries(params).forEach(([k, v]) => {
    if (v === undefined || v === "") return
    sp.set(k, v)
  })
  const s = sp.toString()
  return s ? `?${s}` : ""
}

function pageHref(basePath: string, params: SearchParams, page: number, pageSize: number) {
  return `${basePath}${toSearchParamsString({ ...params, page: String(page), pageSize: String(pageSize) })}`
}

export function PaginationControls({
  basePath,
  searchParams,
  page,
  pageSize,
  totalItems,
}: {
  basePath: string
  searchParams?: SearchParams
  page: number
  pageSize: number
  totalItems: number
}) {
  const totalPages = Math.max(1, Math.ceil(totalItems / pageSize))
  const safePage = Math.min(Math.max(1, page), totalPages)
  const params = searchParams ?? {}

  const windowSize = 5
  const half = Math.floor(windowSize / 2)
  let start = Math.max(1, safePage - half)
  let end = Math.min(totalPages, start + windowSize - 1)
  start = Math.max(1, end - windowSize + 1)

  const showLeftEllipsis = start > 2
  const showRightEllipsis = end < totalPages - 1

  const prevHref = safePage > 1 ? pageHref(basePath, params, safePage - 1, pageSize) : undefined
  const nextHref = safePage < totalPages ? pageHref(basePath, params, safePage + 1, pageSize) : undefined

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          {prevHref ? (
            <Link href={prevHref} legacyBehavior passHref>
              <PaginationPrevious />
            </Link>
          ) : (
            <PaginationPrevious className="pointer-events-none opacity-50" />
          )}
        </PaginationItem>

        <PaginationItem>
          <Link href={pageHref(basePath, params, 1, pageSize)} legacyBehavior passHref>
            <PaginationLink isActive={safePage === 1}>1</PaginationLink>
          </Link>
        </PaginationItem>

        {showLeftEllipsis && (
          <PaginationItem>
            <PaginationEllipsis />
          </PaginationItem>
        )}

        {Array.from({ length: Math.max(0, end - start + 1) }).map((_, i) => {
          const p = start + i
          if (p === 1 || p === totalPages) return null
          return (
            <PaginationItem key={p}>
              <Link href={pageHref(basePath, params, p, pageSize)} legacyBehavior passHref>
                <PaginationLink isActive={safePage === p}>{p}</PaginationLink>
              </Link>
            </PaginationItem>
          )
        })}

        {showRightEllipsis && (
          <PaginationItem>
            <PaginationEllipsis />
          </PaginationItem>
        )}

        {totalPages > 1 && (
          <PaginationItem>
            <Link href={pageHref(basePath, params, totalPages, pageSize)} legacyBehavior passHref>
              <PaginationLink isActive={safePage === totalPages}>{totalPages}</PaginationLink>
            </Link>
          </PaginationItem>
        )}

        <PaginationItem>
          {nextHref ? (
            <Link href={nextHref} legacyBehavior passHref>
              <PaginationNext />
            </Link>
          ) : (
            <PaginationNext className="pointer-events-none opacity-50" />
          )}
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  )
}
