/**
 * Composable for dropdown management in QuickFilters
 * Includes keyboard navigation support (Escape, Arrow keys, Enter)
 */
import { ref, type Ref } from 'vue'
import type { DropdownType } from '../model/types'

export interface DropdownRefs {
    types: Ref<HTMLElement | null>
    langs: Ref<HTMLElement | null>
    smart: Ref<HTMLElement | null>
}

export function useFilterDropdown() {
    const openDropdown = ref<DropdownType | null>(null)
    const dropdownStyle = ref<Record<string, string>>({})
    const focusedIndex = ref(-1)

    const dropdownRefs: DropdownRefs = {
        types: ref<HTMLElement | null>(null),
        langs: ref<HTMLElement | null>(null),
        smart: ref<HTMLElement | null>(null),
    }

    function toggleDropdown(type: DropdownType): void {
        if (openDropdown.value === type) {
            openDropdown.value = null
            focusedIndex.value = -1
            return
        }

        const refElement = dropdownRefs[type].value
        if (refElement) {
            const rect = refElement.getBoundingClientRect()
            dropdownStyle.value = {
                position: 'fixed',
                top: `${rect.bottom + 4}px`,
                left: `${rect.left}px`,
                zIndex: '1100',
            }
        }
        openDropdown.value = type
        focusedIndex.value = -1
    }

    function closeDropdown(): void {
        openDropdown.value = null
        focusedIndex.value = -1
    }

    function handleOutsideClick(e: MouseEvent): void {
        const target = e.target as HTMLElement
        if (!target.closest('.filter-dropdown') && !target.closest('.filter-dropdown-menu')) {
            closeDropdown()
        }
    }

    function handleKeydown(e: KeyboardEvent): void {
        if (!openDropdown.value) return

        const menu = document.querySelector('.filter-dropdown-menu')
        if (!menu) return

        const items = menu.querySelectorAll<HTMLElement>('[role="option"], .filter-item')
        const itemCount = items.length

        switch (e.key) {
            case 'Escape':
                e.preventDefault()
                closeDropdown()
                // Return focus to trigger button
                const triggerRef = dropdownRefs[openDropdown.value]?.value
                triggerRef?.querySelector('button')?.focus()
                break

            case 'ArrowDown':
                e.preventDefault()
                focusedIndex.value = focusedIndex.value < itemCount - 1 ? focusedIndex.value + 1 : 0
                items[focusedIndex.value]?.focus()
                break

            case 'ArrowUp':
                e.preventDefault()
                focusedIndex.value = focusedIndex.value > 0 ? focusedIndex.value - 1 : itemCount - 1
                items[focusedIndex.value]?.focus()
                break

            case 'Home':
                e.preventDefault()
                focusedIndex.value = 0
                items[0]?.focus()
                break

            case 'End':
                e.preventDefault()
                focusedIndex.value = itemCount - 1
                items[itemCount - 1]?.focus()
                break

            case 'Enter':
            case ' ':
                if (focusedIndex.value >= 0 && items[focusedIndex.value]) {
                    e.preventDefault()
                    items[focusedIndex.value].click()
                }
                break

            case 'Tab':
                closeDropdown()
                break
        }
    }

    function setupListeners(): void {
        document.addEventListener('click', handleOutsideClick)
        document.addEventListener('keydown', handleKeydown)
    }

    function cleanupListeners(): void {
        document.removeEventListener('click', handleOutsideClick)
        document.removeEventListener('keydown', handleKeydown)
    }

    return {
        openDropdown,
        dropdownStyle,
        dropdownRefs,
        focusedIndex,
        toggleDropdown,
        closeDropdown,
        setupListeners,
        cleanupListeners,
    }
}
