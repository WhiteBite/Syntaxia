package analyzers

import (
	"regexp"
	"strings"

	"syntaxia/domain/analysis"
)

// VueComponentInfo holds Vue component cross-file symbol information
type VueComponentInfo struct {
	Name      string
	FilePath  string
	Props     []VueProp
	Emits     []VueEmit
	Slots     []VueSlot
	Expose    []string
	ModelDefs []VueModelDef
}

// VueProp represents a Vue component prop
type VueProp struct {
	Name     string
	Type     string
	Required bool
	Default  string
}

// VueEmit represents a Vue component emit
type VueEmit struct {
	Name    string
	Payload string
}

// VueSlot represents a Vue component slot
type VueSlot struct {
	Name  string
	Props []string
}

// VueModelDef represents a Vue v-model definition
type VueModelDef struct {
	Name     string
	PropName string
	Event    string
}

// Precompiled regexes for Vue component parsing
var (
	vueIdxDefinePropsRe  = regexp.MustCompile(`defineProps<\{([^}]+)\}>|defineProps\(\{([^}]+)\}\)|defineProps\(\[([^\]]+)\]\)`)
	vueIdxDefineEmitsRe  = regexp.MustCompile(`defineEmits\(\[([^\]]+)\]\)`)
	vueIdxDefineSlotsRe  = regexp.MustCompile(`defineSlots<\{([^}]+)\}>`)
	vueIdxDefineExposeRe = regexp.MustCompile(`defineExpose\(\{\s*([^}]+)\s*\}\)`)
	vueIdxDefineModelRe  = regexp.MustCompile(`defineModel\(\s*['"](\w+)['"]\s*\)|defineModel\(\)`)
	vueIdxPropNameRe     = regexp.MustCompile(`(\w+)\s*[?:]`)
	vueIdxEmitNameRe     = regexp.MustCompile(`['"]([^'"]+)['"]`)
	vueIdxExposeNameRe   = regexp.MustCompile(`\b(\w+)\b`)
)

// extractVueComponentInfo extracts Vue component props, emits, slots from content.
// Must be called with idx.mu locked.
func (idx *SymbolIndexImpl) extractVueComponentInfo(filePath string, content []byte) {
	text := string(content)

	// Extract component name from file path
	parts := strings.Split(filePath, "/")
	if len(parts) == 0 {
		parts = strings.Split(filePath, "\\")
	}
	fileName := parts[len(parts)-1]
	componentName := strings.TrimSuffix(fileName, ".vue")

	info := &VueComponentInfo{
		Name:     componentName,
		FilePath: filePath,
		Props:    extractVueProps(text),
		Emits:    extractVueEmits(text),
		Slots:    extractVueSlots(text),
		Expose:   extractVueExpose(text),
		ModelDefs: extractVueModels(text),
	}

	idx.vueComponents[filePath] = info
}

// extractVueProps extracts props from Vue component.
func extractVueProps(text string) []VueProp {
	var props []VueProp

	matches := vueIdxDefinePropsRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		propsContent := getFirstNonEmpty(match[1:])
		if propsContent == "" {
			continue
		}

		propMatches := vueIdxPropNameRe.FindAllStringSubmatch(propsContent, -1)
		for _, pm := range propMatches {
			if pm[1] != "" {
				props = append(props, VueProp{
					Name:     pm[1],
					Required: !strings.Contains(propsContent, pm[1]+"?"),
				})
			}
		}
	}

	return props
}

// extractVueEmits extracts emits from Vue component.
func extractVueEmits(text string) []VueEmit {
	var emits []VueEmit

	matches := vueIdxDefineEmitsRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if match[1] == "" {
			continue
		}

		emitMatches := vueIdxEmitNameRe.FindAllStringSubmatch(match[1], -1)
		for _, em := range emitMatches {
			if em[1] != "" {
				emits = append(emits, VueEmit{Name: em[1]})
			}
		}
	}

	return emits
}

// extractVueSlots extracts slots from Vue component.
func extractVueSlots(text string) []VueSlot {
	var slots []VueSlot

	matches := vueIdxDefineSlotsRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if match[1] == "" {
			continue
		}

		slotMatches := vueIdxPropNameRe.FindAllStringSubmatch(match[1], -1)
		for _, sm := range slotMatches {
			if sm[1] != "" {
				slots = append(slots, VueSlot{Name: sm[1]})
			}
		}
	}

	return slots
}

