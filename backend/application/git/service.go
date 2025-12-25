package git

import (
	"fmt"
	"syntaxia/domain"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Service реализует GitService
type Service struct {
	log domain.Logger
}

// NewService создает новый Service
func NewService(log domain.Logger) domain.GitService {
	return &Service{log: log}
}

// GenerateDiff генерирует git diff между HEAD и HEAD~1
func (s *Service) GenerateDiff(projectPath string) (string, error) {
	s.log.Info(fmt.Sprintf("Generating diff for project: %s", projectPath))

	// Открываем репозиторий
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Получаем HEAD
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Получаем коммит, на который указывает HEAD
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get commit object: %w", err)
	}

	// Получаем дерево текущего коммита
	toTree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get commit tree: %w", err)
	}

	// Получаем родительский коммит
	parent, err := commit.Parent(0)
	if err != nil {
		// Если родителя нет (первый коммит), сравниваем с пустым деревом
		patch, err := toTree.Patch(nil)
		if err != nil {
			return "", fmt.Errorf("failed to generate patch against empty tree: %w", err)
		}
		return patch.String(), nil
	}

	// Получаем дерево родительского коммита
	fromTree, err := parent.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get parent tree: %w", err)
	}

	// Получаем изменения
	patch, err := fromTree.Patch(toTree)
	if err != nil {
		return "", fmt.Errorf("failed to generate patch: %w", err)
	}

	return patch.String(), nil
}

// GetAllFiles возвращает список всех файлов в репозитории
func (s *Service) GetAllFiles(projectPath string) ([]string, error) {
	s.log.Info(fmt.Sprintf("Getting all files for project: %s", projectPath))
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree from commit: %w", err)
	}

	var files []string
	err = tree.Files().ForEach(func(f *object.File) error {
		files = append(files, f.Name)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate over files: %w", err)
	}

	return files, nil
}
