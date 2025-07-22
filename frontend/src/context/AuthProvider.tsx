import { createContext, useState } from "react";
import { ReactNode } from "react";


export interface AuthProps{
    username: string;
    password: string;
    accessToken: string;
    roles: string[];
    team: string;
}

const AuthContext = createContext<{
    auth: AuthProps;
    setAuth: React.Dispatch<React.SetStateAction<AuthProps>>;
    }>({
        auth: {
            username: "admin",
            password: "admin",
            accessToken: "helo",
            roles: [],
            team: ""
        },
    setAuth: () => {}
});

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
    const [auth, setAuth] = useState<AuthProps>({
        username: "admin",
        password: "admin",
        accessToken: "helo",
        roles: ["admin"],
        team: ""
    });

    return (
        <AuthContext.Provider value={{auth,setAuth}}>
            {children}
        </AuthContext.Provider>
    )
};

export default AuthContext;