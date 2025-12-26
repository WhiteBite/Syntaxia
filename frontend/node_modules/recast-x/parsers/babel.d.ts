import { parse as babelParse } from "@babel/parser";
import { Overrides } from "./_babel_options";
type BabelParser = {
    parse: typeof babelParse;
};
export declare const parser: BabelParser;
export declare function parse(source: string, options?: Overrides): ReturnType<typeof parser.parse>;
export {};
