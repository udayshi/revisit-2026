'use client'

import {  useEffect,useState } from "react";


type User={
    id:number;
    email:string;
    username:string;
}
const TestApi = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function fetchUsers() {
            try {
                const res = await fetch("/api/users");
                if (!res.ok) {
                    throw new Error("Failed to fetch users");
                }
                const data: User[] = await res.json();
                //setUsers(data);
            } catch (err:any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        fetchUsers();
    }, []);


async function handleClick() {
    const res = await fetch("/api/users", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      }
    });

    const data = await res.json();
    setUsers(data);
    
  }

  return <>
    <h2>User List</h2>
      <ul>
        {users.map((u) => (
          <li key={u.id}>
            {u.id} - {u.username} ({u.email})
          </li>
        ))}
      </ul>
  <button onClick={handleClick}>List Users</button>
  </>;
}

export default TestApi