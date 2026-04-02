// Glyphs/colors derived from https://github.com/file-icons/vscode (MIT).

export type GitFileIconFont = 'fa' | 'fi' | 'mf' | 'octicons'

export interface GitFileIcon {
  font: GitFileIconFont
  glyph: string
  color: string
  scale?: number
}

function glyph(codePoint: number): string {
  return String.fromCodePoint(codePoint)
}

function defineIcon(font: GitFileIconFont, codePoint: number, color: string, scale = 1): GitFileIcon {
  return { font, glyph: glyph(codePoint), color, scale }
}

const iconCSS = defineIcon('fa', 0xf13c, '#be2f31', 0.93)
const iconDocker = defineIcon('fi', 0xf106, '#46788d', 1.29)
const iconGo = defineIcon('fi', 0xeaae, '#6a9fb5', 1.29)
const iconHTML = defineIcon('fa', 0xf13b, '#6a9fb5', 1.07)
const iconImageCool = defineIcon('octicons', 0xf012, '#90a959', 1.14)
const iconImageWarm = defineIcon('octicons', 0xf012, '#d28445', 1.14)
const iconJavaScript = defineIcon('mf', 0xf129, '#f4bf75')
const iconJSON = defineIcon('fi', 0xeabe, '#46788d', 1.07)
const iconMarkdown = defineIcon('octicons', 0xf0c9, '#aa759f', 1.14)
const iconPython = defineIcon('mf', 0xf14c, '#ff00cc')
const iconShell = defineIcon('fi', 0xf0c8, '#aa759f')
const iconSQL = defineIcon('mf', 0xf10e, '#d28445')
const iconSVG = defineIcon('mf', 0xf15c, '#ee9e2e')
const iconTOML = defineIcon('fi', 0x1f143, '#90a959', 1.07)
const iconTypeScript = defineIcon('fi', 0x02a6, '#6a9fb5')
const iconTSX = defineIcon('fi', 0xe9e7, '#9dc0ce')
const iconYAML = defineIcon('fi', 0x0079, '#ac4142', 1.07)
const iconZsh = defineIcon('fi', 0xf0c8, '#6a9fb5')

const FILE_NAME_ICON_MAP: Record<string, GitFileIcon> = {
  '.bash_history': iconShell,
  '.bash_profile': iconShell,
  '.bashrc': iconShell,
  '.dockerignore': iconDocker,
  '.zlogin': iconZsh,
  '.zlogout': iconZsh,
  '.zprofile': iconZsh,
  '.zshenv': iconZsh,
  '.zshrc': iconZsh,
  'containerfile': iconDocker,
  'dockerfile': iconDocker,
  'go.mod': iconGo,
  'go.sum': iconGo,
}

const FILE_EXTENSION_ICON_MAP: Record<string, GitFileIcon> = {
  apng: iconImageWarm,
  avif: iconImageWarm,
  bash: iconShell,
  bmp: iconImageWarm,
  cjs: iconJavaScript,
  css: iconCSS,
  cts: iconTypeScript,
  gif: iconImageWarm,
  go: iconGo,
  htm: iconHTML,
  html: iconHTML,
  ico: iconImageWarm,
  jpeg: iconImageCool,
  jpg: iconImageCool,
  js: iconJavaScript,
  json: iconJSON,
  jsonc: iconJSON,
  jsx: iconJavaScript,
  less: iconCSS,
  markdown: iconMarkdown,
  md: iconMarkdown,
  mdx: iconMarkdown,
  mjs: iconJavaScript,
  mts: iconTypeScript,
  png: iconImageWarm,
  py: iconPython,
  pyi: iconPython,
  sass: iconCSS,
  scss: iconCSS,
  sh: iconShell,
  sql: iconSQL,
  svg: iconSVG,
  toml: iconTOML,
  ts: iconTypeScript,
  tsx: iconTSX,
  webp: iconImageWarm,
  yaml: iconYAML,
  yml: iconYAML,
  zsh: iconZsh,
}

function fileNameFromPath(path: string): string {
  return path.trim().split(/[\\/]/).pop()?.toLowerCase() ?? ''
}

export function resolveGitDiffFileIcon(path: string): GitFileIcon | null {
  const fileName = fileNameFromPath(path)
  if (!fileName) return null

  const directMatch = FILE_NAME_ICON_MAP[fileName]
  if (directMatch) return directMatch

  const extIndex = fileName.lastIndexOf('.')
  if (extIndex === -1 || extIndex === fileName.length - 1) return null

  const extension = fileName.slice(extIndex + 1)
  return FILE_EXTENSION_ICON_MAP[extension] ?? null
}
