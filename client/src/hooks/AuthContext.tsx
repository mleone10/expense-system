import { createContext, useContext, useState, ReactNode, useEffect } from 'react';

interface AuthContextType {
  isSignedIn: boolean;
  handleSignIn(): void;
  signOut(): void;
  userInfo: userInfoType | undefined;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface userInfoType {
  name: string;
  profileUrl: string;
}

const getUserInfo = async () => {
  return fetch(`/api/user`, {
    credentials: "include"
  }).then(response => {
    if (!response.ok) {
      return undefined;
    }
    return response.json().then(data => data as userInfoType);
  })
}

function AuthProvider({ children }: { children: ReactNode }) {
  let [isSignedIn, setIsSignedIn] = useState<boolean>(false);
  let [userInfo, setUserInfo] = useState<userInfoType | undefined>(undefined);

  const handleSignIn = async () => {
    let userInfo = await getUserInfo();
    setUserInfo(userInfo);
    if (userInfo === undefined) {
      setIsSignedIn(false);
    } else {
      setIsSignedIn(true);
    }
  }

  const signOut = async () => {
    console.log("signing out");
    setUserInfo(undefined);
    setIsSignedIn(false);
  }

  useEffect(() => {
    if (!isSignedIn) {
      handleSignIn();
    }
  }, [isSignedIn])

  return (
    <AuthContext.Provider value={{ isSignedIn, handleSignIn, signOut, userInfo }}>
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
