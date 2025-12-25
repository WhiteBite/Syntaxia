import { useState } from 'react'
import { useUserStore } from '../store/userStore'
import './UserForm.css'

interface UserFormProps {
  onClose: () => void
}

export default function UserForm({ onClose }: UserFormProps) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const addUser = useUserStore((state) => state.addUser)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await addUser({ name, email })
      setName('')
      setEmail('')
      onClose()
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="user-form" onSubmit={handleSubmit}>
      <div className="form-group">
        <label htmlFor="name">Name</label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>
      <div className="form-group">
        <label htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
      </div>
      <div className="form-actions">
        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Create'}
        </button>
        <button type="button" onClick={onClose}>
          Cancel
        </button>
      </div>
    </form>
  )
}
