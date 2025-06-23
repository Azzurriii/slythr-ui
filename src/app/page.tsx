"use client";

import { useState } from "react";
import { SlythrActions } from "@/components/SlythrActions";
import { CodeEditor } from "@/components/CodeEditor";
import { AnalysisPanel } from "@/components/AnalysisPanel";
import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
} from "@/lib/types";

const DEFAULT_CONTRACT = `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private storedData;
    address public owner;
    
    constructor() {
        owner = msg.sender;
    }
    
    function set(uint256 x) public {
        storedData = x;
    }
    
    function get() public view returns (uint256) {
        return storedData;
    }
    
    function withdraw() public {
        require(msg.sender == owner, "Only owner can withdraw");
        payable(owner).transfer(address(this).balance);
    }
}`;

export default function HomePage() {
  const [sourceCode, setSourceCode] = useState(DEFAULT_CONTRACT);
  const [isLoading, setIsLoading] = useState(false);
  const [analysisType, setAnalysisType] = useState<"static" | "ai" | "tests">(
    "static"
  );
  const [staticAnalysisResult, setStaticAnalysisResult] =
    useState<StaticAnalysisResponse | null>(null);
  const [aiAnalysisResult, setAiAnalysisResult] =
    useState<AIAnalysisResponse | null>(null);
  const [testCaseResult, setTestCaseResult] = useState<TestCaseResponse | null>(
    null
  );
  const [activeTab, setActiveTab] = useState("static");

  const handleAnalyze = async () => {
    setIsLoading(true);
    setAnalysisType("static");
    setActiveTab("static");

    // Simulate API call
    setTimeout(() => {
      const mockStaticResult: StaticAnalysisResponse = {
        total_issues: 3,
        severity_summary: { high: 1, medium: 1, low: 1, informational: 0 },
        issues: [
          {
            title: "Reentrancy vulnerability in withdraw function",
            severity: "HIGH",
            description:
              "The withdraw function is vulnerable to reentrancy attacks. The balance is transferred before the state is updated.",
            location: "Line 18-21",
          },
          {
            title: "Missing access control",
            severity: "MEDIUM",
            description:
              "The set function lacks proper access control, allowing anyone to modify the stored data.",
            location: "Line 12-14",
          },
          {
            title: "Unused variable",
            severity: "LOW",
            description:
              "The storedData variable is declared as private but could be optimized.",
            location: "Line 4",
          },
        ],
      };
      setStaticAnalysisResult(mockStaticResult);
      setIsLoading(false);
    }, 2000);
  };

  const handleAIAnalyze = async () => {
    setIsLoading(true);
    setAnalysisType("ai");
    setActiveTab("ai");

    setTimeout(() => {
      const mockAIResult: AIAnalysisResponse = {
        security_score: 65,
        risk_level: "MEDIUM",
        vulnerabilities: [
          {
            title: "Potential reentrancy attack vector",
            severity: "HIGH",
            description:
              "The contract's withdraw function follows the check-effects-interactions pattern incorrectly, making it vulnerable to reentrancy attacks.",
            recommendation:
              "Use the checks-effects-interactions pattern or implement a reentrancy guard.",
          },
          {
            title: "Centralization risk",
            severity: "MEDIUM",
            description:
              "The contract has a single owner with significant control over funds.",
            recommendation:
              "Consider implementing multi-signature or decentralized governance.",
          },
        ],
      };
      setAiAnalysisResult(mockAIResult);
      setIsLoading(false);
    }, 3000);
  };

  const handleGenerateTests = async (framework: string, language: string) => {
    setIsLoading(true);
    setAnalysisType("tests");
    setActiveTab("tests");

    setTimeout(() => {
      const mockTestResult: TestCaseResponse = {
        test_framework: framework,
        test_language: language,
        test_code: `const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("SimpleStorage", function () {
  let simpleStorage;
  let owner;
  let addr1;

  beforeEach(async function () {
    [owner, addr1] = await ethers.getSigners();
    const SimpleStorage = await ethers.getContractFactory("SimpleStorage");
    simpleStorage = await SimpleStorage.deploy();
    await simpleStorage.deployed();
  });

  describe("Deployment", function () {
    it("Should set the right owner", async function () {
      expect(await simpleStorage.owner()).to.equal(owner.address);
    });
  });

  describe("Storage", function () {
    it("Should store and retrieve values", async function () {
      await simpleStorage.set(42);
      expect(await simpleStorage.get()).to.equal(42);
    });

    it("Should allow anyone to set values", async function () {
      await simpleStorage.connect(addr1).set(100);
      expect(await simpleStorage.get()).to.equal(100);
    });
  });

  describe("Withdrawal", function () {
    it("Should allow owner to withdraw", async function () {
      // Send some ETH to the contract
      await owner.sendTransaction({
        to: simpleStorage.address,
        value: ethers.utils.parseEther("1.0")
      });

      const initialBalance = await owner.getBalance();
      await simpleStorage.withdraw();
      const finalBalance = await owner.getBalance();
      
      expect(finalBalance).to.be.gt(initialBalance);
    });

    it("Should revert when non-owner tries to withdraw", async function () {
      await expect(
        simpleStorage.connect(addr1).withdraw()
      ).to.be.revertedWith("Only owner can withdraw");
    });
  });
});`,
        warnings_and_recommendations: [
          "Consider adding tests for reentrancy attack scenarios",
          "Add edge case tests for zero values and empty states",
          "Test gas optimization scenarios",
          "Include integration tests with other contracts",
        ],
      };
      setTestCaseResult(mockTestResult);
      setIsLoading(false);
    }, 2500);
  };

  return (
    <div className="flex h-[calc(100vh-3.5rem)]">
      {/* Left Column - Actions */}
      <div className="w-80 border-r border-border bg-card/50 p-4 overflow-y-auto">
        <SlythrActions
          isLoading={isLoading}
          onAnalyze={handleAnalyze}
          onAIAnalyze={handleAIAnalyze}
          onGenerateTests={handleGenerateTests}
        />
      </div>

      {/* Middle Column - Code Editor */}
      <div className="flex-1 min-w-0">
        <CodeEditor value={sourceCode} onChange={setSourceCode} />
      </div>

      {/* Right Column - Analysis Panel */}
      <div className="w-96 border-l border-border bg-card/50">
        <AnalysisPanel
          activeTab={activeTab}
          onTabChange={setActiveTab}
          staticAnalysisResult={staticAnalysisResult}
          aiAnalysisResult={aiAnalysisResult}
          testCaseResult={testCaseResult}
          isLoading={isLoading}
        />
      </div>
    </div>
  );
}
