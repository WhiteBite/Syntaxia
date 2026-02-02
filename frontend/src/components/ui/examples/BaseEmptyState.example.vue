<template>
  <div class="examples-container">
    <h1 class="examples-title">BaseEmptyState Examples</h1>

    <!-- Example 1: Basic with icon component -->
    <section class="example-section">
      <h2 class="example-title">1. Basic Usage</h2>
      <div class="example-preview">
        <BaseEmptyState
          :icon="FolderIcon"
          title="No files selected"
          description="Select files from the tree to build context"
        />
      </div>
    </section>

    <!-- Example 2: With action button -->
    <section class="example-section">
      <h2 class="example-title">2. With Action Button</h2>
      <div class="example-preview">
        <BaseEmptyState
          :icon="RefreshCwIcon"
          title="No data available"
          description="Click refresh to load the latest data"
        >
          <template #action>
            <BaseButton variant="primary" @click="handleRefresh">
              Refresh
            </BaseButton>
          </template>
        </BaseEmptyState>
      </div>
    </section>

    <!-- Example 3: Small size -->
    <section class="example-section">
      <h2 class="example-title">3. Small Size</h2>
      <div class="example-preview example-preview--small">
        <BaseEmptyState
          :icon="InboxIcon"
          title="No messages"
          size="sm"
        />
      </div>
    </section>

    <!-- Example 4: Large size with multiple actions -->
    <section class="example-section">
      <h2 class="example-title">4. Large Size with Multiple Actions</h2>
      <div class="example-preview">
        <BaseEmptyState
          :icon="SearchIcon"
          title="No search results"
          description="We couldn't find any files matching your search criteria. Try adjusting your filters or search terms."
          size="lg"
        >
          <template #action>
            <BaseButton variant="primary" @click="handleClearFilters">
              Clear Filters
            </BaseButton>
            <BaseButton variant="ghost" @click="handleResetSearch">
              Reset Search
            </BaseButton>
          </template>
        </BaseEmptyState>
      </div>
    </section>

    <!-- Example 5: Custom icon slot -->
    <section class="example-section">
      <h2 class="example-title">5. Custom Icon Slot</h2>
      <div class="example-preview">
        <BaseEmptyState
          title="No symbols found"
          description="This file doesn't contain any exportable symbols"
        >
          <template #icon>
            <CodeIcon class="w-12 h-12 text-purple-400" />
          </template>
        </BaseEmptyState>
      </div>
    </section>

    <!-- Example 6: No project state -->
    <section class="example-section">
      <h2 class="example-title">6. No Project State</h2>
      <div class="example-preview">
        <BaseEmptyState
          :icon="FolderOpenIcon"
          title="No project opened"
          description="Open a project to start building context"
        >
          <template #action>
            <BaseButton variant="primary" @click="handleOpenProject">
              Open Project
            </BaseButton>
          </template>
        </BaseEmptyState>
      </div>
    </section>

    <!-- Example 7: Error state -->
    <section class="example-section">
      <h2 class="example-title">7. Error State</h2>
      <div class="example-preview">
        <BaseEmptyState
          :icon="AlertCircleIcon"
          title="Failed to load data"
          description="An error occurred while loading the data. Please try again."
        >
          <template #action>
            <BaseButton variant="danger" @click="handleRetry">
              Retry
            </BaseButton>
          </template>
        </BaseEmptyState>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import BaseEmptyState from '../BaseEmptyState.vue'
import BaseButton from '../BaseButton.vue'
import { 
  FolderIcon, 
  RefreshCwIcon, 
  InboxIcon, 
  SearchIcon,
  CodeIcon,
  FolderOpenIcon,
  AlertCircleIcon
} from 'lucide-vue-next'

function handleRefresh() {
  console.log('Refresh clicked')
}

function handleClearFilters() {
  console.log('Clear filters clicked')
}

function handleResetSearch() {
  console.log('Reset search clicked')
}

function handleOpenProject() {
  console.log('Open project clicked')
}

function handleRetry() {
  console.log('Retry clicked')
}
</script>

<style scoped>
.examples-container {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  background: var(--bg-0);
  min-height: 100vh;
}

.examples-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 32px;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.example-section {
  margin-bottom: 48px;
}

.example-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}

.example-preview {
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-2xl);
  padding: 24px;
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.example-preview--small {
  min-height: 200px;
}
</style>
