import { createContext, useState } from "react";
import { ReactNode } from "react";

const AuthContext = createContext({
    auth: {
        username: "admin",
        password: "admin",
        accessToken: "helo",
        roles:[] as string[]
    },
    setAuth: (_auth: any) => {}
});


interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
    const [auth,setAuth] = useState({
        username: "admin",
        password: "admin",
        accessToken: "helo",
        roles: ["admin"]
    });

    return (
        <AuthContext.Provider value={{auth,setAuth}}>
            {children}
        </AuthContext.Provider>
    )
};

export default AuthContext;