export interface FileTypeInfo {
  icon: string
  colorClass: string
}

// Единый источник иконок и цветов для типов файлов
export const FILE_TYPE_CONFIG: Record<string, FileTypeInfo> = {
  // TypeScript/JavaScript
  ts: { icon: '🔷', colorClass: 'bg-blue-500' },
  tsx: { icon: '⚛️', colorClass: 'bg-cyan-500' },
  js: { icon: '🟨', colorClass: 'bg-yellow-500' },
  jsx: { icon: '⚛️', colorClass: 'bg-cyan-500' },
  mjs: { icon: '🟨', colorClass: 'bg-yellow-500' },
  cjs: { icon: '🟨', colorClass: 'bg-yellow-500' },
  // Frameworks
  vue: { icon: '💚', colorClass: 'bg-emerald-500' },
  svelte: { icon: '🔥', colorClass: 'bg-orange-500' },
  // Backend
  go: { icon: '🐹', colorClass: 'bg-sky-500' },
  py: { icon: '🐍', colorClass: 'bg-yellow-600' },
  java: { icon: '☕', colorClass: 'bg-orange-600' },
  kt: { icon: '🟣', colorClass: 'bg-purple-500' },
  kts: { icon: '🟣', colorClass: 'bg-purple-500' },
  rs: { icon: '🦀', colorClass: 'bg-orange-500' },
  rb: { icon: '💎', colorClass: 'bg-red-500' },
  php: { icon: '🐘', colorClass: 'bg-indigo-400' },
  cs: { icon: '🟦', colorClass: 'bg-purple-600' },
  // C/C++
  cpp: { icon: '⚙️', colorClass: 'bg-blue-600' },
  cc: { icon: '⚙️', colorClass: 'bg-blue-600' },
  cxx: { icon: '⚙️', colorClass: 'bg-blue-600' },
  c: { icon: '⚙️', colorClass: 'bg-blue-500' },
  h: { icon: '📎', colorClass: 'bg-blue-400' },
  hpp: { icon: '📎', colorClass: 'bg-blue-400' },
  // Styles
  css: { icon: '🎨', colorClass: 'bg-pink-500' },
  scss: { icon: '🎨', colorClass: 'bg-pink-600' },
  sass: { icon: '🎨', colorClass: 'bg-pink-600' },
  less: { icon: '🎨', colorClass: 'bg-indigo-500' },
  // Markup
  html: { icon: '🌐', colorClass: 'bg-orange-500' },
  htm: { icon: '🌐', colorClass: 'bg-orange-500' },
  xml: { icon: '📰', colorClass: 'bg-orange-400' },
  svg: { icon: '🖼️', colorClass: 'bg-yellow-500' },
  // Data/Config
  json: { icon: '📋', colorClass: 'bg-yellow-400' },
  yaml: { icon: '⚙️', colorClass: 'bg-red-400' },
  yml: { icon: '⚙️', colorClass: 'bg-red-400' },
  toml: { icon: '⚙️', colorClass: 'bg-gray-500' },
  ini: { icon: '⚙️', colorClass: 'bg-gray-500' },
  env: { icon: '🔐', colorClass: 'bg-yellow-600' },
  // Documentation
  md: { icon: '📝', colorClass: 'bg-gray-400' },
  mdx: { icon: '📝', colorClass: 'bg-gray-400' },
  txt: { icon: '📄', colorClass: 'bg-gray-400' },
  // Database
  sql: { icon: '🗃️', colorClass: 'bg-indigo-500' },
  // Build/Config files
  gradle: { icon: '🐘', colorClass: 'bg-green-600' },
  // Shell
  sh: { icon: '📜', colorClass: 'bg-green-500' },
  bash: { icon: '📜', colorClass: 'bg-green-500' },
  zsh: { icon: '📜', colorClass: 'bg-green-500' },
  ps1: { icon: '📜', colorClass: 'bg-blue-300' },
  bat: { icon: '📜', colorClass: 'bg-gray-500' },
  cmd: { icon: '📜', colorClass: 'bg-gray-500' },
  // Mobile
  dart: { icon: '🎯', colorClass: 'bg-blue-400' },
  swift: { icon: '🍎', colorClass: 'bg-orange-500' },
  // Other
  graphql: { icon: '◈', colorClass: 'bg-pink-500' },
  gql: { icon: '◈', colorClass: 'bg-pink-500' },
  proto: { icon: '📡', colorClass: 'bg-gray-500' },
  lock: { icon: '🔒', colorClass: 'bg-gray-600' },
  default: { icon: '📄', colorClass: 'bg-gray-500' }
}

