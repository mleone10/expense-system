import { Routes, Route, Outlet, Navigate } from 'react-router';
import { AuthProvider, useAuth } from './AuthContext';

import './App.css';
import { useEffect } from 'react';

function App() {
  return (
    <div className="app">
      <AuthProvider>
        <Routes>
          <Route element={<AppLayout />}>
            <Route path="/" element={<AuthenticatedApp />} />
            <Route path="/auth/callback" element={<AuthCallback />} />
          </Route>
        </Routes>
      </AuthProvider>
    </div>
  );
}

function AppLayout() {
  return (
    <div>
      <AppHeader />
      <Outlet />
    </div>
  )
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

function SignInButton() {
  return (
    <a href="https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=https://expense.mleone.dev/auth/callback" className="header-button">
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
          auth.signin(result.token);
        },
        (error) => {
          console.log(`Failed to exchange authorization code: ${error}`)
        }
      )
  }, [auth])

  if (auth.token) {
    return <Navigate to="/" />
  }

  return (
    <div>
      Loading...
    </div>
  )
}

function Footer() {
  return (
    <div>
      <a target="_blank" href="https://icons8.comundefined">Money</a> icon by <a target="_blank" href="https://icons8.com">Icons8</a>
    </div>
  )
}

export default App;
