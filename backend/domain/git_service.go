package domain

// GitService определяет интерфейс для работы с Git
type GitService interface {
	GenerateDiff(projectPath string) (string, error)
}
