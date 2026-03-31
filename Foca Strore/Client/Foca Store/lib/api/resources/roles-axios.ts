import { clientApi } from "../core/client"

export async function getRoles() {
  return clientApi.getArray("/roles", { keys: ["roles"], label: "roles" })
}
