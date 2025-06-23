import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
} from "./types";

export const mockStaticAnalysisResult: StaticAnalysisResponse = {
  success: true,
  issues: [
    {
      type: "detector",
      title: "Reentrancy vulnerability in withdraw function",
      severity: "HIGH",
      confidence: "High",
      description:
        "The withdraw function is vulnerable to reentrancy attacks. The balance is transferred before the state is updated.",
      location: "Contract.sol:L18-L21",
      reference: "Contract.sol#L18-L21",
    },
    {
      type: "detector",
      title: "Missing access control",
      severity: "MEDIUM",
      confidence: "Medium",
      description:
        "The set function lacks proper access control, allowing anyone to modify the stored data.",
      location: "Contract.sol:L12-L14",
      reference: "Contract.sol#L12-L14",
    },
    {
      type: "detector",
      title: "Unused variable",
      severity: "LOW",
      confidence: "Low",
      description:
        "The storedData variable is declared as private but could be optimized.",
      location: "Contract.sol:L4",
      reference: "Contract.sol#L4",
    },
  ],
  total_issues: 3,
  severity_summary: { high: 1, medium: 1, low: 1, informational: 0 },
  analyzed_at: new Date().toISOString(),
  source_hash: "mock_hash_static_123456789",
};

export const mockAIAnalysisResult: AIAnalysisResponse = {
  success: true,
  analysis: {
    security_score: 65,
    risk_level: "MEDIUM",
    summary:
      "The contract has several security issues, including reentrancy and missing access control.",
    vulnerabilities: [
      {
        title: "Potential reentrancy attack vector",
        severity: "HIGH",
        description:
          "The contract's withdraw function follows the check-effects-interactions pattern incorrectly, making it vulnerable to reentrancy attacks.",
        location: {
          function: "withdraw()",
          line_numbers: [18, 21],
        },
        recommendation:
          "Use the checks-effects-interactions pattern or implement a reentrancy guard.",
      },
      {
        title: "Centralization risk",
        severity: "MEDIUM",
        description:
          "The contract has a single owner with significant control over funds.",
        location: {
          function: "constructor()",
          line_numbers: [6, 8],
        },
        recommendation:
          "Consider implementing multi-signature or decentralized governance.",
      },
    ],
    good_practices: [
      "Uses latest Solidity version with proper license identifier",
      "Implements proper owner access control pattern",
    ],
    recommendations: [
      "Add event logging for important state changes",
      "Consider implementing emergency stop functionality",
      "Add input validation for critical functions",
    ],
  },
  total_issues: 2,
  analyzed_at: new Date().toISOString(),
  source_hash: "mock_hash_123456789",
};

export const generateMockTestResult = (
  framework: string,
  language: string
): TestCaseResponse => ({
  success: true,
  test_framework: framework,
  test_language: language,
  file_name: `SimpleStorage.test.${language === "typescript" ? "ts" : "js"}`,
  source_hash: "mock_hash_test_123456789",
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
  generated_at: new Date().toISOString(),
});
