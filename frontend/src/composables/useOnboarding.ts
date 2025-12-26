/**
 * Onboarding Tour Composable
 * Provides guided tour for new users using driver.js
 */

import { useI18n } from './useI18n'

const ONBOARDING_KEY = 'Syntaxia-onboarding-completed'

export function useOnboarding() {
    const { t } = useI18n()

    async function startTour() {
        // Dynamic import to avoid loading driver.js until needed
        const { driver } = await import('driver.js')
        await import('driver.js/dist/driver.css')
        // Custom dark theme styles
        await import('@/assets/styles/driver-theme.css')

        const steps = [
            {
                element: '[data-tour="file-tree"]',
                popover: {
                    title: t('onboarding.step1Title'),
                    description: t('onboarding.step1Desc'),
                    side: 'right' as const,
                },
            },
            {
                element: '[data-tour="build-button"]',
                popover: {
                    title: t('onboarding.step2Title'),
                    description: t('onboarding.step2Desc'),
                    side: 'top' as const,
                },
            },
            {
                element: '[data-tour="context-preview"]',
                popover: {
                    title: t('onboarding.step3Title'),
                    description: t('onboarding.step3Desc'),
                    side: 'left' as const,
                },
            },
            {
                element: '[data-tour="ai-chat"]',
                popover: {
                    title: t('onboarding.step4Title'),
                    description: t('onboarding.step4Desc'),
                    side: 'left' as const,
                },
            },
        ]

        const driverObj = driver({
            showProgress: true,
            steps,
            nextBtnText: t('onboarding.next'),
            prevBtnText: t('onboarding.prev'),
            doneBtnText: t('onboarding.done'),
            onDestroyStarted: () => {
                localStorage.setItem(ONBOARDING_KEY, 'true')
                driverObj.destroy()
            },
        })

        driverObj.drive()
    }

    function shouldShowTour(): boolean {
        return !localStorage.getItem(ONBOARDING_KEY)
    }

    function resetTour() {
        localStorage.removeItem(ONBOARDING_KEY)
    }

    function markTourCompleted() {
        localStorage.setItem(ONBOARDING_KEY, 'true')
    }

    return { startTour, shouldShowTour, resetTour, markTourCompleted }
}
