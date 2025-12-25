package tools

import (
	"encoding/json"
	"fmt"
	"syntaxia/domain"
)

// PreferencesStore interface for user preferences storage
type PreferencesStore interface {
	SetPreference(key, value string) error
	GetPreference(key string) (string, error)
	GetAllPreferences() (map[string]string, error)
}

// PreferencesToolsHandler handles user preferences tools
type PreferencesToolsHandler struct {
	BaseHandler
	Store domain.ContextMemory // Uses ContextMemory which has preference methods
}

// NewPreferencesToolsHandler creates a new preferences tools handler
func NewPreferencesToolsHandler(logger domain.Logger, store domain.ContextMemory) *PreferencesToolsHandler {
	return &PreferencesToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		Store:       store,
	}
}

var preferencesToolNames = map[string]bool{
	"set_preference":  true,
	"get_preferences": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *PreferencesToolsHandler) CanHandle(toolName string) bool {
	return preferencesToolNames[toolName]
}

// GetTools returns the list of preferences tools
func (h *PreferencesToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "set_preference",
			Description: "Set a user preference value",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"key":   {Type: "string", Description: "Preference key"},
					"value": {Type: "string", Description: "Preference value"},
				},
				Required: []string{"key", "value"},
			},
		},
		{
			Name:        "get_preferences",
			Description: "Get user preferences. If key is provided, returns single preference; otherwise returns all.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"key": {Type: "string", Description: "Preference key (optional, omit to get all)"},
				},
			},
		},
	}
}

// Execute executes a preferences tool
func (h *PreferencesToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	switch toolName {
	case "set_preference":
		return h.setPreference(args)
	case "get_preferences":
		return h.getPreferences(args)
	default:
		return "", fmt.Errorf("unknown preferences tool: %s", toolName)
	}
}

func (h *PreferencesToolsHandler) setPreference(args map[string]any) (string, error) {
	if h.Store == nil {
		return "", fmt.Errorf("preferences store not initialized")
	}

	key, _ := args["key"].(string)
	value, _ := args["value"].(string)

	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	if err := h.Store.SetPreference(key, value); err != nil {
		return "", fmt.Errorf("failed to set preference: %w", err)
	}

	return fmt.Sprintf("Preference set: %s = %s", key, value), nil
}

func (h *PreferencesToolsHandler) getPreferences(args map[string]any) (string, error) {
	if h.Store == nil {
		return "", fmt.Errorf("preferences store not initialized")
	}

	key, _ := args["key"].(string)

	if key != "" {
		value, err := h.Store.GetPreference(key)
		if err != nil {
			return "", fmt.Errorf("failed to get preference: %w", err)
		}
		if value == "" {
			return fmt.Sprintf("Preference '%s' not found", key), nil
		}
		return fmt.Sprintf("%s = %s", key, value), nil
	}

	// Get all preferences
	prefs, err := h.Store.GetAllPreferences()
	if err != nil {
		return "", fmt.Errorf("failed to get preferences: %w", err)
	}

	if len(prefs) == 0 {
		return "No preferences set", nil
	}

	jsonBytes, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format preferences: %w", err)
	}

	return fmt.Sprintf("User preferences:\n%s", string(jsonBytes)), nil
}
