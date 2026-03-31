import { clientApi } from "../core/client"

export interface Ownership {
  barang_id: number
  user_id: number
}

export async function getOwnership(): Promise<Ownership[]> {
  return clientApi.getArray<Ownership>("/user/barang", {
    keys: ["ownership"],
    label: "ownership",
  })
}
