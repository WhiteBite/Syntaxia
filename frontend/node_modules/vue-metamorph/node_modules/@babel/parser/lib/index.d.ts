declare const UnparenthesizedPipeBodyDescriptions: Set<"ArrowFunctionExpression" | "AssignmentExpression" | "ConditionalExpression" | "YieldExpression">;
type GetSetMemberType<T extends Set<any>> = T extends Set<infer M> ? M : unknown;
type UnparenthesizedPipeBodyTypes = GetSetMemberType<typeof UnparenthesizedPipeBodyDescriptions>;

declare class Position {
    line: number;
    column: number;
    index: number;
    constructor(line: number, col: number, index: number);
}
declare class SourceLocation {
    start: Position;
    end: Position;
    filename: string;
    identifierName: string | undefined | null;
    constructor(start: Position, end?: Position);
}

type LValAncestor = {
    type: "UpdateExpression";
    prefix: boolean;
} | {
    type: "ArrayPattern" | "AssignmentExpression" | "CatchClause" | "ForOfStatement" | "FormalParameters" | "ForInStatement" | "ForStatement" | "ImportSpecifier" | "ImportNamespaceSpecifier" | "ImportDefaultSpecifier" | "ParenthesizedExpression" | "ObjectPattern" | "RestElement" | "VariableDeclarator";
};

type SyntaxPlugin = "flow" | "typescript" | "jsx" | "pipelineOperator" | "placeholders";
type ParseErrorCode = "BABEL_PARSER_SYNTAX_ERROR" | "BABEL_PARSER_SOURCETYPE_MODULE_REQUIRED";
interface ParseErrorSpecification<ErrorDetails> {
    code: ParseErrorCode;
    reasonCode: string;
    syntaxPlugin?: SyntaxPlugin;
    missingPlugin?: string | string[];
    loc: Position;
    details: ErrorDetails;
    pos: number;
}
type ParseError<ErrorDetails> = SyntaxError & ParseErrorSpecification<ErrorDetails>;
type ParseErrorConstructor<ErrorDetails> = (loc: Position, details: ErrorDetails) => ParseError<ErrorDetails>;
declare const Errors: {
    PipeBodyIsTighter: ParseErrorConstructor<object>;
    PipeTopicRequiresHackPipes: ParseErrorConstructor<object>;
    PipeTopicUnbound: ParseErrorConstructor<object>;
    PipeTopicUnconfiguredToken: ParseErrorConstructor<{
        token: string;
    }>;
    PipeTopicUnused: ParseErrorConstructor<object>;
    PipeUnparenthesizedBody: ParseErrorConstructor<{
        type: UnparenthesizedPipeBodyTypes;
    }>;
    PipelineBodyNoArrow: ParseErrorConstructor<object>;
    PipelineBodySequenceExpression: ParseErrorConstructor<object>;
    PipelineHeadSequenceExpression: ParseErrorConstructor<object>;
    PipelineTopicUnused: ParseErrorConstructor<object>;
    PrimaryTopicNotAllowed: ParseErrorConstructor<object>;
    PrimaryTopicRequiresSmartPipeline: ParseErrorConstructor<object>;
    StrictDelete: ParseErrorConstructor<object>;
    StrictEvalArguments: ParseErrorConstructor<{
        referenceName: string;
    }>;
    StrictEvalArgumentsBinding: ParseErrorConstructor<{
        bindingName: string;
    }>;
    StrictFunction: ParseErrorConstructor<object>;
    StrictNumericEscape: ParseErrorConstructor<object>;
    StrictOctalLiteral: ParseErrorConstructor<object>;
    StrictWith: ParseErrorConstructor<object>;
    AccessorIsGenerator: ParseErrorConstructor<{
        kind: "get" | "set";
    }>;
    ArgumentsInClass: ParseErrorConstructor<object>;
    AsyncFunctionInSingleStatementContext: ParseErrorConstructor<object>;
    AwaitBindingIdentifier: ParseErrorConstructor<object>;
    AwaitBindingIdentifierInStaticBlock: ParseErrorConstructor<object>;
    AwaitExpressionFormalParameter: ParseErrorConstructor<object>;
    AwaitUsingNotInAsyncContext: ParseErrorConstructor<object>;
    AwaitNotInAsyncContext: ParseErrorConstructor<object>;
    AwaitNotInAsyncFunction: ParseErrorConstructor<object>;
    BadGetterArity: ParseErrorConstructor<object>;
    BadSetterArity: ParseErrorConstructor<object>;
    BadSetterRestParameter: ParseErrorConstructor<object>;
    ConstructorClassField: ParseErrorConstructor<object>;
    ConstructorClassPrivateField: ParseErrorConstructor<object>;
    ConstructorIsAccessor: ParseErrorConstructor<object>;
    ConstructorIsAsync: ParseErrorConstructor<object>;
    ConstructorIsGenerator: ParseErrorConstructor<object>;
    DeclarationMissingInitializer: ParseErrorConstructor<{
        kind: "await using" | "const" | "destructuring" | "using";
    }>;
    DecoratorArgumentsOutsideParentheses: ParseErrorConstructor<object>;
    DecoratorBeforeExport: ParseErrorConstructor<object>;
    DecoratorsBeforeAfterExport: ParseErrorConstructor<object>;
    DecoratorConstructor: ParseErrorConstructor<object>;
    DecoratorExportClass: ParseErrorConstructor<object>;
    DecoratorSemicolon: ParseErrorConstructor<object>;
    DecoratorStaticBlock: ParseErrorConstructor<object>;
    DeferImportRequiresNamespace: ParseErrorConstructor<object>;
    DeletePrivateField: ParseErrorConstructor<object>;
    DestructureNamedImport: ParseErrorConstructor<object>;
    DuplicateConstructor: ParseErrorConstructor<object>;
    DuplicateDefaultExport: ParseErrorConstructor<object>;
    DuplicateExport: ParseErrorConstructor<{
        exportName: string;
    }>;
    DuplicateProto: ParseErrorConstructor<object>;
    DuplicateRegExpFlags: ParseErrorConstructor<object>;
    DynamicImportPhaseRequiresImportExpressions: ParseErrorConstructor<{
        phase: string;
    }>;
    ElementAfterRest: ParseErrorConstructor<object>;
    EscapedCharNotAnIdentifier: ParseErrorConstructor<object>;
    ExportBindingIsString: ParseErrorConstructor<{
        localName: string;
        exportName: string;
    }>;
    ExportDefaultFromAsIdentifier: ParseErrorConstructor<object>;
    ForInOfLoopInitializer: ParseErrorConstructor<{
        type: "ForInStatement" | "ForOfStatement";
    }>;
    ForInUsing: ParseErrorConstructor<object>;
    ForOfAsync: ParseErrorConstructor<object>;
    ForOfLet: ParseErrorConstructor<object>;
    GeneratorInSingleStatementContext: ParseErrorConstructor<object>;
    IllegalBreakContinue: ParseErrorConstructor<{
        type: "BreakStatement" | "ContinueStatement";
    }>;
    IllegalLanguageModeDirective: ParseErrorConstructor<object>;
    IllegalReturn: ParseErrorConstructor<object>;
    ImportAttributesUseAssert: ParseErrorConstructor<object>;
    ImportBindingIsString: ParseErrorConstructor<{
        importName: string;
    }>;
    ImportCallArgumentTrailingComma: ParseErrorConstructor<object>;
    ImportCallArity: ParseErrorConstructor<{
        maxArgumentCount: 1 | 2;
    }>;
    ImportCallNotNewExpression: ParseErrorConstructor<object>;
    ImportCallSpreadArgument: ParseErrorConstructor<object>;
    ImportJSONBindingNotDefault: ParseErrorConstructor<object>;
    ImportReflectionHasAssertion: ParseErrorConstructor<object>;
    ImportReflectionNotBinding: ParseErrorConstructor<object>;
    IncompatibleRegExpUVFlags: ParseErrorConstructor<object>;
    InvalidBigIntLiteral: ParseErrorConstructor<object>;
    InvalidCodePoint: ParseErrorConstructor<object>;
    InvalidCoverInitializedName: ParseErrorConstructor<object>;
    InvalidDecimal: ParseErrorConstructor<object>;
    InvalidDigit: ParseErrorConstructor<{
        radix: number;
    }>;
    InvalidEscapeSequence: ParseErrorConstructor<object>;
    InvalidEscapeSequenceTemplate: ParseErrorConstructor<object>;
    InvalidEscapedReservedWord: ParseErrorConstructor<{
        reservedWord: string;
    }>;
    InvalidIdentifier: ParseErrorConstructor<{
        identifierName: string;
    }>;
    InvalidLhs: ParseErrorConstructor<{
        ancestor: LValAncestor;
    }>;
    InvalidLhsBinding: ParseErrorConstructor<{
        ancestor: LValAncestor;
    }>;
    InvalidLhsOptionalChaining: ParseErrorConstructor<{
        ancestor: LValAncestor;
    }>;
    InvalidNumber: ParseErrorConstructor<object>;
    InvalidOrMissingExponent: ParseErrorConstructor<object>;
    InvalidOrUnexpectedToken: ParseErrorConstructor<{
        unexpected: string;
    }>;
    InvalidParenthesizedAssignment: ParseErrorConstructor<object>;
    InvalidPrivateFieldResolution: ParseErrorConstructor<{
        identifierName: string;
    }>;
    InvalidPropertyBindingPattern: ParseErrorConstructor<object>;
    InvalidRecordProperty: ParseErrorConstructor<object>;
    InvalidRestAssignmentPattern: ParseErrorConstructor<object>;
    LabelRedeclaration: ParseErrorConstructor<{
        labelName: string;
    }>;
    LetInLexicalBinding: ParseErrorConstructor<object>;
    LineTerminatorBeforeArrow: ParseErrorConstructor<object>;
    MalformedRegExpFlags: ParseErrorConstructor<object>;
    MissingClassName: ParseErrorConstructor<object>;
    MissingEqInAssignment: ParseErrorConstructor<object>;
    MissingSemicolon: ParseErrorConstructor<object>;
    MissingPlugin: ParseErrorConstructor<{
        missingPlugin: [string];
    }>;
    MissingOneOfPlugins: ParseErrorConstructor<{
        missingPlugin: string[];
    }>;
    MissingUnicodeEscape: ParseErrorConstructor<object>;
    MixingCoalesceWithLogical: ParseErrorConstructor<object>;
    ModuleAttributeDifferentFromType: ParseErrorConstructor<object>;
    ModuleAttributeInvalidValue: ParseErrorConstructor<object>;
    ModuleAttributesWithDuplicateKeys: ParseErrorConstructor<{
        key: string;
    }>;
    ModuleExportNameHasLoneSurrogate: ParseErrorConstructor<{
        surrogateCharCode: number;
    }>;
    ModuleExportUndefined: ParseErrorConstructor<{
        localName: string;
    }>;
    MultipleDefaultsInSwitch: ParseErrorConstructor<object>;
    NewlineAfterThrow: ParseErrorConstructor<object>;
    NoCatchOrFinally: ParseErrorConstructor<object>;
    NumberIdentifier: ParseErrorConstructor<object>;
    NumericSeparatorInEscapeSequence: ParseErrorConstructor<object>;
    ObsoleteAwaitStar: ParseErrorConstructor<object>;
    OptionalChainingNoNew: ParseErrorConstructor<object>;
    OptionalChainingNoTemplate: ParseErrorConstructor<object>;
    OverrideOnConstructor: ParseErrorConstructor<object>;
    ParamDupe: ParseErrorConstructor<object>;
    PatternHasAccessor: ParseErrorConstructor<object>;
    PatternHasMethod: ParseErrorConstructor<object>;
    PrivateInExpectedIn: ParseErrorConstructor<{
        identifierName: string;
    }>;
    PrivateNameRedeclaration: ParseErrorConstructor<{
        identifierName: string;
    }>;
    RecordExpressionBarIncorrectEndSyntaxType: ParseErrorConstructor<object>;
    RecordExpressionBarIncorrectStartSyntaxType: ParseErrorConstructor<object>;
    RecordExpressionHashIncorrectStartSyntaxType: ParseErrorConstructor<object>;
    RecordNoProto: ParseErrorConstructor<object>;
    RestTrailingComma: ParseErrorConstructor<object>;
    SloppyFunction: ParseErrorConstructor<object>;
    SloppyFunctionAnnexB: ParseErrorConstructor<object>;
    SourcePhaseImportRequiresDefault: ParseErrorConstructor<object>;
    StaticPrototype: ParseErrorConstructor<object>;
    SuperNotAllowed: ParseErrorConstructor<object>;
    SuperPrivateField: ParseErrorConstructor<object>;
    TrailingDecorator: ParseErrorConstructor<object>;
    TupleExpressionBarIncorrectEndSyntaxType: ParseErrorConstructor<object>;
    TupleExpressionBarIncorrectStartSyntaxType: ParseErrorConstructor<object>;
    TupleExpressionHashIncorrectStartSyntaxType: ParseErrorConstructor<object>;
    UnexpectedArgumentPlaceholder: ParseErrorConstructor<object>;
    UnexpectedAwaitAfterPipelineBody: ParseErrorConstructor<object>;
    UnexpectedDigitAfterHash: ParseErrorConstructor<object>;
    UnexpectedImportExport: ParseErrorConstructor<object>;
    UnexpectedKeyword: ParseErrorConstructor<{
        keyword: string;
    }>;
    UnexpectedLeadingDecorator: ParseErrorConstructor<object>;
    UnexpectedLexicalDeclaration: ParseErrorConstructor<object>;
    UnexpectedNewTarget: ParseErrorConstructor<object>;
    UnexpectedNumericSeparator: ParseErrorConstructor<object>;
    UnexpectedPrivateField: ParseErrorConstructor<object>;
    UnexpectedReservedWord: ParseErrorConstructor<{
        reservedWord: string;
    }>;
    UnexpectedSuper: ParseErrorConstructor<object>;
    UnexpectedToken: ParseErrorConstructor<{
        expected?: string | null;
        unexpected?: string | null;
    }>;
    UnexpectedTokenUnaryExponentiation: ParseErrorConstructor<object>;
    UnexpectedUsingDeclaration: ParseErrorConstructor<object>;
    UnsupportedBind: ParseErrorConstructor<object>;
    UnsupportedDecoratorExport: ParseErrorConstructor<object>;
    UnsupportedDefaultExport: ParseErrorConstructor<object>;
    UnsupportedImport: ParseErrorConstructor<object>;
    UnsupportedMetaProperty: ParseErrorConstructor<{
        target: string;
        onlyValidPropertyName: string;
    }>;
    UnsupportedParameterDecorator: ParseErrorConstructor<object>;
    UnsupportedPropertyDecorator: ParseErrorConstructor<object>;
    UnsupportedSuper: ParseErrorConstructor<object>;
    UnterminatedComment: ParseErrorConstructor<object>;
    UnterminatedRegExp: ParseErrorConstructor<object>;
    UnterminatedString: ParseErrorConstructor<object>;
    UnterminatedTemplate: ParseErrorConstructor<object>;
    UsingDeclarationExport: ParseErrorConstructor<object>;
    UsingDeclarationHasBindingPattern: ParseErrorConstructor<object>;
    VarRedeclaration: ParseErrorConstructor<{
        identifierName: string;
    }>;
    YieldBindingIdentifier: ParseErrorConstructor<object>;
    YieldInParameter: ParseErrorConstructor<object>;
    ZeroDigitNumericSeparator: ParseErrorConstructor<object>;
    ImportMetaOutsideModule: ParseErrorConstructor<object>;
    ImportOutsideModule: ParseErrorConstructor<object>;
};

