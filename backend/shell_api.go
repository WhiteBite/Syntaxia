package main

// ShellIntegrationStatus returns current shell integration status
type ShellIntegrationStatus struct {
	IsRegistered bool   `json:"isRegistered"`
	CurrentOS    string `json:"currentOS"`
}

// GetShellIntegrationStatus returns whether context menu is registered
func (a *App) GetShellIntegrationStatus() (ShellIntegrationStatus, error) {
	if a.container.ShellIntegration == nil {
		return ShellIntegrationStatus{CurrentOS: "unknown"}, nil
	}

	isRegistered, err := a.container.ShellIntegration.IsRegistered()
	if err != nil {
		a.log.Warning("Failed to check shell integration status: " + err.Error())
	}

	return ShellIntegrationStatus{
		IsRegistered: isRegistered,
		CurrentOS:    a.container.ShellIntegration.GetCurrentOS(),
	}, nil
}

// RegisterShellIntegration adds "Open in Syntaxia" to OS context menu
func (a *App) RegisterShellIntegration() error {
	if a.container.ShellIntegration == nil {
		return nil
	}

	if err := a.container.ShellIntegration.Register(); err != nil {
		a.log.Error("Failed to register shell integration: " + err.Error())
		return err
	}

	a.log.Info("Shell integration registered successfully")
	return nil
}

// UnregisterShellIntegration removes context menu entry
func (a *App) UnregisterShellIntegration() error {
	if a.container.ShellIntegration == nil {
		return nil
	}

	if err := a.container.ShellIntegration.Unregister(); err != nil {
		a.log.Error("Failed to unregister shell integration: " + err.Error())
		return err
	}

	a.log.Info("Shell integration unregistered successfully")
	return nil
}
