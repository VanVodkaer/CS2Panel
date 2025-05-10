import axios from "axios";

// 从 import.meta.env 中读取环境变量
const API_BASE = import.meta.env.VITE_API_BASE_URL;

// 创建一个 axios 实例
const api = axios.create({
  baseURL: API_BASE + "/api",
  timeout: 10000, // 根据需要设置超时时间
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器
// api.interceptors.request.use(
//   (config) => {
//     // 自动携带 token
//     const token = localStorage.getItem("token");
//     if (token) {
//       config.headers.Authorization = `Bearer ${token}`;
//     }
//     return config;
//   },
//   (error) => Promise.reject(error)
// );

// 响应拦截器
// api.interceptors.response.use(
//   (response) => response.data,
//   (error) => {
//     // 全局错误处理
//     console.error("API Error:", error);
//     return Promise.reject(error);
//   }
// );

export default api;
