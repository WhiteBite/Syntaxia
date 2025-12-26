"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = default_1;
var shared_1 = require("../../shared");
function default_1() {
    return {
        BinaryOperators: [
            "==", "!=", "===", "!==",
            "<", "<=", ">", ">=",
            "<<", ">>", ">>>",
            "+", "-", "*", "/", "%",
            "&",
            "|", "^", "in",
            "instanceof",
        ],
        AssignmentOperators: [
            "=", "+=", "-=", "*=", "/=", "%=",
            "<<=", ">>=", ">>>=",
            "|=", "^=", "&=",
        ],
        LogicalOperators: [
            "||", "&&",
        ],
    };
}
(0, shared_1.maybeSetModuleExports)(function () { return module; });
//# sourceMappingURL=core.js.map