import axios from "axios";

const baseURL = "/api";

export const api = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // send cookies
});

// Utility to read cookie by name
function getCookie(name: string): string | null {
  return (
    document.cookie
      .split("; ")
      .find((row) => row.startsWith(name + "="))
      ?.split("=")[1] ?? null
  );
}

// Attach CSRF token automatically
api.interceptors.request.use((config) => {
  const csrf = getCookie("csrf_");
  if (csrf) {
    config.headers["X-CSRF-Token"] = csrf;
  }
  return config;
});

// Optional unauthorized handler storage
let unauthorizedHandler: (() => void) | null = null;

export function setUnauthorizedHandler(handler: () => void) {
  unauthorizedHandler = handler;
}

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && unauthorizedHandler) {
      unauthorizedHandler();
    }
    return Promise.reject(error);
  },
);
