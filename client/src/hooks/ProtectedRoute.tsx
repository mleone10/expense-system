import { Navigate, Outlet } from "react-router";
import { useAuth } from "hooks";

export const ProtectedRoute = () => {
  if (!useAuth().isSignedIn) {
    return <Navigate to="/" />
  }

  return <Outlet />;
}
