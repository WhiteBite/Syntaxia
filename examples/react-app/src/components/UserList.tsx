import { User } from '../store/userStore'
import UserCard from './UserCard'
import './UserList.css'

interface UserListProps {
  users: User[]
}

export default function UserList({ users }: UserListProps) {
  return (
    <div className="user-list">
      <h2>Users ({users.length})</h2>
      <div className="user-grid">
        {users.map((user) => (
          <UserCard key={user.id} user={user} />
        ))}
      </div>
    </div>
  )
}
