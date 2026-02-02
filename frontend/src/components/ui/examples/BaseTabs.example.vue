<template>
  <div class="examples-container">
    <h1 class="examples-title">BaseTabs Examples</h1>

    <!-- Example 1: Basic Tabs -->
    <section class="example-section">
      <h2 class="example-title">1. Basic Tabs</h2>
      <p class="example-description">Simple tabs with text labels</p>
      
      <div class="example-demo">
        <BaseTabs v-model="basicTab">
          <BaseTab name="overview" label="Overview">
            <div class="tab-content">
              <h3>Overview</h3>
              <p>This is the overview tab content. It contains general information about the project.</p>
            </div>
          </BaseTab>
          <BaseTab name="details" label="Details">
            <div class="tab-content">
              <h3>Details</h3>
              <p>Detailed information and specifications are shown here.</p>
            </div>
          </BaseTab>
          <BaseTab name="settings" label="Settings">
            <div class="tab-content">
              <h3>Settings</h3>
              <p>Configure your preferences and options in this tab.</p>
            </div>
          </BaseTab>
        </BaseTabs>
      </div>

      <div class="example-code">
        <pre><code>&lt;BaseTabs v-model="activeTab"&gt;
  &lt;BaseTab name="overview" label="Overview"&gt;
    Content for overview
  &lt;/BaseTab&gt;
  &lt;BaseTab name="details" label="Details"&gt;
    Content for details
  &lt;/BaseTab&gt;
  &lt;BaseTab name="settings" label="Settings"&gt;
    Content for settings
  &lt;/BaseTab&gt;
&lt;/BaseTabs&gt;</code></pre>
      </div>
    </section>

    <!-- Example 2: Tabs with Icons -->
    <section class="example-section">
      <h2 class="example-title">2. Tabs with Icons</h2>
      <p class="example-description">Tabs enhanced with icons for better visual recognition</p>
      
      <div class="example-demo">
        <BaseTabs v-model="iconTab">
          <BaseTab name="general" label="General" :icon="Settings">
            <div class="tab-content">
              <h3>General Settings</h3>
              <p>Configure general application settings and preferences.</p>
              <div class="setting-item">
                <label>Language:</label>
                <select class="setting-select">
                  <option>English</option>
                  <option>Русский</option>
                </select>
              </div>
            </div>
          </BaseTab>
          <BaseTab name="advanced" label="Advanced" :icon="Code">
            <div class="tab-content">
              <h3>Advanced Settings</h3>
              <p>Advanced configuration options for power users.</p>
              <div class="setting-item">
                <label>Debug Mode:</label>
                <input type="checkbox" />
              </div>
            </div>
          </BaseTab>
          <BaseTab name="about" label="About" :icon="Info">
            <div class="tab-content">
              <h3>About</h3>
              <p>Application version: 1.0.0</p>
              <p>Build date: 2024-01-15</p>
            </div>
          </BaseTab>
        </BaseTabs>
      </div>

      <div class="example-code">
        <pre><code>&lt;BaseTabs v-model="activeTab"&gt;
  &lt;BaseTab name="general" label="General" :icon="SettingsIcon"&gt;
    General content
  &lt;/BaseTab&gt;
  &lt;BaseTab name="advanced" label="Advanced" :icon="CodeIcon"&gt;
    Advanced content
  &lt;/BaseTab&gt;
&lt;/BaseTabs&gt;</code></pre>
      </div>
    </section>

    <!-- Example 3: Disabled Tab -->
    <section class="example-section">
      <h2 class="example-title">3. Disabled Tab</h2>
      <p class="example-description">Tabs can be disabled to prevent user interaction</p>
      
      <div class="example-demo">
        <BaseTabs v-model="disabledTab">
          <BaseTab name="available" label="Available" :icon="CheckCircle">
            <div class="tab-content">
              <h3>Available Feature</h3>
              <p>This feature is available and ready to use.</p>
            </div>
          </BaseTab>
          <BaseTab name="locked" label="Locked" :icon="Lock" :disabled="true">
            <div class="tab-content">
              <h3>Locked Feature</h3>
              <p>This content should not be visible.</p>
            </div>
          </BaseTab>
          <BaseTab name="coming-soon" label="Coming Soon" :icon="Clock" :disabled="true">
            <div class="tab-content">
              <h3>Coming Soon</h3>
              <p>This content should not be visible.</p>
            </div>
          </BaseTab>
        </BaseTabs>
      </div>

      <div class="example-code">
        <pre><code>&lt;BaseTabs v-model="activeTab"&gt;
  &lt;BaseTab name="available" label="Available"&gt;
    Available content
  &lt;/BaseTab&gt;
  &lt;BaseTab name="locked" label="Locked" :disabled="true"&gt;
    Locked content (not accessible)
  &lt;/BaseTab&gt;
