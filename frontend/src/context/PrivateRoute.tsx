// components/PrivateRoute.tsx
import React, { useContext } from "react";
import { Navigate } from "react-router-dom";
import AuthContext from "@/context/AuthProvider";
import Loading from "@/components/Loading";

const PrivateRoute = ({ children }: { children: React.ReactElement }) => {
  const { auth, loading } = useContext(AuthContext);

  if (loading) {
    return <Loading />;
  }

  if (!auth) {
    console.log("Unauthorized access attempt");
    console.log(auth);
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default PrivateRoute;
