<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { userApi, type User } from '../api/userApi'

const route = useRoute()
const router = useRouter()
const user = ref<User | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const id = route.params.id as string
    user.value = await userApi.getUserById(id)
  } catch (e) {
    error.value = 'User not found'
  } finally {
    loading.value = false
  }
})

const goBack = () => router.push('/users')
</script>

<template>
  <div class="user-detail">
    <button @click="goBack" class="back-btn">← Back</button>
    
    <div v-if="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="user" class="user-info">
      <h1>{{ user.name }}</h1>
      <p>Email: {{ user.email }}</p>
      <p>Created: {{ new Date(user.createdAt).toLocaleDateString() }}</p>
    </div>
  </div>
</template>

<style scoped>
.user-detail {
  max-width: 600px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 20px;
  padding: 8px 16px;
  cursor: pointer;
}

.user-info {
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
}

.error {
  color: #ff4444;
}
</style>
