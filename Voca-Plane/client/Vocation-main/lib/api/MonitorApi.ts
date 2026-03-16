import axiosInstance from "../axios";
import { getAuthHeaders } from "../getAuth";

export async function getUserAdmin() {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.get("/admin/users", { headers });
  return response.data;
}

export async function getInformation() {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.get("/admin/dashboard", { headers });
  return response.data;
}

export async function updateUserAdmin(id: string, data: any) {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.patch(`/admin/users/${id}/role`, data, { headers });
  return response.data;
}