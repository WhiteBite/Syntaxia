export enum GitStatus {
  Unmodified = "unmodified", // Mapped from empty or " " in Git porcelain
  Modified = "modified", // Mapped from M
  Added = "added", // Mapped from A
  Deleted = "deleted", // Mapped from D
  Renamed = "renamed", // Mapped from R
  Copied = "copied", // Mapped from C
  Unmerged = "unmerged", // Mapped from U / D/A / A/U etc. -> Renamed to UnmergedConflict
  Untracked = "untracked", // Mapped from ??
  Ignored = "ignored", // Files that are ignored by .gitignore or custom rules
  UnmergedConflict = "unmerged-conflict", // Mapped from UM (backend custom)
}

export enum ContextOrigin {
  None = "None",
  Selected = "User Selection",
  Git = "Git Changes",
  AI = "AI Suggestion",
}

export enum ExportMode {
  Clipboard = "clipboard",
  AI = "ai",
  Human = "human",
}