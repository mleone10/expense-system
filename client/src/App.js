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
    <a href="https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=https://expense.mleone.dev" className="header-button">
      Sign In
    </a>
  )
}

export default App;
