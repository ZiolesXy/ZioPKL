import { extractEntries } from "./extract-entries"
import { serverFetchWithAuth } from "./server-fetch-with-auth"

export const serverApi = {
  async getJson<TResponse = unknown>(path: string): Promise<TResponse> {
    const res = await serverFetchWithAuth(path)

    if (!res.ok) {
      throw new Error(`Failed to fetch ${path}: ${res.status} ${res.statusText}`)
    }

    const json = await res.json().catch(() => null)
    return json as TResponse
  },

 
  async getEntries<T>(
    path: string,
    options: { label: string }
  ): Promise<T[]> {
    const payload = await this.getJson<unknown>(path)
    return extractEntries<T>(payload, options.label)
  },
}
