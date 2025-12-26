"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = default_1;
var tslib_1 = require("tslib");
var shared_1 = require("../shared");
var es2021_1 = tslib_1.__importDefault(require("./operators/es2021"));
var es2020_1 = tslib_1.__importDefault(require("./es2020"));
function default_1(fork) {
    // The es2021OpsDef plugin comes before es2020Def so AssignmentOperators will
    // be appropriately augmented before first used.
    fork.use(es2021_1.default);
    fork.use(es2020_1.default);
}
(0, shared_1.maybeSetModuleExports)(function () { return module; });
//# sourceMappingURL=es2021.js.map