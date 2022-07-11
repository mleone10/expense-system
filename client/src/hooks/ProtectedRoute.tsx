import { Navigate, Outlet } from "react-router";
import { useAuth } from "hooks";

export const ProtectedRoute = () => {
  const auth = useAuth();

  if (!auth.isSignedIn) {
    return <Navigate to="/" />
  }

  return <Outlet />;
}
