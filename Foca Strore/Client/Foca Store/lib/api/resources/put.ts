import { putJson } from "../core/http"

export const putData = async <T>(
  endpoint: string,
  id: number | string,
  payload: T
) => {
  return putJson(`${endpoint}/${id}`, payload)
}
