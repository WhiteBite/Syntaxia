// API Types for better type safety

export interface ApiResponse<T = unknown> {
  data?: T;
  error?: string;
  success: boolean;
}

export interface FileNode {
  name: string;
  path: string;
  relPath: string;
  isDir: boolean;
  size: number;
  children?: FileNode[];
  isGitignored: boolean;
  isCustomIgnored: boolean;
  isIgnored: boolean; // Синхронизировано с dto.ts
  // isSelected removed - this is UI state, not domain data
}

export interface FileStatus {
  path: string;
  status: string;
}

export interface CommitWithFiles {
  hash: string;
  subject: string;
  author: string;
  date: string;
  files: string[];
  isMerge: boolean;
}

export interface SettingsDTO {
  customIgnoreRules: string;
  customPromptRules: string;
  openAIAPIKey: string;
  geminiAPIKey: string;
  openRouterAPIKey: string;
  localAIAPIKey: string;
  localAIHost: string;
  localAIModelName: string;
  selectedProvider: string;
  selectedModels: Record<string, string>;
  availableModels: Record<string, string[]>;
  useGitignore: boolean;
  useCustomIgnore: boolean;
}

export interface ExportSettings {
  mode: "clipboard" | "ai" | "human";
  context: string;
  stripComments: boolean;
  includeManifest: boolean;
  exportFormat: "plain" | "manifest" | "json";
  aiProfile: string;
  tokenLimit: number;
  fileSizeLimitKB: number;
  enableAutoSplit: boolean;
  maxTokensPerChunk: number;
  overlapTokens: number;
  splitStrategy: "smart" | "file" | "token";
  theme: string;
  includeLineNumbers: boolean;
  includePageNumbers: boolean;
}

export interface ExportResult {
  mode: string;
  text?: string;
  fileName?: string;
  dataBase64?: string;
  filePath?: string;
  isLarge?: boolean;
  sizeBytes?: number;
}

export interface ContextAnalysisResult {
  task: string;
  taskType: string;
  priority: string;
  selectedFiles: FileNode[];
  dependencyFiles: FileNode[];
  context: string;
  analysisTime: number;
  recommendations: string[];
  estimatedTokens: number;
  confidence: number;
}

export interface Project {
  name: string;
  path: string;
  lastOpened?: string;
}

export interface Context {
  id: string;
  name: string;
  description: string;
  content: string;
  files: string[];
  createdAt: string;
  updatedAt: string;
  projectPath: string;
  tokenCount: number;
}

// Error types
export class ApiError extends Error {
  constructor(
      message: string,
      public statusCode?: number,
      public code?: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export class ValidationError extends Error {
  constructor(
      message: string,
      public field?: string,
  ) {
    super(message);
    this.name = "ValidationError";
  }
}

export interface AutonomousTaskRequest {
  task: string;
  slaPolicy: "lite" | "standard" | "strict";
  projectPath: string;
  options?: {
    maxTokens?: number;
    temperature?: number;
    enableStaticAnalysis?: boolean;
    enableTests?: boolean;
    enableSBOM?: boolean;
  };
}

export interface AutonomousTaskResponse {
  taskId: string;
  status: "accepted" | "rejected";
  message?: string;
}

export interface AutonomousTaskStatus {
  taskId: string;
  status: "pending" | "running" | "completed" | "failed" | "cancelled";
  currentStep?: string;
  progress: number; // 0-100
  estimatedTimeRemaining?: number; // seconds
  error?: string;
  startedAt: string;
  updatedAt: string;
}

export interface TPLPlanStep {
  id: string;
  operation: string;
  description: string;
  status: "pending" | "running" | "completed" | "failed" | "skipped";
  duration?: number; // milliseconds
  error?: string;
  dependencies?: string[]; // IDs of steps this depends on
}

export interface TPLPlan {
  taskId: string;
  steps: TPLPlanStep[];
  totalSteps: number;
  completedSteps: number;
  estimatedTotalTime?: number; // seconds
}

export interface WhyViewReport {
  taskId: string;
  rationale: string;
  decisionFactors: string[];
  impactedFiles: string[];
  riskAssessment: "low" | "medium" | "high";
  confidence: number; // 0-1
  alternatives: string[];
  createdAt: string;
}

export interface TimeToGreenReport {
  taskId: string;
  phases: {
    planning: number; // milliseconds
    generation: number;
    verification: number;
    total: number;
  };
  breakdown: {
    contextAnalysis: number;
    codeGeneration: number;
    compilation: number;
    testing: number;
    staticAnalysis: number;
    sbomScanning: number;
  };
  bottlenecks: string[];
  recommendations: string[];
  createdAt: string;
}

export interface DerivedDiffReport {
  taskId: string;
  summary: {
    filesChanged: number;
    linesAdded: number;
    linesRemoved: number;
    riskLevel: "low" | "medium" | "high";
  };
  changes: {
    filePath: string;
    hunks: Array<{
      header: string;
      lines: string[];
    }>;
    stats: {
      added: number;
      removed: number;
    };
  }[];
  impactAnalysis: {
    criticalFiles: string[];
    testCoverage: number; // 0-1
    breakingChanges: string[];
  };
  createdAt: string;
}

export interface GenericReport {
  id: string;
  type:
      | "why_view"
      | "time_to_green"
      | "derived_diff"
      | "guardrails"
      | "sbom_licensing";
  taskId: string;
  title: string;
  summary: string;
  data: unknown;
  createdAt: string;
  updatedAt: string;
}

export interface SLAPolicy {
  name: string;
  description: string;
  maxExecutionTime: number; // seconds
  requiredPhases: string[];
  qualityGates: {
    staticAnalysis: boolean;
    tests: boolean;
    sbomScanning: boolean;
    guardrails: boolean;
  };
  riskTolerance: "low" | "medium" | "high";
}