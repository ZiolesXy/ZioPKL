import api from "../../axios"
import { extractArrayFromPayload } from "./extract-array"

export const clientApi = {
  async getJson<TResponse = unknown>(path: string): Promise<TResponse> {
    const res = await api.get(path)
    return res.data as TResponse
  },

  async postJson<TResponse = unknown, TBody = unknown>(
    path: string,
    body: TBody
  ): Promise<TResponse> {
    const res = await api.post(path, body)
    return res.data as TResponse
  },

  async putJson<TResponse = unknown, TBody = unknown>(
    path: string,
    body: TBody
  ): Promise<TResponse> {
    const res = await api.put(path, body)
    return res.data as TResponse
  },

  async deleteJson<TResponse = unknown>(path: string): Promise<TResponse> {
    const res = await api.delete(path)
    return res.data as TResponse
  },

  async getArray<T>(
    path: string,
    options: { keys?: string[]; label: string }
  ): Promise<T[]> {
    const payload = await this.getJson<unknown>(path)
    return extractArrayFromPayload<T>(payload, options)
  },
}
