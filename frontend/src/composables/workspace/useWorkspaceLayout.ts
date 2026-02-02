/**
 * Workspace Layout Management
 * 
 * Handles panel resizing and sidebar visibility for the main workspace.
 */

import { useResizablePanel } from '@/composables/useResizablePanel';
import { ref } from 'vue';

export function useWorkspaceLayout() {
    // Right sidebar visibility
    const showRightSidebar = ref(loadSidebarState());

    function loadSidebarState(): boolean {
        try {
            const saved = localStorage.getItem('right-sidebar-visible');
            return saved !== 'false';
        } catch {
            return true;
        }
    }

    function toggleRightSidebar() {
        showRightSidebar.value = !showRightSidebar.value;
        try {
            localStorage.setItem('right-sidebar-visible', String(showRightSidebar.value));
        } catch {
            // Ignore localStorage errors
        }
    }

    // Panel resizing
    const leftResize = useResizablePanel({
        minWidth: 280,
        maxWidth: 700,
        defaultWidth: 380,
        storageKey: 'workspace-left-width'
    });

    const rightResize = useResizablePanel({
        minWidth: 320,
        maxWidth: 700,
        defaultWidth: 380,
        storageKey: 'workspace-right-width',
        invertDirection: true // Тянем влево = увеличиваем ширину
    });

    function resetPanelSizes() {
        leftResize.resetToDefault();
        rightResize.resetToDefault();
    }

    return {
        // Sidebar
        showRightSidebar,
        toggleRightSidebar,

        // Left panel
        leftResize,
        leftWidth: leftResize.width,
        leftPanelRef: leftResize.panelRef,

        // Right panel
        rightResize,
        rightWidth: rightResize.width,
        rightPanelRef: rightResize.panelRef,

        // Actions
        resetPanelSizes
    };
}
