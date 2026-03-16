import axiosInstance from "../axios";
import { ApiResponse, Flight } from "@/lib/type/flight";
import { getAuthHeaders } from "../getAuth";
export const getFlights = async (page: number = 1, limit: number = 10) => {
  const response = await axiosInstance.get(`/flights?page=${page}&limit=${limit}`);
  return response.data; // Mengembalikan { success, data, meta }
};
export const getTickets = getFlights


export async function getFlightById(id: string): Promise<Flight> {
    const response = await axiosInstance.get(`/flights/${id}`)
    const result: ApiResponse<Flight> = response.data
    return result.data
}

export async function createFlight(data: any): Promise<void> {
    const headers = await getAuthHeaders();
    await axiosInstance.post("/admin/flights", data, { headers })
}
export async function deleteFlight(id: number): Promise<void> {
  const headers = await getAuthHeaders();
    await axiosInstance.delete(`/admin/flights/${id}`, { headers })
}

export async function updateFlight(id: number, data: any): Promise<void> {
  const headers = await getAuthHeaders();
  await axiosInstance.put(`/admin/flights/${id}`, data, { headers })
}

