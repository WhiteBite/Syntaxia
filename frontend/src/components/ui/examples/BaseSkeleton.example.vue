<template>
  <div class="example-container">
    <h1 class="example-title">BaseSkeleton Examples</h1>
    <p class="example-description">
      Skeleton loading component with shimmer animation for various content types
    </p>

    <!-- Text Skeleton -->
    <section class="example-section">
      <h2 class="example-section-title">Text Skeleton</h2>
      <div class="example-demo">
        <BaseSkeleton variant="text" width="200px" />
        <BaseSkeleton variant="text" width="150px" />
        <BaseSkeleton variant="text" width="180px" />
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="text" width="200px" /&gt;
&lt;BaseSkeleton variant="text" width="150px" /&gt;
&lt;BaseSkeleton variant="text" width="180px" /&gt;</pre>
    </section>

    <!-- Multiple Lines -->
    <section class="example-section">
      <h2 class="example-section-title">Multiple Lines</h2>
      <div class="example-demo">
        <BaseSkeleton variant="text" :count="3" />
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="text" :count="3" /&gt;</pre>
    </section>

    <!-- Circle (Avatar) -->
    <section class="example-section">
      <h2 class="example-section-title">Circle (Avatar)</h2>
      <div class="example-demo">
        <div style="display: flex; gap: 16px; align-items: center;">
          <BaseSkeleton variant="circle" width="40px" height="40px" />
          <BaseSkeleton variant="circle" width="56px" height="56px" />
          <BaseSkeleton variant="circle" width="72px" height="72px" />
        </div>
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="circle" width="40px" height="40px" /&gt;
&lt;BaseSkeleton variant="circle" width="56px" height="56px" /&gt;
&lt;BaseSkeleton variant="circle" width="72px" height="72px" /&gt;</pre>
    </section>

    <!-- Rectangle -->
    <section class="example-section">
      <h2 class="example-section-title">Rectangle</h2>
      <div class="example-demo">
        <BaseSkeleton variant="rect" width="100%" height="100px" />
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="rect" width="100%" height="100px" /&gt;</pre>
    </section>

    <!-- Card -->
    <section class="example-section">
      <h2 class="example-section-title">Card</h2>
      <div class="example-demo">
        <BaseSkeleton variant="card" height="200px" />
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="card" height="200px" /&gt;</pre>
    </section>

    <!-- Without Animation -->
    <section class="example-section">
      <h2 class="example-section-title">Without Animation</h2>
      <div class="example-demo">
        <BaseSkeleton variant="text" width="200px" :animated="false" />
        <BaseSkeleton variant="circle" width="40px" height="40px" :animated="false" />
      </div>
      <pre class="example-code">
&lt;BaseSkeleton variant="text" width="200px" :animated="false" /&gt;
&lt;BaseSkeleton variant="circle" width="40px" :animated="false" /&gt;</pre>
    </section>

    <!-- Complex Example: User Card -->
    <section class="example-section">
      <h2 class="example-section-title">Complex Example: User Card</h2>
      <div class="example-demo">
        <div class="user-card-skeleton">
          <BaseSkeleton variant="circle" width="64px" height="64px" />
          <div class="user-card-skeleton__content">
            <BaseSkeleton variant="text" width="120px" height="20px" />
            <BaseSkeleton variant="text" width="180px" height="16px" />
            <BaseSkeleton variant="text" width="150px" height="16px" />
          </div>
        </div>
      </div>
      <pre class="example-code">
&lt;div class="user-card-skeleton"&gt;
  &lt;BaseSkeleton variant="circle" width="64px" height="64px" /&gt;
  &lt;div class="user-card-skeleton__content"&gt;
    &lt;BaseSkeleton variant="text" width="120px" height="20px" /&gt;
    &lt;BaseSkeleton variant="text" width="180px" height="16px" /&gt;
    &lt;BaseSkeleton variant="text" width="150px" height="16px" /&gt;
  &lt;/div&gt;
&lt;/div&gt;</pre>
    </section>

    <!-- Complex Example: List -->
    <section class="example-section">
      <h2 class="example-section-title">Complex Example: List Items</h2>
      <div class="example-demo">
        <div class="list-skeleton">
          <div v-for="i in 3" :key="i" class="list-skeleton__item">
            <BaseSkeleton variant="circle" width="32px" height="32px" />
            <div class="list-skeleton__content">
              <BaseSkeleton variant="text" width="140px" height="16px" />
              <BaseSkeleton variant="text" width="100px" height="12px" />
            </div>
          </div>
        </div>
      </div>
      <pre class="example-code">
&lt;div class="list-skeleton"&gt;
  &lt;div v-for="i in 3" :key="i" class="list-skeleton__item"&gt;
    &lt;BaseSkeleton variant="circle" width="32px" height="32px" /&gt;
    &lt;div class="list-skeleton__content"&gt;
      &lt;BaseSkeleton variant="text" width="140px" height="16px" /&gt;
      &lt;BaseSkeleton variant="text" width="100px" height="12px" /&gt;
    &lt;/div&gt;
  &lt;/div&gt;
&lt;/div&gt;</pre>
    </section>

    <!-- Specialized Skeletons -->
    <section class="example-section">
      <h2 class="example-section-title">Specialized: File Tree</h2>
      <div class="example-demo">
        <SkeletonFileTree :rows="7" />
      </div>
      <pre class="example-code">
&lt;SkeletonFileTree :rows="7" /&gt;</pre>
    </section>

    <section class="example-section">
      <h2 class="example-section-title">Specialized: Stats Cards</h2>
      <div class="example-demo">
        <SkeletonStats :cards="3" />
      </div>
      <pre class="example-code">
&lt;SkeletonStats :cards="3" /&gt;</pre>
    </section>
  </div>
</template>

<script setup lang="ts">
import { BaseSkeleton } from '@/components/ui'
import SkeletonFileTree from '@/components/SkeletonFileTree.vue'
import SkeletonStats from '@/components/SkeletonStats.vue'
</script>

<style scoped>
.example-container {
  padding: var(--space-6);
  max-width: 1200px;
  margin: 0 auto;
}

.example-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.example-description {
  font-size: var(--font-size-lg);
  color: var(--text-muted);
  margin-bottom: var(--space-8);
}

.example-section {
  margin-bottom: var(--space-8);
}

.example-section-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-4);
}

.example-demo {
  padding: var(--space-6);
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-3);
}

.example-code {
  padding: var(--space-4);
  background: var(--bg-0);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  overflow-x: auto;
}

.user-card-skeleton {
  display: flex;
  gap: var(--space-4);
  align-items: flex-start;
}

.user-card-skeleton__content {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
}

.list-skeleton {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.list-skeleton__item {
  display: flex;
  gap: var(--space-3);
  align-items: center;
}

.list-skeleton__content {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
}
</style>
