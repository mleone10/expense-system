import './App.css';

function App() {
  return (
    <div className="app">
      <AppHeader />
    </div>
  );
}

function AppHeader() {
  return (
    <header className="app-header">
      <LoginButton />
    </header>
  )
}

function LoginButton() {
  return (
    <a href="https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=https://expense.mleone.dev" className="login-button">
      Login
    </a>
  )
}

export default App;
