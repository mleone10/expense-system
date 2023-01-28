import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import "./Organization.css"

interface orgData {
  name: string
  id: string
  creationDate: string
  members: string[]
}

const getOrg = async (orgId: string) => {
  return fetch(`/api/orgs/${orgId}`, {
    credentials: "include"
  }).then(response => {
    if (response.ok) {
      return response.json().then(body => body as orgData)
    }
  })
}

const Organization = () => {
  const [orgData, setOrg] = useState<orgData>()
  const { orgId } = useParams()

  const loadOrg = () => {
    if (orgId === undefined) {
      return
    }

    getOrg(orgId).then(org => {
      if (org !== undefined) {
        console.log(org)
        setOrg(org)
      }
    })
  }

  useEffect(() => {
    loadOrg()
    // eslint-disable-next-line
  }, [])

  return (
    <React.Fragment>
      <h1>Organization</h1>
      <p>{orgData?.name}</p>
    </React.Fragment>
  )
}

export default Organization;
