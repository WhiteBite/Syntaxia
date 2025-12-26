import type { ContextOrigin, GitStatus } from "./enums";

export interface DomainFileNode {
  name: string;
  path: string;
  relPath: string;
  isDir: boolean;
  children?: DomainFileNode[];
  isGitignored: boolean;
  isCustomIgnored: boolean;
  isIgnored: boolean; // Добавляем для синхронизации с api.ts
  size: number;
}

export interface FileNode {
  name: string;
  path: string;
  relPath: string;
  isDir: boolean;
  children?: FileNode[];
  depth?: number;
  gitStatus?: GitStatus;
  contextOrigin?: ContextOrigin;
  isBinary?: boolean;
  isIgnored: boolean;
  isGitignored: boolean;
  isCustomIgnored: boolean;
  parentPath?: string | null;
  size: number;
  // isSelected removed - this is UI state, managed in stores
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
  isMerge: boolean;
  files: string[];
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
  autonomousMode?: boolean;
  // AI Provider specific settings
  openaiModel?: string;
  geminiModel?: string;
  // Context settings
  maxContextSize?: number;
  maxFilesInContext?: number;
  includeDependencies?: boolean;
  includeTests?: boolean;
  // Context splitting settings
  maxTokensPerChunk?: number;
  overlapTokens?: number;
  splitStrategy?: 'smart' | 'file' | 'token';
  // Code generation settings
  temperature?: number;
  maxTokens?: number;
  autoFormat?: boolean;
  includeComments?: boolean;
  // Safety & validation
  enableGuardrails?: boolean;
  autoTest?: boolean;
  validateSyntax?: boolean;
}

export interface Hunk {
  header: string;
  lines: string[];
}

export interface FileDiff {
  filePath: string;
  hunks: Hunk[];
  stats: { added: number; removed: number };
}

export type ToastType = "info" | "success" | "error" | "warning";
export type LogType = ToastType | "warning" | "debug";

export interface LogEntry {
  id: number;
  message: string;
  type: LogType;
  timestamp: string;
}

export interface AICodeGenerationRequest {
  prompt: string;
  language?: string;
  includeTests?: boolean;
  includeComments?: boolean;
  maxTokens?: number;
  temperature?: number;
}

export interface AICodeGenerationResponse {
  code?: string;
  analysis?: string;
  complexity?: 'low' | 'medium' | 'high';
}

// Autonomous Mode DTOs (синхронизировано с api.ts)
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

export interface GuardrailsReport {
  taskId: string;
  violations: Array<{
    rule: string;
    severity: "error" | "warning" | "info";
    message: string;
    file?: string;
    line?: number;
  }>;
  passed: boolean;
  createdAt: string;
}

export interface SbomLicensingReport {
  taskId: string;
  dependencies: Array<{
    name: string;
    version: string;
    license: string;
    isCompatible: boolean;
  }>;
  incompatibleLicenses: string[];
  createdAt: string;
}

/** Union type for all report data types */
export type ReportData =
  | WhyViewReport
  | TimeToGreenReport
  | DerivedDiffReport
  | GuardrailsReport
  | SbomLicensingReport;

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
  data: ReportData;
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

// Project Structure Types (Phase 6)
export type ArchitectureType =
  | 'clean'
  | 'hexagonal'
  | 'mvc'
  | 'mvvm'
  | 'layered'
  | 'microservices'
  | 'monolith'
  | 'serverless'
  | 'event-driven'
  | 'ddd'
  | 'unknown';

export type NamingStyle = 'camelCase' | 'snake_case' | 'PascalCase' | 'kebab-case' | 'mixed';

export type FolderStructure = 'flat' | 'by-feature' | 'by-type' | 'by-layer' | 'hybrid';

export interface LayerInfo {
  name: string;
  path: string;
  description: string;
  dependencies: string[];
  files: string[];
  patterns: string[];
}

export interface ArchitectureInfo {
  type: ArchitectureType | string;
  confidence: number;
  description: string;
  indicators: string[];
  layers: LayerInfo[];
}

export interface FrameworkInfo {
  name: string;
  version: string;
  category: string;
  language: string;
  configFiles: string[];
  indicators: string[];
  bestPractices: string[];
}

export interface BuildSystemInfo {
  name: string;
  configFile: string;
  scripts: string[];
  commands: string[];
}

export interface LanguageInfo {
  name: string;
  version: string;
  fileCount: number;
  percentage: number;
  primary: boolean;
}

export interface FileNamingStyle {
  style: NamingStyle;
  suffixes: string[];
  prefixes: string[];
  examples: string[];
}

export interface TestConventions {
  location: string;
  fileSuffix: string;
  framework: string;
  patterns: string[];
}

export interface ImportStyle {
  absoluteImports: boolean;
  aliasedImports: boolean;
  groupedImports: boolean;
  importOrder: string[];
}

export interface CodeStyleInfo {
  indentStyle: string;
  indentSize: number;
  maxLineLength: number;
  trailingComma: boolean;
  semicolons: boolean;
  quoteStyle: string;
  configFile: string;
}

export interface ConventionInfo {
  namingStyle: NamingStyle;
  fileNaming: FileNamingStyle;
  folderStructure: FolderStructure;
  testConventions: TestConventions;
  importStyle: ImportStyle;
  codeStyle: CodeStyleInfo;
}

export interface ProjectStructure {
  architecture?: ArchitectureInfo;
  conventions?: ConventionInfo;
  frameworks: FrameworkInfo[];
  buildSystems: BuildSystemInfo[];
  languages: LanguageInfo[];
  layers: LayerInfo[];
  projectType: string;
  confidence: number;
}
