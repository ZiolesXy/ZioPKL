import { clientApi } from "../core/client"

export async function getBarangs() {
  return clientApi.getArray("/barangs", { keys: ["barangs"], label: "barangs" })
}
