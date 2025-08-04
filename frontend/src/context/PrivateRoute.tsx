// components/PrivateRoute.tsx
import React, { useContext } from "react";
import { Navigate } from "react-router-dom";
import { AuthContext } from "@/context/AuthProvider";
import { toast } from "sonner"
import Loading from "@/components/Loading";

const PrivateRoute = ({ children }: { children: React.ReactElement }) => {
  const { auth, loading } = useContext(AuthContext);

  if (loading) {
    return <Loading />;
  }

  if (!auth) {
    toast.error("You must be logged in to view this page.");
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default PrivateRoute;
