import {
  createContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { checkSession, login as loginRequest } from "@/lib/backend-interaction";
import { useNavigate } from "react-router-dom";
import { api } from "@/api/axios";

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

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [auth, setAuth] = useState<AuthProps | null>(null);
  const navigate = useNavigate();

    // Logout function
    const logout = useCallback(async () => {
        setAuth(null);
        const response = await api.post("/logout", {}, { withCredentials: true });
        if (response.status === 200) {
        console.log("Logout successful");
        } else {
        console.error("Logout failed", response);
        }
        navigate("/");
    }, [navigate]);  
  
    // Login function
    const login = useCallback(async (email: string, password: string) => {
        const response = await loginRequest({ email, password });
        if (response.status === 200) {
        const { username, role } = response.data;
        setAuth({ username, roles: [role] });
        navigate("/challenges");
        return true;
        } else {
        logout();
        return false;
        }
    }, [navigate]);

    // On app load: check if session is valid
    useEffect(() => {
        const verifySession = async () => {
        const response = await checkSession();
        if (response.status === 200) {
            const { username, role } = response.data;
            setAuth({ username, roles: [role] });
        } else {
            logout();
        }
        };
        verifySession();
    }, []);

  return (
    <AuthContext.Provider value={{ auth, setAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;