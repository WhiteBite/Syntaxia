/**
 * Error handling utility specifically designed for memory issues and large context handling
 */

import { useLogger } from '@/composables/useLogger';
import { useUIStore } from '@/stores/ui.store';

const logger = useLogger('MemoryErrorHandler');

export interface ErrorOptions {
  silent?: boolean;
  showNotification?: boolean;
  recoverable?: boolean;
  retryCallback?: () => void;
  context?: Record<string, unknown>;
}

export type ErrorCategory =
  'memory' |     // Memory-related errors (OutOfMemory, heap limit, etc.)
  'context' |    // Context building/processing errors
  'file' |       // File-related errors (too large, can't read, etc.)
  'network' |    // Network-related errors
  'api' |        // API errors
  'validation' | // Validation errors
  'unknown';     // Uncategorized errors

export interface ErrorDetails {
  message: string;
  category: ErrorCategory;
  recoverable: boolean;
  code?: string;
  suggestion?: string;
  originalError?: Error;
}

/**
 * Error handler for memory and context errors
 */
export class MemoryErrorHandler {
  private static readonly MEMORY_ERROR_PATTERNS = [
    'out of memory',
    'memory limit exceeded',
    'heap limit',
    'allocation failed',
    'allocation size overflow',
    'insufficient memory',
    'memory overflow',
    'call stack size exceeded',
    'maximum call stack size exceeded'
  ];

  private static readonly CONTEXT_SIZE_ERROR_PATTERNS = [
    'context too large',
    'exceeded maximum context',
    'context size limit',
    'maximum token limit',
    'token limit exceeded',
    'too many tokens',
    'context exceeds limit'
  ];

  /**
   * Categorize an error based on its message and stack trace
   */
  public static categorizeError(error: Error): ErrorCategory {
    const errorText = `${error.message} ${error.stack || ''}`.toLowerCase();

    // Check for memory errors
    if (this.MEMORY_ERROR_PATTERNS.some(pattern => errorText.includes(pattern))) {
      return 'memory';
    }

    // Check for context size errors
    if (this.CONTEXT_SIZE_ERROR_PATTERNS.some(pattern => errorText.includes(pattern))) {
      return 'context';
    }

    // Additional error categories could be checked here

    return 'unknown';
  }

  /**
   * Handle an error with appropriate measures based on type
   */
  public static handle(error: Error | string, options: ErrorOptions = {}): ErrorDetails {
    const uiStore = useUIStore();
    const errorObj = typeof error === 'string' ? new Error(error) : error;
    const category = this.categorizeError(errorObj);

    const details: ErrorDetails = {
      message: errorObj.message,
      category,
      recoverable: options.recoverable !== false,
      originalError: errorObj
    };

    // Handle based on category
    switch (category) {
      case 'memory':
        details.suggestion = 'Try selecting fewer files or freeing up system memory by closing other applications. The application will attempt to recover automatically.';
        details.code = 'ERR_MEMORY';

        if (options.showNotification !== false) {
          uiStore.addToast(
            `Memory error: ${errorObj.message}. ${details.suggestion}`,
            'error'
          );
        }

        // Force garbage collection if possible and attempt recovery
        this.attemptMemoryRecovery();
        break;

      case 'context':
        details.suggestion = 'Try reducing the number of files or selecting smaller files. The application has automatic chunked processing for large contexts.';
        details.code = 'ERR_CONTEXT_SIZE';

        if (options.showNotification !== false) {
          uiStore.addToast(
            `Context error: ${errorObj.message}. ${details.suggestion}`,
            'error'
          );
        }
        break;

      default:
        details.suggestion = 'Check the console for more details.';
        details.code = 'ERR_UNKNOWN';

        if (options.showNotification !== false && !options.silent) {
          uiStore.addToast(
            `An error occurred: ${errorObj.message}`,
            'error'
          );
        }
    }

    // Log to console if not silent
    if (!options.silent) {
      console.error(`[${details.code}] ${errorObj.message}`, errorObj, options.context || {});
    }

    return details;
  }

  /**
   * Attempt to recover from memory issues
   */
  private static attemptMemoryRecovery(): void {
    // Clear any large objects that might be cached
    // This is a placeholder for real memory recovery code
    logger.debug('Attempting memory recovery...');

    // Try to trigger garbage collection in different ways
    if (window.gc) {
      try {
        window.gc();
        logger.debug('Garbage collection triggered successfully');
      } catch (e) {
        console.warn('Failed to trigger garbage collection via window.gc', e);
      }
    }

    // Try to force garbage collection in Node.js environment (if applicable)
    // gc is typed in performance.d.ts as global var
    if (typeof gc !== 'undefined' && gc) {
      try {
        gc();
        logger.debug('Garbage collection triggered successfully via global.gc');
      } catch (e) {
        console.warn('Failed to trigger garbage collection via global.gc', e);
      }
    }

    // Add a small delay to allow garbage collection to work
    setTimeout(() => {
      logger.debug('Memory recovery attempt completed');
    }, 150);
  }

  /**
   * Install a global error handler for uncaught errors
   */
  public static installGlobalHandler(): void {
    window.addEventListener('error', (event) => {
      const error = event.error || new Error(event.message);
      const category = this.categorizeError(error);

      // Only handle memory and context errors globally
      if (category === 'memory' || category === 'context') {
        this.handle(error, { silent: true });

        // For memory errors, we can try to recover
        if (category === 'memory') {
          this.attemptMemoryRecovery();
        }
      }
    });

    // Also handle unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
      const error = event.reason instanceof Error ? event.reason : new Error(String(event.reason));
      const category = this.categorizeError(error);

      // Only handle memory and context errors globally
      if (category === 'memory' || category === 'context') {
        this.handle(error, { silent: true });
      }
    });
  }
}

// Extend window interface to potentially include gc method
declare global {
  interface Window {
    gc?: () => void;
  }
}

export default MemoryErrorHandler;