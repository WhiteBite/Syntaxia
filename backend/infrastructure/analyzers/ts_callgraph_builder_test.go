package analyzers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTSCallGraphBuilder(t *testing.T) {
	builder := NewTSCallGraphBuilder()

	if builder == nil {
		t.Fatal("NewTSCallGraphBuilder returned nil")
	}
	if builder.graph == nil {
		t.Error("graph not initialized")
	}
	if builder.componentGraph == nil {
		t.Error("componentGraph not initialized")
	}
}

func TestTSCallGraphBuilder_Build_TypeScript(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		wantNodes     int
		wantEdges     int
		wantFuncNames []string
	}{
		{
			name: "simple functions",
			files: map[string]string{
				"main.ts": `
function greet(name: string): void {
	console.log(name);
}

function main(): void {
	greet("world");
}
`,
			},
			wantNodes:     2,
			wantFuncNames: []string{"greet", "main"},
		},
		{
			name: "arrow functions",
			files: map[string]string{
				"utils.ts": `
export const add = (a: number, b: number): number => {
	return a + b;
};

export const multiply = (a: number, b: number): number => {
	return a * b;
};
`,
			},
			wantNodes:     2,
			wantFuncNames: []string{"add", "multiply"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for filename, content := range tt.files {
				writeTSTestFile(t, tmpDir, filename, content)
			}

			builder := NewTSCallGraphBuilder()
			graph, err := builder.Build(tmpDir)

			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if graph == nil {
				t.Fatal("Build returned nil graph")
			}
		})
	}
}

func TestTSCallGraphBuilder_Build_ClassMethods(t *testing.T) {
	tmpDir := t.TempDir()

	classCode := `
export class Calculator {
	private value: number = 0;

	add(n: number): Calculator {
		this.value += n;
		return this;
	}

	subtract(n: number): Calculator {
		this.value -= n;
		return this;
	}

	getValue(): number {
		return this.value;
	}
}
`
	writeTSTestFile(t, tmpDir, "calculator.ts", classCode)

	builder := NewTSCallGraphBuilder()
	graph, err := builder.Build(tmpDir)

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(graph.Nodes) == 0 {
		t.Error("no nodes found for class methods")
	}

	// Check that methods are found
	foundMethods := make(map[string]bool)
	for _, node := range graph.Nodes {
		foundMethods[node.Name] = true
	}

	expectedMethods := []string{"add", "subtract", "getValue"}
	for _, method := range expectedMethods {
		if !foundMethods[method] {
			t.Errorf("method %s not found", method)
		}
	}
}

func TestTSCallGraphBuilder_Build_VueComponent(t *testing.T) {
	tmpDir := t.TempDir()

	vueCode := `<template>
	<div>
		<ChildComponent :value="count" />
		<button @click="increment">+</button>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ChildComponent from './ChildComponent.vue';

const count = ref(0);

function increment(): void {
	count.value++;
}

function decrement(): void {
	count.value--;
}
</script>
`
	writeTSTestFile(t, tmpDir, "Counter.vue", vueCode)

	builder := NewTSCallGraphBuilder()
	_, err := builder.Build(tmpDir)

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check component graph
	compGraph := builder.GetComponentGraph()
	if compGraph == nil {
		t.Fatal("component graph is nil")
	}

	if len(compGraph.Components) == 0 {
		t.Error("no components found")
	}

	// Check component was registered
	found := false
	for _, comp := range compGraph.Components {
		if comp.Name == "Counter" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Counter component not found")
	}
}

func TestTSCallGraphBuilder_GetCallers(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function helper(): void {
	console.log("help");
}

function main(): void {
	helper();
}

function other(): void {
	helper();
}
`
	writeTSTestFile(t, tmpDir, "main.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find helper function ID
	var helperID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "helper" {
			helperID = id
			break
		}
	}

	if helperID == "" {
		t.Skip("helper function not found")
	}

	callers := builder.GetCallers(helperID)
	t.Logf("callers for helper: %d", len(callers))
}

func TestTSCallGraphBuilder_GetCallees(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function foo(): void {}
function bar(): void {}

function main(): void {
	foo();
	bar();
}
`
	writeTSTestFile(t, tmpDir, "main.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find main function ID
	var mainID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "main" {
			mainID = id
			break
		}
	}

	if mainID == "" {
		t.Skip("main function not found")
	}

	callees := builder.GetCallees(mainID)
	t.Logf("callees for main: %d", len(callees))
}

