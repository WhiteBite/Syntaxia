<template>
  <div class="examples-container">
    <h1 class="examples-title">BaseAlert Component</h1>
    <p class="examples-description">
      Alert component for displaying informational, success, warning, and error messages with optional dismiss functionality.
    </p>

    <!-- Basic Variants -->
    <section class="example-section">
      <h2 class="section-title">Basic Variants</h2>
      
      <div class="example-item">
        <h3 class="example-label">Info Alert</h3>
        <BaseAlert variant="info">
          <template #title>Information</template>
          This is an informational message to provide context or helpful details.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Success Alert</h3>
        <BaseAlert variant="success">
          <template #title>Success</template>
          Your changes have been saved successfully.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Warning Alert</h3>
        <BaseAlert variant="warning">
          <template #title>Warning</template>
          This action may have unintended consequences. Please review before proceeding.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Error Alert</h3>
        <BaseAlert variant="error">
          <template #title>Error</template>
          An error occurred while processing your request. Please try again.
        </BaseAlert>
      </div>
    </section>

    <!-- Without Title -->
    <section class="example-section">
      <h2 class="section-title">Without Title</h2>
      
      <div class="example-item">
        <BaseAlert variant="info">
          Select files from the tree on the left to build your context.
        </BaseAlert>
      </div>

      <div class="example-item">
        <BaseAlert variant="success">
          Context built successfully with 15 files.
        </BaseAlert>
      </div>
    </section>

    <!-- Dismissible -->
    <section class="example-section">
      <h2 class="section-title">Dismissible Alerts</h2>
      
      <div class="example-item">
        <BaseAlert variant="info" dismissible @dismiss="handleDismiss('info')">
          <template #title>Dismissible Info</template>
          Click the X button to dismiss this alert.
        </BaseAlert>
      </div>

      <div class="example-item">
        <BaseAlert variant="warning" dismissible @dismiss="handleDismiss('warning')">
          <template #title>Dismissible Warning</template>
          This warning can be dismissed if you understand the risks.
        </BaseAlert>
      </div>
    </section>

    <!-- Custom Icons -->
    <section class="example-section">
      <h2 class="section-title">Custom Icons</h2>
      
      <div class="example-item">
        <BaseAlert variant="info">
          <template #icon>
            <Sparkles :size="16" />
          </template>
          <template #title>Custom Icon</template>
          You can override the default icon with your own.
        </BaseAlert>
      </div>

      <div class="example-item">
        <BaseAlert variant="success">
          <template #icon>
            <Rocket :size="16" />
          </template>
          <template #title>Deployment Successful</template>
          Your application has been deployed to production.
        </BaseAlert>
      </div>
    </section>

    <!-- With Actions -->
    <section class="example-section">
      <h2 class="section-title">With Action Buttons</h2>
      
      <div class="example-item">
        <BaseAlert variant="warning">
          <template #title>Unsaved Changes</template>
          You have unsaved changes. Would you like to save them before leaving?
          <template #actions>
            <BaseButton variant="warning" size="sm" @click="handleAction('save')">
              Save Changes
            </BaseButton>
            <BaseButton variant="ghost" size="sm" @click="handleAction('discard')">
              Discard
            </BaseButton>
          </template>
        </BaseAlert>
      </div>

      <div class="example-item">
        <BaseAlert variant="error">
          <template #title>Connection Failed</template>
          Unable to connect to the server. Please check your network connection.
          <template #actions>
            <BaseButton variant="danger" size="sm" @click="handleAction('retry')">
              Retry
            </BaseButton>
          </template>
        </BaseAlert>
      </div>
    </section>

    <!-- Long Content -->
    <section class="example-section">
      <h2 class="section-title">Long Content</h2>
      
      <div class="example-item">
        <BaseAlert variant="info" dismissible>
          <template #title>Getting Started</template>
          Welcome to Syntaxia! To get started, follow these steps:
          <ul style="margin-top: 8px; padding-left: 20px;">
            <li>Select files from the file tree on the left</li>
            <li>Click "Build Context" to generate your context</li>
            <li>Use the AI chat to interact with your codebase</li>
            <li>Export your context when ready</li>
          </ul>
        </BaseAlert>
      </div>
    </section>

    <!-- Real-world Examples -->
    <section class="example-section">
      <h2 class="section-title">Real-world Examples</h2>
      
      <div class="example-item">
        <h3 class="example-label">Preset Info</h3>
        <BaseAlert variant="info">
          Saving preset with 12 selected files.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Token Limit Warning</h3>
        <BaseAlert variant="warning">
          <template #title>Token Limit Exceeded</template>
          Your context exceeds the token limit. Consider reducing the number of files or increasing the limit in settings.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Git Status</h3>
        <BaseAlert variant="success">
          <template #title>Changes Committed</template>
          Successfully committed 5 files to the repository.
        </BaseAlert>
      </div>

      <div class="example-item">
        <h3 class="example-label">Build Error</h3>
        <BaseAlert variant="error" dismissible>
          <template #title>Build Failed</template>
          Failed to build context. Please check that all selected files exist and are accessible.
        </BaseAlert>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import BaseAlert from '../BaseAlert.vue'
import { BaseButton } from '@/components/ui'
import { Sparkles, Rocket } from 'lucide-vue-next'

function handleDismiss(variant: string) {
  console.log(`${variant} alert dismissed`)
}

function handleAction(action: string) {
  console.log(`Action: ${action}`)
}
</script>

<style scoped>
.examples-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  background: var(--bg-0);
  min-height: 100vh;
}

.examples-title {
  font-size: 2rem;
  font-weight: bold;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.examples-description {
  font-size: 1rem;
  color: var(--text-muted);
  margin-bottom: 2rem;
}

.example-section {
  margin-bottom: 3rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 1.5rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--border-default);
}

.example-item {
  margin-bottom: 1.5rem;
}

.example-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>
