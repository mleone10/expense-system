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

function AuthenticatedApp() {
  return (
    <div>
      Welcome known user!
    </div>
  )
}

function UnauthenticatedApp() {
  return (
    <section className="unauthenticated-app">
      <p>Please sign in to continue using the Expense System:</p>
      <SignInButton />
    </section>
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
