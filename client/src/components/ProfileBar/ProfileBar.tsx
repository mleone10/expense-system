import { SignInButton } from "components"
import { useAuth } from "hooks"
import { useEffect, useState } from "react"

import "./ProfileBar.css"

interface Props {
  showMainMenu(): void
}

const ProfileBar = ({ showMainMenu }: Props) => {
  const [profileImageUrl, setProfileImageUrl] = useState<string | undefined>();
  const auth = useAuth();

  useEffect(() => {
    if (auth.userInfo?.profileUrl !== undefined) {
      fetch(auth.userInfo?.profileUrl)
        .then(res => {
          if (res.ok) {
            return res.blob()
          }
        })
        .then(blob => {
          if (blob !== undefined) {
            setProfileImageUrl(URL.createObjectURL(blob))
          }
        })
    }
  }, [auth.userInfo?.profileUrl])

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
        {profileImageUrl !== undefined && <img src={profileImageUrl} alt="Current user" />}
      </span>
    </header>
  )

  return !auth.isSignedIn ?
    unauthenticatedProfileBar :
    authenticatedProfileBar
}

export default ProfileBar;
