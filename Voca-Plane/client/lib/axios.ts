import axios from "axios";

const axiosInstance = axios.create({
    // baseURL: "http://172.16.17.123:8000/api/v1"
    baseURL: "http://localhost:8000/api/v1"
})

export default axiosInstance