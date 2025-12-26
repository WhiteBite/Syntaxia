import type { Code, VueCodeInformation } from '../../types';
export declare function wrapWith(source: string, startOffset: number, endOffset: number, features: VueCodeInformation, ...codes: Code[]): Generator<Code>;
