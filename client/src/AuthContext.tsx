import { createContext, useContext, useState, useCallback, ReactNode } from 'react';

interface AuthContextType {
  isSignedIn: () => boolean;
  signIn: (token: string) => void;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextType>(null!);

function AuthProvider({ children }: { children: ReactNode }) {
  let [token, setToken] = useState<string>("");

  let isSignedIn = useCallback(
    () => {
      return token !== "";
    }, [token]);

  let signIn = useCallback(
    (token: string) => {
      setToken(token);
    }, []);

  let signOut = useCallback(
    () => {
      setToken("");
    }, []);

  let value = { isSignedIn, signIn, signOut };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

function useAuth() {
  return useContext(AuthContext);
}

export { AuthProvider, useAuth };
