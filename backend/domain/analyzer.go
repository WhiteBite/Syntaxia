package domain

import "context"

// ContextAnalyzer определяет интерфейс для сервисов, которые могут
// анализировать проект и предлагать релевантные файлы для задачи.
type ContextAnalyzer interface {
	// SuggestFiles анализирует задачу пользователя и текущее состояние проекта (все файлы)
	// и возвращает список относительных путей к файлам, которые наиболее релевантны.
	SuggestFiles(ctx context.Context, task string, allFiles []*FileNode) ([]string, error)
}
