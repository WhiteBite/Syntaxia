import type { ASTNode, Fork } from "./types";
import type { Path } from "./path";
import type { Scope } from "./scope";
export interface NodePath<N = any, V = any> extends Path<V> {
    node: N;
    parent: any;
    scope: Scope | null;
    replace: Path['replace'];
    prune(...args: any[]): any;
    _computeNode(): any;
    _computeParent(): any;
    _computeScope(): Scope | null;
    getValueProperty(name: any): any;
    needsParens(assumeExpressionContext?: boolean): boolean;
    canBeFirstInStatement(): boolean;
    firstInStatement(): boolean;
}
export interface NodePathConstructor {
    new <N extends ASTNode = any, V = any>(value: any, parentPath?: any, name?: any): NodePath<N, V>;
}
export default function nodePathPlugin(fork: Fork): NodePathConstructor;
