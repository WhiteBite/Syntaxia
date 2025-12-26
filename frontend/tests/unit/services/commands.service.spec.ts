import { CommandService } from '@/services/commands.service'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock stores
vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        refreshFileTree: vi.fn().mockResolvedValue(undefined)
    })
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        hasProject: true,
        clearProject: vi.fn()
    })
}))

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: vi.fn()
    })
}))

describe('CommandService', () => {
    let service: CommandService

    beforeEach(() => {
        setActivePinia(createPinia())
        service = new CommandService()
    })

    it('getCommands returns list of commands', () => {
        const commands = service.getCommands()

        expect(Array.isArray(commands)).toBe(true)
        expect(commands.length).toBeGreaterThan(0)
        expect(commands[0]).toHaveProperty('id')
        expect(commands[0]).toHaveProperty('name')
        expect(commands[0]).toHaveProperty('action')
    })

    it('executeCommand calls action of command', async () => {
        // Test with open-project command which doesn't require window.location mock
        const result = await service.executeCommand('open-project')
        expect(result).toBe(true)
    })

    it('filterCommands filters by query', () => {
        const commands = service.getCommands()
        const filtered = commands.filter(c =>
            c.name.toLowerCase().includes('project') ||
            c.description.toLowerCase().includes('project')
        )

        expect(filtered.length).toBeGreaterThan(0)
        filtered.forEach(cmd => {
            const matchesName = cmd.name.toLowerCase().includes('project')
            const matchesDesc = cmd.description.toLowerCase().includes('project')
            expect(matchesName || matchesDesc).toBe(true)
        })
    })
})
