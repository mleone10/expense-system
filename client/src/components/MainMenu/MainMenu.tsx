import { useAuth } from "hooks"
import React, { useEffect, useRef } from "react"
import { Link } from "react-router-dom"
import "./MainMenu.css"

interface Props {
  isMainMenuVisible: boolean
  clearMainMenu(): void
}

const MainMenu = ({ isMainMenuVisible, clearMainMenu }: Props) => {
  const auth = useAuth();
  const ref = useRef<HTMLElement>(null)

  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      if (isMainMenuVisible && ref.current && e.target instanceof Node && !ref.current.contains(e.target)) {
        clearMainMenu()
      }
    }

    document.addEventListener("mousedown", onClick)
    return () => {
      document.removeEventListener("mousedown", onClick)
    }
  }, [clearMainMenu, isMainMenuVisible])

  const classes = [
    "main-menu",
    isMainMenuVisible && "visible"
  ].filter(e => e).join(" ")

  return auth.isSignedIn ? (
    <nav className={classes} ref={ref}>
      <ul>
        <Link to="/"><li>Home</li></Link>
        <Link reloadDocument to="/"><li>Sign Out</li></Link>
      </ul>
    </nav >
  ) : (<React.Fragment></React.Fragment>)
}

export default MainMenu;