declare class TokContext {
    constructor(token: string, preserveSpace?: boolean);
    token: string;
    preserveSpace: boolean;
}

type TokenOptions = {
    keyword?: string;
    beforeExpr?: boolean;
    startsExpr?: boolean;
    rightAssociative?: boolean;
    isLoop?: boolean;
    isAssign?: boolean;
    prefix?: boolean;
    postfix?: boolean;
    binop?: number | null;
};
type TokenType = number;
declare class ExportedTokenType {
    label: string;
    keyword: string | undefined | null;
    beforeExpr: boolean;
    startsExpr: boolean;
    rightAssociative: boolean;
    isLoop: boolean;
    isAssign: boolean;
    prefix: boolean;
    postfix: boolean;
    binop: number | undefined | null;
    updateContext: ((context: Array<TokContext>) => void) | undefined | null;
    constructor(label: string, conf?: TokenOptions);
}

type DeferredStrictError = typeof Errors.StrictNumericEscape | typeof Errors.StrictOctalLiteral;
type TopicContextState = {
    maxNumOfResolvableTopics: number;
    maxTopicIndex: null | 0;
};
declare const enum LoopLabelKind {
    Loop = 1,
    Switch = 2
}
declare class State {
    flags: number;
    accessor strict: boolean;
    curLine: number;
    lineStart: number;
    startLoc: Position;
    endLoc: Position;
    init({ strictMode, sourceType, startLine, startColumn }: Options): void;
    errors: ParseError<any>[];
    potentialArrowAt: number;
    noArrowAt: number[];
    noArrowParamsConversionAt: number[];
    accessor maybeInArrowParameters: boolean;
    accessor inType: boolean;
    accessor noAnonFunctionType: boolean;
    accessor hasFlowComment: boolean;
    accessor isAmbientContext: boolean;
    accessor inAbstractClass: boolean;
    accessor inDisallowConditionalTypesContext: boolean;
    topicContext: TopicContextState;
    accessor soloAwait: boolean;
    accessor inFSharpPipelineDirectBody: boolean;
    labels: Array<{
        kind: LoopLabelKind;
        name?: string | null;
        statementStart?: number;
    }>;
    commentsLen: number;
    commentStack: Array<CommentWhitespace>;
    pos: number;
    type: TokenType;
    value: any;
    start: number;
    end: number;
    lastTokEndLoc: Position;
    lastTokStartLoc: Position;
    context: Array<TokContext>;
    accessor canStartJSXElement: boolean;
    accessor containsEsc: boolean;
    firstInvalidTemplateEscapePos: null | Position;
    accessor hasTopLevelAwait: boolean;
    strictErrors: Map<number, [DeferredStrictError, Position]>;
    tokensLength: number;
    /**
     * When we add a new property, we must manually update the `clone` method
     * @see State#clone
     */
    curPosition(): Position;
    clone(): State;
}

