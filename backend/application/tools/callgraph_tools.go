package tools

import (
	"fmt"
	"syntaxia/domain"
	"strings"
)

// CallGraphToolsHandler handles call graph related tools
type CallGraphToolsHandler struct {
	BaseHandler
	CallGraph domain.CallGraphBuilder
}

// NewCallGraphToolsHandler creates a new call graph tools handler
func NewCallGraphToolsHandler(logger domain.Logger, callGraph domain.CallGraphBuilder) *CallGraphToolsHandler {
	return &CallGraphToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		CallGraph:   callGraph,
	}
}

var callGraphToolNames = map[string]bool{
	"get_callers":    true,
	"get_callees":    true,
	"get_impact":     true,
	"get_call_chain": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *CallGraphToolsHandler) CanHandle(toolName string) bool {
	return callGraphToolNames[toolName]
}

// GetTools returns the list of call graph tools
func (h *CallGraphToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "get_callers",
			Description: "Get functions that call the specified function",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"function_id": {Type: "string", Description: "Function identifier (e.g., 'pkg.FunctionName')"},
				},
				Required: []string{"function_id"},
			},
		},
		{
			Name:        "get_callees",
			Description: "Get functions called by the specified function",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"function_id": {Type: "string", Description: "Function identifier (e.g., 'pkg.FunctionName')"},
				},
				Required: []string{"function_id"},
			},
		},
		{
			Name:        "get_impact",
			Description: "Get all functions affected if the specified function changes",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"function_id": {Type: "string", Description: "Function identifier"},
					"max_depth":   {Type: "integer", Description: "Maximum depth to search (default: 3)"},
				},
				Required: []string{"function_id"},
			},
		},
		{
			Name:        "get_call_chain",
			Description: "Find call chain between two functions",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"start_function": {Type: "string", Description: "Starting function identifier"},
					"end_function":   {Type: "string", Description: "Target function identifier"},
					"max_depth":      {Type: "integer", Description: "Maximum depth (default: 5)"},
				},
				Required: []string{"start_function", "end_function"},
			},
		},
	}
}

// Execute executes a call graph tool
func (h *CallGraphToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	if h.CallGraph == nil {
		return "", fmt.Errorf("call graph not initialized")
	}

	switch toolName {
	case "get_callers":
		return h.getCallers(args, projectRoot)
	case "get_callees":
		return h.getCallees(args, projectRoot)
	case "get_impact":
		return h.getImpact(args, projectRoot)
	case "get_call_chain":
		return h.getCallChain(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown call graph tool: %s", toolName)
	}
}

func (h *CallGraphToolsHandler) getCallers(args map[string]any, projectRoot string) (string, error) {
	functionID, _ := args["function_id"].(string)
	if functionID == "" {
		return "", fmt.Errorf("function_id is required")
	}

	if _, err := h.CallGraph.Build(projectRoot); err != nil {
		return "", fmt.Errorf("failed to build call graph: %w", err)
	}

	callers := h.CallGraph.GetCallers(functionID)
	if len(callers) == 0 {
		return fmt.Sprintf("No callers found for %s", functionID), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Functions that call %s:\n\n", functionID))
	for _, c := range callers {
		result.WriteString(fmt.Sprintf("  - %s (%s:%d)\n", c.Name, c.FilePath, c.Line))
	}
	return result.String(), nil
}

func (h *CallGraphToolsHandler) getCallees(args map[string]any, projectRoot string) (string, error) {
	functionID, _ := args["function_id"].(string)
	if functionID == "" {
		return "", fmt.Errorf("function_id is required")
	}

	if _, err := h.CallGraph.Build(projectRoot); err != nil {
		return "", fmt.Errorf("failed to build call graph: %w", err)
	}

	callees := h.CallGraph.GetCallees(functionID)
	if len(callees) == 0 {
		return fmt.Sprintf("No callees found for %s", functionID), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Functions called by %s:\n\n", functionID))
	for _, c := range callees {
		result.WriteString(fmt.Sprintf("  - %s (%s:%d)\n", c.Name, c.FilePath, c.Line))
	}
	return result.String(), nil
}

func (h *CallGraphToolsHandler) getImpact(args map[string]any, projectRoot string) (string, error) {
	functionID, _ := args["function_id"].(string)
	if functionID == "" {
		return "", fmt.Errorf("function_id is required")
	}
	maxDepth := 3
	if d, ok := args["max_depth"].(float64); ok {
		maxDepth = int(d)
	}

	if _, err := h.CallGraph.Build(projectRoot); err != nil {
		return "", fmt.Errorf("failed to build call graph: %w", err)
	}

	affected := h.CallGraph.GetImpact(functionID, maxDepth)
	if len(affected) == 0 {
		return fmt.Sprintf("No impact found for %s", functionID), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Functions affected by changes to %s (depth %d):\n\n", functionID, maxDepth))
	for _, a := range affected {
		result.WriteString(fmt.Sprintf("  - %s (%s:%d)\n", a.Name, a.FilePath, a.Line))
	}
	return result.String(), nil
}

func (h *CallGraphToolsHandler) getCallChain(args map[string]any, projectRoot string) (string, error) {
	startID, _ := args["start_function"].(string)
	endID, _ := args["end_function"].(string)
	if startID == "" || endID == "" {
		return "", fmt.Errorf("start_function and end_function are required")
	}
	maxDepth := 5
	if d, ok := args["max_depth"].(float64); ok {
		maxDepth = int(d)
	}

	if _, err := h.CallGraph.Build(projectRoot); err != nil {
		return "", fmt.Errorf("failed to build call graph: %w", err)
	}

	chains := h.CallGraph.GetCallChain(startID, endID, maxDepth)
	if len(chains) == 0 {
		return fmt.Sprintf("No call chain found from %s to %s", startID, endID), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Call chains from %s to %s:\n\n", startID, endID))
	for i, chain := range chains {
		result.WriteString(fmt.Sprintf("  Chain %d: %s\n", i+1, strings.Join(chain, " → ")))
	}
	return result.String(), nil
}
