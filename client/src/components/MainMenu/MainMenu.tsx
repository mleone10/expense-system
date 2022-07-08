import { useEffect, useRef } from "react"
import "./MainMenu.css"

interface Props {
  isMainMenuVisible: boolean
  clearMainMenu(): void
}

const MainMenu = ({ isMainMenuVisible, clearMainMenu }: Props) => {
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

  return (
    <nav className={classes} ref={ref}>
      <ul>
        <li>Foo</li>
        <li>Bar</li>
        <li>Fizz</li>
        <li>Buzz</li>
        <li>VeryLongMainMenuItem</li>
      </ul>
    </nav >
  )
}

export default MainMenu;
