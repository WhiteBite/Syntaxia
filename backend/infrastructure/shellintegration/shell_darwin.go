//go:build darwin

package shellintegration

import (
	"fmt"
	"os"
	"path/filepath"
)

const serviceFileName = "Open in Syntaxia.workflow"

func (s *Service) getServicePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Services", serviceFileName)
}

func (s *Service) isRegisteredOS() (bool, error) {
	_, err := os.Stat(s.getServicePath())
	return err == nil, nil
}

func (s *Service) registerOS(exePath string) error {
	servicePath := s.getServicePath()
	workflowDir := filepath.Join(servicePath, "Contents")

	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>NSServices</key>
	<array>
		<dict>
			<key>NSMenuItem</key>
			<dict>
				<key>default</key>
				<string>Open in Syntaxia</string>
			</dict>
			<key>NSMessage</key>
			<string>runWorkflowAsService</string>
			<key>NSSendFileTypes</key>
			<array>
				<string>public.folder</string>
			</array>
		</dict>
	</array>
</dict>
</plist>`

	if err := os.WriteFile(filepath.Join(workflowDir, "Info.plist"), []byte(infoPlist), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %w", err)
	}

	documentWflow := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AMApplicationBuild</key>
	<string>523</string>
	<key>AMApplicationVersion</key>
	<string>2.10</string>
	<key>actions</key>
	<array>
		<dict>
			<key>action</key>
			<dict>
				<key>AMActionVersion</key>
				<string>1.0.2</string>
				<key>AMApplication</key>
				<array>
					<string>Automator</string>
				</array>
				<key>AMBundleIdentifier</key>
				<string>com.apple.RunShellScript</string>
				<key>AMParameterProperties</key>
				<dict>
					<key>COMMAND_STRING</key>
					<dict/>
					<key>CheckedForUserDefaultShell</key>
					<dict/>
					<key>inputMethod</key>
					<dict/>
					<key>shell</key>
					<dict/>
					<key>source</key>
					<dict/>
				</dict>
				<key>ActionParameters</key>
				<dict>
					<key>COMMAND_STRING</key>
					<string>for f in "$@"; do
	"%s" "$f"
done</string>
					<key>CheckedForUserDefaultShell</key>
					<true/>
					<key>inputMethod</key>
					<integer>1</integer>
					<key>shell</key>
					<string>/bin/zsh</string>
					<key>source</key>
					<string></string>
				</dict>
			</dict>
		</dict>
	</array>
</dict>
</plist>`, exePath)

	if err := os.WriteFile(filepath.Join(workflowDir, "document.wflow"), []byte(documentWflow), 0644); err != nil {
		return fmt.Errorf("failed to write document.wflow: %w", err)
	}

	return nil
}

func (s *Service) unregisterOS() error {
	servicePath := s.getServicePath()
	if err := os.RemoveAll(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service: %w", err)
	}
	return nil
}

// NeedsUpdate checks if registered path differs from current exe path
func (s *Service) NeedsUpdate(currentExePath string) bool {
	// macOS: simplified check - always return false for now
	// Full implementation would parse document.wflow
	return false
}
