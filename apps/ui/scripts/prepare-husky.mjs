// apps/ui is nested inside the matecito-ai git repo (not its own git root), so husky's
// default `.git` lookup (relative to cwd) fails when run from `apps/ui`. This wrapper
// resolves the real git root first, then points husky's hooksPath at `apps/ui/.husky`.
import { execSync } from 'node:child_process'
import { existsSync } from 'node:fs'

const isCi = process.env.CI === 'true'
if (isCi) {
  process.exit(0)
}

const gitRoot = execSync('git rev-parse --show-toplevel', { cwd: import.meta.dirname })
  .toString()
  .trim()
if (!existsSync(`${gitRoot}/.git`)) {
  process.exit(0)
}

process.chdir(gitRoot)
const husky = (await import('husky')).default
process.stdout.write(husky('apps/ui/.husky'))
