import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

const client = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json'
    }
})

export interface User {
    id: string
    name: string
    email: string
    createdAt: string
}

export const userApi = {
    async getUsers(): Promise<User[]> {
        const { data } = await client.get('/users')
        return data
    },

    async getUserById(id: string): Promise<User> {
        const { data } = await client.get(`/users/${id}`)
        return data
    },

    async createUser(user: Omit<User, 'id' | 'createdAt'>): Promise<User> {
        const { data } = await client.post('/users', user)
        return data
    },

    async updateUser(id: string, user: Partial<User>): Promise<User> {
        const { data } = await client.put(`/users/${id}`, user)
        return data
    },

    async deleteUser(id: string): Promise<void> {
        await client.delete(`/users/${id}`)
    }
}
