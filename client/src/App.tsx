import { Routes, Route } from 'react-router';
import { AuthProvider, ProtectedRoute } from 'hooks';
import { ProfileBar, AppFooter, AppHeader } from "components";
import { AuthCallback, UnauthenticatedApp } from 'views';

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
    <main className="bound-width">
      <Routes>
        <Route index element={<UnauthenticatedApp />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/home" element={<AuthenticatedApp />} />
        </Route>
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

export default App;
