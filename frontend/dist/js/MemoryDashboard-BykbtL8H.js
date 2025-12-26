const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["js/index-CCibvNLt.js","js/vue-vendor-DLLjc1XI.js","assets/vue-vendor-Dw2MutP3.css","js/ui-vendor-2hF-XK7S.js","js/utils-DVYXdora.js","js/icons-BGgRd9bo.js","js/highlighter-D41BBaGQ.js","assets/highlighter-DqdTtQ31.css","assets/index-9dBOofo9.css","js/useApiCache-BTgc1mkV.js"])))=>i.map(i=>d[i]);
import { a as useLogger, _ as __vitePreload, b as useMemoryMonitor, u as useUIStore } from "./index-CCibvNLt.js";
import { e as defineComponent, r as ref, c as computed, f as onMounted, g as onUnmounted, y as createElementBlock, z as openBlock, C as createBaseVNode, A as createCommentVNode, D as toDisplayString, L as normalizeStyle, l as normalizeClass } from "./vue-vendor-DLLjc1XI.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./icons-BGgRd9bo.js";
import "./highlighter-D41BBaGQ.js";
const logger = useLogger("MemoryDiagnostics");
class MemoryDiagnostics {
  static instance;
  snapshots = [];
  issues = [];
  startTime = /* @__PURE__ */ new Date();
  isCollecting = false;
  collectionInterval = null;
  isDevMode = false;
  constructor() {
  }
  static getInstance() {
    if (!MemoryDiagnostics.instance) {
      MemoryDiagnostics.instance = new MemoryDiagnostics();
    }
    return MemoryDiagnostics.instance;
  }
  /**
   * Start automatic diagnostics collection
   * In dev mode, use longer intervals to reduce memory overhead from dynamic imports
   */
  startCollection(intervalMs = 1e4) {
    if (this.isCollecting) return;
    this.isCollecting = true;
    this.startTime = /* @__PURE__ */ new Date();
    this.snapshots = [];
    this.issues = [];
    const actualInterval = this.isDevMode ? Math.max(intervalMs, 6e4) : intervalMs;
    logger.debug(`Started collection (interval: ${actualInterval}ms, dev: ${this.isDevMode})`);
    void this.collectSnapshot();
    this.collectionInterval = window.setInterval(() => {
      void this.collectSnapshot().then(() => {
        this.analyzeSnapshots();
      });
    }, actualInterval);
  }
  /**
   * Stop collection
   */
  stopCollection() {
    if (this.collectionInterval) {
      clearInterval(this.collectionInterval);
      this.collectionInterval = null;
    }
    this.isCollecting = false;
  }
  /**
   * Collect a single diagnostic snapshot
   */
  async collectSnapshot() {
    const snapshot = {
      timestamp: (/* @__PURE__ */ new Date()).toISOString(),
      memory: this.getMemoryStats(),
      stores: await this.getStoreStats(),
      performance: this.getPerformanceStats(),
      browser: this.getBrowserInfo()
    };
    this.snapshots.push(snapshot);
    if (this.snapshots.length > 20) {
      this.snapshots.shift();
    }
    return snapshot;
  }
  /**
   * Get memory stats
   */
  getMemoryStats() {
    if (!("performance" in window) || !performance.memory) {
      return null;
    }
    const memory = performance.memory;
    const used = Math.round(memory.usedJSHeapSize / (1024 * 1024));
    const total = Math.round(memory.jsHeapSizeLimit / (1024 * 1024));
    const percentage = Math.round(memory.usedJSHeapSize / memory.jsHeapSizeLimit * 100);
    return { used, total, percentage };
  }
  // Cache store imports to avoid repeated dynamic imports causing memory leaks
  storeImportsCache = {};
  /**
   * Get store statistics
   * In dev mode, skip to reduce dynamic import overhead
   */
  async getStoreStats() {
    const stats = {
      fileStore: {
        nodesCount: 0,
        selectedCount: 0,
        memoryUsage: 0,
        searchQueryLength: 0,
        filterExtensionsCount: 0
      },
      contextStore: {
        hasContext: false,
        fileCount: 0,
        totalSize: 0,
        lineCount: 0,
        cacheSize: 0
      },
      apiCache: {
        entries: 0,
        size: 0,
        sizeMB: "0",
        hitRate: "0%",
        hits: 0,
        misses: 0
      }
    };
    if (this.isDevMode) {
      return stats;
    }
    try {
      if (!this.storeImportsCache.useFileStore) {
        const { useFileStore } = await __vitePreload(async () => {
          const { useFileStore: useFileStore2 } = await import("./index-CCibvNLt.js").then((n) => n.k);
          return { useFileStore: useFileStore2 };
        }, true ? __vite__mapDeps([0,1,2,3,4,5,6,7,8]) : void 0);
        this.storeImportsCache.useFileStore = useFileStore;
      }
      if (!this.storeImportsCache.useContextStore) {
        const { useContextStore } = await __vitePreload(async () => {
          const { useContextStore: useContextStore2 } = await import("./index-CCibvNLt.js").then((n) => n.j);
          return { useContextStore: useContextStore2 };
        }, true ? __vite__mapDeps([0,1,2,3,4,5,6,7,8]) : void 0);
        this.storeImportsCache.useContextStore = useContextStore;
      }
      if (!this.storeImportsCache.getCacheStats) {
        const { getCacheStats } = await __vitePreload(async () => {
          const { getCacheStats: getCacheStats2 } = await import("./useApiCache-BTgc1mkV.js");
          return { getCacheStats: getCacheStats2 };
        }, true ? __vite__mapDeps([9,0,1,2,3,4,5,6,7,8]) : void 0);
        this.storeImportsCache.getCacheStats = getCacheStats;
      }
      const fileStore = this.storeImportsCache.useFileStore?.();
      if (fileStore) {
        stats.fileStore = {
          nodesCount: fileStore.nodes?.length || 0,
          selectedCount: fileStore.selectedCount || 0,
          memoryUsage: fileStore.getMemoryUsage ? fileStore.getMemoryUsage() : 0,
          searchQueryLength: fileStore.searchQuery?.length || 0,
          filterExtensionsCount: fileStore.filterExtensions?.length || 0
        };
      }
      const contextStore = this.storeImportsCache.useContextStore?.();
      if (contextStore) {
        stats.contextStore = {
          hasContext: contextStore.hasContext || false,
          fileCount: contextStore.fileCount || 0,
          totalSize: contextStore.totalSize || 0,
          lineCount: contextStore.lineCount || 0,
          cacheSize: contextStore.getMemoryUsage ? contextStore.getMemoryUsage() : 0
        };
      }
      const cacheStats = this.storeImportsCache.getCacheStats?.();
      if (cacheStats) {
        stats.apiCache = {
          entries: cacheStats.entries || 0,
          size: cacheStats.size || 0,
          sizeMB: cacheStats.sizeMB || "0",
          hitRate: cacheStats.hitRate || "0%",
          hits: cacheStats.hits || 0,
          misses: cacheStats.misses || 0
        };
      }
    } catch (e) {
      console.warn("[MemoryDiagnostics] Error collecting store stats:", e);
    }
    return stats;
  }
  /**
   * Get performance statistics
   */
  getPerformanceStats() {
    const navigation = performance.getEntriesByType("navigation")[0];
    const resources = performance.getEntriesByType("resource").length;
    const measures = performance.getEntriesByType("measure").length;
    return { navigation, resources, measures };
  }
  /**
   * Get browser information
   */
  getBrowserInfo() {
    return {
      userAgent: navigator.userAgent,
      language: navigator.language,
      platform: navigator.platform,
      hardwareConcurrency: navigator.hardwareConcurrency,
      deviceMemory: navigator.deviceMemory
    };
  }
  /**
   * Analyze snapshots and detect issues
   */
  analyzeSnapshots() {
    if (this.snapshots.length < 2) return;
    const latest = this.snapshots[this.snapshots.length - 1];
    const previous = this.snapshots[this.snapshots.length - 2];
    if (latest.memory && previous.memory) {
      const growth = latest.memory.used - previous.memory.used;
      const growthPercent = growth / previous.memory.used * 100;
      if (growth > 50) {
        this.addIssue({
          severity: "warning",
          category: "memory",
          message: `Memory grew by ${growth}MB (${growthPercent.toFixed(1)}%)`,
          details: {
            previous: previous.memory.used,
            current: latest.memory.used,
            growth
          },
          timestamp: latest.timestamp
        });
      }
      if (latest.memory.percentage >= 90) {
        this.addIssue({
          severity: "critical",
          category: "memory",
          message: `Critical memory usage: ${latest.memory.percentage}%`,
          details: latest.memory,
          timestamp: latest.timestamp
        });
      }
    }
    const nodeGrowth = latest.stores.fileStore.nodesCount - previous.stores.fileStore.nodesCount;
    if (nodeGrowth > 1e3) {
      this.addIssue({
        severity: "warning",
        category: "store",
        message: `File store grew by ${nodeGrowth} nodes`,
        details: {
          previous: previous.stores.fileStore.nodesCount,
          current: latest.stores.fileStore.nodesCount
        },
        timestamp: latest.timestamp
      });
    }
    const hitRate = parseFloat(latest.stores.apiCache.hitRate);
    if (hitRate < 50 && latest.stores.apiCache.entries > 10) {
      this.addIssue({
        severity: "info",
        category: "cache",
        message: `Low cache hit rate: ${latest.stores.apiCache.hitRate}`,
        details: latest.stores.apiCache,
        timestamp: latest.timestamp
      });
    }
    if (latest.stores.contextStore.cacheSize > 10 * 1024 * 1024) {
      this.addIssue({
        severity: "warning",
        category: "store",
        message: `Large context cache: ${(latest.stores.contextStore.cacheSize / (1024 * 1024)).toFixed(1)}MB`,
        details: latest.stores.contextStore,
        timestamp: latest.timestamp
      });
    }
  }
  /**
   * Add diagnostic issue
   */
  addIssue(issue) {
    this.issues.push(issue);
    if (this.issues.length > 50) {
      this.issues.shift();
    }
    logger.debug(`${issue.severity.toUpperCase()}: ${issue.message}`, issue.details);
  }
  /**
   * Generate diagnostic report
   */
  generateReport() {
    const endTime = /* @__PURE__ */ new Date();
    const duration = endTime.getTime() - this.startTime.getTime();
    const recommendations = this.generateRecommendations();
    return {
      startTime: this.startTime.toISOString(),
      endTime: endTime.toISOString(),
      duration,
      snapshots: this.snapshots,
      issues: this.issues,
      recommendations
    };
  }
  /**
   * Generate recommendations based on collected data
   */
  generateRecommendations() {
    const recommendations = [];
    if (this.snapshots.length === 0) {
      return ["No data collected yet"];
    }
    const latest = this.snapshots[this.snapshots.length - 1];
    if (latest.memory && latest.memory.percentage > 75) {
      recommendations.push("Memory usage is high. Consider reducing file selection or context size.");
    }
    if (latest.stores.fileStore.nodesCount > 5e3) {
      recommendations.push("Large file tree detected. Consider using file filters or working in subdirectories.");
    }
    if (latest.stores.fileStore.selectedCount > 100) {
      recommendations.push("Many files selected. This may impact performance.");
    }
    if (latest.stores.contextStore.cacheSize > 15 * 1024 * 1024) {
      recommendations.push("Context cache is large. Clear context when not needed.");
    }
    const hitRate = parseFloat(latest.stores.apiCache.hitRate);
    if (hitRate < 50 && latest.stores.apiCache.entries > 10) {
      recommendations.push("Low cache hit rate. Cache may not be effective.");
    }
    if (this.snapshots.length >= 5) {
      const first = this.snapshots[0];
      const memoryGrowth = latest.memory && first.memory ? latest.memory.used - first.memory.used : 0;
      if (memoryGrowth > 100) {
        recommendations.push(`Memory grew by ${memoryGrowth}MB since start. Possible memory leak.`);
      }
    }
    if (recommendations.length === 0) {
      recommendations.push("No issues detected. Memory usage is normal.");
    }
    return recommendations;
  }
  /**
   * Save report to file (as JSON for AI to read)
   */
  async saveReportToFile() {
    const report = this.generateReport();
    const filename = `memory-diagnostics-${Date.now()}.json`;
    try {
      const blob = new Blob([JSON.stringify(report, null, 2)], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      a.click();
      URL.revokeObjectURL(url);
      logger.debug("Report saved:", filename);
      return filename;
    } catch (e) {
      console.error("[MemoryDiagnostics] Failed to save report:", e);
      throw e;
    }
  }
  /**
   * Get current status summary
   */
  getStatusSummary() {
    if (this.snapshots.length === 0) {
      return "No diagnostics collected yet";
    }
    const latest = this.snapshots[this.snapshots.length - 1];
    const criticalIssues = this.issues.filter((i) => i.severity === "critical").length;
    const warnings = this.issues.filter((i) => i.severity === "warning").length;
    let summary = `Diagnostics Summary:
`;
    summary += `- Snapshots collected: ${this.snapshots.length}
`;
    summary += `- Critical issues: ${criticalIssues}
`;
    summary += `- Warnings: ${warnings}
`;
    if (latest.memory) {
      summary += `- Current memory: ${latest.memory.used}MB / ${latest.memory.total}MB (${latest.memory.percentage}%)
`;
    }
    summary += `- File nodes: ${latest.stores.fileStore.nodesCount}
`;
    summary += `- Selected files: ${latest.stores.fileStore.selectedCount}
`;
    summary += `- Context cache: ${(latest.stores.contextStore.cacheSize / (1024 * 1024)).toFixed(1)}MB
`;
    summary += `- API cache: ${latest.stores.apiCache.entries} entries (${latest.stores.apiCache.sizeMB}MB)
`;
    return summary;
  }
  /**
   * Get latest snapshot
   */
  getLatestSnapshot() {
    return this.snapshots.length > 0 ? this.snapshots[this.snapshots.length - 1] : null;
  }
  /**
   * Get all issues
   */
  getIssues() {
    return this.issues;
  }
  /**
   * Clear all data
   */
  clear() {
    this.snapshots = [];
    this.issues = [];
    this.startTime = /* @__PURE__ */ new Date();
  }
}
function useMemoryDiagnostics() {
  return MemoryDiagnostics.getInstance();
}
const _hoisted_1 = { class: "bg-gray-800/90 backdrop-blur p-4 rounded-lg border border-gray-600 max-w-sm shadow-xl" };
const _hoisted_2 = { class: "flex items-center justify-between mb-3" };
const _hoisted_3 = {
  key: 0,
  class: "space-y-3 text-sm"
};
const _hoisted_4 = { class: "flex justify-between text-gray-300 mb-1" };
const _hoisted_5 = { class: "font-mono" };
const _hoisted_6 = { class: "w-full bg-gray-700 rounded-full h-2" };
const _hoisted_7 = { class: "text-xs text-gray-400 mt-1" };
const _hoisted_8 = {
  key: 0,
  class: "border-t border-gray-700 pt-3 space-y-2"
};
const _hoisted_9 = { class: "text-gray-300" };
const _hoisted_10 = { class: "flex justify-between text-xs" };
const _hoisted_11 = { class: "font-mono" };
const _hoisted_12 = { class: "flex justify-between text-xs" };
const _hoisted_13 = { class: "font-mono" };
const _hoisted_14 = { class: "flex justify-between text-xs" };
const _hoisted_15 = { class: "font-mono" };
const _hoisted_16 = {
  key: 1,
  class: "border-t border-gray-700 pt-3 space-y-1 text-xs"
};
const _hoisted_17 = { class: "flex justify-between" };
const _hoisted_18 = { class: "font-mono" };
const _hoisted_19 = { class: "flex justify-between" };
const _hoisted_20 = { class: "text-xs text-gray-400 text-center" };
const _hoisted_21 = {
  key: 1,
  class: "text-gray-400 text-sm text-center py-4"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "MemoryDashboard",
  emits: ["close"],
  setup(__props) {
    const logger2 = useLogger("MemoryDashboard");
    const memoryMonitor = useMemoryMonitor();
    const diagnostics = useMemoryDiagnostics();
    const uiStore = useUIStore();
    const stats = ref(null);
    const lastUpdate = ref(Date.now());
    const timeSinceUpdate = ref(0);
    const storeMetrics = computed(() => stats.value?.storeMetrics);
    const diagnosticsInfo = computed(() => {
      const issues = diagnostics.getIssues();
      return {
        snapshotCount: diagnostics.getLatestSnapshot() ? diagnostics["snapshots"]?.length || 0 : 0,
        criticalIssues: issues.filter((i) => i.severity === "critical").length,
        warnings: issues.filter((i) => i.severity === "warning").length
      };
    });
    let updateInterval = null;
    let timeInterval = null;
    async function updateStats() {
      stats.value = await memoryMonitor.getMemoryStats();
      lastUpdate.value = Date.now();
      timeSinceUpdate.value = 0;
    }
    function dumpSnapshot() {
      memoryMonitor.dumpHeapSnapshot("manual-dashboard");
    }
    function forceCleanup() {
      memoryMonitor.forceCleanup();
      setTimeout(() => void updateStats(), 1e3);
    }
    async function saveReport() {
      try {
        const filename = await diagnostics.saveReportToFile();
        uiStore.addToast(`Diagnostic report saved: ${filename}`, "success");
        logger2.debug("Report saved for AI analysis");
        logger2.debug(diagnostics.getStatusSummary());
      } catch (e) {
        uiStore.addToast("Failed to save diagnostic report", "error");
        logger2.error("Failed to save report:", e);
      }
    }
    function formatBytes(bytes) {
      if (bytes === 0) return "0 B";
      const k = 1024;
      const sizes = ["B", "KB", "MB", "GB"];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return Math.round(bytes / Math.pow(k, i) * 10) / 10 + " " + sizes[i];
    }
    onMounted(() => {
      void updateStats();
      diagnostics.startCollection(3e4);
      const statsInterval = 5e3;
      const timeCounterInterval = 2e3;
      updateInterval = window.setInterval(() => void updateStats(), statsInterval);
      timeInterval = window.setInterval(() => {
        timeSinceUpdate.value = Math.floor((Date.now() - lastUpdate.value) / 1e3);
      }, timeCounterInterval);
    });
    onUnmounted(() => {
      if (updateInterval) clearInterval(updateInterval);
      if (timeInterval) clearInterval(timeInterval);
      diagnostics.stopCollection();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("div", _hoisted_2, [
          _cache[1] || (_cache[1] = createBaseVNode("h3", { class: "text-lg font-semibold text-white flex items-center gap-2" }, " 🧠 Memory Monitor ", -1)),
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close")),
            class: "text-gray-400 hover:text-white transition-colors",
            title: "Close"
          }, " ✕ ")
        ]),
        stats.value ? (openBlock(), createElementBlock("div", _hoisted_3, [
          createBaseVNode("div", null, [
            createBaseVNode("div", _hoisted_4, [
              _cache[2] || (_cache[2] = createBaseVNode("span", null, "Heap Usage", -1)),
              createBaseVNode("span", _hoisted_5, toDisplayString(stats.value.used) + "MB / " + toDisplayString(stats.value.total) + "MB", 1)
            ]),
            createBaseVNode("div", _hoisted_6, [
              createBaseVNode("div", {
                class: normalizeClass([
                  "h-2 rounded-full transition-all duration-300",
                  stats.value.percentage >= 90 ? "bg-red-500" : stats.value.percentage >= 75 ? "bg-yellow-500" : "bg-green-500"
                ]),
                style: normalizeStyle({ width: `${stats.value.percentage}%` })
              }, null, 6)
            ]),
            createBaseVNode("div", _hoisted_7, toDisplayString(stats.value.percentage) + "%", 1)
          ]),
          storeMetrics.value ? (openBlock(), createElementBlock("div", _hoisted_8, [
            createBaseVNode("div", _hoisted_9, [
              _cache[6] || (_cache[6] = createBaseVNode("div", { class: "font-semibold mb-1" }, "Store Metrics", -1)),
              createBaseVNode("div", _hoisted_10, [
                _cache[3] || (_cache[3] = createBaseVNode("span", null, "File Store:", -1)),
                createBaseVNode("span", _hoisted_11, toDisplayString(storeMetrics.value.fileStore.nodesCount) + " nodes", 1)
              ]),
              createBaseVNode("div", _hoisted_12, [
                _cache[4] || (_cache[4] = createBaseVNode("span", null, "Context Cache:", -1)),
                createBaseVNode("span", _hoisted_13, toDisplayString(formatBytes(storeMetrics.value.contextStore.cacheSize)), 1)
              ]),
              createBaseVNode("div", _hoisted_14, [
                _cache[5] || (_cache[5] = createBaseVNode("span", null, "API Cache:", -1)),
                createBaseVNode("span", _hoisted_15, toDisplayString(storeMetrics.value.apiCache.entries) + " entries", 1)
              ])
            ])
          ])) : createCommentVNode("", true),
          diagnosticsInfo.value ? (openBlock(), createElementBlock("div", _hoisted_16, [
            _cache[9] || (_cache[9] = createBaseVNode("div", { class: "text-gray-400" }, "Diagnostics:", -1)),
            createBaseVNode("div", _hoisted_17, [
              _cache[7] || (_cache[7] = createBaseVNode("span", null, "Snapshots:", -1)),
              createBaseVNode("span", _hoisted_18, toDisplayString(diagnosticsInfo.value.snapshotCount), 1)
            ]),
            createBaseVNode("div", _hoisted_19, [
              _cache[8] || (_cache[8] = createBaseVNode("span", null, "Issues:", -1)),
              createBaseVNode("span", {
                class: normalizeClass(["font-mono", diagnosticsInfo.value.criticalIssues > 0 ? "text-red-400" : "text-gray-300"])
              }, toDisplayString(diagnosticsInfo.value.criticalIssues) + " critical, " + toDisplayString(diagnosticsInfo.value.warnings) + " warnings ", 3)
            ])
          ])) : createCommentVNode("", true),
          createBaseVNode("div", { class: "border-t border-gray-700 pt-3 space-y-2" }, [
            createBaseVNode("div", { class: "flex gap-2" }, [
              createBaseVNode("button", {
                onClick: dumpSnapshot,
                class: "flex-1 px-3 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-medium transition-colors",
                title: "Create heap snapshot for Chrome DevTools"
              }, " 📸 Dump Heap "),
              createBaseVNode("button", {
                onClick: forceCleanup,
                class: "flex-1 px-3 py-2 bg-orange-600 hover:bg-orange-500 text-white rounded text-xs font-medium transition-colors",
                title: "Force garbage collection and cleanup"
              }, " 🧹 Cleanup ")
            ]),
            createBaseVNode("button", {
              onClick: saveReport,
              class: "w-full px-3 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded text-xs font-medium transition-colors",
              title: "Save diagnostic report for AI analysis"
            }, " 💾 Save Report (for AI) ")
          ]),
          createBaseVNode("div", _hoisted_20, " Updated " + toDisplayString(timeSinceUpdate.value) + "s ago ", 1)
        ])) : (openBlock(), createElementBlock("div", _hoisted_21, " Memory stats unavailable "))
      ]);
    };
  }
});
export {
  _sfc_main as default
};
