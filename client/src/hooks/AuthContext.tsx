import { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  isSignedIn: boolean;
  setIsSignedIn(state: boolean): void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

function AuthProvider({ children }: { children: ReactNode }) {
  let [isSignedIn, setIsSignedIn] = useState<boolean>(false);

  let value = { isSignedIn, setIsSignedIn };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

function useAuth() {
  const authContext = useContext(AuthContext)
  if (authContext === undefined) {
    throw new Error("auth context is undefined")
  }

  return authContext
}

export { AuthProvider, useAuth };
