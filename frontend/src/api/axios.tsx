import axios from 'axios';
import AuthContext from "@/context/AuthProvider"
import { useContext } from "react";

const baseURL = 'http://127.0.0.1:1337/api/';

export const api = axios.create({
    baseURL: baseURL,
    headers: {
        'Content-type': 'application/json',
    },
    withCredentials: true,
});

api.interceptors.response.use(
    (response) => response,
    (error) => {
      
    const { setAuth } = useContext(AuthContext);
    if (error.response?.status === 401) {
      // Unauthorized → clear auth and redirect to login
      setAuth(null);  // or use a callback to update context
      window.location.href = "/login"; // or use React Router `navigate`
    }
    return Promise.reject(error);
  }
);