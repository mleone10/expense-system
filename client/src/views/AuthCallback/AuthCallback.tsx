import { Navigate } from "react-router";
import { useEffect } from "react";
import { useAuth } from "hooks";

function AuthCallback() {
  const auth = useAuth();
  const code = new URLSearchParams(window.location.search).get("code");

  useEffect(() => {
    fetch(`/api/token?code=${code}`, {
      credentials: "include"
    }).then(response => {
      if (response.ok) {
        auth.setIsSignedIn(true);
      }
    })
  }, [code, auth])

  if (auth.isSignedIn) {
    return <Navigate to="/home" />
  } else {
    return <p>Signing you in...</p>
  }
}

export default AuthCallback;
