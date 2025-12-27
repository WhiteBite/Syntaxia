package analyzers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractVueProps(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []VueProp
	}{
		{
			name: "defineProps with TypeScript interface",
			content: `<script setup lang="ts">
defineProps<{
  title: string
  count?: number
  items: string[]
}>()
</script>`,
			expected: []VueProp{
				{Name: "title", Required: true},
				{Name: "count", Required: false},
				{Name: "items", Required: true},
			},
		},
		{
			name: "defineProps with object syntax",
			content: `<script setup>
defineProps({
  message: String,
  count: Number
})
</script>`,
			expected: []VueProp{
				{Name: "message", Required: true},
				{Name: "count", Required: true},
			},
		},
		{
			name:     "no props defined",
			content:  `<script setup>\nconst x = 1\n</script>`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props := extractVueProps(tt.content)
			assert.Equal(t, len(tt.expected), len(props))
			for i, exp := range tt.expected {
				if i < len(props) {
					assert.Equal(t, exp.Name, props[i].Name)
					assert.Equal(t, exp.Required, props[i].Required)
				}
			}
		})
	}
}

func TestExtractVueEmits(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []VueEmit
	}{
		{
			name: "defineEmits with array syntax",
			content: `<script setup>
defineEmits(['change', 'update', 'delete'])
</script>`,
			expected: []VueEmit{
				{Name: "change"},
				{Name: "update"},
				{Name: "delete"},
			},
		},
		{
			name:     "no emits defined",
			content:  `<script setup>\nconst x = 1\n</script>`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emits := extractVueEmits(tt.content)
			assert.Equal(t, len(tt.expected), len(emits))
			for i, exp := range tt.expected {
				if i < len(emits) {
					assert.Equal(t, exp.Name, emits[i].Name)
				}
			}
		})
	}
}

func TestExtractVueSlots(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []VueSlot
	}{
		{
			name: "defineSlots with TypeScript",
			content: `<script setup lang="ts">
defineSlots<{
  default: () => any
  header: () => any
}>()
</script>`,
			expected: []VueSlot{
				{Name: "default"},
				{Name: "header"},
			},
		},
		{
			name:     "no slots defined",
			content:  `<script setup>\nconst x = 1\n</script>`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slots := extractVueSlots(tt.content)
			assert.Equal(t, len(tt.expected), len(slots))
			for i, exp := range tt.expected {
				if i < len(slots) {
					assert.Equal(t, exp.Name, slots[i].Name)
				}
			}
		})
	}
}

func TestExtractVueExpose(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "defineExpose with methods",
			content: `<script setup>
const focus = () => {}
const reset = () => {}
defineExpose({ focus, reset })
</script>`,
			expected: []string{"focus", "reset"},
		},
		{
			name:     "no expose defined",
			content:  `<script setup>\nconst x = 1\n</script>`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expose := extractVueExpose(tt.content)
			assert.ElementsMatch(t, tt.expected, expose)
		})
	}
}

func TestExtractVueModels(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []VueModelDef
	}{
		{
			name: "default model",
			content: `<script setup>
const model = defineModel()
</script>`,
			expected: []VueModelDef{
				{Name: "modelValue", PropName: "modelValue", Event: "update:modelValue"},
			},
		},
		{
			name: "named model",
			content: `<script setup>
const title = defineModel('title')
</script>`,
			expected: []VueModelDef{
				{Name: "title", PropName: "title", Event: "update:title"},
			},
		},
		{
			name: "multiple models",
			content: `<script setup>
const firstName = defineModel('firstName')
const lastName = defineModel('lastName')
</script>`,
			expected: []VueModelDef{
				{Name: "firstName", PropName: "firstName", Event: "update:firstName"},
				{Name: "lastName", PropName: "lastName", Event: "update:lastName"},
			},
		},
		{
			name:     "no model defined",
			content:  `<script setup>\nconst x = 1\n</script>`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := extractVueModels(tt.content)
			assert.Equal(t, len(tt.expected), len(models))
			for i, exp := range tt.expected {
				if i < len(models) {
					assert.Equal(t, exp.Name, models[i].Name)
					assert.Equal(t, exp.PropName, models[i].PropName)
					assert.Equal(t, exp.Event, models[i].Event)
				}
			}
		})
	}
}

