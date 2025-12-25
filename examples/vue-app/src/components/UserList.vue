<script setup lang="ts">
import type { User } from '../api/userApi';
import UserCard from './UserCard.vue';

defineProps<{
  users: User[]
}>()
</script>

<template>
  <div class="user-list">
    <p v-if="users.length === 0">No users found</p>
    <div v-else class="user-grid">
      <UserCard v-for="user in users" :key="user.id" :user="user" />
    </div>
  </div>
</template>

<style scoped>
.user-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}
</style>
