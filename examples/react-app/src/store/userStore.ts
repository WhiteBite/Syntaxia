import { create } from 'zustand'
import { userApi } from '../api/userApi'

export interface User {
    id: string
    name: string
    email: string
    createdAt: string
}

interface UserStore {
    users: User[]
    loading: boolean
    error: string | null
    fetchUsers: () => Promise<void>
    addUser: (user: Omit<User, 'id' | 'createdAt'>) => Promise<void>
    deleteUser: (id: string) => Promise<void>
}

export const useUserStore = create<UserStore>((set) => ({
    users: [],
    loading: false,
    error: null,

    fetchUsers: async () => {
        set({ loading: true, error: null })
        try {
            const users = await userApi.getUsers()
            set({ users, loading: false })
        } catch (error) {
            set({ error: 'Failed to fetch users', loading: false })
        }
    },

    addUser: async (user) => {
        try {
            const newUser = await userApi.createUser(user)
            set((state) => ({ users: [...state.users, newUser] }))
        } catch (error) {
            set({ error: 'Failed to create user' })
        }
    },

    deleteUser: async (id) => {
        try {
            await userApi.deleteUser(id)
            set((state) => ({ users: state.users.filter((u) => u.id !== id) }))
        } catch (error) {
            set({ error: 'Failed to delete user' })
        }
    },
}))
