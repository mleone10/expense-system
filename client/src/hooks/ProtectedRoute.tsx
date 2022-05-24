import { Navigate, Outlet } from "react-router";
import { useAuth } from "hooks";

export const ProtectedRoute = () => {
  if (!useAuth().getIsSignedIn()) {
    return <Navigate to="/" />
  }

  return <Outlet />;
}
