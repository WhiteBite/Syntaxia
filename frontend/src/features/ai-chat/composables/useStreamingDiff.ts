import { computed, ref, type Ref } from 'vue'

export interface DiffLine {
    type: 'add' | 'remove' | 'context'
    content: string
    lineNumber: number | null
    oldLineNumber: number | null
}

export interface ParsedDiff {
    lines: DiffLine[]
    stats: { additions: number; deletions: number }
    filePath: string
    operation: 'create' | 'modify' | 'delete'
}

export function useStreamingDiff(content: Ref<string>, filePath: Ref<string>, original?: Ref<string | undefined>) {
    const isExpanded = ref(false)
    const operation = computed(() => !original?.value ? 'create' : !content.value ? 'delete' : 'modify')

    const parsedDiff = computed<ParsedDiff>(() => {
        const lines = getDiffLines(content.value, original?.value, operation.value)
        return {
            lines,
            stats: { additions: lines.filter(l => l.type === 'add').length, deletions: lines.filter(l => l.type === 'remove').length },
            filePath: filePath.value,
            operation: operation.value as 'create' | 'modify' | 'delete'
        }
    })

    return {
        parsedDiff,
        previewLines: computed(() => parsedDiff.value.lines.filter(l => l.type !== 'context').slice(0, 8)),
        stats: computed(() => parsedDiff.value.stats),
        operation,
        isExpanded,
        toggleExpand: () => { isExpanded.value = !isExpanded.value }
    }
}

function getDiffLines(content: string, original: string | undefined, op: string): DiffLine[] {
    if (op === 'create') return content.split('\n').map((c, i) => ({ type: 'add', content: c, lineNumber: i + 1, oldLineNumber: null }))
    if (op === 'delete') return (original || '').split('\n').map((c, i) => ({ type: 'remove', content: c, lineNumber: null, oldLineNumber: i + 1 }))

    const orig = original?.split('\n') || [], mod = content.split('\n'), lcs = getLcs(orig, mod), lines: DiffLine[] = []
    let oi = 0, mi = 0, li = 0, nn = 1, on = 1
    while (oi < orig.length || mi < mod.length) {
        if (li < lcs.length && orig[oi] === lcs[li] && mod[mi] === lcs[li]) {
            lines.push({ type: 'context', content: orig[oi], lineNumber: nn++, oldLineNumber: on++ }); oi++; mi++; li++
        } else if (mi < mod.length && (li >= lcs.length || mod[mi] !== lcs[li])) {
            lines.push({ type: 'add', content: mod[mi], lineNumber: nn++, oldLineNumber: null }); mi++
        } else if (oi < orig.length) {
            lines.push({ type: 'remove', content: orig[oi], lineNumber: null, oldLineNumber: on++ }); oi++
        }
    }
    return lines
}

function getLcs(a: string[], b: string[]): string[] {
    const dp: number[][] = Array(a.length + 1).fill(0).map(() => Array(b.length + 1).fill(0))
    for (let i = 1; i <= a.length; i++) for (let j = 1; j <= b.length; j++)
        dp[i][j] = a[i - 1] === b[j - 1] ? dp[i - 1][j - 1] + 1 : Math.max(dp[i - 1][j], dp[i][j - 1])
    const r: string[] = []
    for (let i = a.length, j = b.length; i > 0 && j > 0;) {
        if (a[i - 1] === b[j - 1]) {
            r.unshift(a[--i])
            j--
        } else if (dp[i - 1][j] > dp[i][j - 1]) {
            i--
        } else {
            j--
        }
    }
    return r
}