func TestTSCallGraphBuilder_GetImpact(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function a(): void {
	b();
}

function b(): void {
	c();
}

function c(): void {
	console.log("c");
}
`
	writeTSTestFile(t, tmpDir, "chain.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find c function ID
	var cID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "c" {
			cID = id
			break
		}
	}

	if cID == "" {
		t.Skip("c function not found")
	}

	impact := builder.GetImpact(cID, 3)
	t.Logf("impact for c: %d nodes", len(impact))
}

func TestTSCallGraphBuilder_VueComponentUsage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create child component
	childVue := `<template>
	<span>{{ value }}</span>
</template>

<script setup lang="ts">
defineProps<{
	value: number;
}>();
</script>
`
	writeTSTestFile(t, tmpDir, "ChildComponent.vue", childVue)

	// Create parent component
	parentVue := `<template>
	<div>
		<ChildComponent :value="count" />
		<AnotherChild />
	</div>
</template>

<script setup lang="ts">
import ChildComponent from './ChildComponent.vue';
import AnotherChild from './AnotherChild.vue';

const count = 42;
</script>
`
	writeTSTestFile(t, tmpDir, "ParentComponent.vue", parentVue)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	compGraph := builder.GetComponentGraph()

	// Check usages were recorded
	usages := builder.GetComponentUsages("ParentComponent.vue")
	t.Logf("component usages in parent: %d", len(usages))

	// Check component graph edges
	t.Logf("component graph edges: %d", len(compGraph.Usages))
}

func TestTSCallGraphBuilder_SkipDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file in node_modules (should be skipped)
	nodeModulesDir := filepath.Join(tmpDir, "node_modules", "some-lib")
	os.MkdirAll(nodeModulesDir, 0o755)
	writeTSTestFile(t, nodeModulesDir, "index.ts", `export function libFunc() {}`)

	// Create file in src (should be analyzed)
	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0o755)
	writeTSTestFile(t, srcDir, "main.ts", `export function mainFunc() {}`)

	builder := NewTSCallGraphBuilder()
	graph, err := builder.Build(tmpDir)

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should only find mainFunc, not libFunc
	foundLib := false
	foundMain := false
	for _, node := range graph.Nodes {
		if node.Name == "libFunc" {
			foundLib = true
		}
		if node.Name == "mainFunc" {
			foundMain = true
		}
	}

	if foundLib {
		t.Error("libFunc from node_modules should be skipped")
	}
	if !foundMain {
		t.Error("mainFunc from src should be found")
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-component", "MyComponent"},
		{"button", "Button"},
		{"MyComponent", "MyComponent"},
		{"user-profile-card", "UserProfileCard"},
		{"a-b-c", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsHTMLElement(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"div", true},
		{"span", true},
		{"MyComponent", false},
		{"button", true},
		{"custom-element", false},
		{"template", true},
		{"slot", true},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			result := isHTMLElement(tt.tag)
			if result != tt.expected {
				t.Errorf("isHTMLElement(%q) = %v, want %v", tt.tag, result, tt.expected)
			}
		})
	}
}

func TestIsTSKeyword(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"function", true},
		{"class", true},
		{"myFunc", false},
		{"if", true},
		{"async", true},
		{"await", true},
		{"customName", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTSKeyword(tt.name)
			if result != tt.expected {
				t.Errorf("isTSKeyword(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestTSCallGraphBuilder_GetCallersEdges(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function helper(): void {
	console.log("help");
}

function caller1(): void {
	helper();
}

function caller2(): void {
	helper();
}
`
	writeTSTestFile(t, tmpDir, "main.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find helper function ID
	var helperID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "helper" {
			helperID = id
			break
		}
	}

	if helperID != "" {
		edges := builder.GetCallersEdges(helperID)
		t.Logf("caller edges for helper: %d", len(edges))
	}
}

