"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.parser = void 0;
exports.parse = parse;
var tslib_1 = require("tslib");
var babel_1 = require("./babel");
Object.defineProperty(exports, "parser", { enumerable: true, get: function () { return babel_1.parser; } });
var _babel_options_1 = tslib_1.__importDefault(require("./_babel_options"));
function parse(source, options) {
    var babelOptions = (0, _babel_options_1.default)(options);
    babelOptions.plugins.push("jsx", "typescript");
    return babel_1.parser.parse(source, babelOptions);
}
