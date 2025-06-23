export interface SlitherIssue {
  type: string;
  title: string;
  description: string;
  severity: "HIGH" | "MEDIUM" | "LOW" | "INFORMATIONAL" | "OPTIMIZATION";
  confidence: string;
  location: string;
  reference: string;
}

export interface StaticAnalysisResponse {
  success: boolean;
  issues: SlitherIssue[];
  total_issues: number;
  severity_summary: {
    high: number;
    medium: number;
    low: number;
    informational: number;
    optimization?: number;
  };
  analyzed_at: string;
  source_hash: string;
}

export interface AIVulnerability {
  title: string;
  severity: "HIGH" | "MEDIUM" | "LOW" | "INFORMATIONAL";
  description: string;
  location: {
    function: string;
    line_numbers: number[];
  };
  recommendation: string;
}

export interface AIAnalysisResponse {
  success: boolean;
  analysis: {
    security_score: number;
    risk_level: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
    summary: string;
    vulnerabilities: AIVulnerability[];
    good_practices: string[] | null;
    recommendations: string[] | null;
  };
  total_issues: number;
  analyzed_at: string;
  source_hash: string;
}

export interface TestCaseResponse {
  success: boolean;
  test_code: string;
  test_framework: string;
  test_language: string;
  file_name: string;
  source_hash: string;
  warnings_and_recommendations: string[];
  generated_at: string;
}

// Types for Fetch Contract API
export interface ContractSourceResponse {
  address: string;
  source_code: string;
  source_hash: string;
  network: string;
}

export type NetworkType =
  | "ethereum"
  | "polygon"
  | "bsc"
  | "base"
  | "arbitrum"
  | "avalanche"
  | "optimism"
  | "gnosis"
  | "fantom"
  | "celo";
