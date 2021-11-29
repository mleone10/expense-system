import { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  token: string;
  signin: (token: string) => void;
}

const AuthContext = createContext<AuthContextType>(null!);

function AuthProvider({ children }: { children: ReactNode }) {
  let [token, setToken] = useState<string>("");

  let signin = (token: string) => {
    setToken(token);
  }

  let value = { token, signin };

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
