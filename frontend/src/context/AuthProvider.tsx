import {
  createContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { getSessionInfo, login as loginRequest } from "@/lib/backend-interaction";
import { api, setUnauthorizedHandler } from "@/api/axios";

// --- Types ---
type Team = {
  id: number;
  name: string;
  logo?: string;
  country?: string;
  members: TeamMember[];
  badges?: Badge[];
  score: number;
  solves?: TeamSolves[];
  bio?: string;

};

type Solve = {
  challengeId: number;
  solveTimestamp: string;
}

type Badge = {
  description: string;
  name: string;
}

type TeamSolves = {
  category: string;
  id: number;
  name: string;
  timestamp: string;
}

type User = {
  id: number;
  username: string;
  role: string;
  profilePicture?: string;
  score: number;
  email: string;
  country: string;
  joinedAt: string;
  solves: Solve[];
  teamId: number | null;
}

type TeamMember = {
  id: number;
  username: string;
  role: string;
  score: number;
}

interface AuthProps {
  id: number;
  username: string;
  role: string;
  teamId: number | null;
}

function isAuthProps(obj: any): obj is AuthProps {
  return (
    typeof obj === "object" &&
    obj !== null &&
    "id" in obj &&
    "username" in obj &&
    "role" in obj
  );
}

function isTeam(obj: any): obj is Team {
  return (
    typeof obj === "object" &&
    obj !== null &&
    typeof obj.id === "number" &&
    typeof obj.name === "string" &&
    typeof obj.score === "number" &&
    (
      obj.members === undefined ||
      (
        Array.isArray(obj.members) &&
        obj.members.every(
          (member: any) =>
            typeof member === "object" &&
            member !== null &&
            typeof member.id === "number" &&
            typeof member.username === "string" &&
            typeof member.role === "string" &&
            typeof member.score === "number"
        )
      )
    )
  );
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
const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [auth, setAuth] = useState<AuthProps | null>(null);
  const [loading, setLoading] = useState(true);

  const logout = useCallback(async () => {
    setAuth(null);
    try {
      await api.post("/logout");
      console.log("Logged out successfully");
    } catch (err) {
      console.warn("Logout request error:", err);
    }
  }, []);

  const login = useCallback(
    async (email: string, password: string) => {
      const loginResponse = await loginRequest({ email, password });
      switch (loginResponse) {
        case 200:
          return true;
        case 401:
          console.warn("Unauthorized: Invalid credentials");
          return false;
        default:
          console.error("Unexpected login response:", loginResponse);
          return false;
      }
    },
    []
  );

  useEffect(() => {
    const verifySession = async () => {
      try {
        const response = await getSessionInfo();
        isAuthProps(response) ? setAuth(response) : setAuth(null);
      } catch {
        setAuth(null);
      } finally {
        setLoading(false);
      }
    };
    verifySession();
  }, []);

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

export { AuthContext, AuthProvider, isTeam, isAuthProps };
export type { Team, AuthContextType, AuthProps, User, TeamMember, Solve };
