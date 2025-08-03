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

// --- Types ---
export type Team = {
  id: number;
  name: string;
};

export interface AuthProps {
  username: string;
  roles: string[];
  team: Team | null;
}

interface AuthContextType {
  auth: AuthProps | null;
  setAuth: React.Dispatch<React.SetStateAction<AuthProps | null>>;
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
  loading: boolean;
}

// --- Context ---
const AuthContext = createContext<AuthContextType>({
  auth: null,
  setAuth: () => {},
  login: async () => false,
  logout: () => {},
  loading: true,
});

// --- Provider ---
export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [auth, setAuth] = useState<AuthProps | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const logout = useCallback(async () => {
    setAuth(null);
    try {
      await api.post("/logout");
    } catch (err) {
      console.warn("Logout request error:", err);
    }
  }, []);

  const login = useCallback(
    async (email: string, password: string) => {
      const response = await loginRequest({ email, password });
      if (response.status === 200) {
        const { username, role, teamId, teamName } = response.data;
        const team = teamId && teamName ? { id: teamId, name: teamName } : null;
        setAuth({ username, roles: [role], team });
        return true;
      } else {
        await logout();
        return false;
      }
    },
    [logout]
  );

  useEffect(() => {
    const verifySession = async () => {
      try {
        const response = await checkSession();
        if (response.status === 200) {
          const { username, role, team_id, team_name } = response.data;
          console.log("Session verified:", response.data);
          const team = team_id && team_name ? { id: team_id, name: team_name } : null;
          setAuth({ username, roles: [role], team });
        } else {
          await logout();
        }
      } catch {
        await logout();
      } finally {
        setLoading(false);
      }
    };

    verifySession();
  }, [logout]);

  useEffect(() => {
    setUnauthorizedHandler(() => {
      setAuth(null);
    });
  }, []);

  return (
    <AuthContext.Provider value={{ auth, setAuth, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
