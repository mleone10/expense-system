import { Routes, Route, Navigate } from 'react-router';
import { AuthProvider, useAuth } from './AuthContext';

import './App.css';
import { useEffect } from 'react';

function App() {
  return (
    <div className="app">
      <AuthProvider>
        <AppHeader />
        <AppContent />
        <AppFooter />
      </AuthProvider>
    </div>
  );
}

function AppHeader() {
  return (
    <header className="header">
      <div className="header-block">
        <h1 className="header-title">Expense System</h1>
        <p className="header-subtitle">
          Reimbursement simplified.
        </p>
      </div>
      <div className="header-block">
        <SignInButton />
      </div>
    </header>
  )
}

function AppContent() {
  return (
    <div className="app-content">
      <Routes>
        <Route path="/" element={<AuthenticatedApp />} />
        <Route path="/auth/callback" element={<AuthCallback />} />
      </Routes>
    </div>
  )
}

function AppFooter() {
  return (
    <footer className="app-footer">
      <p>Copyright &copy; 2021 <a href="https://twitter.com/mleone5244">Mario Leone</a></p>
      <p>Money icon by <a href="https://icons8.com">Icons8</a></p>
    </footer>
  )
}

function SignInButton() {
  const signInLink = `https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=${process.env.NODE_ENV === "development" ? 'http://localhost:3000' : 'https://expense.mleone.dev'}/auth/callback`
  return (
    <a href={signInLink} className="header-button" >
      Sign In
    </a>
  )
}

function AuthenticatedApp() {
  return (
    <div>
      Welcome
    </div>
  )
}

function AuthCallback() {
  let auth = useAuth();

  useEffect(() => {
    let code = new URLSearchParams(window.location.search).get("code");
    fetch(`/api/token?code=${code}`)
      .then(res => res.json())
      .then(
        (result) => {
          auth.signIn(result.token);
        },
        (error) => {
          console.log(`Failed to exchange authorization code: ${error}`)
        }
      )
  }, [auth])

  if (auth.isSignedIn()) {
    return <Navigate to="/" />
  }

  return (
    <div>
      Loading...
    </div>
  )
}

export default App;
