"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateSrc = generateSrc;
exports.resolveSrcPath = resolveSrcPath;
const codeFeatures_1 = require("../codeFeatures");
const utils_1 = require("../utils");
const boundary_1 = require("../utils/boundary");
function* generateSrc(src) {
    if (src === true) {
        return;
    }
    let { text } = src;
    text = resolveSrcPath(text);
    yield `export * from `;
    const wrapCodeFeatures = {
        ...codeFeatures_1.codeFeatures.all,
        ...text !== src.text ? codeFeatures_1.codeFeatures.navigationWithoutRename : {},
    };
    const token = yield* (0, boundary_1.startBoundary)('main', src.offset, wrapCodeFeatures);
    yield `'`;
    yield [text.slice(0, src.text.length), 'main', src.offset, { __combineToken: token }];
    yield text.slice(src.text.length);
    yield `'`;
    yield (0, boundary_1.endBoundary)(token, src.offset + src.text.length);
    yield utils_1.endOfLine;
    yield `export { default } from '${text}'${utils_1.endOfLine}`;
}
function resolveSrcPath(text) {
    if (text.endsWith('.d.ts')) {
        text = text.slice(0, -'.d.ts'.length);
    }
    else if (text.endsWith('.ts')) {
        text = text.slice(0, -'.ts'.length);
    }
    else if (text.endsWith('.tsx')) {
        text = text.slice(0, -'.tsx'.length) + '.jsx';
    }
    if (!text.endsWith('.js') && !text.endsWith('.jsx')) {
        text = text + '.js';
    }
    return text;
}
//# sourceMappingURL=src.js.map