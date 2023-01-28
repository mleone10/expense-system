import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import "./Organization.css"

interface orgData {
  name: string
  id: string
  creationDate: string
  members: member[]
}

interface member {
  name: string
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
      <OrgTitle orgData={orgData} />
      <Members members={orgData?.members} />
    </React.Fragment>
  )
}

interface orgTitleProps {
  orgData: orgData | undefined
}

const OrgTitle = (props: orgTitleProps) => {
  if (props.orgData !== undefined) {
    return <h1>{props.orgData.name}</h1>
  }
  return <h1>Organization</h1>
}

interface membersProps {
  members: member[] | undefined
}

const Members = (props: membersProps) => {
  return (
    <section>
      <h2>Members</h2>
      <ul>
        {props?.members?.map((member) => {
          return <li>{member.name}</li>
        })
        }
      </ul>
    </section>
  )
}

export default Organization;
