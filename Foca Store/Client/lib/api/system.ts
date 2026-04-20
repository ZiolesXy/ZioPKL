import api from "../axios"
import { getAuthHeaders } from "../auth-server"
import { clientApi } from "../client-api"
import { handleApiError } from "../utils"

export async function ResetAllNoProCa(password : string) {
    try{
      const formdata = new FormData()
      formdata.append("password", password)

      const response = await clientApi.post("/system/reset/catalog", formdata)
      return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function ResetCache(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/redis", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}
export async function Seed(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/seed/all", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function SeedProduct(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/seed/all-product", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function SeedProductWithAssets(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/seed/assets", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function SeedProductWithoutAssets(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/seed/products", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function SyncProductAssets(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/seed/sync", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function DeleteSyncAssets(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.request({
            url: "/api/admin/products/assets",
            method: "DELETE",
            data: formdata,
            headers: {}
        })
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function MigrateDatabase(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/migrate", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function ResetDatabase(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/reset", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function ResetDatabaseFull(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

        const response = await clientApi.post("/system/reset/product", formdata)
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function DeleteCloudinaryAssets(password : string) {
    try{
        const formdata = new FormData()
        formdata.append("password", password)

      const response = await clientApi.request({
            url: "/system/reset/assets",
            method: "DELETE",
            data: formdata,
            headers: {
               
            }
        })
        return response
    } catch(error: unknown){
        handleApiError(error)
    }
}

export async function GetOverview() {
    try{
        const headers = await getAuthHeaders()
        const response = await api.get("/api/admin/overview", { headers })
        return response.data
    } catch(error: unknown){
        handleApiError(error)
    }
}