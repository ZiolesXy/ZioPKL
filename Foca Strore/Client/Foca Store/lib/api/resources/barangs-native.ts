import { serverApi } from "../core/server"

/**
 * Get barangs from server API
 */
async function getBarangs() {
  return serverApi.getEntries("/barangs", { label: "barangs" })
}

export default getBarangs