&lt;/BaseTabs&gt;</code></pre>
      </div>
    </section>

    <!-- Example 4: Many Tabs -->
    <section class="example-section">
      <h2 class="example-title">4. Many Tabs</h2>
      <p class="example-description">Tabs with keyboard navigation support (Arrow keys, Home, End)</p>
      
      <div class="example-demo">
        <BaseTabs v-model="manyTab">
          <BaseTab name="tab1" label="Tab 1">
            <div class="tab-content">
              <h3>Tab 1</h3>
              <p>Use Arrow Left/Right keys to navigate between tabs.</p>
            </div>
          </BaseTab>
          <BaseTab name="tab2" label="Tab 2">
            <div class="tab-content">
              <h3>Tab 2</h3>
              <p>Press Home to go to the first tab.</p>
            </div>
          </BaseTab>
          <BaseTab name="tab3" label="Tab 3">
            <div class="tab-content">
              <h3>Tab 3</h3>
              <p>Press End to go to the last tab.</p>
            </div>
          </BaseTab>
          <BaseTab name="tab4" label="Tab 4">
            <div class="tab-content">
              <h3>Tab 4</h3>
              <p>Keyboard navigation skips disabled tabs.</p>
            </div>
          </BaseTab>
          <BaseTab name="tab5" label="Tab 5">
            <div class="tab-content">
              <h3>Tab 5</h3>
              <p>This is the last tab. Try pressing End key!</p>
            </div>
          </BaseTab>
        </BaseTabs>
      </div>

      <div class="example-code">
        <pre><code>&lt;!-- Keyboard Navigation --&gt;
&lt;!-- Arrow Left/Right: Navigate between tabs --&gt;
&lt;!-- Home: Jump to first tab --&gt;
&lt;!-- End: Jump to last tab --&gt;</code></pre>
      </div>
    </section>

    <!-- Example 5: Lazy Loading -->
    <section class="example-section">
      <h2 class="example-title">5. Lazy Loading Content</h2>
      <p class="example-description">Tab content is only rendered when the tab becomes active</p>
      
      <div class="example-demo">
        <BaseTabs v-model="lazyTab">
          <BaseTab name="light" label="Light Content">
            <div class="tab-content">
              <h3>Light Content</h3>
              <p>This is a simple tab with minimal content.</p>
            </div>
          </BaseTab>
          <BaseTab name="heavy" label="Heavy Content">
            <div class="tab-content">
              <h3>Heavy Content</h3>
              <p>This tab contains heavy content that is only rendered when active.</p>
              <div v-for="i in 20" :key="i" class="heavy-item">
                Item {{ i }} - This content is lazy loaded
              </div>
            </div>
          </BaseTab>
          <BaseTab name="dynamic" label="Dynamic Content">
            <div class="tab-content">
              <h3>Dynamic Content</h3>
              <p>Mounted at: {{ mountTime }}</p>
              <p>This content is created when you first visit this tab.</p>
            </div>
          </BaseTab>
        </BaseTabs>
      </div>

      <div class="example-code">
        <pre><code>&lt;!-- Content is only rendered when tab is active --&gt;
&lt;BaseTab name="heavy" label="Heavy Content"&gt;
  &lt;ExpensiveComponent /&gt;
&lt;/BaseTab&gt;</code></pre>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseTabs from '../BaseTabs.vue'
import BaseTab from '../BaseTab.vue'
import { Settings, Code, Info, CheckCircle, Lock, Clock } from 'lucide-vue-next'

// Tab states
const basicTab = ref('overview')
const iconTab = ref('general')
const disabledTab = ref('available')
const manyTab = ref('tab1')
const lazyTab = ref('light')

// Dynamic content example
const mountTime = computed(() => new Date().toLocaleTimeString())
</script>

<style scoped>
.examples-container {
  padding: var(--space-8);
  max-width: 1200px;
  margin: 0 auto;
  background: var(--bg-0);
  min-height: 100vh;
}

.examples-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--space-8);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.example-section {
  margin-bottom: var(--space-8);
  padding: var(--space-6);
  background: var(--bg-1);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-2xl);
}

.example-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-2);
}

.example-description {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
  margin-bottom: var(--space-4);
}

.example-demo {
  padding: var(--space-6);
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-4);
}

.example-code {
  background: var(--bg-0);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  overflow-x: auto;
}

.example-code pre {
  margin: 0;
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.6;
}

.example-code code {
  color: var(--text-secondary);
}

.tab-content {
  padding: var(--space-4);
}

.tab-content h3 {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-3);
}

.tab-content p {
  font-size: var(--font-size-md);
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: var(--space-2);
}

.setting-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-4);
  padding: var(--space-3);
  background: var(--bg-1);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
}

.setting-item label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.setting-select {
  padding: var(--space-2) var(--space-3);
  background: var(--bg-0);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  cursor: pointer;
}

.setting-select:focus {
  outline: none;
  border-color: var(--accent-indigo);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.setting-item input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.heavy-item {
  padding: var(--space-2);
  margin-bottom: var(--space-1);
  background: var(--bg-1);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}
</style>
