import axiosInstance from "../axios";
import { getAuthHeaders } from "../getAuth";
export const getAirports = async (page: number = 1, limit: number = 10) => {
    const response = await axiosInstance.get(`/airports?page=${page}&limit=${limit}`);
    return response.data;
};

export async function createAirport(data: any): Promise<void> {
    const headers = await getAuthHeaders();
    await axiosInstance.post("/admin/airports", data, { headers });
}

export async function updateAirport(id: number, data: any): Promise<void> {
    const headers = await getAuthHeaders();
    await axiosInstance.put(`/admin/airports/${id}`, data, { headers });
}

export async function deleteAirport(id: number): Promise<void> {
    const headers = await getAuthHeaders();
    await axiosInstance.delete(`/admin/airports/${id}`, { headers })
}