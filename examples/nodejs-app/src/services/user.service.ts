import axios from 'axios';

export interface User {
    id: string;
    name: string;
    email: string;
    createdAt: Date;
}

export class UserService {
    private apiUrl = process.env.API_URL || 'http://localhost:3000';

    async getAllUsers(): Promise<User[]> {
        try {
            const response = await axios.get(`${this.apiUrl}/users`);
            return response.data;
        } catch (error) {
            console.error('Error fetching users:', error);
            throw error;
        }
    }

    async createUser(userData: Partial<User>): Promise<User> {
        try {
            const response = await axios.post(`${this.apiUrl}/users`, userData);
            return response.data;
        } catch (error) {
            console.error('Error creating user:', error);
            throw error;
        }
    }

    async getUserById(id: string): Promise<User> {
        try {
            const response = await axios.get(`${this.apiUrl}/users/${id}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching user ${id}:`, error);
            throw error;
        }
    }
}
