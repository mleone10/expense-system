import { useEffect, useState } from "react";
import { SignInButton } from "components"
import { useAuth } from "hooks"

import "./ProfileBar.css"

interface Props {
  showMainMenu(): void
}

const ProfileBar = ({ showMainMenu }: Props) => {
  interface userInfoType {
    name: string;
    profileUrl: string;
  }

  const [userInfo, setUserInfo] = useState<userInfoType | undefined>(undefined)
  const isSignedIn = useAuth().getIsSignedIn();

  useEffect(() => {
    if (!isSignedIn) {
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
  }, [isSignedIn])

  if (userInfo === undefined) {
    return (
      <nav className="unauthenticated-nav">
        <SignInButton />
      </nav>
    )
  }

  return (
    <nav className="authenticated-nav">
      <svg
        className='main-menu-selector'
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 50 50"
        onClick={showMainMenu}>
        <path d="M 0 7.5 L 0 12.5 L 50 12.5 L 50 7.5 Z M 0 22.5 L 0 27.5 L 50 27.5 L 50 22.5 Z M 0 37.5 L 0 42.5 L 50 42.5 L 50 37.5 Z"></path>
      </svg>
      <span className="right-side">
        <span className="username">{userInfo.name}</span>
        <div className='dropdown-selector'>
          <img src={userInfo.profileUrl} alt="Current user" />
          <div className='dropdown-content'>
            <a href="/"><p>Sign Out</p></a>
          </div>
        </div>
      </span>
    </nav >
  )
}

export default ProfileBar;
