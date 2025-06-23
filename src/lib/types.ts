export interface SlitherIssue {
  title: string;
  severity: "HIGH" | "MEDIUM" | "LOW" | "INFORMATIONAL";
  description: string;
  location: string;
}

export interface StaticAnalysisResponse {
  total_issues: number;
  severity_summary: {
    high: number;
    medium: number;
    low: number;
    informational: number;
  };
  issues: SlitherIssue[];
}

export interface AIVulnerability {
  title: string;
  severity: "HIGH" | "MEDIUM" | "LOW" | "INFORMATIONAL";
  description: string;
  recommendation: string;
}

export interface AIAnalysisResponse {
  security_score: number;
  risk_level: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  vulnerabilities: AIVulnerability[];
}

export interface TestCaseResponse {
  test_framework: string;
  test_language: string;
  test_code: string;
  warnings_and_recommendations: string[];
}
