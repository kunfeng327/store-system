import axios from "axios";

const api = axios.create({ baseURL: "http://localhost:8080/api" });

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

export const getUser = () => JSON.parse(localStorage.getItem("user") || "null");
export const setUser = (u) => localStorage.setItem("user", JSON.stringify(u));
export const logout = () => {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
};

export default api;
