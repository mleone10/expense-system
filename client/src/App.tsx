import { Routes, Route, Navigate } from 'react-router';
import { AuthProvider, useAuth } from './AuthContext';

import './App.css';
import { useEffect } from 'react';

function App() {
  return (
    <AuthProvider>
      <AppHeader />
      <AppContent />
      <AppFooter />
    </AuthProvider>
  );
}

function AppHeader() {
  return (
    <header className="bound-width">
      <div className="header-block">
        <h1 className="header-title">Expense System</h1>
        <p className="header-subtitle">
          Reimbursement simplified.
        </p>
      </div>
      <div className="header-block">
        {useAuth().getIsSignedIn() ? <SignOutButton /> : <SignInButton />}
      </div>
    </header>
  )
}

function AppContent() {
  return (
    <div className="app-content bound-width">
      <Routes>
        <Route path="/" element={useAuth().getIsSignedIn() ? <AuthenticatedApp /> : <UnauthenticatedApp />} />
        <Route path="/auth/callback" element={<AuthCallback />} />
      </Routes>
    </div>
  )
}

function AppFooter() {
  return (
    <footer className="bound-width">
      <p>Copyright &copy; 2021 <a href="https://twitter.com/mleone5244">Mario Leone</a></p>
      <p>Money icon by <a href="https://icons8.com">Icons8</a></p>
    </footer>
  )
}

function SignInButton() {
  const signInLink = `https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=${process.env.NODE_ENV === "development" ? 'http://localhost:3000' : 'https://expense.mleone.dev'}/auth/callback`
  return (
    <a href={signInLink} className="header-button">
      Sign In
    </a>
  )
}

function SignOutButton() {
  return (
    <button className="header-button" onClick={useAuth().signOut}>
      Sign Out
    </button>
  )
}

function AuthenticatedApp() {
  return (
    <div>
      Welcome known user!
    </div>
  )
}

function UnauthenticatedApp() {
  return (
    <div>
      Welcome stranger!
    </div>
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
