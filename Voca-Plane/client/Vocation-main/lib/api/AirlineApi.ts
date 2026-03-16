import axiosInstance from "../axios";
import { getAuthHeaders } from "../getAuth";
export const getAirlines = async (page: number = 1, limit: number = 10) => {
  const response = await axiosInstance.get(`/airlines?page=${page}&limit=${limit}`);
  return response.data; // Mengembalikan { success, data, meta }
};

export async function createAirline(data: any): Promise<void> {
  const headers = await getAuthHeaders();
   const formData = new FormData();
    formData.append("name", data.name);
    formData.append("code", data.code);
    if(data.logo instanceof File){
        formData.append("logo", data.logo); 
    }
    await axiosInstance.post("/admin/airlines", formData, { headers: { ...headers, "Content-Type": "multipart/form-data" } });
  }

export async function updateAirline(id: number, data: any): Promise<void> {
  const headers = await getAuthHeaders();
  const formData = new FormData();
    formData.append("name", data.name);
    formData.append("code", data.code);
    if(data.logo instanceof File){
        formData.append("logo", data.logo); 
    }
    await axiosInstance.put(`/admin/airlines/${id}`, formData, { headers: { ...headers, "Content-Type": "multipart/form-data" } });
  }

export async function deleteAirlines(id: number) {
  const headers = await getAuthHeaders();    
  const response = await axiosInstance.delete(`/admin/airlines/${id}`, { headers });
  return response.data;
}

