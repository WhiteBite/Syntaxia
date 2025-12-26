/**
 * Type annotation defs shared between Flow and TypeScript.
 * These defs could not be defined in ./flow.ts or ./typescript.ts directly
 * because they use the same name.
 */
import type { Fork } from "../types";
export default function (fork: Fork): void;
