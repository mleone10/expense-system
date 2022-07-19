import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Organizations.css"

interface orgData {
  orgs: org[]
}

interface org {
  name: string
  id: string
  admin: boolean
}

const getOrgs = async () => {
  return fetch(`/api/orgs`, {
    credentials: "include"
  }).then(response => {
    if (response.ok) {
      return response.json().then(body => body as orgData)
    }
  })
}

interface OrgsTableProps {
  orgData: orgData
}

const OrgsTable = (props: OrgsTableProps) => {
  const navigate = useNavigate();
  const handleRowClick = (id: string) => {
    return navigate(`/orgs/${id}`)
  }

  return (
    <table className="orgs-table">
      <thead>
        <tr>
          <th>Org Name</th>
          <th>Role</th>
        </tr>
      </thead>
      <tbody>
        {props.orgData.orgs.sort((a, b) => a.name > b.name ? 1 : -1).map((org) => {
          return (
            <tr key={org.id} onClick={() => handleRowClick(org.id)}>
              <td>{org.name}</td>
              <td>{org.admin ? "Admin" : "Member"}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  )
}

interface newOrg {
  name: string
}

interface newOrgResponse {
  id: string
}

const createOrg = async (org: newOrg) => {
  return fetch(`/api/orgs`, {
    method: "POST",
    body: JSON.stringify({
      name: org.name
    })
  }).then(response => {
    if (response.ok) {
      return response.json().then(res => res as newOrgResponse)
    }
  })
}

interface CreateOrgFormProps {
  onCreateOrg(): void
}

const CreateOrgForm = (props: CreateOrgFormProps) => {
  const [orgName, setOrgName] = useState<string>("")

  const handleOrgNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setOrgName(event.target.value)
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    createOrg({ name: orgName }).then(res => {
      console.log(res)
      if (res !== undefined) {
        setOrgName("")
      }
    }).then(() => {
      props.onCreateOrg()
    })
    event.preventDefault()
  }

  return (
    <section className="create-org-form">
      <h2>Create New Organization</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Organization Name:
          <input type="text" value={orgName} onChange={handleOrgNameChange} />
        </label>
        <input className="submit" type="submit" value="Create Org" />
      </form>
    </section>
  )
}

const Organizations = () => {
  const [orgData, setOrgs] = useState<orgData>({ "orgs": [] })

  const loadOrgs = () => {
    getOrgs().then(orgs => {
      if (orgs !== undefined) {
        setOrgs(orgs)
      }
    })
  }

  useEffect(() => {
    loadOrgs()
  }, [])

  return (
    <React.Fragment>
      <h1>Organizations</h1>
      <OrgsTable orgData={orgData} />
      <CreateOrgForm onCreateOrg={loadOrgs} />
    </React.Fragment>
  )
}

export default Organizations;
