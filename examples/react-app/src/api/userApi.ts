import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

const client = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
})

export interface UserData {
    name: string
    email: string
}

export interface User extends UserData {
    id: string
    createdAt: string
}

export const userApi = {
    getUsers: async (): Promise<User[]> => {
        const response = await client.get('/users')
        return response.data
    },

    getUserById: async (id: string): Promise<User> => {
        const response = await client.get(`/users/${id}`)
        return response.data
    },

    createUser: async (user: UserData): Promise<User> => {
        const response = await client.post('/users', user)
        return response.data
    },

    updateUser: async (id: string, user: Partial<UserData>): Promise<User> => {
        const response = await client.put(`/users/${id}`, user)
        return response.data
    },

    deleteUser: async (id: string): Promise<void> => {
        await client.delete(`/users/${id}`)
    },
}
