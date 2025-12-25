import { useEffect, useState } from 'react'
import './App.css'
import UserForm from './components/UserForm'
import UserList from './components/UserList'
import { useUserStore } from './store/userStore'

function App() {
  const { users, loading, error, fetchUsers } = useUserStore()
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  return (
    <div className="app">
      <header>
        <h1>User Management</h1>
      </header>
      
      <main>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : 'Add User'}
        </button>

        {showForm && <UserForm onClose={() => setShowForm(false)} />}

        {loading && <p>Loading...</p>}
        {error && <p className="error">{error}</p>}
        {!loading && users.length > 0 && <UserList users={users} />}
        {!loading && users.length === 0 && <p>No users found</p>}
      </main>
    </div>
  )
}

export default App
