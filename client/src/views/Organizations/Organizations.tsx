import { useEffect, useState } from "react";
import "./Organizations.css"

interface orgs {
  orgs: org[]
}

interface org {
  name: string
  id: string
  admin: boolean
}

const getOrgs = () => {
  return fetch(`/api/orgs`, {
    credentials: "include"
  }).then(response => {
    if (response.ok) {
      return response.json().then(body => body as orgs)
    }
  })
}

const Organizations = () => {
  const [orgs, setOrgs] = useState<orgs>({ "orgs": [] })

  useEffect(() => {
    getOrgs().then(orgs => {
      if (orgs !== undefined) {
        setOrgs(orgs)
      }
    })
  }, [])

  return (
    <section className="orgs">
      <h1>Organizations</h1>
    </section>
  )
}

export default Organizations;