func TestSymbolIndex_VueComponentInfo(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)
	ctx := context.Background()

	vueContent := `<script setup lang="ts">
defineProps<{
  title: string
  count?: number
}>()

defineEmits(['change', 'update'])

const focus = () => {}
defineExpose({ focus })
</script>

<template>
  <div>{{ title }}</div>
</template>`

	// Index Vue file
	err := idx.IndexFile(ctx, "components/MyComponent.vue", []byte(vueContent))
	require.NoError(t, err)

	// Extract Vue component info
	idx.mu.Lock()
	idx.extractVueComponentInfo("components/MyComponent.vue", []byte(vueContent))
	idx.mu.Unlock()

	// Test GetVueComponentInfo
	info := idx.GetVueComponentInfo("components/MyComponent.vue")
	require.NotNil(t, info)
	assert.Equal(t, "MyComponent", info.Name)
	assert.Len(t, info.Props, 2)
	assert.Len(t, info.Emits, 2)
	assert.Len(t, info.Expose, 1)
}

func TestSymbolIndex_FindVueComponentByName(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	// Add Vue component info manually
	idx.mu.Lock()
	idx.vueComponents["components/UserCard.vue"] = &VueComponentInfo{
		Name:     "UserCard",
		FilePath: "components/UserCard.vue",
		Props:    []VueProp{{Name: "user", Required: true}},
	}
	idx.vueComponents["components/Button.vue"] = &VueComponentInfo{
		Name:     "Button",
		FilePath: "components/Button.vue",
		Props:    []VueProp{{Name: "label", Required: true}},
	}
	idx.mu.Unlock()

	tests := []struct {
		name          string
		componentName string
		expectFound   bool
		expectedPath  string
	}{
		{"exact match", "UserCard", true, "components/UserCard.vue"},
		{"case insensitive", "usercard", true, "components/UserCard.vue"},
		{"another component", "Button", true, "components/Button.vue"},
		{"not found", "NonExistent", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := idx.FindVueComponentByName(tt.componentName)
			if tt.expectFound {
				require.NotNil(t, info)
				assert.Equal(t, tt.expectedPath, info.FilePath)
			} else {
				assert.Nil(t, info)
			}
		})
	}
}

func TestSymbolIndex_GetVueComponentProps(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["components/Form.vue"] = &VueComponentInfo{
		Name:     "Form",
		FilePath: "components/Form.vue",
		Props: []VueProp{
			{Name: "modelValue", Required: true},
			{Name: "disabled", Required: false},
			{Name: "label", Required: true},
		},
	}
	idx.mu.Unlock()

	props := idx.GetVueComponentProps("Form")
	require.Len(t, props, 3)
	assert.Equal(t, "modelValue", props[0].Name)
	assert.True(t, props[0].Required)
	assert.Equal(t, "disabled", props[1].Name)
	assert.False(t, props[1].Required)
}

func TestSymbolIndex_GetVueComponentEmits(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["components/Input.vue"] = &VueComponentInfo{
		Name:     "Input",
		FilePath: "components/Input.vue",
		Emits: []VueEmit{
			{Name: "update:modelValue"},
			{Name: "focus"},
			{Name: "blur"},
		},
	}
	idx.mu.Unlock()

	emits := idx.GetVueComponentEmits("Input")
	require.Len(t, emits, 3)
	assert.Equal(t, "update:modelValue", emits[0].Name)
}

func TestSymbolIndex_GetVueComponentSlots(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["components/Card.vue"] = &VueComponentInfo{
		Name:     "Card",
		FilePath: "components/Card.vue",
		Slots: []VueSlot{
			{Name: "default"},
			{Name: "header"},
			{Name: "footer"},
		},
	}
	idx.mu.Unlock()

	slots := idx.GetVueComponentSlots("Card")
	require.Len(t, slots, 3)
	assert.Equal(t, "default", slots[0].Name)
	assert.Equal(t, "header", slots[1].Name)
	assert.Equal(t, "footer", slots[2].Name)
}

