<script setup lang="ts">
import { useRouter } from 'vue-router';
import type { User } from '../api/userApi';
import { useUserStore } from '../stores/userStore';

const props = defineProps<{
  user: User
}>()

const router = useRouter()
const userStore = useUserStore()

const viewDetails = () => {
  router.push(`/users/${props.user.id}`)
}

const handleDelete = async () => {
  if (confirm('Are you sure?')) {
    await userStore.deleteUser(props.user.id)
  }
}
</script>

<template>
  <div class="user-card" @click="viewDetails">
    <h3>{{ user.name }}</h3>
    <p>{{ user.email }}</p>
    <small>{{ new Date(user.createdAt).toLocaleDateString() }}</small>
    <button @click.stop="handleDelete" class="delete-btn">Delete</button>
  </div>
</template>

<style scoped>
.user-card {
  background: #1a1a1a;
  padding: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
}

.user-card:hover {
  transform: translateY(-2px);
}

.user-card h3 {
  margin: 0 0 10px;
  color: #42b883;
}

.user-card p {
  margin: 0 0 5px;
  color: #888;
}

.delete-btn {
  margin-top: 10px;
  background: #ff4444;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn:hover {
  background: #cc0000;
}
</style>