// extractVueExpose extracts exposed methods/properties from Vue component.
func extractVueExpose(text string) []string {
	var expose []string

	matches := vueIdxDefineExposeRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if match[1] == "" {
			continue
		}

		exposeMatches := vueIdxExposeNameRe.FindAllStringSubmatch(match[1], -1)
		for _, em := range exposeMatches {
			if em[1] != "" {
				expose = append(expose, em[1])
			}
		}
	}

	return expose
}

// extractVueModels extracts defineModel definitions.
func extractVueModels(text string) []VueModelDef {
	var models []VueModelDef

	matches := vueIdxDefineModelRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		name := "modelValue"
		if len(match) > 1 && match[1] != "" {
			name = match[1]
		}
		models = append(models, VueModelDef{
			Name:     name,
			PropName: name,
			Event:    "update:" + name,
		})
	}

	return models
}

// getFirstNonEmpty returns the first non-empty string from slice.
func getFirstNonEmpty(strs []string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}

// GetVueComponentInfo returns Vue component info for cross-file resolution.
func (idx *SymbolIndexImpl) GetVueComponentInfo(filePath string) *VueComponentInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.vueComponents[filePath]
}

// GetAllVueComponents returns all indexed Vue components.
func (idx *SymbolIndexImpl) GetAllVueComponents() map[string]*VueComponentInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	result := make(map[string]*VueComponentInfo, len(idx.vueComponents))
	for k, v := range idx.vueComponents {
		result[k] = v
	}
	return result
}

// FindVueComponentByName finds a Vue component by its name.
func (idx *SymbolIndexImpl) FindVueComponentByName(name string) *VueComponentInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	nameLower := strings.ToLower(name)
	for _, info := range idx.vueComponents {
		if strings.ToLower(info.Name) == nameLower {
			return info
		}
	}
	return nil
}

// ResolveVueComponentUsage resolves a Vue component usage to its definition.
func (idx *SymbolIndexImpl) ResolveVueComponentUsage(componentName string) *VueComponentInfo {
	return idx.FindVueComponentByName(componentName)
}

// GetVueComponentProps returns props for a Vue component.
func (idx *SymbolIndexImpl) GetVueComponentProps(componentName string) []VueProp {
	info := idx.FindVueComponentByName(componentName)
	if info == nil {
		return nil
	}
	return info.Props
}

// GetVueComponentEmits returns emits for a Vue component.
func (idx *SymbolIndexImpl) GetVueComponentEmits(componentName string) []VueEmit {
	info := idx.FindVueComponentByName(componentName)
	if info == nil {
		return nil
	}
	return info.Emits
}

// GetVueComponentSlots returns slots for a Vue component.
func (idx *SymbolIndexImpl) GetVueComponentSlots(componentName string) []VueSlot {
	info := idx.FindVueComponentByName(componentName)
	if info == nil {
		return nil
	}
	return info.Slots
}

// ResolveComponentUsage finds all files that use a given Vue component.
// Returns list of file paths where the component is imported or used in template.
func (idx *SymbolIndexImpl) ResolveComponentUsage(componentPath string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	targetInfo := idx.vueComponents[componentPath]
	if targetInfo == nil {
		return nil
	}

	var usages []string
	componentName := targetInfo.Name
	componentNameLower := strings.ToLower(componentName)
	kebabName := toKebabCase(componentName)

	for filePath, symbols := range idx.byFile {
		if filePath == componentPath {
			continue
		}

		// Check imports in the file
		for _, symIdx := range symbols {
			sym := idx.symbols[symIdx]
			if sym.Kind == analysis.KindImport {
				// Check if import references the component
				if strings.Contains(strings.ToLower(sym.Name), componentNameLower) {
					usages = append(usages, filePath)
					break
				}
			}
		}
	}

	// Also check Vue component template usages
	for filePath, info := range idx.vueComponents {
		if filePath == componentPath {
			continue
		}
		// Check if this component's template uses the target component
		if containsComponentUsage(info, componentName, kebabName) {
			if !containsString(usages, filePath) {
				usages = append(usages, filePath)
			}
		}
	}

	return usages
}

// containsComponentUsage checks if a Vue component uses another component.
func containsComponentUsage(info *VueComponentInfo, pascalName, kebabName string) bool {
	// This would require storing template component usages in VueComponentInfo
	// For now, return false - full implementation would parse template
	_ = info
	_ = pascalName
	_ = kebabName
	return false
}

// containsString checks if slice contains a string.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// toKebabCase converts PascalCase to kebab-case.
func toKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
