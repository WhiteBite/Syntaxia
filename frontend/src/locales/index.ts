// Locales index - exports all translations
import en from './en'
import ru from './ru'

export type Locale = 'ru' | 'en'

export const translations: Record<Locale, Record<string, string>> = {
    ru,
    en,
}

export { en, ru }

