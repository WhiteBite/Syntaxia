import { plural, useI18n } from '@/composables/useI18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock localStorage
const localStorageMock = {
    store: {} as Record<string, string>,
    getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
        localStorageMock.store[key] = value
    }),
    clear: vi.fn(() => {
        localStorageMock.store = {}
    }),
}

Object.defineProperty(global, 'localStorage', { value: localStorageMock })

describe('useI18n', () => {
    beforeEach(() => {
        localStorageMock.clear()
        vi.clearAllMocks()
    })

    it('t() returns Russian translation when locale=ru', () => {
        localStorageMock.store['app-locale'] = 'ru'
        const { t, setLocale } = useI18n()
        setLocale('ru')

        expect(t('welcome.title')).toBe('Syntaxia')
        expect(t('files.title')).toBe('Файлы проекта')
    })

    it('t() returns English translation when locale=en', () => {
        const { t, setLocale } = useI18n()
        setLocale('en')

        expect(t('welcome.title')).toBe('Syntaxia')
        expect(t('files.title')).toBe('Project Files')
    })

    it('t() returns key if translation not found', () => {
        const { t } = useI18n()

        expect(t('nonexistent.key')).toBe('nonexistent.key')
        expect(t('another.missing.key')).toBe('another.missing.key')
    })

    it('setLocale changes language and saves to localStorage', () => {
        const { setLocale, locale } = useI18n()

        setLocale('en')
        expect(locale.value).toBe('en')
        expect(localStorageMock.setItem).toHaveBeenCalledWith('app-locale', 'en')

        setLocale('ru')
        expect(locale.value).toBe('ru')
        expect(localStorageMock.setItem).toHaveBeenCalledWith('app-locale', 'ru')
    })
})

describe('plural', () => {
    it('returns "файл" for 1', () => {
        expect(plural(1, 'файл', 'файла', 'файлов')).toBe('файл')
        expect(plural(21, 'файл', 'файла', 'файлов')).toBe('файл')
        expect(plural(101, 'файл', 'файла', 'файлов')).toBe('файл')
    })

    it('returns "файла" for 2, 3, 4', () => {
        expect(plural(2, 'файл', 'файла', 'файлов')).toBe('файла')
        expect(plural(3, 'файл', 'файла', 'файлов')).toBe('файла')
        expect(plural(4, 'файл', 'файла', 'файлов')).toBe('файла')
        expect(plural(22, 'файл', 'файла', 'файлов')).toBe('файла')
        expect(plural(34, 'файл', 'файла', 'файлов')).toBe('файла')
    })

    it('returns "файлов" for 5, 11, 100', () => {
        expect(plural(5, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(10, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(100, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(0, 'файл', 'файла', 'файлов')).toBe('файлов')
    })

    it('correctly handles 11-14 (special case)', () => {
        expect(plural(11, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(12, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(13, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(14, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(111, 'файл', 'файла', 'файлов')).toBe('файлов')
        expect(plural(112, 'файл', 'файла', 'файлов')).toBe('файлов')
    })
})
