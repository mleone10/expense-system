import { Routes, Route } from 'react-router';
import { AuthProvider, ProtectedRoute } from 'hooks';
import { ProfileBar, AppFooter, AppHeader } from "components";
import { AuthCallback, Home, UnauthenticatedApp } from 'views';

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
          <Route path="/home" element={<Home />} />
        </Route>
        <Route path="/auth/callback" element={<AuthCallback />} />
      </Routes>
    </main>
  )
}

export default App;
