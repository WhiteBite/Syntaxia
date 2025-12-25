import { translations, type Locale } from '@/locales'
import { computed, ref } from 'vue'

export type { Locale }

const currentLocale = ref<Locale>((localStorage.getItem('app-locale') as Locale) || 'ru')

/**
 * Pluralization function for Russian language
 * @param count - number to pluralize
 * @param one - form for 1 (файл)
 * @param few - form for 2-4 (файла)
 * @param many - form for 5+ (файлов)
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

export function useI18n() {
    const t = (key: string, params?: Record<string, string | number>): string => {
        const locale = currentLocale.value
        const dict = translations[locale] as Record<string, string>
        let result = dict[key] || key

        // Support interpolation: {param} -> value
        if (params) {
            for (const [paramKey, paramValue] of Object.entries(params)) {
                result = result.replace(new RegExp(`\\{${paramKey}\\}`, 'g'), String(paramValue))
            }
        }

        return result
    }

    const setLocale = (locale: Locale) => {
        currentLocale.value = locale
        localStorage.setItem('app-locale', locale)
    }

    const locale = computed(() => currentLocale.value)

    return {
        t,
        locale,
        setLocale,
        plural,
    }
}
