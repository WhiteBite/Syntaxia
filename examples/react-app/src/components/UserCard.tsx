import { User, useUserStore } from '../store/userStore'
import './UserCard.css'

interface UserCardProps {
  user: User
}

export default function UserCard({ user }: UserCardProps) {
  const deleteUser = useUserStore((state) => state.deleteUser)

  const handleDelete = async () => {
    if (confirm('Are you sure?')) {
      await deleteUser(user.id)
    }
  }

  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <small>{new Date(user.createdAt).toLocaleDateString()}</small>
      <button onClick={handleDelete} className="delete-btn">
        Delete
      </button>
    </div>
  )
}
