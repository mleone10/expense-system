import { SignInButton } from "components"
import { useAuth } from "hooks"

import "./ProfileBar.css"

interface Props {
  showMainMenu(): void
}

const ProfileBar = ({ showMainMenu }: Props) => {
  const auth = useAuth();

  const unauthenticatedProfileBar = (
    <header className="profile-bar unauthenticated-profile-bar">
      <SignInButton />
    </header>
  )

  const authenticatedProfileBar = (
    <header className="profile-bar authenticated-profile-bar">
      <svg
        className='main-menu-selector'
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 50 50"
        onClick={showMainMenu}>
        <path d="M 0 7.5 L 0 12.5 L 50 12.5 L 50 7.5 Z M 0 22.5 L 0 27.5 L 50 27.5 L 50 22.5 Z M 0 37.5 L 0 42.5 L 50 42.5 L 50 37.5 Z"></path>
      </svg>
      <span className="right-side">
        <span className="username">{auth.userInfo?.name}</span>
        <img src={auth.userInfo?.profileUrl} alt="Current user" />
      </span>
    </header>
  )

  return auth.handleSignIn === undefined ?
    unauthenticatedProfileBar :
    authenticatedProfileBar
}

export default ProfileBar;
