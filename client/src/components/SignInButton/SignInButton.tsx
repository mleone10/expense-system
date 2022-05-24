import "./SignInButton.css";

function SignInButton() {
  const signInLink = `https://auth.expense.mleone.dev/login?client_id=6ka3m790cv5hrhjdqt2ju89v45&response_type=code&scope=email+openid+profile&redirect_uri=${process.env.NODE_ENV === "development" ? 'http://localhost:3000' : 'https://expense.mleone.dev'}/auth/callback`
  return (
    <a href={signInLink} className="sign-in-button">
      Sign In
    </a>
  )
}

export default SignInButton;
