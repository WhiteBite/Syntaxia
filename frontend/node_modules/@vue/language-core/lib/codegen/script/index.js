"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateScript = generate;
const path = require("path-browserify");
const codeFeatures_1 = require("../codeFeatures");
const names = require("../names");
const utils_1 = require("../utils");
const boundary_1 = require("../utils/boundary");
const context_1 = require("./context");
const scriptSetup_1 = require("./scriptSetup");
const src_1 = require("./src");
const template_1 = require("./template");
const exportExpression = `{} as typeof ${names._export}`;
function generate(options) {
    const ctx = (0, context_1.createScriptCodegenContext)(options);
    const codeGenerator = generateWorker(options, ctx);
    return { ...ctx, codes: [...codeGenerator] };
}
function* generateWorker(options, ctx) {
    yield* generateGlobalTypesReference(options);
    const { script, scriptRanges, scriptSetup, scriptSetupRanges, vueCompilerOptions } = options;
    if (scriptSetup && scriptSetupRanges) {
        yield* (0, scriptSetup_1.generateScriptSetupImports)(scriptSetup, scriptSetupRanges);
    }
    if (script?.src) {
        yield* (0, src_1.generateSrc)(script.src);
    }
    // <script> + <script setup>
    if (script && scriptRanges && scriptSetup && scriptSetupRanges) {
        // <script>
        let selfType;
        const { exportDefault } = scriptRanges;
        if (exportDefault) {
            yield* generateScriptWithExportDefault(ctx, script, scriptRanges, exportDefault, vueCompilerOptions, selfType = '__VLS_self');
        }
        else {
            yield* (0, utils_1.generateSfcBlockSection)(script, 0, script.content.length, codeFeatures_1.codeFeatures.all);
            yield `export default ${exportExpression}${utils_1.endOfLine}`;
        }
        // <script setup>
        yield* generateExportDeclareEqual(scriptSetup, names._export);
        if (scriptSetup.generic) {
            yield* (0, scriptSetup_1.generateGeneric)(options, ctx, scriptSetup, scriptSetupRanges, scriptSetup.generic, (0, scriptSetup_1.generateSetupFunction)(options, ctx, scriptSetup, scriptSetupRanges, (0, template_1.generateTemplate)(options, ctx, selfType)));
        }
        else {
            yield `await (async () => {${utils_1.newLine}`;
            yield* (0, scriptSetup_1.generateSetupFunction)(options, ctx, scriptSetup, scriptSetupRanges, (0, template_1.generateTemplate)(options, ctx, selfType), [`return `]);
            yield `})()${utils_1.endOfLine}`;
        }
    }
    // only <script setup>
    else if (scriptSetup && scriptSetupRanges) {
        if (scriptSetup.generic) {
            yield* generateExportDeclareEqual(scriptSetup, names._export);
            yield* (0, scriptSetup_1.generateGeneric)(options, ctx, scriptSetup, scriptSetupRanges, scriptSetup.generic, (0, scriptSetup_1.generateSetupFunction)(options, ctx, scriptSetup, scriptSetupRanges, (0, template_1.generateTemplate)(options, ctx)));
        }
        else {
            // no script block, generate script setup code at root
            yield* (0, scriptSetup_1.generateSetupFunction)(options, ctx, scriptSetup, scriptSetupRanges, (0, template_1.generateTemplate)(options, ctx), generateExportDeclareEqual(scriptSetup, names._export));
        }
        yield `export default ${exportExpression}${utils_1.endOfLine}`;
    }
    // only <script>
    else if (script && scriptRanges) {
        const { exportDefault } = scriptRanges;
        if (exportDefault) {
            yield* generateScriptWithExportDefault(ctx, script, scriptRanges, exportDefault, vueCompilerOptions, names._export, (0, template_1.generateTemplate)(options, ctx, names._export));
        }
        else {
            yield* (0, utils_1.generateSfcBlockSection)(script, 0, script.content.length, codeFeatures_1.codeFeatures.all);
            yield* generateExportDeclareEqual(script, names._export);
            yield `(await import('${vueCompilerOptions.lib}')).defineComponent({})${utils_1.endOfLine}`;
            yield* (0, template_1.generateTemplate)(options, ctx, names._export);
            yield `export default ${exportExpression}${utils_1.endOfLine}`;
        }
    }
    yield* ctx.localTypes.generate();
}
function* generateScriptWithExportDefault(ctx, script, scriptRanges, exportDefault, vueCompilerOptions, varName, templateGenerator) {
    const { componentOptions } = scriptRanges;
    const { expression, isObjectLiteral } = componentOptions ?? exportDefault;
    let wrapLeft;
    let wrapRight;
    if (isObjectLiteral
        && vueCompilerOptions.optionsWrapper.length) {
        [wrapLeft, wrapRight] = vueCompilerOptions.optionsWrapper;
        ctx.inlayHints.push({
            blockName: script.name,
            offset: expression.start,
            setting: 'vue.inlayHints.optionsWrapper',
            label: wrapLeft || '[Missing optionsWrapper[0]]',
            tooltip: [
                'This is virtual code that is automatically wrapped for type support, it does not affect your runtime behavior, you can customize it via `vueCompilerOptions.optionsWrapper` option in tsconfig / jsconfig.',
                'To hide it, you can set `"vue.inlayHints.optionsWrapper": false` in IDE settings.',
            ].join('\n\n'),
        }, {
            blockName: script.name,
            offset: expression.end,
            setting: 'vue.inlayHints.optionsWrapper',
            label: wrapRight || '[Missing optionsWrapper[1]]',
        });
    }
    yield* (0, utils_1.generateSfcBlockSection)(script, 0, expression.start, codeFeatures_1.codeFeatures.all);
    yield exportExpression;
    yield* (0, utils_1.generateSfcBlockSection)(script, expression.end, exportDefault.end, codeFeatures_1.codeFeatures.all);
    yield utils_1.endOfLine;
    if (templateGenerator) {
        yield* templateGenerator;
    }
    yield* generateExportDeclareEqual(script, varName);
    if (wrapLeft && wrapRight) {
        yield wrapLeft;
        yield* (0, utils_1.generateSfcBlockSection)(script, expression.start, expression.end, codeFeatures_1.codeFeatures.all);
        yield wrapRight;
    }
    else {
        yield* (0, utils_1.generateSfcBlockSection)(script, expression.start, expression.end, codeFeatures_1.codeFeatures.all);
    }
    yield utils_1.endOfLine;
    yield* (0, utils_1.generateSfcBlockSection)(script, exportDefault.end, script.content.length, codeFeatures_1.codeFeatures.all);
}
function* generateGlobalTypesReference(options) {
    const globalTypesPath = options.vueCompilerOptions.globalTypesPath(options.fileName);
    if (!globalTypesPath) {
        yield `/* placeholder */${utils_1.newLine}`;
    }
    else if (path.isAbsolute(globalTypesPath)) {
        let relativePath = path.relative(path.dirname(options.fileName), globalTypesPath);
        if (relativePath !== globalTypesPath
            && !relativePath.startsWith('./')
            && !relativePath.startsWith('../')) {
            relativePath = './' + relativePath;
        }
        yield `/// <reference types="${relativePath}" />${utils_1.newLine}`;
    }
    else {
        yield `/// <reference types="${globalTypesPath}" />${utils_1.newLine}`;
    }
}
function* generateExportDeclareEqual(block, name) {
    yield `const `;
    const token = yield* (0, boundary_1.startBoundary)(block.name, 0, codeFeatures_1.codeFeatures.doNotReportTs6133);
    yield name;
    yield (0, boundary_1.endBoundary)(token, block.content.length);
    yield ` = `;
}
//# sourceMappingURL=index.js.map