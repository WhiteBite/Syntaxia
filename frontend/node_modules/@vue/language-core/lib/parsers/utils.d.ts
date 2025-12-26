import type * as ts from 'typescript';
import type { TextRange } from '../types';
export declare function parseBindingRanges(ts: typeof import('typescript'), ast: ts.SourceFile, componentExtsensions: string[]): {
    bindings: TextRange[];
    components: TextRange[];
};
export declare function getClosestMultiLineCommentRange(ts: typeof import('typescript'), node: ts.Node, parents: ts.Node[], ast: ts.SourceFile): {
    start: number;
    end: number;
} | undefined;
