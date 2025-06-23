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
  isLoading,
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
    <div className="h-full flex flex-col">
      <div className="border-b border-border bg-muted/50 px-4 py-2">
        <h2 className="text-sm font-medium">Analysis Results</h2>
      </div>

      <Tabs
        value={activeTab}
        onValueChange={onTabChange}
        className="flex-1 flex flex-col"
      >
        <TabsList className="grid w-full grid-cols-3 m-4 mb-0">
          <TabsTrigger value="static" className="flex items-center gap-2">
            <Shield className="h-4 w-4" />
            Static
          </TabsTrigger>
          <TabsTrigger value="ai" className="flex items-center gap-2">
            <Zap className="h-4 w-4" />
            AI Audit
          </TabsTrigger>
          <TabsTrigger value="tests" className="flex items-center gap-2">
            <TestTube className="h-4 w-4" />
            Tests
          </TabsTrigger>
        </TabsList>

        <div className="flex-1 overflow-y-auto p-4 pt-2">
          <TabsContent value="static" className="mt-2 space-y-4">
            {!staticAnalysisResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-8">
                  <Shield className="h-12 w-12 text-muted-foreground mb-4" />
                  <p className="text-center text-muted-foreground">
                    Run a static analysis to see the results here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader>
                    <CardTitle>Slither Analysis Summary</CardTitle>
                    <CardDescription>
                      Found {staticAnalysisResult.total_issues} total issues
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-wrap gap-2">
                      {staticAnalysisResult.severity_summary.high > 0 && (
                        <Badge className={getSeverityColor("HIGH")}>
                          {staticAnalysisResult.severity_summary.high} High
                        </Badge>
                      )}
                      {staticAnalysisResult.severity_summary.medium > 0 && (
                        <Badge className={getSeverityColor("MEDIUM")}>
                          {staticAnalysisResult.severity_summary.medium} Medium
                        </Badge>
                      )}
                      {staticAnalysisResult.severity_summary.low > 0 && (
                        <Badge className={getSeverityColor("LOW")}>
                          {staticAnalysisResult.severity_summary.low} Low
                        </Badge>
                      )}
                      {staticAnalysisResult.severity_summary.informational >
                        0 && (
                        <Badge className={getSeverityColor("INFORMATIONAL")}>
                          {staticAnalysisResult.severity_summary.informational}{" "}
                          Info
                        </Badge>
                      )}
                    </div>
                  </CardContent>
                </Card>

                <Accordion type="single" collapsible className="space-y-2">
                  {staticAnalysisResult.issues.map((issue, index) => (
                    <AccordionItem
                      key={index}
                      value={`item-${index}`}
                      className="border rounded-lg px-4"
                    >
                      <AccordionTrigger className="hover:no-underline">
                        <div className="flex items-center justify-between w-full mr-4">
                          <span className="text-left">{issue.title}</span>
                          <Badge className={getSeverityColor(issue.severity)}>
                            {issue.severity}
                          </Badge>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent className="space-y-2">
                        <p className="text-sm text-muted-foreground">
                          {issue.description}
                        </p>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <span>Location:</span>
                          <code className="bg-muted px-1 py-0.5 rounded">
                            {issue.location}
                          </code>
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  ))}
                </Accordion>
              </>
            )}
          </TabsContent>

          <TabsContent value="ai" className="mt-2 space-y-4">
            {!aiAnalysisResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-8">
                  <Zap className="h-12 w-12 text-muted-foreground mb-4" />
                  <p className="text-center text-muted-foreground">
                    Run an AI audit to see the results here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader>
                    <CardTitle>AI Security Assessment</CardTitle>
                    <CardDescription>
                      Overall risk level: {aiAnalysisResult.risk_level}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <div className="flex justify-between text-sm mb-2">
                        <span>Security Score</span>
                        <span>{aiAnalysisResult.security_score}/100</span>
                      </div>
                      <Progress
                        value={aiAnalysisResult.security_score}
                        className="h-2"
                      />
                    </div>
                  </CardContent>
                </Card>

                <Accordion type="single" collapsible className="space-y-2">
                  {aiAnalysisResult.vulnerabilities.map((vuln, index) => (
                    <AccordionItem
                      key={index}
                      value={`ai-item-${index}`}
                      className="border rounded-lg px-4"
                    >
                      <AccordionTrigger className="hover:no-underline">
                        <div className="flex items-center justify-between w-full mr-4">
                          <span className="text-left">{vuln.title}</span>
                          <Badge className={getSeverityColor(vuln.severity)}>
                            {vuln.severity}
                          </Badge>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent className="space-y-3">
                        <p className="text-sm text-muted-foreground">
                          {vuln.description}
                        </p>
                        <Alert>
                          <CheckCircle className="h-4 w-4" />
                          <AlertDescription className="text-sm">
                            <strong>Recommendation:</strong>{" "}
                            {vuln.recommendation}
                          </AlertDescription>
                        </Alert>
                      </AccordionContent>
                    </AccordionItem>
                  ))}
                </Accordion>
              </>
            )}
          </TabsContent>

          <TabsContent value="tests" className="mt-2 space-y-4">
            {!testCaseResult ? (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-8">
                  <TestTube className="h-12 w-12 text-muted-foreground mb-4" />
                  <p className="text-center text-muted-foreground">
                    Generate test cases to see them here.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <div>
                      <CardTitle>Generated Test Suite</CardTitle>
                      <CardDescription>
                        {testCaseResult.test_framework} •{" "}
                        {testCaseResult.test_language}
                      </CardDescription>
                    </div>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => copyToClipboard(testCaseResult.test_code)}
                    >
                      <Copy className="h-4 w-4 mr-2" />
                      {copiedCode ? "Copied!" : "Copy"}
                    </Button>
                  </CardHeader>
                  <CardContent>
                    <div className="bg-muted/50 rounded-lg p-4 overflow-x-auto">
                      <pre className="text-sm">
                        <code>{testCaseResult.test_code}</code>
                      </pre>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <AlertTriangle className="h-5 w-5 text-yellow-500" />
                      Recommendations
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <ul className="space-y-2">
                      {testCaseResult.warnings_and_recommendations.map(
                        (rec, index) => (
                          <li
                            key={index}
                            className="flex items-start gap-2 text-sm"
                          >
                            <AlertTriangle className="h-4 w-4 text-yellow-500 mt-0.5 flex-shrink-0" />
                            {rec}
                          </li>
                        )
                      )}
                    </ul>
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
