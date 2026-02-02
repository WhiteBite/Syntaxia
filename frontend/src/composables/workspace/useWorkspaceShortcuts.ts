/**
 * Workspace Global Keyboard Shortcuts
 * 
 * Handles global keyboard shortcuts for workspace actions.
 */

import { onMounted, onUnmounted } from 'vue';

export function useWorkspaceShortcuts(handlers: {
    onBuildContext: () => void;
    onOpenExport: () => void;
    onCopyContext: () => void;
}) {
    const handleGlobalBuildContext = () => handlers.onBuildContext();
    const handleGlobalOpenExport = () => handlers.onOpenExport();
    const handleGlobalCopyContext = () => handlers.onCopyContext();

    onMounted(() => {
        window.addEventListener('global-build-context', handleGlobalBuildContext);
        window.addEventListener('global-open-export', handleGlobalOpenExport);
        window.addEventListener('global-copy-context', handleGlobalCopyContext);
    });

    onUnmounted(() => {
        window.removeEventListener('global-build-context', handleGlobalBuildContext);
        window.removeEventListener('global-open-export', handleGlobalOpenExport);
        window.removeEventListener('global-copy-context', handleGlobalCopyContext);
    });
}
