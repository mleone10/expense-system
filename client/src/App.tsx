import { Routes, Route } from 'react-router';
import { AuthProvider, ProtectedRoute } from 'hooks';
import { ProfileBar, MainMenu, AppFooter } from "components";
import { AuthCallback, Home, UnauthenticatedApp } from 'views';

import './App.css';
import { useState } from 'react';

function App() {
  const [isMainMenuVisible, setIsMainMenuVisible] = useState(false)

  function showMainMenu() {
    setIsMainMenuVisible(true)
  }

  function clearMainMenu() {
    setIsMainMenuVisible(false)
  }

  return (
    <AuthProvider>
      <ProfileBar showMainMenu={showMainMenu} />
      <div className="page-body">
        <MainMenu isMainMenuVisible={isMainMenuVisible} clearMainMenu={clearMainMenu} />
        <main className="bound-width">
          <Routes>
            <Route index element={<UnauthenticatedApp />} />
            <Route element={<ProtectedRoute />}>
              <Route path="/home" element={<Home />} />
            </Route>
          </Routes>
        </main>
      </div>
      <AppFooter />
    </AuthProvider>
  );
}

export default App;
