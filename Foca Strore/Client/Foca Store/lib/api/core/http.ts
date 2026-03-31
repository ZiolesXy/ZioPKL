import api from "../../axios"
import { extractArrayFromPayload } from "./extract-array"

export { clientApi } from "./client"
export { serverApi } from "./server"

export async function getJson<TResponse = unknown>(path: string): Promise<TResponse> {
  const res = await api.get(path)
  return res.data as TResponse
}

export async function postJson<TResponse = unknown, TBody = unknown>(
  path: string,
  body: TBody
): Promise<TResponse> {
  const res = await api.post(path, body)
  return res.data as TResponse
}

export async function putJson<TResponse = unknown, TBody = unknown>(
  path: string,
  body: TBody
): Promise<TResponse> {
  const res = await api.put(path, body)
  return res.data as TResponse
}

export async function deleteJson<TResponse = unknown>(path: string): Promise<TResponse> {
  const res = await api.delete(path)
  return res.data as TResponse
}

export async function getArray<T>(
  path: string,
  options: { keys?: string[]; label: string }
): Promise<T[]> {
  const payload = await getJson<unknown>(path)
  return extractArrayFromPayload<T>(payload, options)
}
