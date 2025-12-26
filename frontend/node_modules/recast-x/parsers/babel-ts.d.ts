import { parser } from "./babel";
import { Overrides } from "./_babel_options";
export { parser };
export declare function parse(source: string, options?: Overrides): ReturnType<typeof parser.parse>;
