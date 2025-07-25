import { createContext, useState } from "react";
import { ReactNode } from "react";


export interface AuthProps{
    username: string;
    roles: string[];
}

const AuthContext = createContext<{
    auth: AuthProps;
    setAuth: React.Dispatch<React.SetStateAction<AuthProps>>;
    }>({
        auth: {
            username: "",
            roles: [],
        },
    setAuth: () => {}
});

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
    const [auth, setAuth] = useState<AuthProps>({
        username: "",
        roles: [],
    });

    return (
        <AuthContext.Provider value={{auth,setAuth}}>
            {children}
        </AuthContext.Provider>
    )
};

export default AuthContext;