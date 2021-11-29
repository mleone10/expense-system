import { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  isSignedIn: () => boolean;
  signIn: (token: string) => void;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextType>(null!);

function AuthProvider({ children }: { children: ReactNode }) {
  let [token, setToken] = useState<string>("");

  let isSignedIn = () => {
    console.log("Evaluating token")
    return token !== "";
  }

  let signIn = (token: string) => {
    setToken(token);
  }

  let signOut = () => {
    console.log("Signing out user")
    setToken("");
  }

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
