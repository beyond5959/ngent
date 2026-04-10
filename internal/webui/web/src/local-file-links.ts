export interface LocalFileLinkTarget {
  path: string
  line?: number
}

export type LocalFilePreviewHint = 'text' | 'image' | 'unsupported' | 'unknown'

const IMAGE_EXTENSIONS = new Set([
  'apng', 'avif', 'bmp', 'gif', 'heic', 'heif', 'ico', 'jpeg', 'jpg', 'png', 'svg', 'tif', 'tiff', 'webp',
])

const TEXT_EXTENSIONS = new Set([
  'bash', 'c', 'cc', 'cfg', 'conf', 'cpp', 'cs', 'css', 'csv', 'env', 'go', 'graphql', 'h', 'hpp', 'htm', 'html',
  'ini', 'java', 'js', 'json', 'jsonc', 'jsx', 'kt', 'less', 'log', 'lua', 'markdown', 'md', 'mdx', 'mjs', 'py',
  'rb', 'rs', 'sass', 'scss', 'sh', 'sql', 'swift', 'toml', 'ts', 'tsx', 'txt', 'vue', 'xml', 'yaml', 'yml', 'zsh',
])

const TEXT_FILE_NAMES = new Set([
  '.bash_history',
  '.bash_profile',
  '.bashrc',
  '.dockerignore',
  '.env',
  '.gitignore',
  '.npmrc',
  '.prettierrc',
  '.zlogin',
  '.zlogout',
  '.zprofile',
  '.zshenv',
  '.zshrc',
  'containerfile',
  'dockerfile',
  'go.mod',
  'go.sum',
  'license',
  'makefile',
  'readme',
])

const UNSUPPORTED_EXTENSIONS = new Set([
  '7z', 'aac', 'avi', 'bin', 'class', 'doc', 'docx', 'dylib', 'eot', 'exe', 'flac', 'gz', 'jar', 'm4a', 'mkv',
  'mov', 'mp3', 'mp4', 'o', 'obj', 'otf', 'pdf', 'ppt', 'pptx', 'rar', 'so', 'tar', 'tgz', 'ttf', 'wav', 'woff',
  'woff2', 'xls', 'xlsx', 'zip',
])

function isAbsoluteLocalPath(value: string): boolean {
  return value.startsWith('/') || /^[A-Za-z]:[\\/]/.test(value)
}

function fileNameFromPath(path: string): string {
  return path.trim().split(/[\\/]/).pop()?.toLowerCase() ?? ''
}

function fileExtension(path: string): string {
  const fileName = fileNameFromPath(path)
  const extIndex = fileName.lastIndexOf('.')
  if (extIndex <= 0 || extIndex === fileName.length - 1) return ''
  return fileName.slice(extIndex + 1)
}

function parseLineFragment(fragment: string): number | undefined {
  const match = fragment.trim().match(/^L(\d+)(?:C\d+)?$/i)
  if (!match) return undefined
  const line = Number.parseInt(match[1], 10)
  return Number.isFinite(line) && line > 0 ? line : undefined
}

export function parseLocalFileLinkHref(href: string): LocalFileLinkTarget | null {
  let raw = href.trim()
  if (!raw || raw.startsWith('#')) return null
  if (!isAbsoluteLocalPath(raw) && /^[a-z][a-z0-9+.-]*:/i.test(raw)) return null
  if (!isAbsoluteLocalPath(raw)) return null

  let fragment = ''
  const hashIndex = raw.indexOf('#')
  if (hashIndex >= 0) {
    fragment = raw.slice(hashIndex + 1)
    raw = raw.slice(0, hashIndex)
  }

  try {
    raw = decodeURIComponent(raw)
  } catch {
    // Keep the original path if it is not valid percent-encoding.
  }
  if (!isAbsoluteLocalPath(raw)) return null

  return {
    path: raw,
    line: parseLineFragment(fragment),
  }
}

export function localFilePreviewHint(path: string): LocalFilePreviewHint {
  const normalizedPath = path.trim()
  if (!normalizedPath) return 'unknown'

  const fileName = fileNameFromPath(normalizedPath)
  const extension = fileExtension(normalizedPath)
  if (IMAGE_EXTENSIONS.has(extension)) return 'image'
  if (TEXT_FILE_NAMES.has(fileName) || TEXT_EXTENSIONS.has(extension)) return 'text'
  if (UNSUPPORTED_EXTENSIONS.has(extension)) return 'unsupported'
  if (!extension && (fileName.startsWith('.') || /^[A-Z][A-Z0-9_-]*$/.test(fileName))) return 'text'
  return 'unknown'
}
