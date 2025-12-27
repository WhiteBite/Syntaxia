/**
 * Singleton i18n module for use outside Vue components (stores, services)
 * For Vue components, use useI18n() composable instead
 */
import { translations, type Locale } from './index'

let currentLocale: Locale = (typeof localStorage !== 'undefined'
    ? localStorage.getItem('app-locale') as Locale
    : null) || 'ru'

/**
 * Translate a key to current locale
 * @param key - Translation key (e.g., 'common.save')
 * @param params - Optional interpolation params
 */
export function t(key: string, params?: Record<string, string | number>): string {
    const dict = translations[currentLocale]
    let result = dict[key] || key

    if (params) {
        for (const [paramKey, paramValue] of Object.entries(params)) {
            result = result.replace(new RegExp(`\\{${paramKey}\\}`, 'g'), String(paramValue))
        }
    }

    return result
}

/**
 * Set current locale
 */
export function setLocale(locale: Locale): void {
    currentLocale = locale
    if (typeof localStorage !== 'undefined') {
        localStorage.setItem('app-locale', locale)
    }
}

/**
 * Get current locale
 */
export function getLocale(): Locale {
    return currentLocale
}

/**
 * Pluralization for Russian language
 */
export function plural(count: number, one: string, few: string, many: string): string {
    const absCount = Math.abs(count)
    const mod10 = absCount % 10
    const mod100 = absCount % 100

    if (mod10 === 1 && mod100 !== 11) {
        return one
    }
    if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
        return few
    }
    return many
}
