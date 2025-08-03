import {
  createContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { checkSession, login as loginRequest } from "@/lib/backend-interaction";
import { useNavigate } from "react-router-dom";
import { api, setUnauthorizedHandler } from "@/api/axios";

export interface AuthProps {
  username: string;
  roles: string[];
}

interface AuthContextType {
  auth: AuthProps | null;
  setAuth: React.Dispatch<React.SetStateAction<AuthProps | null>>;
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  auth: null,
  setAuth: () => {},
  login: async () => false,
  logout: () => {},
});

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [auth, setAuth] = useState<AuthProps | null>(null);
  const navigate = useNavigate();

  // 🔓 Logout function
  const logout = useCallback(async () => {
    setAuth(null);
    try {
      await api.post("/logout");
    } catch (err) {
      console.warn("Logout request error:", err);
    }
    navigate("/");
  }, [navigate]);

  // 🔐 Login function
  const login = useCallback(
    async (email: string, password: string) => {
      const response = await loginRequest({ email, password });
      if (response.status === 200) {
        const { username, role } = response.data;
        setAuth({ username, roles: [role] });
        return true;
      } else {
        await logout();
        return false;
      }
    },
    [logout]
  );

  // 📡 Session check on mount
  useEffect(() => {
    const verifySession = async () => {
      const response = await checkSession();
      if (response.status === 200) {
        const { username, role } = response.data;
        setAuth({ username, roles: [role] });
      } else {
        await logout();
      }
    };
    verifySession();
  }, [logout]);

  // 🧩 Register global 401 handler
  useEffect(() => {
    setUnauthorizedHandler(() => {
      setAuth(null);
      navigate("/login");
    });
  }, [navigate]);

  return (
    <AuthContext.Provider value={{ auth, setAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