// Special filename mappings (exact match)
const SPECIAL_FILES: Record<string, FileTypeInfo> = {
  'pom.xml': { icon: '🏺', colorClass: 'bg-red-600' },
  'build.gradle': { icon: '🐘', colorClass: 'bg-green-600' },
  'build.gradle.kts': { icon: '🐘', colorClass: 'bg-green-600' },
  'settings.gradle': { icon: '🐘', colorClass: 'bg-green-600' },
  'settings.gradle.kts': { icon: '🐘', colorClass: 'bg-green-600' },
  'package.json': { icon: '📦', colorClass: 'bg-red-500' },
  'package-lock.json': { icon: '🔒', colorClass: 'bg-gray-500' },
  'yarn.lock': { icon: '🔒', colorClass: 'bg-blue-400' },
  'pnpm-lock.yaml': { icon: '🔒', colorClass: 'bg-orange-400' },
  'tsconfig.json': { icon: '🔷', colorClass: 'bg-blue-500' },
  'jsconfig.json': { icon: '🟨', colorClass: 'bg-yellow-500' },
  '.gitignore': { icon: '🚫', colorClass: 'bg-gray-500' },
  '.gitattributes': { icon: '🔧', colorClass: 'bg-gray-500' },
  '.env': { icon: '🔐', colorClass: 'bg-yellow-600' },
  '.env.local': { icon: '🔐', colorClass: 'bg-yellow-600' },
  '.env.development': { icon: '🔐', colorClass: 'bg-yellow-600' },
  '.env.production': { icon: '🔐', colorClass: 'bg-yellow-600' },
  'Dockerfile': { icon: '🐳', colorClass: 'bg-blue-400' },
  'docker-compose.yml': { icon: '🐳', colorClass: 'bg-blue-400' },
  'docker-compose.yaml': { icon: '🐳', colorClass: 'bg-blue-400' },
  'Makefile': { icon: '🔨', colorClass: 'bg-gray-600' },
  'CMakeLists.txt': { icon: '🔨', colorClass: 'bg-blue-500' },
  'go.mod': { icon: '🐹', colorClass: 'bg-sky-500' },
  'go.sum': { icon: '🔒', colorClass: 'bg-sky-400' },
  'Cargo.toml': { icon: '🦀', colorClass: 'bg-orange-500' },
  'Cargo.lock': { icon: '🔒', colorClass: 'bg-orange-400' },
  'requirements.txt': { icon: '🐍', colorClass: 'bg-yellow-600' },
  'Pipfile': { icon: '🐍', colorClass: 'bg-yellow-600' },
  'pyproject.toml': { icon: '🐍', colorClass: 'bg-yellow-600' },
  'Gemfile': { icon: '💎', colorClass: 'bg-red-500' },
  'Gemfile.lock': { icon: '🔒', colorClass: 'bg-red-400' },
  'composer.json': { icon: '🐘', colorClass: 'bg-indigo-400' },
  'README.md': { icon: '📖', colorClass: 'bg-blue-400' },
  'LICENSE': { icon: '📜', colorClass: 'bg-gray-500' },
  'LICENSE.md': { icon: '📜', colorClass: 'bg-gray-500' },
  '.prettierrc': { icon: '✨', colorClass: 'bg-pink-400' },
  '.eslintrc': { icon: '🔍', colorClass: 'bg-purple-500' },
  '.eslintrc.js': { icon: '🔍', colorClass: 'bg-purple-500' },
  '.eslintrc.json': { icon: '🔍', colorClass: 'bg-purple-500' },
  'vite.config.ts': { icon: '⚡', colorClass: 'bg-purple-500' },
  'vite.config.js': { icon: '⚡', colorClass: 'bg-purple-500' },
  'webpack.config.js': { icon: '📦', colorClass: 'bg-blue-500' },
  'rollup.config.js': { icon: '📦', colorClass: 'bg-red-500' },
  'tailwind.config.js': { icon: '🎨', colorClass: 'bg-cyan-500' },
  'tailwind.config.ts': { icon: '🎨', colorClass: 'bg-cyan-500' },
  'postcss.config.js': { icon: '🎨', colorClass: 'bg-red-500' },
  'Application.java': { icon: '🚀', colorClass: 'bg-green-500' },
}

export function getFileIcon(name: string): string {
  // Check special files first (exact match)
  if (SPECIAL_FILES[name]) {
    return SPECIAL_FILES[name].icon
  }
  // Then check by extension
  const ext = name.split('.').pop()?.toLowerCase() || ''
  return FILE_TYPE_CONFIG[ext]?.icon || FILE_TYPE_CONFIG.default.icon
}

export function getFileTypeInfo(name: string): FileTypeInfo {
  // Check special files first (exact match)
  if (SPECIAL_FILES[name]) {
    return SPECIAL_FILES[name]
  }
  // Then check by extension
  const ext = name.split('.').pop()?.toLowerCase() || ''
  return FILE_TYPE_CONFIG[ext] || FILE_TYPE_CONFIG.default
}