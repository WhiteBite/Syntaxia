/**
 * Мок для wailsjs/go/models
 */

export namespace domain {
    export class FileNode {
        path: string = ''
        name: string = ''
        isDir: boolean = false
        children: FileNode[] = []
        size: number = 0
        modTime: string = ''
        extension: string = ''
        selected: boolean = false
        expanded: boolean = false
    }

    export class ProjectInfo {
        path: string = ''
        name: string = ''
        lastOpened: string = ''
    }

    export class ContextItem {
        id: string = ''
        type: string = ''
        path: string = ''
        content: string = ''
        tokens: number = 0
    }

    export class BuildResult {
        success: boolean = false
        output: string = ''
        errors: string[] = []
    }

    export class TestResult {
        success: boolean = false
        passed: number = 0
        failed: number = 0
        total: number = 0
    }

    export class VersionInfo {
        version: string = ''
        buildDate: string = ''
        commit: string = ''
    }

    export class Release {
        tagName: string = ''
        name: string = ''
        body: string = ''
        publishedAt: string = ''
        htmlUrl: string = ''
    }

    export class FileStats {
        size: number = 0
        modTime: string = ''
        isDir: boolean = false
    }

    export class ProjectStructure {
        frameworks: string[] = []
        languages: string[] = []
        architecture: string = ''
        patterns: string[] = []
    }
}
