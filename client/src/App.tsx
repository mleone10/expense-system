import { Routes, Route } from 'react-router';
import { AuthProvider, useAuth } from 'hooks';
import { ProfileBar, AppFooter, AppHeader } from "components";
import { AuthCallback } from 'views';

import './App.css';

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

export default App;
