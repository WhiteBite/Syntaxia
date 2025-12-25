/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string;
  readonly DEV: boolean;
  readonly PROD: boolean;
  // more env variables...
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// Wails Go bindings
interface Window {
  go?: {
    main?: {
      App?: {
        ResetWindowState?: () => Promise<void>;
        [key: string]: unknown;
      };
      [key: string]: unknown;
    };
    [key: string]: unknown;
  };
}
