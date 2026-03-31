import { serverApi } from "../core/server"

/**
 * Get users from server API
 */
async function getUsers() {
  return serverApi.getEntries("/users/barangs", { label: "users" })
}

export default getUsers
