import { useEffect, useState } from "react";
import { SignInButton } from "components"
import { useAuth } from "hooks"

import "./AppHeader.css"

interface Props {
  showMainMenu(): void
}

// TODO: Break ProfileBar and Header back into two components
// TODO: Separate main page content into MainMenu and PageContent so that both are visible simultaneously on large screens
// TODO: Configure mobile experience such that MainMenu appears on top of PageContent on mobile
// TODO: Configure MainMenu icon to disappear on larger screens so that MainMenu is always visible
const AppHeader = ({ showMainMenu }: Props) => {
  interface userInfoType {
    name: string;
    profileUrl: string;
  }

  const auth = useAuth();
  const [userInfo, setUserInfo] = useState<userInfoType | undefined>(undefined)

  useEffect(() => {
    if (!auth.isSignedIn) {
      return
    }

    fetch(`/api/user`, {
      credentials: "include"
    }).then(response => {
      if (response.ok) {
        return response.json().then(res => res as userInfoType)
      }
    }).then(data => {
      setUserInfo(data)
    })
  }, [auth.isSignedIn])

  const unauthenticatedProfileBar = (
    <nav className="profile-bar unauthenticated-profile-bar">
      <SignInButton />
    </nav>
  )

  const authenticatedProfileBar = (
    <nav className="profile-bar authenticated-profile-bar">
      <svg
        className='main-menu-selector'
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 50 50"
        onClick={showMainMenu}>
        <path d="M 0 7.5 L 0 12.5 L 50 12.5 L 50 7.5 Z M 0 22.5 L 0 27.5 L 50 27.5 L 50 22.5 Z M 0 37.5 L 0 42.5 L 50 42.5 L 50 37.5 Z"></path>
      </svg>
      <span className="right-side">
        <span className="username">{userInfo?.name}</span>
        <img src={userInfo?.profileUrl} alt="Current user" />
      </span>
    </nav >
  )

  return (
    <header>
      {userInfo === undefined ?
        unauthenticatedProfileBar :
        authenticatedProfileBar}
    </header>
  )
}

export default AppHeader;
