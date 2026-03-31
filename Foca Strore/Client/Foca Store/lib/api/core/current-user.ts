import {serverApi} from "./server"

export interface CurrentUser  {
  id: number,
  name: string,
  email: string

}

export const getCurrentUser = () => {
    return serverApi.getJson<{data: CurrentUser }>("/me")
}
