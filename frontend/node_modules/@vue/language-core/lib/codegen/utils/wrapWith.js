"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.wrapWith = wrapWith;
function* wrapWith(source, startOffset, endOffset, features, ...codes) {
    yield ['', source, startOffset, features];
    let offset = 1;
    for (const code of codes) {
        if (typeof code !== 'string') {
            offset++;
        }
        yield code;
    }
    yield ['', source, endOffset, { __combineOffset: offset }];
}
//# sourceMappingURL=wrapWith.js.map