type Plugin =
  | "asyncDoExpressions"
  | "asyncGenerators"
  | "bigInt"
  | "classPrivateMethods"
  | "classPrivateProperties"
  | "classProperties"
  | "classStaticBlock" // Enabled by default
  | "decimal"
  | "decorators-legacy"
  | "deferredImportEvaluation"
  | "decoratorAutoAccessors"
  | "destructuringPrivate"
  | "doExpressions"
  | "dynamicImport"
  | "explicitResourceManagement"
  | "exportDefaultFrom"
  | "exportNamespaceFrom" // deprecated
  | "flow"
  | "flowComments"
  | "functionBind"
  | "functionSent"
  | "importMeta"
  | "jsx"
  | "logicalAssignment"
  | "importAssertions" // deprecated
  | "importAttributes"
  | "importReflection"
  | "moduleBlocks"
  | "moduleStringNames"
  | "nullishCoalescingOperator"
  | "numericSeparator"
  | "objectRestSpread"
  | "optionalCatchBinding"
  | "optionalChaining"
  | "partialApplication"
  | "placeholders"
  | "privateIn" // Enabled by default
  | "regexpUnicodeSets" // Enabled by default
  | "sourcePhaseImports"
  | "throwExpressions"
  | "topLevelAwait"
  | "v8intrinsic"
  | ParserPluginWithOptions[0];

type ParserPluginWithOptions =
  | ["decorators", DecoratorsPluginOptions]
  | ["estree", { classFeatures?: boolean }]
  | ["importAttributes", { deprecatedAssertSyntax: boolean }]
  // @deprecated
  | ["moduleAttributes", { version: "may-2020" }]
  | ["optionalChainingAssign", { version: "2023-07" }]
  | ["pipelineOperator", PipelineOperatorPluginOptions]
  | ["recordAndTuple", RecordAndTuplePluginOptions]
  | ["flow", FlowPluginOptions]
  | ["typescript", TypeScriptPluginOptions];

type PluginConfig = Plugin | ParserPluginWithOptions;

interface DecoratorsPluginOptions {
  decoratorsBeforeExport?: boolean;
  allowCallParenthesized?: boolean;
}

interface PipelineOperatorPluginOptions {
  proposal: "minimal" | "fsharp" | "hack" | "smart";
  topicToken?: "%" | "#" | "@@" | "^^" | "^";
}

interface RecordAndTuplePluginOptions {
  syntaxType: "bar" | "hash";
}

interface FlowPluginOptions {
  all?: boolean;
  enums?: boolean;
}

interface TypeScriptPluginOptions {
  dts?: boolean;
  disallowAmbiguousJSXLike?: boolean;
}

/**
 * A whitespace token containing comments
 */
type CommentWhitespace = {
    /**
     * the start of the whitespace token.
     */
    start: number;
    /**
     * the end of the whitespace token.
     */
    end: number;
    /**
     * the containing comments
     */
    comments: Array<Comment>;
    /**
     * the immediately preceding AST node of the whitespace token
     */
    leadingNode: Node | null;
    /**
     * the immediately following AST node of the whitespace token
     */
    trailingNode: Node | null;
    /**
     * the innermost AST node containing the whitespace with minimal size (|end - start|)
     */
    containingNode: Node | null;
};

declare class Token {
    constructor(state: State);
    type: TokenType;
    value: any;
    start: number;
    end: number;
    loc: SourceLocation;
}

