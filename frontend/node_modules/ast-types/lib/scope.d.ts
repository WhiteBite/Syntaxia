import type { namedTypes } from "./main";
import type { NodePath } from "./node-path";
import type { Fork } from "./types";
export type ScopeType = "global" | "function" | "block" | "switch" | "catch" | "with" | "for" | "class" | "type";
export type ScopeBinding = Record<string, NodePath<namedTypes.Identifier>[]>;
export type ScopeTypes = Record<string, NodePath[]>;
export interface Scope {
    type: ScopeType;
    path: NodePath;
    node: NodePath['value'];
    isGlobal: boolean;
    depth: number;
    parent: Scope | null;
    bindings: ScopeBinding;
    types: ScopeTypes;
    didScan: boolean;
    declares(name: string): boolean;
    declaresType(name: string): boolean;
    markAsStale(): void;
    scan(force?: boolean): void;
    getBindings(): ScopeBinding;
    getTypes(): ScopeTypes;
    lookup(name: string): Scope;
    lookupType(name: string): Scope;
    getGlobalScope(): Scope | null;
    scanScope: (path: NodePath, bindings: ScopeBinding, scopeTypes: ScopeTypes) => void;
    recursiveScanScope: (path: NodePath, bindings: ScopeBinding, scopeTypes: ScopeTypes) => void;
    recursiveScanChild: (path: NodePath, bindings: ScopeBinding, scopeTypes: ScopeTypes) => void;
}
export interface ScopeConstructor {
    new (path: NodePath, parentScope?: Scope | null): Scope;
    isEstablishedBy(node: NodePath['node']): boolean;
}
export default function scopePlugin(fork: Fork): ScopeConstructor;
