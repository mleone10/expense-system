import { Routes, Route, Navigate } from 'react-router';
import { AuthProvider, useAuth } from './AuthContext';

import './App.css';
import { useEffect, useState } from 'react';

function App() {
  return (
    <AuthProvider>
      <ProfileBar />
      <AppHeader />
      <AppContent />
      <AppFooter />
    </AuthProvider>
  );
}

function AppHeader() {
  return (
    <header className="bound-width">
      <h1>Expense System</h1>
      <p>Reimbursement simplified.</p>
    </header>
  )
}

function AppContent() {
  return (
    <main className=" bound-width">
      <Routes>
        <Route path="/" element={useAuth().getIsSignedIn() ? <AuthenticatedApp /> : <UnauthenticatedApp />} />
        <Route path="/auth/callback" element={<AuthCallback />} />
      </Routes>
    </main>
  )
}

function AppFooter() {
  return (
    <footer>
      <p>&copy; <a href="https://twitter.com/mleone5244">Mario Leone</a></p>
      <p>Money icon by <a href="https://icons8.com">Icons8</a></p>
    </footer>
  )
}

function SignInButton() {
  const signInLink = `https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=${process.env.NODE_ENV === "development" ? 'http://localhost:3000' : 'https://expense.mleone.dev'}/auth/callback`
  return (
    <a href={signInLink} className="sign-in-button">
      Sign In
    </a>
  )
}

function AuthenticatedApp() {
  return (
    <p>
      Welcome known user!
    </p>
  )
}

function UnauthenticatedApp() {
  return (
    <section className="unauthenticated-app">
      <p>Please sign in to continue.</p>
    </section>
  )
}

function ProfileBar() {
  interface userInfoType {
    name: string;
    profileUrl: string;
  }

  const [userInfo, setUserInfo] = useState<userInfoType | undefined>(undefined)
  const isSignedIn = useAuth().getIsSignedIn();

  useEffect(() => {
    if (!isSignedIn) {
      return
    }

    fetch(`/api/user`, {
      credentials: "include"
    }).then(response => {
      if (response.ok) {
        return response.json().then(res => res as userInfoType)
      }
    }).then(data => {
      setUserInfo(data)
    })
  }, [isSignedIn])

  if (userInfo === undefined) {
    return (
      <nav>
        <SignInButton />
      </nav>
    )
  }

  return (
    <nav>
      <span>{userInfo.name}</span>
      <div className='dropdown-selector'>
        <img src={userInfo.profileUrl} alt="Current user" />
        <div className='dropdown-content'>
          <a href="/"><p>Sign Out</p></a>
        </div>
      </div>
    </nav >
  )
}

function AuthCallback() {
  const auth = useAuth();
  const signIn = auth.signIn;
  const code = new URLSearchParams(window.location.search).get("code");

  useEffect(() => {
    fetch(`/api/token?code=${code}`, {
      credentials: "include"
    }).then(response => {
      if (response.ok) {
        signIn();
      }
    })
  }, [code, signIn])

  if (auth.getIsSignedIn()) {
    return <Navigate to="/" />
  } else {
    return <div>Something went wrong.  Please try again.</div>
  }
}

export default App;
