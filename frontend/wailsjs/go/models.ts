export namespace domain {
	
	export class AffectedGraph {
	    changedFiles: string[];
	    affectedFiles: string[];
	    dependencies: Record<string, Array<string>>;
	    testMapping: Record<string, Array<string>>;
	
	    static createFrom(source: any = {}) {
	        return new AffectedGraph(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.changedFiles = source["changedFiles"];
	        this.affectedFiles = source["affectedFiles"];
	        this.dependencies = source["dependencies"];
	        this.testMapping = source["testMapping"];
	    }
	}
	export class ApplyResult {
	    success: boolean;
	    path: string;
	    operationId: string;
	    error?: string;
	    appliedLines: number;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.path = source["path"];
	        this.operationId = source["operationId"];
	        this.error = source["error"];
	        this.appliedLines = source["appliedLines"];
	        this.metadata = source["metadata"];
	    }
	}
	export class LayerInfo {
	    name: string;
	    path: string;
	    description: string;
	    dependencies: string[];
	    files: string[];
	    patterns: string[];
	
	    static createFrom(source: any = {}) {
	        return new LayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.description = source["description"];
	        this.dependencies = source["dependencies"];
	        this.files = source["files"];
	        this.patterns = source["patterns"];
	    }
	}
	export class ArchitectureInfo {
	    type: string;
	    confidence: number;
	    description: string;
	    indicators: string[];
	    layers: LayerInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ArchitectureInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.confidence = source["confidence"];
	        this.description = source["description"];
	        this.indicators = source["indicators"];
	        this.layers = this.convertValues(source["layers"], LayerInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Bottleneck {
	    Type: string;
	    Description: string;
	    Duration: number;
	    Impact: string;
	    Suggestions: string[];
	
	    static createFrom(source: any = {}) {
	        return new Bottleneck(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Description = source["Description"];
	        this.Duration = source["Duration"];
	        this.Impact = source["Impact"];
	        this.Suggestions = source["Suggestions"];
	    }
	}
	export class BudgetPolicy {
	    ID: string;
	    Name: string;
	    Description: string;
	    Type: string;
	    Limit: number;
	    Unit: string;
	    TimeWindow: number;
	    Enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BudgetPolicy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Type = source["Type"];
	        this.Limit = source["Limit"];
	        this.Unit = source["Unit"];
	        this.TimeWindow = source["TimeWindow"];
	        this.Enabled = source["Enabled"];
	    }
	}
	export class BudgetViolation {
	    PolicyID: string;
	    Type: string;
	    Current: number;
	    Limit: number;
	    Unit: string;
	    Message: string;
	    // Go type: time
	    Timestamp: any;
	    Context: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new BudgetViolation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PolicyID = source["PolicyID"];
	        this.Type = source["Type"];
	        this.Current = source["Current"];
	        this.Limit = source["Limit"];
	        this.Unit = source["Unit"];
	        this.Message = source["Message"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Context = source["Context"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BuildResult {
	    success: boolean;
	    language: string;
	    projectPath: string;
	    output: string;
	    error?: string;
	    duration: number;
	    artifacts?: string[];
	    warnings?: string[];
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new BuildResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.language = source["language"];
	        this.projectPath = source["projectPath"];
	        this.output = source["output"];
	        this.error = source["error"];
	        this.duration = source["duration"];
	        this.artifacts = source["artifacts"];
	        this.warnings = source["warnings"];
	        this.metadata = source["metadata"];
	    }
	}
	export class BuildSystemInfo {
	    name: string;
	    configFile: string;
	    scripts: string[];
	    commands: string[];
	
	    static createFrom(source: any = {}) {
	        return new BuildSystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.configFile = source["configFile"];
	        this.scripts = source["scripts"];
	        this.commands = source["commands"];
	    }
	}
	export class CodeStyleInfo {
	    indentStyle: string;
	    indentSize: number;
	    maxLineLength: number;
	    trailingComma: boolean;
	    semicolons: boolean;
	    quoteStyle: string;
	    configFile: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeStyleInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.indentStyle = source["indentStyle"];
	        this.indentSize = source["indentSize"];
	        this.maxLineLength = source["maxLineLength"];
	        this.trailingComma = source["trailingComma"];
	        this.semicolons = source["semicolons"];
	        this.quoteStyle = source["quoteStyle"];
	        this.configFile = source["configFile"];
	    }
	}
	export class CommitWithFiles {
	    hash: string;
	    subject: string;
	    author: string;
	    date: string;
	    files: string[];
	    isMerge: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CommitWithFiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.subject = source["subject"];
	        this.author = source["author"];
	        this.date = source["date"];
	        this.files = source["files"];
	        this.isMerge = source["isMerge"];
	    }
	}
	export class ComplianceIssue {
	    type: string;
	    severity: string;
	    description: string;
	    component?: string;
	    recommendation?: string;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.description = source["description"];
	        this.component = source["component"];
	        this.recommendation = source["recommendation"];
	    }
	}
	export class LicenseConflict {
	    license1: string;
	    license2: string;
	    description: string;
	    severity: string;
	
	    static createFrom(source: any = {}) {
	        return new LicenseConflict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.license1 = source["license1"];
	        this.license2 = source["license2"];
	        this.description = source["description"];
	        this.severity = source["severity"];
	    }
	}
	export class LicenseSummary {
	    totalLicenses: number;
	    byType: Record<string, number>;
	    byLicense: Record<string, number>;
	    conflicts?: LicenseConflict[];
	
	    static createFrom(source: any = {}) {
	        return new LicenseSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalLicenses = source["totalLicenses"];
	        this.byType = source["byType"];
	        this.byLicense = source["byLicense"];
	        this.conflicts = this.convertValues(source["conflicts"], LicenseConflict);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LicenseInfo {
	    name: string;
	    spdxId: string;
	    type: string;
	    files: string[];
	    confidence: number;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new LicenseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.spdxId = source["spdxId"];
	        this.type = source["type"];
	        this.files = source["files"];
	        this.confidence = source["confidence"];
	        this.description = source["description"];
	    }
	}
	export class LicenseScanResult {
	    success: boolean;
	    projectPath: string;
	    licenses: LicenseInfo[];
	    summary?: LicenseSummary;
	    metadata?: Record<string, any>;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new LicenseScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.projectPath = source["projectPath"];
	        this.licenses = this.convertValues(source["licenses"], LicenseInfo);
	        this.summary = this.convertValues(source["summary"], LicenseSummary);
	        this.metadata = source["metadata"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VulnerabilitySummary {
	    total: number;
	    critical: number;
	    high: number;
	    medium: number;
	    low: number;
	    fixed: number;
	
	    static createFrom(source: any = {}) {
	        return new VulnerabilitySummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.critical = source["critical"];
	        this.high = source["high"];
	        this.medium = source["medium"];
	        this.low = source["low"];
	        this.fixed = source["fixed"];
	    }
	}
	export class VulnerabilityScanResult {
	    success: boolean;
	    projectPath: string;
	    vulnerabilities: Vulnerability[];
	    summary?: VulnerabilitySummary;
	    metadata?: Record<string, any>;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new VulnerabilityScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.projectPath = source["projectPath"];
	        this.vulnerabilities = this.convertValues(source["vulnerabilities"], Vulnerability);
	        this.summary = this.convertValues(source["summary"], VulnerabilitySummary);
	        this.metadata = source["metadata"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Vulnerability {
	    id: string;
	    severity: string;
	    description: string;
	    cvss?: number;
	    fixedIn?: string;
	
	    static createFrom(source: any = {}) {
	        return new Vulnerability(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.severity = source["severity"];
	        this.description = source["description"];
	        this.cvss = source["cvss"];
	        this.fixedIn = source["fixedIn"];
	    }
	}
	export class SBOMComponent {
	    name: string;
	    version: string;
	    type: string;
	    purl?: string;
	    license?: string;
	    vulnerabilities?: Vulnerability[];
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SBOMComponent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.type = source["type"];
	        this.purl = source["purl"];
	        this.license = source["license"];
	        this.vulnerabilities = this.convertValues(source["vulnerabilities"], Vulnerability);
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SBOMResult {
	    success: boolean;
	    projectPath: string;
	    format: string;
	    outputPath: string;
	    components: SBOMComponent[];
	    metadata?: Record<string, any>;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new SBOMResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.projectPath = source["projectPath"];
	        this.format = source["format"];
	        this.outputPath = source["outputPath"];
	        this.components = this.convertValues(source["components"], SBOMComponent);
	        this.metadata = source["metadata"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ComplianceSummary {
	    totalIssues: number;
	    critical: number;
	    high: number;
	    medium: number;
	    low: number;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalIssues = source["totalIssues"];
	        this.critical = source["critical"];
	        this.high = source["high"];
	        this.medium = source["medium"];
	        this.low = source["low"];
	    }
	}
	export class ComplianceReport {
	    success: boolean;
	    projectPath: string;
	    compliant: boolean;
	    issues: ComplianceIssue[];
	    summary?: ComplianceSummary;
	    metadata?: Record<string, any>;
	    // Go type: time
	    generatedAt: any;
	    sbomResult?: SBOMResult;
	    vulnerabilityResult?: VulnerabilityScanResult;
	    licenseResult?: LicenseScanResult;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.projectPath = source["projectPath"];
	        this.compliant = source["compliant"];
	        this.issues = this.convertValues(source["issues"], ComplianceIssue);
	        this.summary = this.convertValues(source["summary"], ComplianceSummary);
	        this.metadata = source["metadata"];
	        this.generatedAt = this.convertValues(source["generatedAt"], null);
	        this.sbomResult = this.convertValues(source["sbomResult"], SBOMResult);
	        this.vulnerabilityResult = this.convertValues(source["vulnerabilityResult"], VulnerabilityScanResult);
	        this.licenseResult = this.convertValues(source["licenseResult"], LicenseScanResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ComplianceRequirements {
	    allowedLicenses: string[];
	    forbiddenLicenses: string[];
	    maxVulnerabilities: number;
	    maxCriticalVulnerabilities: number;
	    maxHighVulnerabilities: number;
	    maxCVSS: number;
	    requireSBOM: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ComplianceRequirements(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.allowedLicenses = source["allowedLicenses"];
	        this.forbiddenLicenses = source["forbiddenLicenses"];
	        this.maxVulnerabilities = source["maxVulnerabilities"];
	        this.maxCriticalVulnerabilities = source["maxCriticalVulnerabilities"];
	        this.maxHighVulnerabilities = source["maxHighVulnerabilities"];
	        this.maxCVSS = source["maxCVSS"];
	        this.requireSBOM = source["requireSBOM"];
	    }
	}
	
	export class ContextBuildOptions {
	    stripComments: boolean;
	    includeManifest: boolean;
	    includeLineNumbers: boolean;
	    maxTokens: number;
	    maxMemoryMB: number;
	    includeTests: boolean;
	    splitStrategy: string;
	    forceStream: boolean;
	    enableProgressEvents: boolean;
	    outputFormat: string;
	    excludeTests: boolean;
	    collapseEmptyLines: boolean;
	    stripLicense: boolean;
	    compactDataFiles: boolean;
	    skeletonMode: boolean;
	    trimWhitespace: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ContextBuildOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stripComments = source["stripComments"];
	        this.includeManifest = source["includeManifest"];
	        this.includeLineNumbers = source["includeLineNumbers"];
	        this.maxTokens = source["maxTokens"];
	        this.maxMemoryMB = source["maxMemoryMB"];
	        this.includeTests = source["includeTests"];
	        this.splitStrategy = source["splitStrategy"];
	        this.forceStream = source["forceStream"];
	        this.enableProgressEvents = source["enableProgressEvents"];
	        this.outputFormat = source["outputFormat"];
	        this.excludeTests = source["excludeTests"];
	        this.collapseEmptyLines = source["collapseEmptyLines"];
	        this.stripLicense = source["stripLicense"];
	        this.compactDataFiles = source["compactDataFiles"];
	        this.skeletonMode = source["skeletonMode"];
	        this.trimWhitespace = source["trimWhitespace"];
	    }
	}
	export class ContextMetadata {
	    buildDuration: number;
	    // Go type: time
	    lastModified: any;
	    selectedFiles: string[];
	    buildOptions?: ContextBuildOptions;
	    warnings: string[];
	    errors: string[];
	    contentPath: string;
	    skippedFiles?: string[];
	    skippedReasons?: Record<string, string>;
	    projectPath?: string;
	    chunksCount?: number;
	    splitStrategy?: string;
	
	    static createFrom(source: any = {}) {
	        return new ContextMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.buildDuration = source["buildDuration"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.selectedFiles = source["selectedFiles"];
	        this.buildOptions = this.convertValues(source["buildOptions"], ContextBuildOptions);
	        this.warnings = source["warnings"];
	        this.errors = source["errors"];
	        this.contentPath = source["contentPath"];
	        this.skippedFiles = source["skippedFiles"];
	        this.skippedReasons = source["skippedReasons"];
	        this.projectPath = source["projectPath"];
	        this.chunksCount = source["chunksCount"];
	        this.splitStrategy = source["splitStrategy"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ContextSummary {
	    id: string;
	    projectPath: string;
	    fileCount: number;
	    totalSize: number;
	    tokenCount: number;
	    lineCount: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    status: string;
	    metadata: ContextMetadata;
	
	    static createFrom(source: any = {}) {
	        return new ContextSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectPath = source["projectPath"];
	        this.fileCount = source["fileCount"];
	        this.totalSize = source["totalSize"];
	        this.tokenCount = source["tokenCount"];
	        this.lineCount = source["lineCount"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.status = source["status"];
	        this.metadata = this.convertValues(source["metadata"], ContextMetadata);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportStyle {
	    absoluteImports: boolean;
	    aliasedImports: boolean;
	    groupedImports: boolean;
	    importOrder: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImportStyle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.absoluteImports = source["absoluteImports"];
	        this.aliasedImports = source["aliasedImports"];
	        this.groupedImports = source["groupedImports"];
	        this.importOrder = source["importOrder"];
	    }
	}
	export class TestConventions {
	    location: string;
	    fileSuffix: string;
	    framework: string;
	    patterns: string[];
	
	    static createFrom(source: any = {}) {
	        return new TestConventions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.location = source["location"];
	        this.fileSuffix = source["fileSuffix"];
	        this.framework = source["framework"];
	        this.patterns = source["patterns"];
	    }
	}
	export class FileNamingStyle {
	    style: string;
	    suffixes: string[];
	    prefixes: string[];
	    examples: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileNamingStyle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.style = source["style"];
	        this.suffixes = source["suffixes"];
	        this.prefixes = source["prefixes"];
	        this.examples = source["examples"];
	    }
	}
	export class ConventionInfo {
	    namingStyle: string;
	    fileNaming: FileNamingStyle;
	    folderStructure: string;
	    testConventions: TestConventions;
	    importStyle: ImportStyle;
	    codeStyle: CodeStyleInfo;
	
	    static createFrom(source: any = {}) {
	        return new ConventionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.namingStyle = source["namingStyle"];
	        this.fileNaming = this.convertValues(source["fileNaming"], FileNamingStyle);
	        this.folderStructure = source["folderStructure"];
	        this.testConventions = this.convertValues(source["testConventions"], TestConventions);
	        this.importStyle = this.convertValues(source["importStyle"], ImportStyle);
	        this.codeStyle = this.convertValues(source["codeStyle"], CodeStyleInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DiffImpact {
	    RiskLevel: string;
	    AffectedTests: string[];
	    BreakingChanges: string[];
	    PerformanceImpact: string;
	    SecurityImpact: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffImpact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RiskLevel = source["RiskLevel"];
	        this.AffectedTests = source["AffectedTests"];
	        this.BreakingChanges = source["BreakingChanges"];
	        this.PerformanceImpact = source["PerformanceImpact"];
	        this.SecurityImpact = source["SecurityImpact"];
	    }
	}
	export class RiskAssessment {
	    level: string;
	    risks: string[];
	    mitigations: string[];
	    testCoverage: string;
	    reviewNeeded: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RiskAssessment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.risks = source["risks"];
	        this.mitigations = source["mitigations"];
	        this.testCoverage = source["testCoverage"];
	        this.reviewNeeded = source["reviewNeeded"];
	    }
	}
	export class ImpactAnalysis {
	    level: string;
	    affectedAreas: string[];
	    breaking: boolean;
	    performance: string;
	    security: string;
	    maintainability: string;
	
	    static createFrom(source: any = {}) {
	        return new ImpactAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.affectedAreas = source["affectedAreas"];
	        this.breaking = source["breaking"];
	        this.performance = source["performance"];
	        this.security = source["security"];
	        this.maintainability = source["maintainability"];
	    }
	}
	export class WhyView {
	    reason: string;
	    taskId: string;
	    stepId: string;
	    confidence: number;
	    relatedFiles: string[];
	    context?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new WhyView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reason = source["reason"];
	        this.taskId = source["taskId"];
	        this.stepId = source["stepId"];
	        this.confidence = source["confidence"];
	        this.relatedFiles = source["relatedFiles"];
	        this.context = source["context"];
	    }
	}
	export class DiffSummary {
	    totalFiles: number;
	    addedFiles: number;
	    modifiedFiles: number;
	    deletedFiles: number;
	    totalLines: number;
	    addedLines: number;
	    removedLines: number;
	    whyView?: WhyView;
	    impact?: ImpactAnalysis;
	    risk?: RiskAssessment;
	
	    static createFrom(source: any = {}) {
	        return new DiffSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalFiles = source["totalFiles"];
	        this.addedFiles = source["addedFiles"];
	        this.modifiedFiles = source["modifiedFiles"];
	        this.deletedFiles = source["deletedFiles"];
	        this.totalLines = source["totalLines"];
	        this.addedLines = source["addedLines"];
	        this.removedLines = source["removedLines"];
	        this.whyView = this.convertValues(source["whyView"], WhyView);
	        this.impact = this.convertValues(source["impact"], ImpactAnalysis);
	        this.risk = this.convertValues(source["risk"], RiskAssessment);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DiffChange {
	    Type: string;
	    FilePath: string;
	    LineNumber: number;
	    OldContent: string;
	    NewContent: string;
	    Reason: string;
	    Confidence: number;
	
	    static createFrom(source: any = {}) {
	        return new DiffChange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.FilePath = source["FilePath"];
	        this.LineNumber = source["LineNumber"];
	        this.OldContent = source["OldContent"];
	        this.NewContent = source["NewContent"];
	        this.Reason = source["Reason"];
	        this.Confidence = source["Confidence"];
	    }
	}
	export class DerivedDiffReport {
	    TaskID: string;
	    OriginalDiff: string;
	    DerivedDiff: string;
	    Changes: DiffChange[];
	    Summary?: DiffSummary;
	    Impact: DiffImpact;
	
	    static createFrom(source: any = {}) {
	        return new DerivedDiffReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskID = source["TaskID"];
	        this.OriginalDiff = source["OriginalDiff"];
	        this.DerivedDiff = source["DerivedDiff"];
	        this.Changes = this.convertValues(source["Changes"], DiffChange);
	        this.Summary = this.convertValues(source["Summary"], DiffSummary);
	        this.Impact = this.convertValues(source["Impact"], DiffImpact);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DiffHunk {
	    oldStart: number;
	    oldCount: number;
	    newStart: number;
	    newCount: number;
	    lines: string[];
	
	    static createFrom(source: any = {}) {
	        return new DiffHunk(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.oldStart = source["oldStart"];
	        this.oldCount = source["oldCount"];
	        this.newStart = source["newStart"];
	        this.newCount = source["newCount"];
	        this.lines = source["lines"];
	    }
	}
	export class DiffEntry {
	    path: string;
	    operation: string;
	    oldContent?: string;
	    newContent?: string;
	    hunks?: DiffHunk[];
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new DiffEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.operation = source["operation"];
	        this.oldContent = source["oldContent"];
	        this.newContent = source["newContent"];
	        this.hunks = this.convertValues(source["hunks"], DiffHunk);
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class DiffResult {
	    id: string;
	    format: string;
	    content: string;
	    entries: DiffEntry[];
	    summary?: DiffSummary;
	    generatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.format = source["format"];
	        this.content = source["content"];
	        this.entries = this.convertValues(source["entries"], DiffEntry);
	        this.summary = this.convertValues(source["summary"], DiffSummary);
	        this.generatedAt = source["generatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Edit {
	    id: string;
	    atomicGroup?: string;
	    dependsOn?: string[];
	    kind: string;
	    op: string;
	    path: string;
	    filePath: string;
	    language: string;
	    content?: string;
	    anchor?: any;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Edit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.atomicGroup = source["atomicGroup"];
	        this.dependsOn = source["dependsOn"];
	        this.kind = source["kind"];
	        this.op = source["op"];
	        this.path = source["path"];
	        this.filePath = source["filePath"];
	        this.language = source["language"];
	        this.content = source["content"];
	        this.anchor = source["anchor"];
	        this.metadata = source["metadata"];
	    }
	}
	export class EditsMetadata {
	    reason: string;
	    taskId: string;
	    stepId: string;
	    confidence: number;
	    estimatedImpact: string;
	
	    static createFrom(source: any = {}) {
	        return new EditsMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reason = source["reason"];
	        this.taskId = source["taskId"];
	        this.stepId = source["stepId"];
	        this.confidence = source["confidence"];
	        this.estimatedImpact = source["estimatedImpact"];
	    }
	}
	export class EditsJSON {
	    schemaVersion: string;
	    toolchainVersion: string;
	    metadata?: EditsMetadata;
	    edits: Edit[];
	
	    static createFrom(source: any = {}) {
	        return new EditsJSON(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schemaVersion = source["schemaVersion"];
	        this.toolchainVersion = source["toolchainVersion"];
	        this.metadata = this.convertValues(source["metadata"], EditsMetadata);
	        this.edits = this.convertValues(source["edits"], Edit);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ExportResult {
	    mode: string;
	    text?: string;
	    fileName?: string;
	    dataBase64?: string;
	    filePath?: string;
	    isLarge?: boolean;
	    sizeBytes?: number;
	
	    static createFrom(source: any = {}) {
	        return new ExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.text = source["text"];
	        this.fileName = source["fileName"];
	        this.dataBase64 = source["dataBase64"];
	        this.filePath = source["filePath"];
	        this.isLarge = source["isLarge"];
	        this.sizeBytes = source["sizeBytes"];
	    }
	}
	
	export class FileNode {
	    name: string;
	    path: string;
	    relPath: string;
	    isDir: boolean;
	    size: number;
	    contentType: string;
	    children?: FileNode[];
	    isGitignored: boolean;
	    isCustomIgnored: boolean;
	    isIgnored: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.relPath = source["relPath"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.contentType = source["contentType"];
	        this.children = this.convertValues(source["children"], FileNode);
	        this.isGitignored = source["isGitignored"];
	        this.isCustomIgnored = source["isCustomIgnored"];
	        this.isIgnored = source["isIgnored"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileReason {
	    FilePath: string;
	    Reason: string;
	    Impact: string;
	    Confidence: number;
	    RelatedFiles: string[];
	    Context: Record<string, any>;
	    category: string;
	    importance: string;
	    suggestions: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileReason(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FilePath = source["FilePath"];
	        this.Reason = source["Reason"];
	        this.Impact = source["Impact"];
	        this.Confidence = source["Confidence"];
	        this.RelatedFiles = source["RelatedFiles"];
	        this.Context = source["Context"];
	        this.category = source["category"];
	        this.importance = source["importance"];
	        this.suggestions = source["suggestions"];
	    }
	}
	export class FileStatus {
	    path: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new FileStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.status = source["status"];
	    }
	}
	export class FrameworkInfo {
	    name: string;
	    version: string;
	    category: string;
	    language: string;
	    configFiles: string[];
	    indicators: string[];
	    bestPractices: string[];
	
	    static createFrom(source: any = {}) {
	        return new FrameworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.category = source["category"];
	        this.language = source["language"];
	        this.configFiles = source["configFiles"];
	        this.indicators = source["indicators"];
	        this.bestPractices = source["bestPractices"];
	    }
	}
	export class GuardrailRule {
	    ID: string;
	    Pattern: string;
	    Description: string;
	    Action: string;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new GuardrailRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Pattern = source["Pattern"];
	        this.Description = source["Description"];
	        this.Action = source["Action"];
	        this.Message = source["Message"];
	    }
	}
	export class GuardrailPolicy {
	    ID: string;
	    Name: string;
	    Description: string;
	    Type: string;
	    Severity: string;
	    Enabled: boolean;
	    Rules: GuardrailRule[];
	
	    static createFrom(source: any = {}) {
	        return new GuardrailPolicy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Type = source["Type"];
	        this.Severity = source["Severity"];
	        this.Enabled = source["Enabled"];
	        this.Rules = this.convertValues(source["Rules"], GuardrailRule);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GuardrailViolation {
	    PolicyID: string;
	    RuleID: string;
	    Severity: string;
	    Message: string;
	    FilePath: string;
	    LineNumber: number;
	    // Go type: time
	    Timestamp: any;
	    Context: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new GuardrailViolation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PolicyID = source["PolicyID"];
	        this.RuleID = source["RuleID"];
	        this.Severity = source["Severity"];
	        this.Message = source["Message"];
	        this.FilePath = source["FilePath"];
	        this.LineNumber = source["LineNumber"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Context = source["Context"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class LanguageAnalysisValidation {
	    language: string;
	    success: boolean;
	    issueCount: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new LanguageAnalysisValidation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.success = source["success"];
	        this.issueCount = source["issueCount"];
	        this.error = source["error"];
	    }
	}
	export class LanguageInfo {
	    name: string;
	    version: string;
	    fileCount: number;
	    percentage: number;
	    primary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LanguageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.fileCount = source["fileCount"];
	        this.percentage = source["percentage"];
	        this.primary = source["primary"];
	    }
	}
	export class TypeIssue {
	    file: string;
	    line: number;
	    column: number;
	    severity: string;
	    message: string;
	    code?: string;
	
	    static createFrom(source: any = {}) {
	        return new TypeIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.column = source["column"];
	        this.severity = source["severity"];
	        this.message = source["message"];
	        this.code = source["code"];
	    }
	}
	export class TypeCheckResult {
	    success: boolean;
	    language: string;
	    projectPath: string;
	    output: string;
	    error?: string;
	    duration: number;
	    issues?: TypeIssue[];
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new TypeCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.language = source["language"];
	        this.projectPath = source["projectPath"];
	        this.output = source["output"];
	        this.error = source["error"];
	        this.duration = source["duration"];
	        this.issues = this.convertValues(source["issues"], TypeIssue);
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LanguageValidationResult {
	    success: boolean;
	    language: string;
	    typeCheckResult?: TypeCheckResult;
	    buildResult?: BuildResult;
	    typeCheckError?: string;
	    buildError?: string;
	
	    static createFrom(source: any = {}) {
	        return new LanguageValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.language = source["language"];
	        this.typeCheckResult = this.convertValues(source["typeCheckResult"], TypeCheckResult);
	        this.buildResult = this.convertValues(source["buildResult"], BuildResult);
	        this.typeCheckError = source["typeCheckError"];
	        this.buildError = source["buildError"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class PerformanceMetrics {
	    TaskID: string;
	    MemoryUsage: number;
	    CPUUsage: number;
	    DiskIO: number;
	    NetworkIO: number;
	    FileOperations: number;
	    APIRequests: number;
	    CacheHits: number;
	    CacheMisses: number;
	    Timestamps: string[];
	    Values: number[];
	
	    static createFrom(source: any = {}) {
	        return new PerformanceMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskID = source["TaskID"];
	        this.MemoryUsage = source["MemoryUsage"];
	        this.CPUUsage = source["CPUUsage"];
	        this.DiskIO = source["DiskIO"];
	        this.NetworkIO = source["NetworkIO"];
	        this.FileOperations = source["FileOperations"];
	        this.APIRequests = source["APIRequests"];
	        this.CacheHits = source["CacheHits"];
	        this.CacheMisses = source["CacheMisses"];
	        this.Timestamps = source["Timestamps"];
	        this.Values = source["Values"];
	    }
	}
	export class ProjectStructure {
	    architecture?: ArchitectureInfo;
	    conventions?: ConventionInfo;
	    frameworks: FrameworkInfo[];
	    buildSystems: BuildSystemInfo[];
	    languages: LanguageInfo[];
	    layers: LayerInfo[];
	    projectType: string;
	    confidence: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectStructure(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.architecture = this.convertValues(source["architecture"], ArchitectureInfo);
	        this.conventions = this.convertValues(source["conventions"], ConventionInfo);
	        this.frameworks = this.convertValues(source["frameworks"], FrameworkInfo);
	        this.buildSystems = this.convertValues(source["buildSystems"], BuildSystemInfo);
	        this.languages = this.convertValues(source["languages"], LanguageInfo);
	        this.layers = this.convertValues(source["layers"], LayerInfo);
	        this.projectType = source["projectType"];
	        this.confidence = source["confidence"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectValidationResult {
	    success: boolean;
	    projectPath: string;
	    languages: string[];
	    results: Record<string, LanguageValidationResult>;
	
	    static createFrom(source: any = {}) {
	        return new ProjectValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.projectPath = source["projectPath"];
	        this.languages = source["languages"];
	        this.results = this.convertValues(source["results"], LanguageValidationResult, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RecentProjectInfo {
	    path: string;
	    name: string;
	    lastOpenedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new RecentProjectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.lastOpenedAt = source["lastOpenedAt"];
	    }
	}
	export class RepairResult {
	    Success: boolean;
	    RuleID: string;
	    FixedFiles: string[];
	    Error: string;
	    Duration: number;
	    Attempts: number;
	
	    static createFrom(source: any = {}) {
	        return new RepairResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Success = source["Success"];
	        this.RuleID = source["RuleID"];
	        this.FixedFiles = source["FixedFiles"];
	        this.Error = source["Error"];
	        this.Duration = source["Duration"];
	        this.Attempts = source["Attempts"];
	    }
	}
	export class RepairRule {
	    ID: string;
	    Name: string;
	    Description: string;
	    Pattern: string;
	    Fix: string;
	    Priority: number;
	    Language: string;
	    Category: string;
	
	    static createFrom(source: any = {}) {
	        return new RepairRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Pattern = source["Pattern"];
	        this.Fix = source["Fix"];
	        this.Priority = source["Priority"];
	        this.Language = source["Language"];
	        this.Category = source["Category"];
	    }
	}
	
	
	
	export class SettingsDTO {
	    customIgnoreRules: string;
	    customPromptRules: string;
	    openAIAPIKey: string;
	    geminiAPIKey: string;
	    openRouterAPIKey: string;
	    localAIAPIKey: string;
	    localAIHost: string;
	    localAIModelName: string;
	    qwenAPIKey: string;
	    qwenHost: string;
	    selectedProvider: string;
	    selectedModels: Record<string, string>;
	    availableModels: Record<string, Array<string>>;
	    useGitignore: boolean;
	    useCustomIgnore: boolean;
	    recentProjects?: RecentProjectInfo[];
	
	    static createFrom(source: any = {}) {
	        return new SettingsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customIgnoreRules = source["customIgnoreRules"];
	        this.customPromptRules = source["customPromptRules"];
	        this.openAIAPIKey = source["openAIAPIKey"];
	        this.geminiAPIKey = source["geminiAPIKey"];
	        this.openRouterAPIKey = source["openRouterAPIKey"];
	        this.localAIAPIKey = source["localAIAPIKey"];
	        this.localAIHost = source["localAIHost"];
	        this.localAIModelName = source["localAIModelName"];
	        this.qwenAPIKey = source["qwenAPIKey"];
	        this.qwenHost = source["qwenHost"];
	        this.selectedProvider = source["selectedProvider"];
	        this.selectedModels = source["selectedModels"];
	        this.availableModels = source["availableModels"];
	        this.useGitignore = source["useGitignore"];
	        this.useCustomIgnore = source["useCustomIgnore"];
	        this.recentProjects = this.convertValues(source["recentProjects"], RecentProjectInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaticAnalysisReportSummary {
	    totalIssues: number;
	    totalErrors: number;
	    totalWarnings: number;
	    languagesAnalyzed: string[];
	    analyzersUsed: string[];
	    criticalIssues?: StaticIssue[];
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StaticAnalysisReportSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalIssues = source["totalIssues"];
	        this.totalErrors = source["totalErrors"];
	        this.totalWarnings = source["totalWarnings"];
	        this.languagesAnalyzed = source["languagesAnalyzed"];
	        this.analyzersUsed = source["analyzersUsed"];
	        this.criticalIssues = this.convertValues(source["criticalIssues"], StaticIssue);
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaticAnalysisSummary {
	    totalIssues: number;
	    errorCount: number;
	    warningCount: number;
	    infoCount: number;
	    hintCount: number;
	    severityBreakdown: Record<string, number>;
	    categoryBreakdown: Record<string, number>;
	    filesAnalyzed: number;
	    filesWithIssues: number;
	
	    static createFrom(source: any = {}) {
	        return new StaticAnalysisSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalIssues = source["totalIssues"];
	        this.errorCount = source["errorCount"];
	        this.warningCount = source["warningCount"];
	        this.infoCount = source["infoCount"];
	        this.hintCount = source["hintCount"];
	        this.severityBreakdown = source["severityBreakdown"];
	        this.categoryBreakdown = source["categoryBreakdown"];
	        this.filesAnalyzed = source["filesAnalyzed"];
	        this.filesWithIssues = source["filesWithIssues"];
	    }
	}
	export class StaticIssue {
	    file: string;
	    line: number;
	    column: number;
	    severity: string;
	    message: string;
	    code?: string;
	    category?: string;
	    confidence?: string;
	    suggestions?: string[];
	
	    static createFrom(source: any = {}) {
	        return new StaticIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.column = source["column"];
	        this.severity = source["severity"];
	        this.message = source["message"];
	        this.code = source["code"];
	        this.category = source["category"];
	        this.confidence = source["confidence"];
	        this.suggestions = source["suggestions"];
	    }
	}
	export class StaticAnalysisResult {
	    success: boolean;
	    language: string;
	    projectPath: string;
	    analyzer: string;
	    issues: StaticIssue[];
	    summary?: StaticAnalysisSummary;
	    duration: number;
	    error?: string;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new StaticAnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.language = source["language"];
	        this.projectPath = source["projectPath"];
	        this.analyzer = source["analyzer"];
	        this.issues = this.convertValues(source["issues"], StaticIssue);
	        this.summary = this.convertValues(source["summary"], StaticAnalysisSummary);
	        this.duration = source["duration"];
	        this.error = source["error"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StaticAnalysisReport {
	    projectPath: string;
	    timestamp: string;
	    totalDuration: number;
	    results: Record<string, StaticAnalysisResult>;
	    summary?: StaticAnalysisReportSummary;
	    recommendations?: string[];
	
	    static createFrom(source: any = {}) {
	        return new StaticAnalysisReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectPath = source["projectPath"];
	        this.timestamp = source["timestamp"];
	        this.totalDuration = source["totalDuration"];
	        this.results = this.convertValues(source["results"], StaticAnalysisResult, true);
	        this.summary = this.convertValues(source["summary"], StaticAnalysisReportSummary);
	        this.recommendations = source["recommendations"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class StaticAnalysisValidationResult {
	    success: boolean;
	    totalLanguages: number;
	    successCount: number;
	    failureCount: number;
	    successRate: number;
	    totalIssues: number;
	    totalErrors: number;
	    totalWarnings: number;
	    failedLanguages: string[];
	    criticalIssues: StaticIssue[];
	    languages: Record<string, LanguageAnalysisValidation>;
	
	    static createFrom(source: any = {}) {
	        return new StaticAnalysisValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.totalLanguages = source["totalLanguages"];
	        this.successCount = source["successCount"];
	        this.failureCount = source["failureCount"];
	        this.successRate = source["successRate"];
	        this.totalIssues = source["totalIssues"];
	        this.totalErrors = source["totalErrors"];
	        this.totalWarnings = source["totalWarnings"];
	        this.failedLanguages = source["failedLanguages"];
	        this.criticalIssues = this.convertValues(source["criticalIssues"], StaticIssue);
	        this.languages = this.convertValues(source["languages"], LanguageAnalysisValidation, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SymbolEdge {
	    from: string;
	    to: string;
	    type: string;
	    weight: number;
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SymbolEdge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	        this.type = source["type"];
	        this.weight = source["weight"];
	        this.metadata = source["metadata"];
	    }
	}
	export class SymbolNode {
	    id: string;
	    name: string;
	    type: string;
	    path: string;
	    line: number;
	    column: number;
	    package?: string;
	    visibility: string;
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SymbolNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.line = source["line"];
	        this.column = source["column"];
	        this.package = source["package"];
	        this.visibility = source["visibility"];
	        this.metadata = source["metadata"];
	    }
	}
	export class SymbolGraph {
	    nodes: SymbolNode[];
	    edges: SymbolEdge[];
	
	    static createFrom(source: any = {}) {
	        return new SymbolGraph(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodes = this.convertValues(source["nodes"], SymbolNode);
	        this.edges = this.convertValues(source["edges"], SymbolEdge);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TaskBudgets {
	    maxFiles: number;
	    maxChangedLines: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskBudgets(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxFiles = source["maxFiles"];
	        this.maxChangedLines = source["maxChangedLines"];
	    }
	}
	export class Task {
	    ID: string;
	    Name: string;
	    Description: string;
	    State: string;
	    DependsOn: string[];
	    StepFile: string;
	    Budgets: TaskBudgets;
	    Status: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: time
	    StartedAt?: any;
	    // Go type: time
	    CompletedAt?: any;
	    Error: string;
	    Metadata: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.State = source["State"];
	        this.DependsOn = source["DependsOn"];
	        this.StepFile = source["StepFile"];
	        this.Budgets = this.convertValues(source["Budgets"], TaskBudgets);
	        this.Status = source["Status"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.StartedAt = this.convertValues(source["StartedAt"], null);
	        this.CompletedAt = this.convertValues(source["CompletedAt"], null);
	        this.Error = source["Error"];
	        this.Metadata = source["Metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TaskStatus {
	    TaskID: string;
	    State: string;
	    Progress: number;
	    Message: string;
	    Error: string;
	    // Go type: time
	    StartedAt?: any;
	    // Go type: time
	    CompletedAt?: any;
	    // Go type: time
	    UpdatedAt: any;
	    Duration: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskID = source["TaskID"];
	        this.State = source["State"];
	        this.Progress = source["Progress"];
	        this.Message = source["Message"];
	        this.Error = source["Error"];
	        this.StartedAt = this.convertValues(source["StartedAt"], null);
	        this.CompletedAt = this.convertValues(source["CompletedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Duration = source["Duration"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TestConfig {
	    language: string;
	    projectPath: string;
	    scope: string;
	    parallel: boolean;
	    timeout: number;
	    coverage: boolean;
	    verbose: boolean;
	    envVars?: Record<string, string>;
	    testPatterns?: string[];
	    excludePatterns?: string[];
	
	    static createFrom(source: any = {}) {
	        return new TestConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.projectPath = source["projectPath"];
	        this.scope = source["scope"];
	        this.parallel = source["parallel"];
	        this.timeout = source["timeout"];
	        this.coverage = source["coverage"];
	        this.verbose = source["verbose"];
	        this.envVars = source["envVars"];
	        this.testPatterns = source["testPatterns"];
	        this.excludePatterns = source["excludePatterns"];
	    }
	}
	
	export class TestCoverage {
	    percentage: number;
	    lines: number;
	    functions: number;
	    branches: number;
	    files?: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new TestCoverage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.percentage = source["percentage"];
	        this.lines = source["lines"];
	        this.functions = source["functions"];
	        this.branches = source["branches"];
	        this.files = source["files"];
	    }
	}
	export class TestInfo {
	    path: string;
	    name: string;
	    type: string;
	    targetFiles?: string[];
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new TestInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.targetFiles = source["targetFiles"];
	        this.metadata = source["metadata"];
	    }
	}
	export class TestResult {
	    success: boolean;
	    testPath: string;
	    testName: string;
	    language: string;
	    duration: number;
	    output: string;
	    error?: string;
	    coverage?: TestCoverage;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new TestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.testPath = source["testPath"];
	        this.testName = source["testName"];
	        this.language = source["language"];
	        this.duration = source["duration"];
	        this.output = source["output"];
	        this.error = source["error"];
	        this.coverage = this.convertValues(source["coverage"], TestCoverage);
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TestSuite {
	    name: string;
	    language: string;
	    projectPath: string;
	    tests: TestInfo[];
	    config?: TestConfig;
	
	    static createFrom(source: any = {}) {
	        return new TestSuite(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.language = source["language"];
	        this.projectPath = source["projectPath"];
	        this.tests = this.convertValues(source["tests"], TestInfo);
	        this.config = this.convertValues(source["config"], TestConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TestValidationResult {
	    success: boolean;
	    totalTests: number;
	    passedTests: number;
	    failedTests: number;
	    skippedTests: number;
	    successRate: number;
	    totalDuration: number;
	    failedTestPaths: string[];
	
	    static createFrom(source: any = {}) {
	        return new TestValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.totalTests = source["totalTests"];
	        this.passedTests = source["passedTests"];
	        this.failedTests = source["failedTests"];
	        this.skippedTests = source["skippedTests"];
	        this.successRate = source["successRate"];
	        this.totalDuration = source["totalDuration"];
	        this.failedTestPaths = source["failedTestPaths"];
	    }
	}
	export class TimeToGreenMetrics {
	    TaskID: string;
	    // Go type: time
	    StartTime: any;
	    // Go type: time
	    EndTime: any;
	    Duration: number;
	    Attempts: number;
	    RepairAttempts: number;
	    BuildTime: number;
	    TestTime: number;
	    StaticAnalysisTime: number;
	    TotalTime: number;
	    Success: boolean;
	    Bottlenecks: Bottleneck[];
	
	    static createFrom(source: any = {}) {
	        return new TimeToGreenMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskID = source["TaskID"];
	        this.StartTime = this.convertValues(source["StartTime"], null);
	        this.EndTime = this.convertValues(source["EndTime"], null);
	        this.Duration = source["Duration"];
	        this.Attempts = source["Attempts"];
	        this.RepairAttempts = source["RepairAttempts"];
	        this.BuildTime = source["BuildTime"];
	        this.TestTime = source["TestTime"];
	        this.StaticAnalysisTime = source["StaticAnalysisTime"];
	        this.TotalTime = source["TotalTime"];
	        this.Success = source["Success"];
	        this.Bottlenecks = this.convertValues(source["Bottlenecks"], Bottleneck);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class UXReport {
	    ID: string;
	    Type: string;
	    Title: string;
	    Description: string;
	    Content: any;
	    // Go type: time
	    CreatedAt: any;
	    Metadata: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new UXReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Type = source["Type"];
	        this.Title = source["Title"];
	        this.Description = source["Description"];
	        this.Content = source["Content"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.Metadata = source["Metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class WhyViewReport {
	    TaskID: string;
	    Files: FileReason[];
	    Context: string;
	    Explanation: string;
	    Confidence: number;
	    Suggestions: string[];
	
	    static createFrom(source: any = {}) {
	        return new WhyViewReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskID = source["TaskID"];
	        this.Files = this.convertValues(source["Files"], FileReason);
	        this.Context = source["Context"];
	        this.Explanation = source["Explanation"];
	        this.Confidence = source["Confidence"];
	        this.Suggestions = source["Suggestions"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class AffectedFile {
	    path: string;
	    type: string;
	    dependents: number;
	
	    static createFrom(source: any = {}) {
	        return new AffectedFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.type = source["type"];
	        this.dependents = source["dependents"];
	    }
	}
	export class ContextMemoryEntry {
	    id: string;
	    topic: string;
	    summary: string;
	    files: string[];
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ContextMemoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.topic = source["topic"];
	        this.summary = source["summary"];
	        this.files = source["files"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class DirCount {
	    dir: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new DirCount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	        this.count = source["count"];
	    }
	}
	export class FileQuickInfo {
	    symbolCount: number;
	    importCount: number;
	    dependentCount: number;
	    changeRisk: number;
	    riskLevel: string;
	
	    static createFrom(source: any = {}) {
	        return new FileQuickInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.symbolCount = source["symbolCount"];
	        this.importCount = source["importCount"];
	        this.dependentCount = source["dependentCount"];
	        this.changeRisk = source["changeRisk"];
	        this.riskLevel = source["riskLevel"];
	    }
	}
	export class IgnorePreviewResult {
	    totalFiles: number;
	    byDirectory: Record<string, number>;
	    byRule: Record<string, number>;
	    sampleFiles: string[];
	    topDirs: DirCount[];
	
	    static createFrom(source: any = {}) {
	        return new IgnorePreviewResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalFiles = source["totalFiles"];
	        this.byDirectory = source["byDirectory"];
	        this.byRule = source["byRule"];
	        this.sampleFiles = source["sampleFiles"];
	        this.topDirs = this.convertValues(source["topDirs"], DirCount);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImpactPreviewResult {
	    totalDependents: number;
	    aggregateRisk: number;
	    riskLevel: string;
	    affectedFiles: AffectedFile[];
	    relatedTests: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImpactPreviewResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalDependents = source["totalDependents"];
	        this.aggregateRisk = source["aggregateRisk"];
	        this.riskLevel = source["riskLevel"];
	        this.affectedFiles = this.convertValues(source["affectedFiles"], AffectedFile);
	        this.relatedTests = source["relatedTests"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ShellIntegrationStatus {
	    isRegistered: boolean;
	    currentOS: string;
	
	    static createFrom(source: any = {}) {
	        return new ShellIntegrationStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isRegistered = source["isRegistered"];
	        this.currentOS = source["currentOS"];
	    }
	}
	export class SmartSuggestion {
	    path: string;
	    source: string;
	    reason: string;
	    confidence: number;
	
	    static createFrom(source: any = {}) {
	        return new SmartSuggestion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.source = source["source"];
	        this.reason = source["reason"];
	        this.confidence = source["confidence"];
	    }
	}
	export class SmartSuggestionsResult {
	    suggestions: SmartSuggestion[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new SmartSuggestionsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.suggestions = this.convertValues(source["suggestions"], SmartSuggestion);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace version {
	
	export class Info {
	    version: string;
	    gitCommit: string;
	    buildDate: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.gitCommit = source["gitCommit"];
	        this.buildDate = source["buildDate"];
	    }
	}
	export class Release {
	    tag_name: string;
	    name: string;
	    body: string;
	    // Go type: time
	    published_at: any;
	    html_url: string;
	    prerelease: boolean;
	    draft: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Release(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.name = source["name"];
	        this.body = source["body"];
	        this.published_at = this.convertValues(source["published_at"], null);
	        this.html_url = source["html_url"];
	        this.prerelease = source["prerelease"];
	        this.draft = source["draft"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ReleasesResponse {
	    currentVersion: string;
	    latestVersion: string;
	    hasUpdate: boolean;
	    releases: Release[];
	    // Go type: time
	    lastChecked: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ReleasesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.hasUpdate = source["hasUpdate"];
	        this.releases = this.convertValues(source["releases"], Release);
	        this.lastChecked = this.convertValues(source["lastChecked"], null);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

