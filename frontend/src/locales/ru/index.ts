// Russian locale - merges all JSON modules
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
import task from './task.json'
import templates from './templates.json'
import welcome from './welcome.json'

const ru: Record<string, string> = {
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
}

export default ru
