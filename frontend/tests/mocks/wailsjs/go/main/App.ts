/**
 * Мок для wailsjs/go/main/App
 * Используется в тестах вместо реального Wails API
 */
import { vi } from 'vitest'

export const SelectDirectory = vi.fn().mockResolvedValue('/mocked/path')
export const GetRecentProjects = vi.fn().mockResolvedValue([])
export const AddRecentProject = vi.fn().mockResolvedValue(undefined)
export const RemoveRecentProject = vi.fn().mockResolvedValue(undefined)
export const GetCurrentDirectory = vi.fn().mockResolvedValue('/home/user')
export const GetProjectContext = vi.fn().mockResolvedValue({})
export const BuildProject = vi.fn().mockResolvedValue({ success: true })
export const RunTests = vi.fn().mockResolvedValue({ success: true })
export const AnalyzeProject = vi.fn().mockResolvedValue({})
export const ExportProject = vi.fn().mockResolvedValue({ success: true })
export const GetVersionInfo = vi.fn().mockResolvedValue({ version: '1.0.0', buildDate: '2024-01-01' })
export const GetReleases = vi.fn().mockResolvedValue([])
export const OpenProject = vi.fn().mockResolvedValue(undefined)
export const CloseProject = vi.fn().mockResolvedValue(undefined)
export const GetSettings = vi.fn().mockResolvedValue({})
export const SaveSettings = vi.fn().mockResolvedValue(undefined)
export const ReadFileContent = vi.fn().mockResolvedValue('')
export const GetFileStats = vi.fn().mockResolvedValue({ size: 0, modTime: '', isDir: false })
export const GetProjectStructure = vi.fn().mockResolvedValue({ frameworks: [], languages: [], architecture: '', patterns: [] })
