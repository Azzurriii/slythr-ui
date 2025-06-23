"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Copy,
  AlertTriangle,
  Shield,
  Zap,
  TestTube,
  CheckCircle,
} from "lucide-react";
import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
} from "@/lib/types";

interface AnalysisPanelProps {
  activeTab: string;
  onTabChange: (tab: string) => void;
  staticAnalysisResult: StaticAnalysisResponse | null;
  aiAnalysisResult: AIAnalysisResponse | null;
  testCaseResult: TestCaseResponse | null;
  isLoading: boolean;
}

export function AnalysisPanel({
  activeTab,
  onTabChange,
  staticAnalysisResult,
  aiAnalysisResult,
  testCaseResult,
}: AnalysisPanelProps) {
  const [copiedCode, setCopiedCode] = useState(false);

  const getSeverityColor = (severity: string) => {
    switch (severity.toUpperCase()) {
      case "HIGH":
      case "CRITICAL":
        return "bg-red-500/10 text-red-500 border-red-500/20";
      case "MEDIUM":
        return "bg-orange-500/10 text-orange-500 border-orange-500/20";
      case "LOW":
        return "bg-yellow-500/10 text-yellow-500 border-yellow-500/20";
      case "OPTIMIZATION":
        return "bg-purple-500/10 text-purple-500 border-purple-500/20";
      default:
        return "bg-blue-500/10 text-blue-500 border-blue-500/20";
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedCode(true);
      setTimeout(() => setCopiedCode(false), 2000);
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }
  };

  return (
    <div className="h-full flex flex-col min-h-0">
      {/* Header - Fixed */}
      <div className="flex-shrink-0 border-b border-border bg-muted/50 px-3 sm:px-4 py-2">
        <h2 className="text-sm font-medium truncate">Analysis Results</h2>
      </div>

      <Tabs
        value={activeTab}
        onValueChange={onTabChange}
        className="flex-1 flex flex-col min-h-0"
      >
        {/* Tabs List - Fixed */}
        <div className="flex-shrink-0 px-2 sm:px-4 pt-2 sm:pt-4">
          <TabsList className="grid w-full grid-cols-3 h-8 sm:h-10">
            <TabsTrigger
              value="static"
              className="flex items-center gap-1 px-1 sm:px-2 text-xs sm:text-sm min-w-0"
            >
              <Shield className="h-3 w-3 sm:h-4 sm:w-4 flex-shrink-0" />
              <span className="truncate">Static</span>
            </TabsTrigger>
            <TabsTrigger
              value="ai"
              className="flex items-center gap-1 px-1 sm:px-2 text-xs sm:text-sm min-w-0"
            >
              <Zap className="h-3 w-3 sm:h-4 sm:w-4 flex-shrink-0" />
              <span className="truncate">
                <span className="hidden sm:inline">AI Audit</span>
                <span className="sm:hidden">AI</span>
              </span>
            </TabsTrigger>
            <TabsTrigger
              value="tests"
              className="flex items-center gap-1 px-1 sm:px-2 text-xs sm:text-sm min-w-0"
            >
              <TestTube className="h-3 w-3 sm:h-4 sm:w-4 flex-shrink-0" />
              <span className="truncate">Tests</span>
            </TabsTrigger>
          </TabsList>
        </div>

        {/* Content Area - Scrollable */}
        <div className="flex-1 overflow-y-auto px-2 sm:px-4 pb-4 min-h-0">
          <TabsContent value="static" className="mt-2 space-y-3 sm:space-y-4">
            {!staticAnalysisResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-6 sm:py-8">
                  <Shield className="h-10 w-10 sm:h-12 sm:w-12 text-muted-foreground mb-3 sm:mb-4" />
                  <p className="text-center text-sm text-muted-foreground px-4">
                    Run a static analysis to see the results here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader className="pb-3 sm:pb-4">
                    <CardTitle className="text-base sm:text-lg">
                      Slither Static Analysis
                    </CardTitle>
                    <CardDescription className="text-xs sm:text-sm space-y-1">
                      <div>
                        Found {staticAnalysisResult.total_issues} issues
                      </div>
                      <div className="break-all">
                        Analyzed at{" "}
                        {new Date(
                          staticAnalysisResult.analyzed_at
                        ).toLocaleString()}
                      </div>
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3 sm:space-y-4">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-3 sm:gap-4">
                      <div className="space-y-2">
                        <h4 className="text-sm font-medium">
                          Issues by Severity
                        </h4>
                        <div className="space-y-2">
                          {staticAnalysisResult.severity_summary.high > 0 && (
                            <div className="flex items-center justify-between">
                              <Badge
                                className={getSeverityColor("HIGH")}
                                variant="outline"
                              >
                                High
                              </Badge>
                              <span className="text-sm">
                                {staticAnalysisResult.severity_summary.high}
                              </span>
                            </div>
                          )}
                          {staticAnalysisResult.severity_summary.medium > 0 && (
                            <div className="flex items-center justify-between">
                              <Badge
                                className={getSeverityColor("MEDIUM")}
                                variant="outline"
                              >
                                Medium
                              </Badge>
                              <span className="text-sm">
                                {staticAnalysisResult.severity_summary.medium}
                              </span>
                            </div>
                          )}
                          {staticAnalysisResult.severity_summary.low > 0 && (
                            <div className="flex items-center justify-between">
                              <Badge
                                className={getSeverityColor("LOW")}
                                variant="outline"
                              >
                                Low
                              </Badge>
                              <span className="text-sm">
                                {staticAnalysisResult.severity_summary.low}
                              </span>
                            </div>
                          )}
                          {staticAnalysisResult.severity_summary.informational >
                            0 && (
                            <div className="flex items-center justify-between">
                              <Badge
                                className={getSeverityColor("INFORMATIONAL")}
                                variant="outline"
                              >
                                Info
                              </Badge>
                              <span className="text-sm">
                                {
                                  staticAnalysisResult.severity_summary
                                    .informational
                                }
                              </span>
                            </div>
                          )}
                          {staticAnalysisResult.severity_summary.optimization &&
                            staticAnalysisResult.severity_summary.optimization >
                              0 && (
                              <div className="flex items-center justify-between">
                                <Badge
                                  className={getSeverityColor("OPTIMIZATION")}
                                  variant="outline"
                                >
                                  Optimization
                                </Badge>
                                <span className="text-sm">
                                  {
                                    staticAnalysisResult.severity_summary
                                      .optimization
                                  }
                                </span>
                              </div>
                            )}
                        </div>
                      </div>
                      <div className="space-y-2">
                        <h4 className="text-sm font-medium">Analysis Info</h4>
                        <div className="space-y-1 text-xs text-muted-foreground">
                          <div className="break-all">
                            Source Hash:{" "}
                            <code className="bg-muted px-1 py-0.5 rounded text-xs">
                              {staticAnalysisResult.source_hash.substring(0, 8)}
                              ...
                            </code>
                          </div>
                          <div>
                            Total Issues:{" "}
                            <span className="font-medium">
                              {staticAnalysisResult.total_issues}
                            </span>
                          </div>
                          <div>
                            Status:{" "}
                            <span className="text-green-600 font-medium">
                              ✓ Success
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Accordion type="single" collapsible className="space-y-2">
                  {staticAnalysisResult.issues.map((issue, index) => (
                    <AccordionItem
                      key={index}
                      value={`item-${index}`}
                      className="border rounded-lg px-3 sm:px-4"
                    >
                      <AccordionTrigger className="hover:no-underline py-3">
                        <div className="flex items-center justify-between w-full mr-2 sm:mr-4 min-w-0">
                          <span className="text-left text-sm truncate pr-2">
                            {issue.title}
                          </span>
                          <Badge
                            className={getSeverityColor(issue.severity)}
                            variant="outline"
                          >
                            {issue.severity}
                          </Badge>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent className="space-y-2 pb-3">
                        <div className="flex flex-wrap items-center gap-2 sm:gap-4 text-xs text-muted-foreground mb-2">
                          <div className="flex items-center gap-1">
                            <span>Type:</span>
                            <Badge variant="outline" className="text-xs">
                              {issue.type}
                            </Badge>
                          </div>
                          <div className="flex items-center gap-1">
                            <span>Confidence:</span>
                            <Badge variant="outline" className="text-xs">
                              {issue.confidence}
                            </Badge>
                          </div>
                        </div>
                        <p className="text-sm text-muted-foreground leading-relaxed">
                          {issue.description}
                        </p>
                        <div className="space-y-1">
                          <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2 text-xs text-muted-foreground">
                            <span className="flex-shrink-0">Location:</span>
                            <code className="bg-muted px-1 py-0.5 rounded break-all">
                              {issue.location}
                            </code>
                          </div>
                          <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2 text-xs text-muted-foreground">
                            <span className="flex-shrink-0">Reference:</span>
                            <code className="bg-muted px-1 py-0.5 rounded break-all">
                              {issue.reference}
                            </code>
                          </div>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  ))}
                </Accordion>
              </>
            )}
          </TabsContent>

          <TabsContent value="ai" className="mt-2 space-y-3 sm:space-y-4">
            {!aiAnalysisResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-6 sm:py-8">
                  <Zap className="h-10 w-10 sm:h-12 sm:w-12 text-muted-foreground mb-3 sm:mb-4" />
                  <p className="text-center text-sm text-muted-foreground px-4">
                    Run an AI audit to see the results here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader className="pb-3 sm:pb-4">
                    <CardTitle className="text-base sm:text-lg">
                      AI Dynamic Analysis
                    </CardTitle>
                    <CardDescription className="text-xs sm:text-sm break-words">
                      Risk Level: {aiAnalysisResult.analysis.risk_level} •{" "}
                      {aiAnalysisResult.total_issues} vulnerabilities found •
                      Analyzed at{" "}
                      {new Date(aiAnalysisResult.analyzed_at).toLocaleString()}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3 sm:space-y-4">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-3 sm:gap-4">
                      <div className="space-y-3">
                        <div>
                          <div className="flex justify-between text-sm mb-2">
                            <span>Security Score</span>
                            <span className="font-medium">
                              {aiAnalysisResult.analysis.security_score}/100
                            </span>
                          </div>
                          <Progress
                            value={aiAnalysisResult.analysis.security_score}
                            className="h-2"
                          />
                        </div>
                        <div className="space-y-2">
                          <h4 className="text-sm font-medium">
                            Analysis Summary
                          </h4>
                          <p className="text-sm text-muted-foreground leading-relaxed">
                            {aiAnalysisResult.analysis.summary}
                          </p>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <h4 className="text-sm font-medium">
                          Analysis Details
                        </h4>
                        <div className="space-y-1 text-xs text-muted-foreground">
                          <div>
                            Risk Level:{" "}
                            <span
                              className={`font-medium ${
                                aiAnalysisResult.analysis.risk_level === "HIGH"
                                  ? "text-red-600"
                                  : aiAnalysisResult.analysis.risk_level ===
                                    "MEDIUM"
                                  ? "text-yellow-600"
                                  : "text-green-600"
                              }`}
                            >
                              {aiAnalysisResult.analysis.risk_level}
                            </span>
                          </div>
                          <div className="break-all">
                            Source Hash:{" "}
                            <code className="bg-muted px-1 py-0.5 rounded text-xs">
                              {aiAnalysisResult.source_hash.substring(0, 8)}...
                            </code>
                          </div>
                          <div>
                            Total Vulnerabilities:{" "}
                            <span className="font-medium">
                              {aiAnalysisResult.total_issues}
                            </span>
                          </div>
                          <div>
                            Status:{" "}
                            <span className="text-green-600 font-medium">
                              ✓ Success
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Accordion type="single" collapsible className="space-y-2">
                  {aiAnalysisResult.analysis.vulnerabilities.map(
                    (vuln, index) => (
                      <AccordionItem
                        key={index}
                        value={`ai-item-${index}`}
                        className="border rounded-lg px-3 sm:px-4"
                      >
                        <AccordionTrigger className="hover:no-underline py-3">
                          <div className="flex items-center justify-between w-full mr-2 sm:mr-4 min-w-0">
                            <span className="text-left text-sm truncate pr-2">
                              {vuln.title}
                            </span>
                            <Badge
                              className={getSeverityColor(vuln.severity)}
                              variant="outline"
                            >
                              {vuln.severity}
                            </Badge>
                          </div>
                        </AccordionTrigger>
                        <AccordionContent className="space-y-3 pb-3">
                          <p className="text-sm text-muted-foreground leading-relaxed">
                            {vuln.description}
                          </p>
                          <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2 text-xs text-muted-foreground">
                            <span className="flex-shrink-0">Location:</span>
                            <code className="bg-muted px-1 py-0.5 rounded break-all">
                              {vuln.location.function} (Lines{" "}
                              {vuln.location.line_numbers.join(", ")})
                            </code>
                          </div>
                          <Alert>
                            <CheckCircle className="h-4 w-4" />
                            <AlertDescription className="text-sm">
                              <strong>Recommendation:</strong>{" "}
                              {vuln.recommendation}
                            </AlertDescription>
                          </Alert>
                        </AccordionContent>
                      </AccordionItem>
                    )
                  )}
                </Accordion>

                {/* Good Practices */}
                {aiAnalysisResult.analysis.good_practices &&
                  aiAnalysisResult.analysis.good_practices.length > 0 && (
                    <div className="space-y-3">
                      <div className="flex items-center gap-2 text-sm font-medium text-green-600">
                        <CheckCircle className="h-4 w-4" />
                        Good Practices Found
                      </div>
                      <div className="bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800/30 rounded-lg p-3 sm:p-4">
                        <ul className="space-y-2">
                          {aiAnalysisResult.analysis.good_practices.map(
                            (practice, index) => (
                              <li
                                key={index}
                                className="flex items-start gap-2 text-sm text-green-800 dark:text-green-200"
                              >
                                <CheckCircle className="h-3 w-3 text-green-600 dark:text-green-400 mt-0.5 flex-shrink-0" />
                                <span className="break-words">{practice}</span>
                              </li>
                            )
                          )}
                        </ul>
                      </div>
                    </div>
                  )}

                {/* General Recommendations */}
                {aiAnalysisResult.analysis.recommendations &&
                  aiAnalysisResult.analysis.recommendations.length > 0 && (
                    <div className="space-y-3">
                      <div className="flex items-center gap-2 text-sm font-medium text-blue-600">
                        <AlertTriangle className="h-4 w-4" />
                        General Recommendations
                      </div>
                      <div className="bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800/30 rounded-lg p-3 sm:p-4">
                        <ul className="space-y-2">
                          {aiAnalysisResult.analysis.recommendations.map(
                            (rec, index) => (
                              <li
                                key={index}
                                className="flex items-start gap-2 text-sm text-blue-800 dark:text-blue-200"
                              >
                                <AlertTriangle className="h-3 w-3 text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" />
                                <span className="break-words">{rec}</span>
                              </li>
                            )
                          )}
                        </ul>
                      </div>
                    </div>
                  )}
              </>
            )}
          </TabsContent>

          <TabsContent value="tests" className="mt-2 space-y-3 sm:space-y-4">
            {!testCaseResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-6 sm:py-8">
                  <TestTube className="h-10 w-10 sm:h-12 sm:w-12 text-muted-foreground mb-3 sm:mb-4" />
                  <p className="text-center text-sm text-muted-foreground px-4">
                    Generate test cases to see them here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader className="flex flex-col sm:flex-row sm:items-start sm:justify-between space-y-2 sm:space-y-0 pb-2">
                    <div className="space-y-1 min-w-0 flex-1">
                      <CardTitle className="text-base sm:text-lg">
                        Generated Test Cases
                      </CardTitle>
                      <CardDescription className="text-xs sm:text-sm break-words">
                        <div className="space-y-1">
                          <div>File: {testCaseResult.file_name}</div>
                          <div>Framework: {testCaseResult.test_framework}</div>
                          <div>Language: {testCaseResult.test_language}</div>
                        </div>
                      </CardDescription>
                      <div className="text-xs text-muted-foreground space-y-1">
                        <div>
                          Generated:{" "}
                          {new Date(
                            testCaseResult.generated_at
                          ).toLocaleString()}
                        </div>
                        <div className="break-all">
                          Source Hash:{" "}
                          <code className="bg-muted px-1 py-0.5 rounded text-xs">
                            {testCaseResult.source_hash.substring(0, 8)}...
                          </code>
                        </div>
                      </div>
                    </div>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => copyToClipboard(testCaseResult.test_code)}
                      className="flex-shrink-0"
                    >
                      <Copy className="h-4 w-4 mr-2" />
                      {copiedCode ? "Copied!" : "Copy"}
                    </Button>
                  </CardHeader>
                  <CardContent className="space-y-4 sm:space-y-6">
                    <div className="relative">
                      <div className="bg-slate-950 rounded-lg p-3 sm:p-4 overflow-x-auto max-h-80 sm:max-h-96 overflow-y-auto">
                        <pre className="text-xs sm:text-sm text-slate-100 font-mono leading-relaxed">
                          <code className="language-javascript break-words whitespace-pre-wrap">
                            {testCaseResult.test_code}
                          </code>
                        </pre>
                      </div>
                      {/* Code language indicator */}
                      <div className="absolute top-2 right-2">
                        <div className="bg-slate-800 text-slate-300 px-2 py-1 rounded text-xs font-mono">
                          {testCaseResult.test_language}
                        </div>
                      </div>
                    </div>

                    {/* Recommendations section */}
                    {testCaseResult.warnings_and_recommendations &&
                      testCaseResult.warnings_and_recommendations.length >
                        0 && (
                        <div className="space-y-3">
                          <div className="flex items-center gap-2 text-sm font-medium text-yellow-600">
                            <AlertTriangle className="h-4 w-4" />
                            Recommendations
                          </div>
                          <div className="bg-yellow-50 dark:bg-yellow-950/20 border border-yellow-200 dark:border-yellow-800/30 rounded-lg p-3 sm:p-4">
                            <ul className="space-y-2">
                              {testCaseResult.warnings_and_recommendations.map(
                                (rec, index) => (
                                  <li
                                    key={index}
                                    className="flex items-start gap-2 text-sm text-yellow-800 dark:text-yellow-200"
                                  >
                                    <AlertTriangle className="h-3 w-3 text-yellow-600 dark:text-yellow-400 mt-0.5 flex-shrink-0" />
                                    <span className="break-words">{rec}</span>
                                  </li>
                                )
                              )}
                            </ul>
                          </div>
                        </div>
                      )}
                  </CardContent>
                </Card>
              </>
            )}
          </TabsContent>
        </div>
      </Tabs>
    </div>
  );
}
