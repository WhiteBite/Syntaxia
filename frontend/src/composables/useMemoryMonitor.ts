/**
 * Composable for Memory Monitoring
 * 
 * Provides easy access to memory monitoring functionality
 */

import type { MemoryWarningOptions } from '@/utils/memory-monitor'
import { MemoryMonitor } from '@/utils/memory-monitor'

export function useMemoryMonitor(options?: MemoryWarningOptions): MemoryMonitor {
    return MemoryMonitor.getInstance(options)
}

export default useMemoryMonitor
