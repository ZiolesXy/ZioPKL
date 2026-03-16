import axiosInstance from "../axios";
import { getAuthHeaders } from "../getAuth";

export async function createTransaction(transaction: any) {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.post("/transactions", transaction, { headers });
  return response.data;
}

export async function getTransactionByUser() {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.get(`/transactions`, { headers });
  return response.data;
}

export async function getTransactionByCode(code: string) {
  const headers = await getAuthHeaders();
  const response = await axiosInstance.get(`/transactions/${code}`, { headers });
  return response.data;
}

export async function getAllTransaction(page: number = 1, limit: number = 10){
  const headers = await getAuthHeaders();
  const response = await axiosInstance.get(`/admin/transactions?page=${page}&limit=${limit}`, { headers });
  return response.data;
}