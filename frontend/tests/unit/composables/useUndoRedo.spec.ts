import { useUndoRedo } from '@/composables/useUndoRedo'
import { describe, expect, it } from 'vitest'

describe('useUndoRedo', () => {
    it('push adds state to history', () => {
        const { push, historyLength } = useUndoRedo<string[]>([])

        expect(historyLength.value).toBe(1) // Initial state

        push(['file1.ts'])
        expect(historyLength.value).toBe(2)

        push(['file1.ts', 'file2.ts'])
        expect(historyLength.value).toBe(3)
    })

    it('undo returns previous state', () => {
        const { push, undo } = useUndoRedo<string[]>([])

        push(['a'])
        push(['a', 'b'])
        push(['a', 'b', 'c'])

        const prevState = undo()
        expect(prevState).toEqual(['a', 'b'])

        const prevState2 = undo()
        expect(prevState2).toEqual(['a'])
    })

    it('redo returns next state', () => {
        const { push, undo, redo } = useUndoRedo<string[]>([])

        push(['a'])
        push(['a', 'b'])

        undo() // Go back to ['a']

        const nextState = redo()
        expect(nextState).toEqual(['a', 'b'])
    })

    it('canUndo returns false at start', () => {
        const { canUndo } = useUndoRedo<string[]>([])

        expect(canUndo.value).toBe(false)
    })

    it('canRedo returns false at end', () => {
        const { push, canRedo } = useUndoRedo<string[]>([])

        push(['a'])
        expect(canRedo.value).toBe(false)
    })

    it('history is limited to maxHistory states', () => {
        const { push, historyLength } = useUndoRedo<number>(0, 20)

        for (let i = 1; i <= 30; i++) {
            push(i)
        }

        expect(historyLength.value).toBe(20)
    })
})
