import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { userApi, type User } from '../api/userApi'

export const useUserStore = defineStore('user', () => {
    const users = ref<User[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    const userCount = computed(() => users.value.length)

    async function fetchUsers() {
        loading.value = true
        error.value = null
        try {
            users.value = await userApi.getUsers()
        } catch (e) {
            error.value = 'Failed to fetch users'
        } finally {
            loading.value = false
        }
    }

    async function addUser(userData: Omit<User, 'id' | 'createdAt'>) {
        try {
            const newUser = await userApi.createUser(userData)
            users.value.push(newUser)
            return newUser
        } catch (e) {
            error.value = 'Failed to create user'
            throw e
        }
    }

    async function deleteUser(id: string) {
        try {
            await userApi.deleteUser(id)
            users.value = users.value.filter(u => u.id !== id)
        } catch (e) {
            error.value = 'Failed to delete user'
            throw e
        }
    }

    function getUserById(id: string) {
        return users.value.find(u => u.id === id)
    }

    return {
        users,
        loading,
        error,
        userCount,
        fetchUsers,
        addUser,
        deleteUser,
        getUserById
    }
})
