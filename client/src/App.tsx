import { Routes, Route } from 'react-router';
import { AuthProvider, ProtectedRoute } from 'hooks';
import { AppHeader, MainMenu, AppFooter } from "components";
import { AuthCallback, Home, UnauthenticatedApp } from 'views';

import './App.css';
import { useState } from 'react';

function App() {
  const [isMainMenuVisible, setIsMainMenuVisible] = useState(false)

  function showMainMenu() {
    setIsMainMenuVisible(true)
    console.log("clicked")
  }

  return (
    <AuthProvider>
      <AppHeader showMainMenu={showMainMenu} />
      <MainMenu isMainMenuVisible={isMainMenuVisible} />
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
