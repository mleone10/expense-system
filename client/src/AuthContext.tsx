import { createContext, useContext, useState, useCallback, ReactNode } from 'react';

interface AuthContextType {
  getIsSignedIn: () => boolean;
  signIn: () => void;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextType>(null!);

function AuthProvider({ children }: { children: ReactNode }) {
  let [isSignedIn, setIsSignedIn] = useState<boolean>(false);

  let getIsSignedIn: () => boolean = useCallback(
    () => {
      return isSignedIn
    }, [isSignedIn]
  )

  let signIn = useCallback(
    () => {
      setIsSignedIn(true);
    }, []);

  let signOut = useCallback(
    () => {
      setIsSignedIn(false);
    }, []);

  let value = { getIsSignedIn, signIn, signOut };

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
