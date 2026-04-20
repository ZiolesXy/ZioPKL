import axios from "axios"

// Server-only Axios instance — digunakan oleh API routes dan server components
// Client-side code menggunakan /api/proxy/* sebagai gateway
const api = axios.create({
  baseURL: process.env.API_URL || "http://localhost:8080",
  withCredentials: true,
})

export default api
