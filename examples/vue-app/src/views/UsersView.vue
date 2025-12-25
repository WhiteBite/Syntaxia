<script setup lang="ts">
import { onMounted } from 'vue'
import UserList from '../components/UserList.vue'
import { useUserStore } from '../stores/userStore'

const userStore = useUserStore()

onMounted(() => {
  userStore.fetchUsers()
})
</script>

<template>
  <div class="users-view">
    <h1>Users</h1>
    <p v-if="userStore.loading">Loading...</p>
    <p v-else-if="userStore.error" class="error">{{ userStore.error }}</p>
    <UserList v-else :users="userStore.users" />
  </div>
</template>

<style scoped>
.users-view h1 {
  color: #646cff;
}

.error {
  color: #ff4444;
}
</style>
