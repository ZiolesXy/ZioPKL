import { clientApi } from "../core/client"

export async function getUsers() {
  return clientApi.getArray("/users", { keys: ["users"], label: "users" })
}
