import { useState } from "react";
import { toast } from "sonner";
import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
  ContractSourceResponse,
  NetworkType,
} from "../types";
import { AnalysisService } from "../services/analysis";
import { ContractService } from "../services/contract";
import { NOTIFICATION_TIMEOUTS } from "../constants";

type LoadingType = "static" | "ai" | "tests" | "fetch" | undefined;

export function useAnalysis() {
  const [isLoading, setIsLoading] = useState(false);
  const [loadingType, setLoadingType] = useState<LoadingType>(undefined);
  const [staticAnalysisResult, setStaticAnalysisResult] =
    useState<StaticAnalysisResponse | null>(null);
  const [aiAnalysisResult, setAiAnalysisResult] =
    useState<AIAnalysisResponse | null>(null);
  const [testCaseResult, setTestCaseResult] = useState<TestCaseResponse | null>(
    null
  );
  const [activeTab, setActiveTab] = useState("static");

  const handleStaticAnalysis = async (sourceCode: string) => {
    setIsLoading(true);
    setLoadingType("static");
    setActiveTab("static");

    try {
      const result = await AnalysisService.performStaticAnalysis(sourceCode);
      setStaticAnalysisResult(result);
      toast.success("Static analysis completed successfully!", {
        description: `Found ${result.total_issues} security issues`,
      });
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Static analysis failed";
      toast.error("Static analysis failed", {
        description: errorMessage,
      });
    } finally {
      setIsLoading(false);
      setLoadingType(undefined);
    }
  };

  const handleAIAnalysis = async (sourceCode: string) => {
    setIsLoading(true);
    setLoadingType("ai");
    setActiveTab("ai");

    try {
      const result = await AnalysisService.performAIAnalysis(sourceCode);
      setAiAnalysisResult(result);
      toast.success("AI analysis completed successfully!", {
        description: `Security score: ${result.security_score}/100 (${result.risk_level} risk)`,
      });
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "AI analysis failed";
      toast.error("AI analysis failed", {
        description: errorMessage,
      });
    } finally {
      setIsLoading(false);
      setLoadingType(undefined);
    }
  };

  const handleGenerateTests = async (
    sourceCode: string,
    framework: string,
    language: string
  ) => {
    setIsLoading(true);
    setLoadingType("tests");
    setActiveTab("tests");

    try {
      const result = await AnalysisService.generateTestCases(
        sourceCode,
        framework,
        language
      );
      setTestCaseResult(result);
      toast.success("Test cases generated successfully!", {
        description: `Generated tests using ${framework} with ${language}`,
      });
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Test generation failed";
      toast.error("Test generation failed", {
        description: errorMessage,
      });
    } finally {
      setIsLoading(false);
      setLoadingType(undefined);
    }
  };

  const handleFetchContract = async (address: string, network: NetworkType) => {
    setIsLoading(true);
    setLoadingType("fetch");

    try {
      const contractData = await ContractService.fetchSourceCode(
        address,
        network
      );

      // Clear previous analysis results when fetching new contract
      clearAnalysisResults();
      setActiveTab("static");

      toast.success("Contract fetched successfully!", {
        description: `Loaded contract ${address.slice(0, 8)}...${address.slice(
          -6
        )} from ${network}`,
      });

      return contractData.source_code;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to fetch contract";
      toast.error("Failed to fetch contract", {
        description: errorMessage,
      });
      throw err;
    } finally {
      setIsLoading(false);
      setLoadingType(undefined);
    }
  };

  const clearAnalysisResults = () => {
    setStaticAnalysisResult(null);
    setAiAnalysisResult(null);
    setTestCaseResult(null);
  };

  return {
    // State
    isLoading,
    loadingType,
    staticAnalysisResult,
    aiAnalysisResult,
    testCaseResult,
    activeTab,

    // Actions
    handleStaticAnalysis,
    handleAIAnalysis,
    handleGenerateTests,
    handleFetchContract,
    clearAnalysisResults,
    setActiveTab,
  };
}
