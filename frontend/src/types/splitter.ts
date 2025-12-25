/* Auto-generated: types for content splitting and copy actions */
export type ClipboardFormat = "plain" | "manifest" | "json";
export type SplitStrategy = "smart" | "file" | "token";

export interface SplitSettings {
  enableAutoSplit: boolean;
  maxTokensPerChunk: number;
  overlapTokens: number;
  splitStrategy: SplitStrategy;
}

export interface TokenEstimateOptions {
  estimator?: (text: string) => number; // optional precise estimator (tiktoken etc.)
}

export interface Segment {
  id: string;
  title?: string;
  start: number; // inclusive
  end: number; // exclusive
  isFile?: boolean;
  relPath?: string;
}

export interface SplitInput {
  text: string;
  segments?: Segment[];
}

export interface SplitChunkRef {
  segmentId: string;
  from: number;
  to: number;
}

export interface SplitChunk {
  index: number;
  start: number;
  end: number;
  text: string;
  tokens: number;
  chars: number;
  refs?: SplitChunkRef[];
  note?: string;
}

export interface SplitPreview {
  totalTokens: number;
  totalChars: number;
  chunkCount: number;
  chunks: SplitChunk[];
  warnings: string[];
}

export interface CopyRequest {
  target: "all" | "chunk";
  chunkIndex?: number;
  format: ClipboardFormat;
  stripComments: boolean;
}

export interface CopyPlan {
  request: CopyRequest;
  parts: string[];
}