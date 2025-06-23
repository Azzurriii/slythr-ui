"use client";

import { useState } from "react";
import { SlythrActions, CodeEditor, AnalysisPanel } from "@/components";
import { useAnalysis, DEFAULT_CONTRACT, type NetworkType } from "@/lib";

export default function HomePage() {
  const [sourceCode, setSourceCode] = useState(DEFAULT_CONTRACT);
  const {
    isLoading,
    loadingType,
    staticAnalysisResult,
    aiAnalysisResult,
    testCaseResult,
    activeTab,
    handleStaticAnalysis,
    handleAIAnalysis,
    handleGenerateTests,
    handleFetchContract,
    setActiveTab,
    checkCachedResults,
  } = useAnalysis();

  // Handler wrapper functions to connect source code with analysis hook
  const handleAnalyze = () => handleStaticAnalysis(sourceCode);
  const handleAIAnalyze = () => handleAIAnalysis(sourceCode);
  const handleGenerateTestsWrapper = (framework: string, language: string) => {
    handleGenerateTests(sourceCode, framework, language);
  };

  const handleFetchContractWrapper = async (
    address: string,
    network: NetworkType
  ) => {
    try {
      const newSourceCode = await handleFetchContract(address, network);
      setSourceCode(newSourceCode);
      // Check for cached results with the new source code
      checkCachedResults(newSourceCode);
    } catch (error) {
      // Error is already handled in the hook
      console.error("Error fetching contract:", error);
    }
  };

  // Check for cached results when source code changes
  const handleSourceCodeChange = (newSourceCode: string) => {
    setSourceCode(newSourceCode);
    checkCachedResults(newSourceCode);
  };

  return (
    <div className="flex flex-col lg:flex-row h-[calc(100vh-3.5rem)] min-h-0">
      {/* Left Column - Actions */}
      <div className="w-full lg:w-72 xl:w-80 border-b lg:border-b-0 lg:border-r border-border bg-card/50 p-3 sm:p-4 overflow-y-auto flex-shrink-0 max-h-[40vh] lg:max-h-none">
        <SlythrActions
          isLoading={isLoading}
          loadingType={loadingType}
          onAnalyze={handleAnalyze}
          onAIAnalyze={handleAIAnalyze}
          onGenerateTests={handleGenerateTestsWrapper}
          onFetchContract={handleFetchContractWrapper}
        />
      </div>

      {/* Middle Column - Code Editor */}
      <div className="flex-1 min-w-0 min-h-0">
        <CodeEditor value={sourceCode} onChange={handleSourceCodeChange} />
      </div>

      {/* Right Column - Analysis Panel */}
      <div className="w-full lg:w-[32rem] xl:w-[36rem] 2xl:w-[40rem] border-t lg:border-t-0 lg:border-l border-border bg-card/50 flex-shrink-0 min-h-0">
        <AnalysisPanel
          isLoading={isLoading}
          staticAnalysisResult={staticAnalysisResult}
          aiAnalysisResult={aiAnalysisResult}
          testCaseResult={testCaseResult}
          activeTab={activeTab}
          onTabChange={setActiveTab}
        />
      </div>
    </div>
  );
}