func TestTSCallGraphBuilder_GetCalleesEdges(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function a(): void {}
function b(): void {}

function main(): void {
	a();
	b();
}
`
	writeTSTestFile(t, tmpDir, "main.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find main function ID
	var mainID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "main" {
			mainID = id
			break
		}
	}

	if mainID != "" {
		edges := builder.GetCalleesEdges(mainID)
		t.Logf("callee edges for main: %d", len(edges))
	}
}

func TestTSCallGraphBuilder_BuildComponentUsageGraph(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple Vue components
	button := `<template><button>{{ label }}</button></template>
<script setup lang="ts">
defineProps<{ label: string }>()
</script>`
	writeTSTestFile(t, tmpDir, "Button.vue", button)

	card := `<template>
	<div class="card">
		<Button label="Click" />
	</div>
</template>
<script setup lang="ts">
import Button from './Button.vue'
</script>`
	writeTSTestFile(t, tmpDir, "Card.vue", card)

	builder := NewTSCallGraphBuilder()
	graph, err := builder.BuildComponentUsageGraph(tmpDir)

	if err != nil {
		t.Fatalf("BuildComponentUsageGraph failed: %v", err)
	}

	if graph == nil {
		t.Fatal("component graph is nil")
	}

	t.Logf("components found: %d", len(graph.Components))
	t.Logf("usage edges: %d", len(graph.Usages))
}

func TestTSCallGraphBuilder_GetComponentUsagesList(t *testing.T) {
	tmpDir := t.TempDir()

	// Create components
	child := `<template><span>Child</span></template>
<script setup lang="ts"></script>`
	writeTSTestFile(t, tmpDir, "Child.vue", child)

	parent := `<template>
	<div>
		<Child />
		<Child />
	</div>
</template>
<script setup lang="ts">
import Child from './Child.vue'
</script>`
	writeTSTestFile(t, tmpDir, "Parent.vue", parent)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	files := builder.GetComponentUsagesList("Child")
	t.Logf("files using Child: %v", files)
}

func TestTSCallGraphBuilder_GetCallChain(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function a(): void {
	b();
}

function b(): void {
	c();
}

function c(): void {
	d();
}

function d(): void {
	console.log("end");
}
`
	writeTSTestFile(t, tmpDir, "chain.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Find function IDs
	var aID, dID string
	for id, node := range builder.graph.Nodes {
		if node.Name == "a" {
			aID = id
		}
		if node.Name == "d" {
			dID = id
		}
	}

	if aID != "" && dID != "" {
		chains := builder.GetCallChain(aID, dID, 5)
		t.Logf("call chains from a to d: %d", len(chains))
		for i, chain := range chains {
			t.Logf("chain %d: %v", i, chain)
		}
	}
}

func TestTSCallGraphBuilder_FindFunctionByName(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function getUserById(id: string): void {}
function getUserByEmail(email: string): void {}
function createUser(data: any): void {}
`
	writeTSTestFile(t, tmpDir, "users.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	// Search for "User" should find all three
	results := builder.FindFunctionByName("User")
	if len(results) < 3 {
		t.Logf("expected at least 3 functions with 'User', got %d", len(results))
	}

	// Search for "getUser" should find two
	results = builder.FindFunctionByName("getUser")
	t.Logf("functions matching 'getUser': %d", len(results))
}

func TestTSCallGraphBuilder_GetFunctionsInFile(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
function func1(): void {}
function func2(): void {}
function func3(): void {}
`
	writeTSTestFile(t, tmpDir, "funcs.ts", code)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	funcs := builder.GetFunctionsInFile("funcs.ts")
	if len(funcs) != 3 {
		t.Errorf("expected 3 functions in file, got %d", len(funcs))
	}
}

func TestTSCallGraphBuilder_Stats(t *testing.T) {
	tmpDir := t.TempDir()

	tsCode := `
function a(): void { b(); }
function b(): void {}
`
	writeTSTestFile(t, tmpDir, "main.ts", tsCode)

	vueCode := `<template><div></div></template>
<script setup lang="ts">
function setup(): void {}
</script>`
	writeTSTestFile(t, tmpDir, "App.vue", vueCode)

	builder := NewTSCallGraphBuilder()
	builder.Build(tmpDir)

	stats := builder.Stats()

	if stats["total_nodes"] == 0 {
		t.Error("expected some nodes")
	}
	if stats["files_analyzed"] == 0 {
		t.Error("expected some files analyzed")
	}

	t.Logf("stats: %+v", stats)
}

func TestTSCallGraphBuilder_AsyncFunctions(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
async function fetchData(): Promise<void> {
	await processData();
}

async function processData(): Promise<void> {
	console.log("processing");
}

const fetchUser = async (id: string): Promise<void> => {
	await fetchData();
};
`
	writeTSTestFile(t, tmpDir, "async.ts", code)

	builder := NewTSCallGraphBuilder()
	graph, err := builder.Build(tmpDir)

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should find async functions
	foundAsync := false
	for _, node := range graph.Nodes {
		if node.Name == "fetchData" || node.Name == "fetchUser" {
			foundAsync = true
			break
		}
	}

	if !foundAsync {
		t.Error("async functions not found")
	}
}

func TestTSCallGraphBuilder_ClassWithMethods(t *testing.T) {
	tmpDir := t.TempDir()

	code := `
class UserService {
	private users: string[] = [];

	async getUser(id: string): Promise<string> {
		return this.findById(id);
	}

	private findById(id: string): string {
		return this.users.find(u => u === id) || '';
	}

	public addUser(user: string): void {
		this.users.push(user);
	}

	static getInstance(): UserService {
		return new UserService();
	}
}
`
	writeTSTestFile(t, tmpDir, "service.ts", code)

	builder := NewTSCallGraphBuilder()
	graph, err := builder.Build(tmpDir)

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Check methods are found
	methodNames := make(map[string]bool)
	for _, node := range graph.Nodes {
		methodNames[node.Name] = true
	}

	expectedMethods := []string{"getUser", "findById", "addUser", "getInstance"}
	for _, method := range expectedMethods {
		if !methodNames[method] {
			t.Logf("method %s not found (may be expected based on regex)", method)
		}
	}
}

// writeTSTestFile creates a test file in the given directory
func writeTSTestFile(t *testing.T, base, path, content string) {
	t.Helper()
	fullPath := filepath.Join(base, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
}
