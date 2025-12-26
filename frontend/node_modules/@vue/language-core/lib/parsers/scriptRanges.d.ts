import type * as ts from 'typescript';
import type { TextRange, VueCompilerOptions } from '../types';
export interface ScriptRanges extends ReturnType<typeof parseScriptRanges> {
}
export declare function parseScriptRanges(ts: typeof import('typescript'), ast: ts.SourceFile, vueCompilerOptions: VueCompilerOptions): {
    exportDefault: (TextRange & {
        expression: TextRange;
        isObjectLiteral: boolean;
    }) | undefined;
    componentOptions: {
        isObjectLiteral: boolean;
        expression: TextRange;
        args: TextRange;
        argsNode: ts.ObjectLiteralExpression;
        components: TextRange | undefined;
        componentsNode: ts.ObjectLiteralExpression | undefined;
        directives: TextRange | undefined;
        name: TextRange | undefined;
        inheritAttrs: string | undefined;
    } | undefined;
    bindings: TextRange[];
    components: TextRange[];
};
