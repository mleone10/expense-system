import { useEffect, useState } from "react";
import { SignInButton } from "components"
import { useAuth } from "hooks"

import "./ProfileBar.css"

function ProfileBar() {
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
      <nav>
        <SignInButton />
      </nav>
    )
  }

  return (
    <nav>
      <span>{userInfo.name}</span>
      <div className='dropdown-selector'>
        <img src={userInfo.profileUrl} alt="Current user" />
        <div className='dropdown-content'>
          <a href="/"><p>Sign Out</p></a>
        </div>
      </div>
    </nav >
  )
}

export default ProfileBar;
