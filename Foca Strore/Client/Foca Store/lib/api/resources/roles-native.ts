import { serverApi } from "../core/server"

/**
 * Get roles from server API
 */
async function getRoles() {
  return serverApi.getEntries("/roles", { label: "roles" })
}

export default getRoles
