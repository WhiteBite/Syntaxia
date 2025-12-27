// English locale - merges all JSON modules
import build from './build.json'
import chat from './chat.json'
import commands from './commands.json'
import common from './common.json'
import context from './context.json'
import errors from './errors.json'
import exportModule from './export.json'
import files from './files.json'
import git from './git.json'
import onboarding from './onboarding.json'
import settings from './settings.json'
import symbols from './symbols.json'
import task from './task.json'
import templates from './templates.json'
import testing from './testing.json'
import verification from './verification.json'
import welcome from './welcome.json'

const en: Record<string, string> = {
    ...common,
    ...files,
    ...context,
    ...chat,
    ...git,
    ...settings,
    ...errors,
    ...templates,
    ...welcome,
    ...onboarding,
    ...exportModule,
    ...task,
    ...commands,
    ...verification,
    ...build,
    ...symbols,
    ...testing,
}

export default en
