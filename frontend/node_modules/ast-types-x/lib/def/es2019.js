"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = default_1;
var tslib_1 = require("tslib");
var types_1 = tslib_1.__importDefault(require("../types"));
var shared_1 = tslib_1.__importStar(require("../shared"));
var es2018_1 = tslib_1.__importDefault(require("./es2018"));
function default_1(fork) {
    fork.use(es2018_1.default);
    var types = fork.use(types_1.default);
    var def = types.Type.def;
    var or = types.Type.or;
    var defaults = fork.use(shared_1.default).defaults;
    def("CatchClause")
        .field("param", or(def("Pattern"), null), defaults["null"]);
}
;
(0, shared_1.maybeSetModuleExports)(function () { return module; });
//# sourceMappingURL=es2019.js.map