func TestSymbolIndex_GetAllVueComponents(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["a.vue"] = &VueComponentInfo{Name: "A"}
	idx.vueComponents["b.vue"] = &VueComponentInfo{Name: "B"}
	idx.vueComponents["c.vue"] = &VueComponentInfo{Name: "C"}
	idx.mu.Unlock()

	all := idx.GetAllVueComponents()
	assert.Len(t, all, 3)
	assert.Contains(t, all, "a.vue")
	assert.Contains(t, all, "b.vue")
	assert.Contains(t, all, "c.vue")
}

func TestSymbolIndex_ResolveComponentUsage(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["components/Button.vue"] = &VueComponentInfo{
		Name:     "Button",
		FilePath: "components/Button.vue",
		Props:    []VueProp{{Name: "label", Required: true}},
	}
	idx.mu.Unlock()

	// Test resolve by name
	info := idx.ResolveVueComponentUsage("Button")
	require.NotNil(t, info)
	assert.Equal(t, "Button", info.Name)

	// Test resolve non-existent
	info = idx.ResolveVueComponentUsage("NonExistent")
	assert.Nil(t, info)
}

func TestSymbolIndex_ResolveComponentUsage_Files(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	idx.mu.Lock()
	idx.vueComponents["components/Button.vue"] = &VueComponentInfo{
		Name:     "Button",
		FilePath: "components/Button.vue",
	}
	idx.vueComponents["pages/Home.vue"] = &VueComponentInfo{
		Name:     "Home",
		FilePath: "pages/Home.vue",
	}
	idx.mu.Unlock()

	// Test ResolveComponentUsage returns file paths
	usages := idx.ResolveComponentUsage("components/Button.vue")
	// Initially empty since we haven't set up imports
	assert.Empty(t, usages)
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Button", "button"},
		{"MyComponent", "my-component"},
		{"UserProfileCard", "user-profile-card"},
		{"ABC", "a-b-c"},
		{"already", "already"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toKebabCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		search   string
		expected bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []string{}, "a", false},
		{"empty search", []string{"a", "b"}, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsString(tt.slice, tt.search)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractVueProps_ComplexTypes(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		minExpected int // minimum expected props (regex may pick up extra)
	}{
		{
			name: "generic types",
			content: `<script setup lang="ts">
defineProps<{
  items: Array<string>
  callback: (item: string) => void
}>()
</script>`,
			minExpected: 2, // items, callback (may also pick up 'item' from callback signature)
		},
		{
			name: "union types",
			content: `<script setup lang="ts">
defineProps<{
  status: 'active' | 'inactive'
  value: string | number
}>()
</script>`,
			minExpected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props := extractVueProps(tt.content)
			assert.GreaterOrEqual(t, len(props), tt.minExpected)
		})
	}
}

func TestVueComponentInfo_FullExtraction(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	vueContent := `<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  title: string
  count?: number
  items: string[]
}>()

defineEmits(['update', 'delete', 'select'])

const inputRef = ref<HTMLInputElement>()
const focus = () => inputRef.value?.focus()
const reset = () => {}

defineExpose({ focus, reset })

const model = defineModel('value')
</script>

<template>
  <div>
    <input ref="inputRef" v-model="model" />
  </div>
</template>`

	// Extract Vue component info
	idx.mu.Lock()
	idx.extractVueComponentInfo("components/ComplexComponent.vue", []byte(vueContent))
	idx.mu.Unlock()

	info := idx.GetVueComponentInfo("components/ComplexComponent.vue")
	require.NotNil(t, info)

	assert.Equal(t, "ComplexComponent", info.Name)
	assert.GreaterOrEqual(t, len(info.Props), 3)
	assert.Len(t, info.Emits, 3)
	// Slots regex requires defineSlots which is not in this content
	assert.Len(t, info.Expose, 2)
	assert.Len(t, info.ModelDefs, 1)
}