interface CommentBase {
    type: "CommentBlock" | "CommentLine";
    value: string;
    start: number;
    end: number;
    loc: SourceLocation;
}
interface CommentBlock extends CommentBase {
    type: "CommentBlock";
}
interface CommentLine extends CommentBase {
    type: "CommentLine";
}
type Comment = CommentBlock | CommentLine;
interface NodeBase {
    start: number;
    end: number;
    loc: SourceLocation;
    range?: [number, number];
    leadingComments?: Array<Comment>;
    trailingComments?: Array<Comment>;
    innerComments?: Array<Comment>;
    extra?: {
        [key: string]: any;
    };
}
type NodeAny<T extends string, KnownProps = object> = NodeBase & {
    type: T;
    [key: string]: any;
} & KnownProps;
type Expression = ArrayExpression | AssignmentExpression | BinaryExpression | CallExpression | ConditionalExpression | FunctionExpression | Identifier | StringLiteral | NumericLiteral | NullLiteral | BooleanLiteral | RegExpLiteral | LogicalExpression | MemberExpression | NewExpression | ObjectExpression | SequenceExpression | ParenthesizedExpression | ThisExpression | UnaryExpression | UpdateExpression | ArrowFunctionExpression | ClassExpression | ImportExpression | MetaProperty | Super | TaggedTemplateExpression | TemplateLiteral | YieldExpression | AwaitExpression | Import | BigIntLiteral | OptionalMemberExpression | OptionalCallExpression | TypeCastExpression | JSXElement | JSXFragment | BindExpression | DoExpression | RecordExpression | TupleExpression | DecimalLiteral | ModuleExpression | TopicReference | PipelineTopicExpression | PipelineBareFunction | PipelinePrimaryTopicReference | TsInstantiationExpression | TsAsExpression | TsSatisfiesExpression | TsTypeAssertion | TsTypeCastExpression | TsNonNullExpression | EstreeChainExpression | EstreeLiteral;
type Statement = BlockStatement | BreakStatement | ContinueStatement | DebuggerStatement | DoWhileStatement | EmptyStatement | ExpressionStatement | ForInStatement | ForStatement | FunctionDeclaration | IfStatement | LabeledStatement | ReturnStatement | SwitchStatement | ThrowStatement | TryStatement | VariableDeclaration | WhileStatement | WithStatement | ClassDeclaration | ExportAllDeclaration | ExportDefaultDeclaration | ExportNamedDeclaration | ForOfStatement | ImportDeclaration | FlowDeclareClass | FlowDeclareFunction | FlowDeclareInterface | FlowDeclareModule | FlowDeclareModuleExports | FlowDeclareTypeAlias | FlowDeclareOpaqueType | FlowDeclareVariable | FlowDeclareExportDeclaration | FlowEnumDeclaration | FlowInterface | FlowOpaqueType | FlowTypeAlias | TSDeclareFunction | TsInterfaceDeclaration | TsTypeAliasDeclaration | TsEnumDeclaration | TsModuleDeclaration | TsImportEqualsDeclaration | TsExportAssignment | TsNamespaceExportDeclaration;
type Pattern = Identifier | ObjectPattern | ArrayPattern | RestElement | AssignmentPattern;
type Declaration = VariableDeclaration | ClassDeclaration | FunctionDeclaration | TsInterfaceDeclaration | TsTypeAliasDeclaration | TsEnumDeclaration | TsModuleDeclaration;
interface DeclarationBase extends NodeBase {
    declare?: true;
}
interface HasDecorators extends NodeBase {
    decorators?: Decorator[];
}
interface InterpreterDirective extends NodeBase {
    type: "InterpreterDirective";
    value: string;
}
interface Identifier extends PatternBase {
    type: "Identifier";
    name: string;
    __clone(): Identifier;
    optional?: true | null;
}
interface PrivateName extends NodeBase {
    type: "PrivateName";
    id: Identifier;
}
type Literal = RegExpLiteral | NullLiteral | StringLiteral | BooleanLiteral | NumericLiteral | BigIntLiteral | DecimalLiteral;
type RegExpFlag = "g" | "i" | "m" | "u" | "s" | "y" | "v";
interface RegExpLiteral extends NodeBase {
    type: "RegExpLiteral";
    pattern: string;
    flags: RegExpFlag[];
}
interface NullLiteral extends NodeBase {
    type: "NullLiteral";
}
interface StringLiteral extends NodeBase {
    type: "StringLiteral";
    value: string;
}
interface BooleanLiteral extends NodeBase {
    type: "BooleanLiteral";
    value: boolean;
}
interface NumericLiteral extends NodeBase {
    type: "NumericLiteral";
    value: number;
}
interface BigIntLiteral extends NodeBase {
    type: "BigIntLiteral";
    value: number;
}
interface DecimalLiteral extends NodeBase {
    type: "DecimalLiteral";
    value: number;
}
interface ParserOutput {
    comments: Comment[];
    errors: Array<ParseError<any>>;
    tokens?: Array<Token | Comment>;
}
interface File extends NodeBase, ParserOutput {
    type: "File";
    program: Program;
}
interface Program extends NodeBase {
    type: "Program";
    sourceType: SourceType;
    body: Array<Statement | ModuleDeclaration>;
    directives: Directive[];
    interpreter: InterpreterDirective | null;
}
interface BodilessFunctionOrMethodBase extends HasDecorators {
    id: Identifier | undefined | null;
    params: Array<Pattern | TSParameterProperty>;
    generator: boolean;
    async: boolean;
    expression: boolean;
    typeParameters?: TypeParameterDeclarationBase | null;
    returnType?: TypeAnnotationBase | null;
}
interface FunctionBase extends BodilessFunctionOrMethodBase {
    body: BlockStatement;
}
interface ExpressionStatement extends NodeBase {
    type: "ExpressionStatement";
    expression: Expression;
}
interface BlockStatement extends NodeBase {
    type: "BlockStatement";
    body: Array<Statement>;
    directives: Directive[];
}
interface EmptyStatement extends NodeBase {
    type: "EmptyStatement";
}
interface DebuggerStatement extends NodeBase {
    type: "DebuggerStatement";
}
interface WithStatement extends NodeBase {
    type: "WithStatement";
    object: Expression;
    body: Statement;
}
interface ReturnStatement extends NodeBase {
    type: "ReturnStatement";
    argument: Expression | undefined | null;
}
interface LabeledStatement extends NodeBase {
    type: "LabeledStatement";
    label: Identifier;
    body: Statement;
}
interface BreakStatement extends NodeBase {
    type: "BreakStatement";
    label: Identifier | undefined | null;
}
interface ContinueStatement extends NodeBase {
    type: "ContinueStatement";
    label: Identifier | undefined | null;
}
interface IfStatement extends NodeBase {
    type: "IfStatement";
    test: Expression;
    consequent: Statement;
    alternate: Statement | undefined | null;
}
interface SwitchStatement extends NodeBase {
    type: "SwitchStatement";
    discriminant: Expression;
    cases: SwitchCase[];
}
interface SwitchCase extends NodeBase {
    type: "SwitchCase";
    test: Expression | undefined | null;
    consequent: Statement[];
}
interface ThrowStatement extends NodeBase {
    type: "ThrowStatement";
    argument: Expression;
}
interface TryStatement extends NodeBase {
    type: "TryStatement";
    block: BlockStatement;
    handler: CatchClause | null;
    finalizer: BlockStatement | null;
}
interface CatchClause extends NodeBase {
    type: "CatchClause";
    param: Pattern;
    body: BlockStatement;
}
interface WhileStatement extends NodeBase {
    type: "WhileStatement";
    test: Expression;
    body: Statement;
}
interface DoWhileStatement extends NodeBase {
    type: "DoWhileStatement";
    body: Statement;
    test: Expression;
}
interface ForStatement extends NodeBase {
    type: "ForStatement";
    init: VariableDeclaration | Expression | undefined | null;
    test: Expression | undefined | null;
    update: Expression | undefined | null;
    body: Statement;
}
interface ForInOfBase extends NodeBase {
    left: VariableDeclaration | Assignable;
    right: Expression;
    body: Statement;
}
interface ForInStatement extends ForInOfBase {
    type: "ForInStatement";
    await: boolean;
}
interface ForOfStatement extends ForInOfBase {
    type: "ForOfStatement";
    await: boolean;
}
interface OptFunctionDeclaration extends FunctionBase, DeclarationBase {
    type: "FunctionDeclaration";
}
interface FunctionDeclaration extends OptFunctionDeclaration {
    id: Identifier;
}
interface VariableDeclaration extends DeclarationBase, HasDecorators {
    type: "VariableDeclaration";
    declarations: VariableDeclarator[];
    kind: "var" | "let" | "const" | "using" | "await using";
}
interface VariableDeclarator extends NodeBase {
    type: "VariableDeclarator";
    id: Pattern;
    init: Expression | undefined | null;
    definite?: true;
}
interface ArgumentPlaceholder extends NodeBase {
    type: "ArgumentPlaceholder";
}
interface Decorator extends NodeBase {
    type: "Decorator";
    expression: Expression;
    arguments?: Array<Expression | SpreadElement>;
}
interface Directive extends NodeBase {
    type: "Directive";
    value: DirectiveLiteral;
}
interface DirectiveLiteral extends NodeBase {
    type: "DirectiveLiteral";
    value: string;
}
interface ImportAttribute extends NodeBase {
    type: "ImportAttribute";
    key: Identifier | StringLiteral;
    value: StringLiteral;
}
interface Super extends NodeBase {
    type: "Super";
}
interface Import extends NodeBase {
    type: "Import";
}
interface ThisExpression extends NodeBase {
    type: "ThisExpression";
}
interface ArrowFunctionExpression extends BodilessFunctionOrMethodBase {
    type: "ArrowFunctionExpression";
    body: BlockStatement | Expression;
}
interface YieldExpression extends NodeBase {
    type: "YieldExpression";
    argument: Expression | undefined | null;
    delegate: boolean;
}
interface AwaitExpression extends NodeBase {
    type: "AwaitExpression";
    argument: Expression;
}
interface ArrayExpression extends NodeBase {
    type: "ArrayExpression";
    elements: Array<Expression | SpreadElement | undefined | null>;
}
interface DoExpression extends NodeBase {
    type: "DoExpression";
    body: BlockStatement | undefined | null;
    async: boolean;
}
interface TupleExpression extends NodeBase {
    type: "TupleExpression";
    elements: Array<Expression | SpreadElement | undefined | null>;
}
interface ObjectExpression extends NodeBase {
    type: "ObjectExpression";
    properties: Array<ObjectProperty | ObjectMethod | SpreadElement>;
}
interface RecordExpression extends NodeBase {
    type: "RecordExpression";
    properties: Array<ObjectProperty | ObjectMethod | SpreadElement>;
}
interface ObjectMemberBase extends NodeBase {
    computed: boolean;
    value: Expression | Pattern;
    decorators?: Decorator[];
    kind?: "get" | "set" | "method";
    method: boolean;
    typeParameters?: TypeParameterDeclarationBase | null;
    variance?: FlowVariance | null;
}
interface ObjectProperty extends ObjectMemberBase {
    type: "ObjectProperty";
    shorthand: boolean;
    key: Expression | PrivateName;
    value: Expression | Pattern;
}
interface ObjectMethod extends ObjectMemberBase, FunctionBase {
    type: "ObjectMethod";
    kind: "get" | "set" | "method";
    key: Expression;
    value: Expression;
}
interface FunctionExpression extends FunctionBase {
    kind?: void;
    type: "FunctionExpression";
}
interface UnaryExpression extends NodeBase {
    type: "UnaryExpression";
    operator: UnaryOperator;
    prefix: boolean;
    argument: Expression;
}
type UnaryOperator = "-" | "+" | "!" | "~" | "typeof" | "void" | "delete" | "throw";
interface UpdateExpression extends NodeBase {
    type: "UpdateExpression";
    operator: UpdateOperator;
    argument: Expression;
    prefix: boolean;
}
type UpdateOperator = "++" | "--";
interface BinaryExpression extends NodeBase {
    type: "BinaryExpression";
    operator: BinaryOperator;
    left: Expression | PrivateName;
    right: Expression;
}
type BinaryOperator = "==" | "!=" | "===" | "!==" | "<" | "<=" | ">" | ">=" | "<<" | ">>" | ">>>" | "+" | "-" | "*" | "/" | "%" | "|" | "^" | "&" | "in" | "instanceof";
type Assignable = Pattern | MemberExpression | ParenthesizedExpression | TsTypeCastExpression | TypeCastExpression;
interface AssignmentExpression extends NodeBase {
    type: "AssignmentExpression";
    operator: AssignmentOperator;
    left: Assignable;
    right: Expression;
}
type AssignmentOperator = "=" | "+=" | "-=" | "*=" | "/=" | "%=" | "<<=" | ">>=" | ">>>=" | "|=" | "^=" | "&=";
interface LogicalExpression extends NodeBase {
    type: "LogicalExpression";
    operator: LogicalOperator;
    left: Expression;
    right: Expression;
}
type LogicalOperator = "||" | "&&";
interface SpreadElement extends NodeBase {
    type: "SpreadElement";
    argument: Expression;
}
interface MemberExpression extends NodeBase {
    type: "MemberExpression";
    object: Expression | Super;
    property: Expression | PrivateName;
    computed: boolean;
}
interface OptionalMemberExpression extends NodeBase {
    type: "OptionalMemberExpression";
    object: Expression | Super;
    property: Expression | PrivateName;
    computed: boolean;
    optional: boolean;
}
interface OptionalCallExpression extends CallOrNewBase {
    type: "OptionalCallExpression";
    optional: boolean;
}
interface BindExpression extends NodeBase {
    type: "BindExpression";
    object: Expression | undefined | null;
    callee: Expression;
}
interface ConditionalExpression extends NodeBase {
    type: "ConditionalExpression";
    test: Expression;
    alternate: Expression;
    consequent: Expression;
}
interface CallOrNewBase extends NodeBase {
    callee: Expression | Super | Import;
    arguments: Array<Expression | SpreadElement>;
    typeArguments: TypeParameterInstantiationBase | undefined | null;
    typeParameters?: TypeParameterInstantiationBase | null;
}
interface CallExpression extends CallOrNewBase {
    type: "CallExpression";
}
interface NewExpression extends CallOrNewBase {
    type: "NewExpression";
    optional?: boolean;
}
interface ImportExpression extends NodeBase {
    type: "ImportExpression";
    source: Expression;
    phase?: null | "source" | "defer";
    options: Expression | null;
}
interface SequenceExpression extends NodeBase {
    type: "SequenceExpression";
    expressions: Expression[];
}
interface ParenthesizedExpression extends NodeBase {
    type: "ParenthesizedExpression";
    expression: Expression;
}
interface TopicReference extends NodeBase {
    type: "TopicReference";
}
interface PipelineBody extends NodeBase {
    type: "PipelineBody";
}
interface PipelineBareFunctionBody extends NodeBase {
    type: "PipelineBareFunctionBody";
    callee: Expression;
}
interface PipelineBareConstructorBody extends NodeBase {
    type: "PipelineBareConstructorBody";
    callee: Expression;
}
interface PipelineBareAwaitedFunctionBody extends NodeBase {
    type: "PipelineBareAwaitedFunctionBody";
    callee: Expression;
}
interface PipelineTopicBody extends NodeBase {
    type: "PipelineTopicBody";
    expression: Expression;
}
interface PipelinePrimaryTopicReference extends NodeBase {
    type: "PipelinePrimaryTopicReference";
}
interface TemplateLiteral extends NodeBase {
    type: "TemplateLiteral";
    quasis: TemplateElement[];
    expressions: Expression[] | TsType[];
}
interface TaggedTemplateExpression extends NodeBase {
    type: "TaggedTemplateExpression";
    tag: Expression;
    quasi: TemplateLiteral;
    typeParameters?: TypeParameterInstantiationBase | null;
}
interface TemplateElement extends NodeBase {
    type: "TemplateElement";
    tail: boolean;
    value: {
        cooked: string;
        raw: string;
    };
}
interface ModuleExpression extends NodeBase {
    type: "ModuleExpression";
    body: Program;
}
type Accessibility = "public" | "protected" | "private";
interface PatternBase extends HasDecorators {
    typeAnnotation?: TypeAnnotationBase | null;
}
interface AssignmentProperty extends ObjectProperty {
    value: Pattern;
}
interface ObjectPattern extends PatternBase {
    type: "ObjectPattern";
    properties: (AssignmentProperty | RestElement)[];
}
interface ArrayPattern extends PatternBase {
    type: "ArrayPattern";
    elements: (Pattern | undefined | null)[];
}
interface RestElement extends PatternBase {
    type: "RestElement";
    argument: Pattern;
}
interface AssignmentPattern extends PatternBase {
    type: "AssignmentPattern";
    left: Pattern;
    right: Expression;
}
interface ClassBase extends HasDecorators {
    id: Identifier | undefined | null;
    superClass: Expression | undefined | null;
    body: ClassBody;
    decorators?: Decorator[];
    typeParameters?: TypeParameterDeclarationBase | null;
    superTypeParameters?: TypeParameterInstantiationBase | null;
    abstract?: boolean;
    implements?: TsExpressionWithTypeArguments[] | undefined | null | FlowClassImplements[];
}
interface ClassBody extends NodeBase {
    type: "ClassBody";
    body: Array<ClassMember | StaticBlock | TsIndexSignature>;
}
interface ClassMemberBase extends NodeBase, HasDecorators {
    static: boolean;
    computed: boolean;
    accessibility?: Accessibility | null;
    override?: true | null;
    abstract?: true | null;
    optional?: true | null;
}
interface StaticBlock extends NodeBase {
    type: "StaticBlock";
    body: Array<Statement>;
}
type ClassMember = ClassMethod | ClassPrivateMethod | ClassProperty | ClassPrivateProperty | ClassAccessorProperty;
interface MethodBase extends FunctionBase {
    kind: MethodKind;
}
type MethodKind = "constructor" | "method" | "get" | "set";
interface ClassMethodOrDeclareMethodCommon extends ClassMemberBase {
    key: Expression | PrivateName;
    kind: MethodKind;
    static: boolean;
    decorators?: Decorator[];
}
interface ClassMethod extends MethodBase, ClassMethodOrDeclareMethodCommon {
    type: "ClassMethod";
    variance?: FlowVariance | null;
}
interface ClassPrivateMethod extends NodeBase, ClassMethodOrDeclareMethodCommon, MethodBase {
    type: "ClassPrivateMethod";
    key: PrivateName;
    computed: false;
    variance?: FlowVariance | null;
}
interface ClassProperty extends ClassMemberBase, DeclarationBase {
    type: "ClassProperty";
    key: Expression;
    value: Expression | undefined | null;
    typeAnnotation?: TypeAnnotationBase | null;
    variance?: FlowVariance | null;
    readonly?: true;
    definite?: true;
}
interface ClassPrivateProperty extends NodeBase {
    type: "ClassPrivateProperty";
    key: PrivateName;
    value: Expression | undefined | null;
    static: boolean;
    computed: false;
    typeAnnotation?: TypeAnnotationBase | null;
    optional?: true;
    definite?: true;
    readonly?: true;
    override?: true;
    variance?: FlowVariance | null;
}
interface ClassAccessorProperty extends ClassMemberBase, DeclarationBase {
    type: "ClassAccessorProperty";
    key: Expression | PrivateName;
    value: Expression | undefined | null;
    typeAnnotation?: TypeAnnotationBase | null;
    variance?: FlowVariance | null;
    readonly?: true;
    definite?: true;
}
interface OptClassDeclaration extends ClassBase, DeclarationBase, HasDecorators {
    type: "ClassDeclaration";
    abstract?: true | null;
}
interface ClassDeclaration extends OptClassDeclaration {
    id: Identifier;
}
interface ClassExpression extends ClassBase {
    type: "ClassExpression";
}
interface MetaProperty extends NodeBase {
    type: "MetaProperty";
    meta: Identifier;
    property: Identifier;
}
type ModuleDeclaration = AnyImport | AnyExport;
type AnyImport = ImportDeclaration | TsImportEqualsDeclaration;
type AnyExport = ExportNamedDeclaration | ExportDefaultDeclaration | ExportAllDeclaration | TsExportAssignment | TsImportEqualsDeclaration | TsNamespaceExportDeclaration;
interface ModuleSpecifier extends NodeBase {
    local: Identifier;
}
interface ImportDeclaration extends NodeBase {
    type: "ImportDeclaration";
    specifiers: Array<ImportSpecifier | ImportDefaultSpecifier | ImportNamespaceSpecifier>;
    source: Literal;
    importKind?: "type" | "typeof" | "value";
    phase?: null | "source" | "defer";
    attributes?: ImportAttribute[];
    assertions?: ImportAttribute[];
    module?: boolean;
}
interface ImportSpecifier extends ModuleSpecifier {
    type: "ImportSpecifier";
    imported: Identifier | StringLiteral;
    importKind?: "type" | "value";
}
interface ImportDefaultSpecifier extends ModuleSpecifier {
    type: "ImportDefaultSpecifier";
}
interface ImportNamespaceSpecifier extends ModuleSpecifier {
    type: "ImportNamespaceSpecifier";
}
interface ExportNamedDeclaration extends NodeBase {
    type: "ExportNamedDeclaration";
    declaration: Declaration | undefined | null;
    specifiers: Array<ExportSpecifier | ExportDefaultSpecifier | ExportNamespaceSpecifier>;
    source: Literal | undefined | null;
    exportKind?: "type" | "value";
    attributes?: ImportAttribute[];
    assertions?: ImportAttribute[];
}
interface ExportSpecifier extends NodeBase {
    type: "ExportSpecifier";
    exported: Identifier | StringLiteral;
    local: Identifier;
    exportKind?: "type" | "value";
}
interface ExportDefaultSpecifier extends NodeBase {
    type: "ExportDefaultSpecifier";
    exported: Identifier;
}
interface ExportNamespaceSpecifier extends NodeBase {
    type: "ExportNamespaceSpecifier";
    exported: Identifier | StringLiteral;
}
interface ExportDefaultDeclaration extends NodeBase {
    type: "ExportDefaultDeclaration";
    declaration: OptFunctionDeclaration | OptTSDeclareFunction | OptClassDeclaration | FlowEnumDeclaration | TsInterfaceDeclaration | Expression;
}
interface ExportAllDeclaration extends NodeBase {
    type: "ExportAllDeclaration";
    source: Literal;
    exportKind?: "type" | "value";
    assertions?: ImportAttribute[];
    attributes?: ImportAttribute[];
}
interface PipelineTopicExpression extends NodeBase {
    type: "PipelineTopicExpression";
    expression: Expression;
}
interface PipelineBareFunction extends NodeBase {
    type: "PipelineBareFunction";
    callee: Expression;
}
type JSXIdentifier = NodeAny<"JSXIdentifier">;
type JSXNamespacedName = NodeAny<"JSXNamespacedName">;
type JSXMemberExpression = NodeAny<"JSXMemberExpression">;
type JSXEmptyExpression = NodeAny<"JSXEmptyExpression">;
type JSXSpreadChild = NodeAny<"JSXSpreadChild">;
type JSXExpressionContainer = NodeAny<"JSXExpressionContainer">;
type JSXAttribute = NodeAny<"JSXAttribute">;
type JSXSpreadAttribute = NodeAny<"JSXSpreadAttribute">;
interface JSXOpeningElement extends NodeBase {
    type: "JSXOpeningElement";
    name: JSXNamespacedName | JSXMemberExpression;
    typeParameters?: TypeParameterInstantiationBase | null;
    attributes: (JSXAttribute | JSXSpreadAttribute)[];
    selfClosing: boolean;
}
type JSXClosingElement = NodeAny<"JSXClosingElement">;
type JSXElement = NodeAny<"JSXElement">;
type JSXOpeningFragment = NodeAny<"JSXOpeningFragment">;
type JSXClosingFragment = NodeAny<"JSXClosingFragment">;
type JSXFragment = NodeAny<"JSXFragment">;
type JSXText = NodeAny<"JSXText">;
interface TypeAnnotationBase extends NodeBase {
    typeAnnotation: Node;
}
interface TypeAnnotation extends NodeBase {
    type: "TypeAnnotation";
    typeAnnotation: FlowType;
}
interface TsTypeAnnotation extends NodeBase {
    type: "TSTypeAnnotation";
    typeAnnotation: TsType;
}
interface TypeParameterDeclarationBase extends NodeBase {
    params: Array<TypeParameter | TsTypeParameter>;
}
interface TypeParameterDeclaration extends TypeParameterDeclarationBase {
    type: "TypeParameterDeclaration";
    params: TypeParameter[];
}
interface TsTypeParameterDeclaration extends TypeParameterDeclarationBase {
    type: "TSTypeParameterDeclaration";
    params: TsTypeParameter[];
}
interface TypeParameter extends NodeBase {
    type: "TypeParameter";
    name: string;
    default?: TypeAnnotation;
}
interface TsTypeParameter extends NodeBase {
    type: "TSTypeParameter";
    name: string | Identifier;
    in?: boolean;
    out?: boolean;
    const?: boolean;
    constraint?: TsType;
    default?: TsType;
}
interface TypeParameterInstantiationBase extends NodeBase {
    params: Node[];
}
interface TypeParameterInstantiation extends TypeParameterInstantiationBase {
    type: "TypeParameterInstantiation";
    params: FlowType[];
}
interface TsTypeParameterInstantiation extends TypeParameterInstantiationBase {
    type: "TSTypeParameterInstantiation";
    params: TsType[];
}
interface TypeCastExpression extends NodeBase {
    type: "TypeCastExpression";
    expression: Expression;
    typeAnnotation: TypeAnnotation;
}
interface TsTypeCastExpression extends NodeBase {
    type: "TSTypeCastExpression";
    expression: Expression;
    typeAnnotation: TsTypeAnnotation;
}
type FlowPredicate = NodeAny<"DeclaredPredicate"> | NodeAny<"InferredPredicate">;
type FlowDeclareClass = NodeAny<"DeclareClass">;
type FlowDeclareExportDeclaration = NodeAny<"DeclareExportDeclaration"> | NodeAny<"DeclareExportDefaultDeclaration"> | NodeAny<"DeclareExportAllDeclaration">;
type FlowDeclareFunction = NodeAny<"DeclareFunction">;
type FlowDeclareVariable = NodeAny<"DeclareVariable">;
type FlowDeclareModule = NodeAny<"DeclareModule">;
type FlowDeclareModuleExports = NodeAny<"DeclareModuleExports">;
type FlowDeclareTypeAlias = NodeAny<"DeclareTypeAlias">;
type FlowDeclareOpaqueType = NodeAny<"DeclareOpaqueType">;
type FlowDeclareInterface = NodeAny<"DeclareInterface">;
type FlowInterface = NodeAny<"InterfaceDeclaration">;
type FlowInterfaceExtends = NodeAny<"InterfaceExtends">;
type FlowTypeAlias = NodeAny<"TypeAlias">;
type FlowOpaqueType = NodeAny<"OpaqueType">;
type FlowObjectTypeIndexer = NodeAny<"ObjectTypeIndexer">;
type FlowObjectTypeInternalSlot = NodeAny<"ObjectTypeInternalSlot">;
type FlowEnumDeclaration = NodeAny<"EnumDeclaration">;
type FlowEnumBody = NodeAny<"EnumBooleanBody" | "EnumNumberBody" | "EnumStringBody" | "EnumSymbolBody">;
type FlowEnumMember = NodeAny<"EnumBooleanMember"> | NodeAny<"EnumNumberMember"> | NodeAny<"EnumStringMember"> | NodeAny<"EnumDefaultedMember">;
type FlowFunctionTypeAnnotation = NodeAny<"FunctionTypeAnnotation">;
type FlowObjectTypeProperty = NodeAny<"ObjectTypeProperty">;
type FlowObjectTypeSpreadProperty = NodeAny<"ObjectTypeSpreadProperty">;
type FlowObjectTypeCallProperty = NodeAny<"ObjectTypeCallProperty">;
type FlowObjectTypeAnnotation = NodeAny<"ObjectTypeAnnotation">;
type FlowQualifiedTypeIdentifier = NodeAny<"QualifiedTypeIdentifier">;
type FlowGenericTypeAnnotation = NodeAny<"GenericTypeAnnotation">;
type FlowTypeofTypeAnnotation = NodeAny<"TypeofTypeAnnotation">;
type FlowTupleTypeAnnotation = NodeAny<"TupleTypeAnnotation">;
type FlowFunctionTypeParam = NodeAny<"FunctionTypeParam">;
type FlowOtherTypeAnnotation = NodeAny<"AnyTypeAnnotation" | "BooleanTypeAnnotation" | "MixedTypeAnnotation" | "EmptyTypeAnnotation" | "ExistsTypeAnnotation" | "NumberTypeAnnotation" | "StringTypeAnnotation" | "SymbolTypeAnnotation" | "NullLiteralTypeAnnotation" | "VoidTypeAnnotation" | "ThisTypeAnnotation" | "ArrayTypeAnnotation" | "NullableTypeAnnotation" | "IntersectionTypeAnnotation" | "UnionTypeAnnotation">;
type FlowType = FlowFunctionTypeAnnotation | FlowObjectTypeAnnotation | FlowGenericTypeAnnotation | FlowTypeofTypeAnnotation | FlowTupleTypeAnnotation | FlowInterfaceType | FlowIndexedAccessType | FlowOptionalIndexedAccessType | FlowOtherTypeAnnotation | StringLiteralTypeAnnotation | BooleanLiteralTypeAnnotation | NumberLiteralTypeAnnotation | BigIntLiteralTypeAnnotation | Identifier;
type FlowVariance = NodeAny<"Variance">;
type FlowClassImplements = NodeAny<"ClassImplements">;
interface FlowInterfaceType extends NodeBase {
    type: "InterfaceTypeAnnotation";
    extends: FlowInterfaceExtends[];
    body: FlowObjectTypeAnnotation;
}
interface FlowIndexedAccessType extends NodeBase {
    type: "IndexedAccessType";
    objectType: FlowType;
    indexType: FlowType;
}
interface FlowOptionalIndexedAccessType extends NodeBase {
    type: "OptionalIndexedAccessType";
    objectType: FlowType;
    indexType: FlowType;
    optional: boolean;
}
interface StringLiteralTypeAnnotation extends NodeBase {
    type: "StringLiteralTypeAnnotation";
    value: string;
}
interface BooleanLiteralTypeAnnotation extends NodeBase {
    type: "BooleanLiteralTypeAnnotation";
    value: boolean;
}
interface NumberLiteralTypeAnnotation extends NodeBase {
    type: "NumberLiteralTypeAnnotation";
    value: number;
}
interface BigIntLiteralTypeAnnotation extends NodeBase {
    type: "BigIntLiteralTypeAnnotation";
    value: number;
}
interface EstreeLiteral extends NodeBase {
    type: "Literal";
    value: any;
    decimal?: string;
}
interface EstreeProperty extends NodeBase {
    type: "Property";
    method: boolean;
    shorthand: boolean;
    key: Expression | EstreePrivateIdentifier;
    computed: boolean;
    value: Expression;
    decorators: Decorator[];
    kind?: "get" | "set" | "init";
    variance?: FlowVariance | null;
}
interface EstreeMethodDefinition extends NodeBase {
    type: "MethodDefinition";
    static: boolean;
    key: Expression;
    computed: boolean;
    value: FunctionExpression;
    decorators: Decorator[];
    kind?: "get" | "set" | "method";
    variance?: FlowVariance | null;
}
interface EstreePrivateIdentifier extends NodeBase {
    type: "PrivateIdentifier";
    name: string;
}
interface EstreePropertyDefinition extends NodeBase {
    type: "PropertyDefinition";
    static: boolean;
    key: Expression | EstreePrivateIdentifier;
    computed: boolean;
    value: Expression;
}
interface EstreeChainExpression extends NodeBase {
    type: "ChainExpression";
    expression: Expression;
}
interface TSParameterProperty extends HasDecorators {
    type: "TSParameterProperty";
    accessibility?: Accessibility | null;
    readonly?: true | null;
    override?: true | null;
    parameter: Identifier | AssignmentPattern;
}
interface OptTSDeclareFunction extends FunctionBase, DeclarationBase {
    type: "TSDeclareFunction";
}
interface TSDeclareFunction extends OptTSDeclareFunction {
    id: Identifier;
}
interface TSDeclareMethod extends FunctionBase, ClassMethodOrDeclareMethodCommon {
    type: "TSDeclareMethod";
    kind: MethodKind;
}
interface TsQualifiedName extends NodeBase {
    type: "TSQualifiedName";
    left: TsEntityName;
    right: Identifier;
}
type TsEntityName = Identifier | TsQualifiedName;
interface TsSignatureDeclarationOrIndexSignatureBase extends NodeBase {
    params: Array<Identifier | RestElement | ObjectPattern | ArrayPattern>;
    returnType: TsTypeAnnotation | undefined | null;
    parameters: Array<Identifier | RestElement | ObjectPattern | ArrayPattern>;
    typeAnnotation: TsTypeAnnotation | undefined | null;
}
interface TsSignatureDeclarationBase extends TsSignatureDeclarationOrIndexSignatureBase {
    typeParameters: TsTypeParameterDeclaration | undefined | null;
}
type TsTypeElement = TsCallSignatureDeclaration | TsConstructSignatureDeclaration | TsPropertySignature | TsMethodSignature | TsIndexSignature;
interface TsCallSignatureDeclaration extends TsSignatureDeclarationBase {
    type: "TSCallSignatureDeclaration";
}
interface TsConstructSignatureDeclaration extends TsSignatureDeclarationBase {
    type: "TSConstructSignatureDeclaration";
}
interface TsNamedTypeElementBase extends NodeBase {
    key: Expression;
    computed: boolean;
    optional?: true;
}
interface TsPropertySignature extends TsNamedTypeElementBase {
    type: "TSPropertySignature";
    readonly?: true;
    typeAnnotation?: TsTypeAnnotation;
}
interface TsMethodSignature extends TsSignatureDeclarationBase, TsNamedTypeElementBase {
    type: "TSMethodSignature";
    kind: "method" | "get" | "set";
}
interface TsIndexSignature extends TsSignatureDeclarationOrIndexSignatureBase {
    readonly?: true;
    static?: true;
    type: "TSIndexSignature";
}
type TsType = TsKeywordType | TsThisType | TsFunctionOrConstructorType | TsTypeReference | TsTypeQuery | TsTypeLiteral | TsArrayType | TsTupleType | TsOptionalType | TsRestType | TsUnionOrIntersectionType | TsConditionalType | TsInferType | TsParenthesizedType | TsTypeOperator | TsIndexedAccessType | TsMappedType | TsLiteralType | TsImportType | TsTypePredicate;
type TsTypeBase = NodeBase;
type TsKeywordTypeType = "TSAnyKeyword" | "TSUnknownKeyword" | "TSNumberKeyword" | "TSObjectKeyword" | "TSBooleanKeyword" | "TSBigIntKeyword" | "TSStringKeyword" | "TSSymbolKeyword" | "TSVoidKeyword" | "TSUndefinedKeyword" | "TSNullKeyword" | "TSNeverKeyword" | "TSIntrinsicKeyword";
interface TsKeywordType extends TsTypeBase {
    type: TsKeywordTypeType;
}
interface TsThisType extends TsTypeBase {
    type: "TSThisType";
}
type TsFunctionOrConstructorType = TsFunctionType | TsConstructorType;
interface TsFunctionType extends TsTypeBase, TsSignatureDeclarationBase {
    type: "TSFunctionType";
    typeAnnotation: TsTypeAnnotation;
}
interface TsConstructorType extends TsTypeBase, TsSignatureDeclarationBase {
    type: "TSConstructorType";
    typeAnnotation: TsTypeAnnotation;
    abstract: boolean;
}
interface TsTypeReference extends TsTypeBase {
    type: "TSTypeReference";
    typeName: TsEntityName;
    typeParameters?: TsTypeParameterInstantiation;
}
interface TsTypePredicate extends TsTypeBase {
    type: "TSTypePredicate";
    parameterName: Identifier | TsThisType;
    typeAnnotation: TsTypeAnnotation | null;
    asserts: boolean;
}
interface TsTypeQuery extends TsTypeBase {
    type: "TSTypeQuery";
    exprName: TsEntityName | TsImportType;
    typeParameters?: TsTypeParameterInstantiation;
}
interface TsTypeLiteral extends TsTypeBase {
    type: "TSTypeLiteral";
    members: TsTypeElement[];
}
interface TsArrayType extends TsTypeBase {
    type: "TSArrayType";
    elementType: TsType;
}
interface TsTupleType extends TsTypeBase {
    type: "TSTupleType";
    elementTypes: Array<TsType | TsNamedTupleMember>;
}
interface TsNamedTupleMember extends TsTypeBase {
    type: "TSNamedTupleMember";
    label: Identifier;
    optional: boolean;
    elementType: TsType;
}
interface TsOptionalType extends TsTypeBase {
    type: "TSOptionalType";
    typeAnnotation: TsType;
}
interface TsRestType extends TsTypeBase {
    type: "TSRestType";
    typeAnnotation: TsType | TsNamedTupleMember;
}
type TsUnionOrIntersectionType = TsUnionType | TsIntersectionType;
interface TsUnionOrIntersectionTypeBase extends TsTypeBase {
    types: TsType[];
}
interface TsUnionType extends TsUnionOrIntersectionTypeBase {
    type: "TSUnionType";
}
interface TsIntersectionType extends TsUnionOrIntersectionTypeBase {
    type: "TSIntersectionType";
}
interface TsConditionalType extends TsTypeBase {
    type: "TSConditionalType";
    checkType: TsType;
    extendsType: TsType;
    trueType: TsType;
    falseType: TsType;
}
interface TsInferType extends TsTypeBase {
    type: "TSInferType";
    typeParameter: TsTypeParameter;
}
interface TsParenthesizedType extends TsTypeBase {
    type: "TSParenthesizedType";
    typeAnnotation: TsType;
}
interface TsTypeOperator extends TsTypeBase {
    type: "TSTypeOperator";
    operator: "keyof" | "unique" | "readonly";
    typeAnnotation: TsType;
}
interface TsIndexedAccessType extends TsTypeBase {
    type: "TSIndexedAccessType";
    objectType: TsType;
    indexType: TsType;
}
interface TsMappedType extends TsTypeBase {
    type: "TSMappedType";
    readonly?: true | "+" | "-";
    typeParameter: TsTypeParameter;
    optional?: true | "+" | "-";
    typeAnnotation: TsType | undefined | null;
    nameType: TsType | undefined | null;
}
interface TsLiteralType extends TsTypeBase {
    type: "TSLiteralType";
    literal: NumericLiteral | StringLiteral | BooleanLiteral | TemplateLiteral;
}
interface TsImportType extends TsTypeBase {
    type: "TSImportType";
    argument: StringLiteral;
    qualifier?: TsEntityName;
    typeParameters?: TsTypeParameterInstantiation;
    options?: Expression | null;
}
interface TsInterfaceDeclaration extends DeclarationBase {
    type: "TSInterfaceDeclaration";
    id: Identifier | undefined | null;
    typeParameters: TsTypeParameterDeclaration | undefined | null;
    extends?: TsExpressionWithTypeArguments[];
    body: TSInterfaceBody;
}
interface TSInterfaceBody extends NodeBase {
    type: "TSInterfaceBody";
    body: TsTypeElement[];
}
interface TsExpressionWithTypeArguments extends TsTypeBase {
    type: "TSExpressionWithTypeArguments";
    expression: TsEntityName;
    typeParameters?: TsTypeParameterInstantiation;
}
interface TsTypeAliasDeclaration extends DeclarationBase {
    type: "TSTypeAliasDeclaration";
    id: Identifier;
    typeParameters: TsTypeParameterDeclaration | undefined | null;
    typeAnnotation: TsType;
}
interface TsEnumDeclaration extends DeclarationBase {
    type: "TSEnumDeclaration";
    const?: true;
    id: Identifier;
    members: TsEnumMember[];
}
interface TsEnumMember extends NodeBase {
    type: "TSEnumMember";
    id: Identifier | StringLiteral;
    initializer?: Expression;
}
interface TsModuleDeclaration extends DeclarationBase {
    type: "TSModuleDeclaration";
    global?: true;
    id: TsModuleName;
    body: TsNamespaceBody;
}
type TsNamespaceBody = TsModuleBlock | TsNamespaceDeclaration;
interface TsModuleBlock extends NodeBase {
    type: "TSModuleBlock";
    body: Statement[];
}
interface TsNamespaceDeclaration extends TsModuleDeclaration {
    id: Identifier;
    body: TsNamespaceBody;
}
type TsModuleName = Identifier | StringLiteral;
interface TsImportEqualsDeclaration extends NodeBase {
    type: "TSImportEqualsDeclaration";
    isExport: boolean;
    id: Identifier;
    importKind: "type" | "value";
    moduleReference: TsModuleReference;
}
type TsModuleReference = TsEntityName | TsExternalModuleReference;
interface TsExternalModuleReference extends NodeBase {
    type: "TSExternalModuleReference";
    expression: StringLiteral;
}
interface TsExportAssignment extends NodeBase {
    type: "TSExportAssignment";
    expression: Expression;
}
interface TsNamespaceExportDeclaration extends NodeBase {
    type: "TSNamespaceExportDeclaration";
    id: Identifier;
}
interface TsTypeAssertionLikeBase extends NodeBase {
    expression: Expression;
    typeAnnotation: TsType;
}
interface TsAsExpression extends TsTypeAssertionLikeBase {
    type: "TSAsExpression";
}
interface TsTypeAssertion extends TsTypeAssertionLikeBase {
    type: "TSTypeAssertion";
}
type TsSatisfiesExpression = TsTypeAssertionLikeBase & {
    type: "TSSatisfiesExpression";
};
interface TsNonNullExpression extends NodeBase {
    type: "TSNonNullExpression";
    expression: Expression;
}
interface TsInstantiationExpression extends NodeBase {
    type: "TSInstantiationExpression";
    expression: Expression;
    typeParameters: TsTypeParameterInstantiation;
}
interface Placeholder<N extends PlaceholderTypes = PlaceholderTypes> extends NodeBase {
    type: "Placeholder";
    name: Identifier;
    expectedNode: N;
}
type Node = ArgumentPlaceholder | ArrayExpression | ArrayPattern | ArrowFunctionExpression | AssignmentExpression | AssignmentPattern | AwaitExpression | BigIntLiteral | BigIntLiteralTypeAnnotation | BinaryExpression | BindExpression | BlockStatement | BooleanLiteral | BooleanLiteralTypeAnnotation | BreakStatement | CallExpression | CatchClause | ClassAccessorProperty | ClassBody | ClassDeclaration | ClassExpression | ClassMethod | ClassPrivateMethod | ClassPrivateProperty | ClassProperty | ConditionalExpression | ContinueStatement | DebuggerStatement | DecimalLiteral | Decorator | Directive | DirectiveLiteral | DoExpression | DoWhileStatement | EmptyStatement | EstreeChainExpression | EstreeLiteral | EstreeMethodDefinition | EstreePrivateIdentifier | EstreeProperty | EstreePropertyDefinition | ExportAllDeclaration | ExportDefaultDeclaration | ExportDefaultSpecifier | ExportNamedDeclaration | ExportNamespaceSpecifier | ExportSpecifier | ExpressionStatement | File | FlowClassImplements | FlowDeclareClass | FlowDeclareExportDeclaration | FlowDeclareFunction | FlowDeclareInterface | FlowDeclareModule | FlowDeclareModuleExports | FlowDeclareOpaqueType | FlowDeclareTypeAlias | FlowDeclareVariable | FlowEnumBody | FlowEnumDeclaration | FlowEnumMember | FlowFunctionTypeAnnotation | FlowFunctionTypeParam | FlowGenericTypeAnnotation | FlowIndexedAccessType | FlowInterface | FlowInterfaceExtends | FlowInterfaceType | FlowObjectTypeAnnotation | FlowObjectTypeCallProperty | FlowObjectTypeIndexer | FlowObjectTypeInternalSlot | FlowObjectTypeProperty | FlowObjectTypeSpreadProperty | FlowOpaqueType | FlowOptionalIndexedAccessType | FlowOtherTypeAnnotation | FlowPredicate | FlowPredicate | FlowQualifiedTypeIdentifier | FlowTupleTypeAnnotation | FlowTypeAlias | FlowTypeofTypeAnnotation | FlowVariance | ForInStatement | ForOfStatement | ForStatement | FunctionDeclaration | FunctionExpression | Identifier | IfStatement | Import | ImportAttribute | ImportDeclaration | ImportDefaultSpecifier | ImportExpression | ImportExpression | ImportNamespaceSpecifier | ImportSpecifier | InterpreterDirective | JSXAttribute | JSXClosingElement | JSXClosingFragment | JSXElement | JSXEmptyExpression | JSXExpressionContainer | JSXFragment | JSXIdentifier | JSXMemberExpression | JSXNamespacedName | JSXOpeningElement | JSXOpeningFragment | JSXSpreadAttribute | JSXSpreadChild | JSXText | LabeledStatement | Literal | LogicalExpression | MemberExpression | MetaProperty | ModuleExpression | NewExpression | NullLiteral | NumberLiteralTypeAnnotation | NumericLiteral | ObjectExpression | ObjectMethod | ObjectPattern | ObjectProperty | OptionalCallExpression | OptionalMemberExpression | ParenthesizedExpression | PipelineBareAwaitedFunctionBody | PipelineBareConstructorBody | PipelineBareFunction | PipelineBareFunctionBody | PipelineBody | PipelinePrimaryTopicReference | PipelineTopicBody | PipelineTopicExpression | Placeholder | PrivateName | Program | RecordExpression | RegExpLiteral | RestElement | ReturnStatement | SequenceExpression | SpreadElement | StaticBlock | StringLiteral | StringLiteralTypeAnnotation | Super | SwitchCase | SwitchStatement | TSDeclareFunction | TSDeclareMethod | TSInterfaceBody | TSParameterProperty | TaggedTemplateExpression | TemplateElement | TemplateLiteral | ThisExpression | ThrowStatement | TopicReference | TryStatement | TsArrayType | TsAsExpression | TsCallSignatureDeclaration | TsConditionalType | TsConstructSignatureDeclaration | TsConstructorType | TsEnumDeclaration | TsEnumMember | TsExportAssignment | TsExpressionWithTypeArguments | TsExternalModuleReference | TsFunctionType | TsImportEqualsDeclaration | TsImportType | TsIndexSignature | TsIndexedAccessType | TsInferType | TsInstantiationExpression | TsInterfaceDeclaration | TsIntersectionType | TsKeywordType | TsLiteralType | TsMappedType | TsMethodSignature | TsModuleBlock | TsModuleDeclaration | TsNamedTupleMember | TsNamespaceExportDeclaration | TsNonNullExpression | TsOptionalType | TsParenthesizedType | TsPropertySignature | TsQualifiedName | TsRestType | TsSatisfiesExpression | TsThisType | TsTupleType | TsTypeAliasDeclaration | TsTypeAnnotation | TsTypeAssertion | TsTypeCastExpression | TsTypeLiteral | TsTypeOperator | TsTypeParameter | TsTypeParameterDeclaration | TsTypeParameterInstantiation | TsTypePredicate | TsTypeQuery | TsTypeReference | TsUnionType | TupleExpression | TypeAnnotation | TypeCastExpression | TypeParameter | TypeParameterDeclaration | TypeParameterInstantiation | UnaryExpression | UpdateExpression | VariableDeclaration | VariableDeclarator | WhileStatement | WithStatement | YieldExpression;

