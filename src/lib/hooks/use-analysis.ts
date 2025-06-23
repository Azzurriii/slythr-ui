import { useState, useEffect } from "react";
import { toast } from "sonner";
import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
  NetworkType,
} from "../types";
import { AnalysisService } from "../services/analysis";
import { ContractService } from "../services/contract";
import { generateSourceHash } from "../utils";

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
  const [currentSourceHash, setCurrentSourceHash] = useState<string | null>(
    null
  );

  // Load cached results on mount if URL has hash params
  useEffect(() => {
    const loadCachedResults = async () => {
      if (typeof window === "undefined") return;

      const searchParams = new URLSearchParams(window.location.search);
      const hash = searchParams.get("hash");
      const tab = searchParams.get("tab");

      // Restore active tab
      if (tab && ["static", "ai", "tests"].includes(tab)) {
        setActiveTab(tab);
      }

      // Load cached analysis results if hash exists
      if (hash) {
        setCurrentSourceHash(hash);

        try {
          // Try to load all analysis results with the same hash
          const [staticResult, aiResult, testResult] = await Promise.allSettled(
            [
              AnalysisService.getCachedStaticAnalysis(hash),
              AnalysisService.getCachedAIAnalysis(hash),
              AnalysisService.getCachedTestCases(hash),
            ]
          );

          if (staticResult.status === "fulfilled") {
            setStaticAnalysisResult(staticResult.value);
          }

          if (aiResult.status === "fulfilled") {
            setAiAnalysisResult(aiResult.value);
          }

          if (testResult.status === "fulfilled") {
            setTestCaseResult(testResult.value);
          }
        } catch (error) {
          console.error("⚠️ Failed to restore analysis results:", error);
        }
      }
    };

    loadCachedResults();
  }, []);

  // Update URL when analysis results change
  const updateURL = (hash?: string, tab?: string) => {
    if (typeof window === "undefined") return;

    const searchParams = new URLSearchParams(window.location.search);

    if (hash) {
      searchParams.set("hash", hash);
      setCurrentSourceHash(hash);
    }
    if (tab) {
      searchParams.set("tab", tab);
    }

    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.replaceState(null, "", newUrl);
  };

  const handleStaticAnalysis = async (sourceCode: string) => {
    setIsLoading(true);
    setLoadingType("static");
    setActiveTab("static");

    try {
      const result = await AnalysisService.performStaticAnalysis(sourceCode);
      setStaticAnalysisResult(result);

      // Update URL with hash for persistence
      updateURL(result.source_hash, "static");

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

      // Update URL with hash for persistence
      updateURL(result.source_hash, "ai");

      toast.success("AI analysis completed successfully!", {
        description: `Security score: ${result.analysis.security_score}/100 (${result.analysis.risk_level} risk)`,
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

      // Update URL with hash for persistence
      updateURL(result.source_hash, "tests");

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

    // Clear URL params when clearing results
    if (typeof window !== "undefined") {
      const newUrl = window.location.pathname;
      window.history.replaceState(null, "", newUrl);
    }
  };

  // Enhanced setActiveTab that also updates URL
  const setActiveTabWithURL = (tab: string) => {
    setActiveTab(tab);
    updateURL(undefined, tab);
  };

  // Check if we can instantly load results for current source code
  const checkCachedResults = async (sourceCode: string) => {
    const hash = await generateSourceHash(sourceCode);

    if (hash === currentSourceHash) {
      // Already have results for this source code
      return;
    }

    // Try to load cached results for this new source code
    setCurrentSourceHash(hash);
    updateURL(hash);

    try {
      const [staticResult, aiResult, testResult] = await Promise.allSettled([
        AnalysisService.getCachedStaticAnalysis(hash),
        AnalysisService.getCachedAIAnalysis(hash),
        AnalysisService.getCachedTestCases(hash),
      ]);

      if (staticResult.status === "fulfilled") {
        setStaticAnalysisResult(staticResult.value);
      }

      if (aiResult.status === "fulfilled") {
        setAiAnalysisResult(aiResult.value);
      }

      if (testResult.status === "fulfilled") {
        setTestCaseResult(testResult.value);
      }
    } catch (error) {
      console.error(`Error checking cached results: ${error}`);
    }
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
    setActiveTab: setActiveTabWithURL,
    checkCachedResults,
  };
}
