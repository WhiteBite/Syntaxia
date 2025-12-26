import { Options } from 'prettier';
export { Options as PrettierOptions } from 'prettier';
import { TransformOptions as TransformOptions$1 } from '@babel/core';

interface RemoveTypeOptions {
    /** Whether to remove ts-ignore and ts-expect-error comments */
    removeTsComments?: boolean;
    /** Escape hatch for customizing Babel configuration */
    customizeBabelConfig?: (config: TransformOptions$1) => void;
}
interface TransformOptions extends RemoveTypeOptions {
    /** Prettier options */
    prettierOptions?: Options | null;
}
/**
 * Transform TypeScript code into vanilla JavaScript without affecting the formatting
 * @param code            Source coude
 * @param fileName        File name for the source
 * @param options         Options
 */
declare function transform(code: string, fileName: string, options?: TransformOptions): Promise<string>;
/**
 * Removes magic comments without performing the TS to JS transform
 * @param code            Source coude
 * @param fileName        File name for the source
 * @param prettierOptions Options to pass to prettier
 */
declare function removeMagicComments(code: string, fileName: string, prettierOptions?: Options | null): Promise<string>;

/**
 * Transform the input file and write the output to another file
 * @param inputFileName
 * @param outputFileName
 */
declare function transformFile(inputFileName: string, outputFileName: string, options?: RemoveTypeOptions): Promise<void>;
/**
 * Remove magic comments from the input file and write the output to another file
 * @param inputFileName
 * @param outputFileName
 */
declare function removeMagicCommentsFromFile(inputFileName: string, outputFileName: string): Promise<void>;

export { removeMagicComments, removeMagicCommentsFromFile, transform, transformFile };
