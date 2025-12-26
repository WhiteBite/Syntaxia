/**
 * Memory Profiling Script
 * 
 * Automatically creates heap snapshots for memory leak analysis.
 * Usage: npm run memory:profile
 * 
 * Features:
 * - Auto-dump heap snapshots every 30 seconds
 * - Manual dump on SIGUSR2 signal (kill -USR2 <pid>)
 * - Automatic dump when heap exceeds 2GB
 */

const fs = require('fs');
const path = require('path');
const v8 = require('v8');

console.log('[MemoryProfile] Starting memory profiling...');
console.log('[MemoryProfile] PID:', process.pid);
console.log('[MemoryProfile] Send SIGUSR2 to manually trigger snapshot: kill -USR2', process.pid);

// Manual snapshot trigger
process.on('SIGUSR2', () => {
    const filename = `heap-manual-${Date.now()}.heapsnapshot`;
    console.log('[MemoryProfile] Manual snapshot triggered:', filename);
    v8.writeHeapSnapshot(filename);
    console.log('[MemoryProfile] Snapshot saved:', filename);
});

// Auto-snapshot on high memory usage
let snapshotCount = 0;
const MAX_SNAPSHOTS = 10; // Limit to prevent disk space issues

setInterval(() => {
    const stats = process.memoryUsage();
    const heapUsedMB = Math.round(stats.heapUsed / (1024 * 1024));
    const heapTotalMB = Math.round(stats.heapTotal / (1024 * 1024));
    const rss = Math.round(stats.rss / (1024 * 1024));

    console.log(`[MemoryProfile] Heap: ${heapUsedMB}MB / ${heapTotalMB}MB | RSS: ${rss}MB`);

    // Auto-dump if heap exceeds 2GB
    if (stats.heapUsed > 2e9 && snapshotCount < MAX_SNAPSHOTS) {
        const filename = `heap-auto-${Date.now()}.heapsnapshot`;
        console.log('[MemoryProfile] High memory detected! Creating snapshot:', filename);
        v8.writeHeapSnapshot(filename);
        console.log('[MemoryProfile] Snapshot saved:', filename);
        snapshotCount++;

        if (snapshotCount >= MAX_SNAPSHOTS) {
            console.warn('[MemoryProfile] Max snapshots reached. No more auto-snapshots will be created.');
        }
    }
}, 30000); // Check every 30 seconds

// Start Vite dev server
console.log('[MemoryProfile] Starting Vite dev server...');
require('../node_modules/vite/bin/vite.js');