type PossiblePlaceholders = {
    Identifier: Identifier;
    StringLiteral: StringLiteral;
    Expression: Expression;
    Statement: Statement;
    Declaration: Declaration;
    BlockStatement: BlockStatement;
    ClassBody: ClassBody;
    Pattern: Pattern;
};
type PlaceholderTypes = keyof PossiblePlaceholders;

type PluginList = PluginConfig[];

type SourceType = "script" | "module" | "unambiguous";
interface Options {
    sourceType?: SourceType;
    sourceFilename?: string;
    startColumn?: number;
    startLine?: number;
    allowAwaitOutsideFunction?: boolean;
    allowReturnOutsideFunction?: boolean;
    allowNewTargetOutsideFunction?: boolean;
    allowImportExportEverywhere?: boolean;
    allowSuperOutsideMethod?: boolean;
    allowUndeclaredExports?: boolean;
    plugins?: PluginList;
    strictMode?: boolean | undefined | null;
    ranges?: boolean;
    tokens?: boolean;
    createImportExpressions?: boolean;
    createParenthesizedExpressions?: boolean;
    errorRecovery?: boolean;
    attachComment?: boolean;
    annexB?: boolean;
}

declare function parse(input: string, options?: Options): File;
declare function parseExpression(input: string, options?: Options): Expression;
declare const tokTypes: Record<string, ExportedTokenType>;

type ParserOptions = Partial<Options>;

export { type Expression, type File, type FlowPluginOptions, type ParserOptions, type PluginConfig as ParserPlugin, type PipelineOperatorPluginOptions, type RecordAndTuplePluginOptions, parse, parseExpression, tokTypes };
