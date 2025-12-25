<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <div id="app">
    <header>
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/users">Users</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </header>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
header {
  background: #1a1a1a;
  padding: 1rem;
}

nav {
  display: flex;
  gap: 1rem;
}

nav a {
  color: #646cff;
  text-decoration: none;
}

nav a:hover {
  text-decoration: underline;
}

nav a.router-link-active {
  color: #42b883;
}

main {
  padding: 2rem;
}
</style>
