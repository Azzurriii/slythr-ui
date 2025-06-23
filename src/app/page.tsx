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
    } catch (error) {
      // Error is already handled in the hook
      console.error("Error fetching contract:", error);
    }
  };

  return (
    <div className="flex h-[calc(100vh-3.5rem)]">
      {/* Left Column - Actions */}
      <div className="w-80 border-r border-border bg-card/50 p-4 overflow-y-auto">
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
      <div className="flex-1 min-w-0">
        <CodeEditor value={sourceCode} onChange={setSourceCode} />
      </div>

      {/* Right Column - Analysis Panel */}
      <div className="w-96 border-l border-border bg-card/50">
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
