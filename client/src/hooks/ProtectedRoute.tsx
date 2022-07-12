import { Outlet } from "react-router";
import { useAuth } from "hooks";
import { UnauthenticatedApp } from "views";

export const ProtectedRoute = () => {
  const auth = useAuth();

  if (!auth.isSignedIn) {
    return <UnauthenticatedApp />
  }

  return <Outlet />;
}
