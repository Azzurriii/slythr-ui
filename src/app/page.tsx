"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import {
  Loader2,
  Search,
  Code,
  AlertTriangle,
  CheckCircle,
  XCircle,
  ChevronDown,
  FileText,
  Shield,
  TrendingUp,
  Info,
  Copy,
  ExternalLink,
} from "lucide-react";

interface ContractInfo {
  name: string;
  compiler: string;
  optimization: boolean;
  runs: number;
  evmVersion: string;
  sourceCode: string;
  abi: string;
}

interface SlitherDetector {
  check: string;
  impact: "Critical" | "High" | "Medium" | "Low" | "Informational";
  confidence: "High" | "Medium" | "Low";
  description: string;
  elements: Array<{
    type: string;
    name: string;
    source_mapping: {
      start: number;
      length: number;
      filename_relative: string;
      filename_absolute: string;
      lines: number[];
    };
  }>;
  additional_fields?: {
    underlying_type?: string;
    variable_name?: string;
  };
}

interface SlitherResults {
  success: boolean;
  error: string | null;
  results: {
    detectors: SlitherDetector[];
    printers: any[];
  };
  version: string;
}

interface AnalysisStats {
  totalIssues: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  informational: number;
  gasOptimizations: number;
}

const mockContractInfo: ContractInfo = {
  name: "SimpleToken",
  compiler: "0.8.19+commit.7dd6d404",
  optimization: true,
  runs: 200,
  evmVersion: "default",
  sourceCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract SimpleToken {
    mapping(address => uint256) public balances;
    mapping(address => mapping(address => uint256)) public allowances;
    
    uint256 public totalSupply;
    string public name = "SimpleToken";
    string public symbol = "STK";
    uint8 public decimals = 18;
    
    address public owner;
    bool private locked;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    constructor(uint256 _totalSupply) {
        totalSupply = _totalSupply;
        balances[msg.sender] = _totalSupply;
        owner = msg.sender;
    }
    
    function transfer(address to, uint256 amount) public returns (bool) {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        
        // Vulnerable to reentrancy
        if (to.code.length > 0) {
            (bool success,) = to.call("");
            require(success, "Transfer failed");
        }
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
        
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    function withdraw() public {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");
        
        // Vulnerable: external call before state change
        (bool success,) = msg.sender.call{value: amount}("");
        require(success, "Withdrawal failed");
        
        balances[msg.sender] = 0;
    }
    
    function updateConfig(uint256 newSupply) public {
        // Missing access control
        totalSupply = newSupply;
    }
    
    // Gas inefficient loop
    function distributeTokens(address[] memory recipients, uint256 amount) public {
        for (uint256 i = 0; i < recipients.length; i++) {
            balances[recipients[i]] += amount;
        }
    }
}`,
  abi: `[{"inputs":[{"internalType":"uint256","name":"_totalSupply","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"}]`,
};

const mockSlitherResults: SlitherResults = {
  success: true,
  error: null,
  results: {
    detectors: [
      {
        check: "reentrancy-eth",
        impact: "Critical",
        confidence: "High",
        description:
          "Reentrancy in SimpleToken.withdraw() (contracts/SimpleToken.sol#25-32):\n\tExternal calls:\n\t- (success) = msg.sender.call{value: amount}() (contracts/SimpleToken.sol#29)\n\tState variables written after the call(s):\n\t- balances[msg.sender] = 0 (contracts/SimpleToken.sol#31)",
        elements: [
          {
            type: "function",
            name: "withdraw",
            source_mapping: {
              start: 1205,
              length: 185,
              filename_relative: "contracts/SimpleToken.sol",
              filename_absolute: "/home/user/contracts/SimpleToken.sol",
              lines: [25, 26, 27, 28, 29, 30, 31, 32],
            },
          },
        ],
      },
      {
        check: "missing-zero-check",
        impact: "Low",
        confidence: "Medium",
        description:
          "SimpleToken.constructor(uint256)._totalSupply (contracts/SimpleToken.sol#15) lacks a zero-check on :\n\t\t- totalSupply = _totalSupply (contracts/SimpleToken.sol#16)",
        elements: [
          {
            type: "parameter",
            name: "_totalSupply",
            source_mapping: {
              start: 890,
              length: 18,
              filename_relative: "contracts/SimpleToken.sol",
              filename_absolute: "/home/user/contracts/SimpleToken.sol",
              lines: [15],
            },
          },
        ],
      },
      {
        check: "unprotected-upgrade",
        impact: "High",
        confidence: "High",
        description:
          "SimpleToken.updateConfig(uint256) (contracts/SimpleToken.sol#35-37) should be protected:\n\t- State variables written:\n\t\t- totalSupply (contracts/SimpleToken.sol#36)",
        elements: [
          {
            type: "function",
            name: "updateConfig",
            source_mapping: {
              start: 1450,
              length: 95,
              filename_relative: "contracts/SimpleToken.sol",
              filename_absolute: "/home/user/contracts/SimpleToken.sol",
              lines: [35, 36, 37],
            },
          },
        ],
      },
      {
        check: "costly-loop",
        impact: "Informational",
        confidence: "Medium",
        description:
          "SimpleToken.distributeTokens(address[],uint256) (contracts/SimpleToken.sol#40-44) has costly operations inside a loop:\n\t- balances[recipients[i]] += amount (contracts/SimpleToken.sol#42)",
        elements: [
          {
            type: "function",
            name: "distributeTokens",
            source_mapping: {
              start: 1580,
              length: 165,
              filename_relative: "contracts/SimpleToken.sol",
              filename_absolute: "/home/user/contracts/SimpleToken.sol",
              lines: [40, 41, 42, 43, 44],
            },
          },
        ],
      },
      {
        check: "pragma",
        impact: "Informational",
        confidence: "High",
        description:
          "Different versions of Solidity are used:\n\t- Version used: ['^0.8.19']\n\t- ^0.8.19 (contracts/SimpleToken.sol#2)",
        elements: [
          {
            type: "pragma",
            name: "0.8.19",
            source_mapping: {
              start: 32,
              length: 23,
              filename_relative: "contracts/SimpleToken.sol",
              filename_absolute: "/home/user/contracts/SimpleToken.sol",
              lines: [2],
            },
          },
        ],
      },
    ],
    printers: [],
  },
  version: "0.9.6",
};

export default function SlitherAnalyzer() {
  const [contractAddress, setContractAddress] = useState("");
  const [solidityCode, setSolidityCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingContract, setIsFetchingContract] = useState(false);
  const [contractInfo, setContractInfo] = useState<ContractInfo | null>(null);
  const [slitherResults, setSlitherResults] = useState<SlitherResults | null>(
    null
  );
  const [error, setError] = useState("");
  const [inputMethod, setInputMethod] = useState<"address" | "code">("address");
  const [analysisProgress, setAnalysisProgress] = useState(0);

  const validateEthereumAddress = (address: string): boolean => {
    return /^0x[a-fA-F0-9]{40}$/.test(address);
  };

  const validateSolidityCode = (code: string): boolean => {
    return code.includes("pragma solidity") || code.includes("contract ");
  };

  const fetchContractSource = async () => {
    if (!validateEthereumAddress(contractAddress)) {
      setError("Please enter a valid Ethereum address");
      return;
    }

    setIsFetchingContract(true);
    setError("");

    try {
      // Simulate API call to Etherscan
      await new Promise((resolve) => setTimeout(resolve, 2000));
      setContractInfo(mockContractInfo);
      setSolidityCode(mockContractInfo.sourceCode);
    } catch (err) {
      setError("Failed to fetch contract source code");
    } finally {
      setIsFetchingContract(false);
    }
  };

  const handleAnalyze = async () => {
    setError("");
    setSlitherResults(null);
    setAnalysisProgress(0);

    if (inputMethod === "address" && !contractInfo) {
      setError("Please fetch contract source code first");
      return;
    }

    if (inputMethod === "code") {
      if (!solidityCode.trim()) {
        setError("Please enter Solidity code");
        return;
      }
      if (!validateSolidityCode(solidityCode)) {
        setError("Please enter valid Solidity code");
        return;
      }
    }

    setIsLoading(true);

    try {
      // Simulate analysis progress
      const progressSteps = [
        { step: "Parsing contract...", progress: 20 },
        { step: "Running detectors...", progress: 50 },
        { step: "Analyzing vulnerabilities...", progress: 80 },
        { step: "Generating report...", progress: 100 },
      ];

      for (const { progress } of progressSteps) {
        await new Promise((resolve) => setTimeout(resolve, 800));
        setAnalysisProgress(progress);
      }

      setSlitherResults(mockSlitherResults);
    } catch (err) {
      setError("Analysis failed. Please try again.");
    } finally {
      setIsLoading(false);
      setAnalysisProgress(0);
    }
  };

  const getImpactColor = (impact: string) => {
    switch (impact.toLowerCase()) {
      case "critical":
        return "bg-red-100 text-red-800 border-red-200";
      case "high":
        return "bg-orange-100 text-orange-800 border-orange-200";
      case "medium":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "low":
        return "bg-blue-100 text-blue-800 border-blue-200";
      case "informational":
        return "bg-gray-100 text-gray-800 border-gray-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const getImpactIcon = (impact: string) => {
    switch (impact.toLowerCase()) {
      case "critical":
        return <XCircle className="h-4 w-4 text-red-600" />;
      case "high":
        return <AlertTriangle className="h-4 w-4 text-orange-600" />;
      case "medium":
        return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
      case "low":
        return <Info className="h-4 w-4 text-blue-600" />;
      case "informational":
        return <CheckCircle className="h-4 w-4 text-gray-600" />;
      default:
        return <CheckCircle className="h-4 w-4 text-gray-600" />;
    }
  };

  const getConfidenceColor = (confidence: string) => {
    switch (confidence.toLowerCase()) {
      case "high":
        return "bg-green-100 text-green-800";
      case "medium":
        return "bg-yellow-100 text-yellow-800";
      case "low":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const calculateStats = (): AnalysisStats => {
    if (!slitherResults) {
      return {
        totalIssues: 0,
        critical: 0,
        high: 0,
        medium: 0,
        low: 0,
        informational: 0,
        gasOptimizations: 0,
      };
    }

    const stats = slitherResults.results.detectors.reduce(
      (acc, detector) => {
        acc.totalIssues++;
        switch (detector.impact.toLowerCase()) {
          case "critical":
            acc.critical++;
            break;
          case "high":
            acc.high++;
            break;
          case "medium":
            acc.medium++;
            break;
          case "low":
            acc.low++;
            break;
          case "informational":
            acc.informational++;
            if (
              detector.check.includes("gas") ||
              detector.check.includes("costly")
            ) {
              acc.gasOptimizations++;
            }
            break;
        }
        return acc;
      },
      {
        totalIssues: 0,
        critical: 0,
        high: 0,
        medium: 0,
        low: 0,
        informational: 0,
        gasOptimizations: 0,
      }
    );

    return stats;
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const stats = calculateStats();

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center gap-3 mb-4">
            <Shield className="h-8 w-8 text-blue-600" />
            <h1 className="text-4xl font-bold text-gray-900">
              Slither Security Analyzer
            </h1>
          </div>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Advanced static analysis for Solidity smart contracts. Detect
            vulnerabilities, optimize gas usage, and ensure security best
            practices.
          </p>
        </div>

        {/* Input Section */}
        <Card className="mb-8 shadow-lg">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-xl">
              <Code className="h-6 w-6 text-blue-600" />
              Contract Input
            </CardTitle>
            <CardDescription className="text-base">
              Enter a verified contract address to fetch from Etherscan or paste
              Solidity code directly
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 p-6">
            <Tabs
              value={inputMethod}
              onValueChange={(value) =>
                setInputMethod(value as "address" | "code")
              }
            >
              <TabsList className="grid w-full grid-cols-2 h-12">
                <TabsTrigger value="address" className="text-base">
                  Contract Address
                </TabsTrigger>
                <TabsTrigger value="code" className="text-base">
                  Solidity Code
                </TabsTrigger>
              </TabsList>

              <TabsContent value="address" className="space-y-6">
                <div className="flex gap-3">
                  <Input
                    placeholder="0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
                    value={contractAddress}
                    onChange={(e) => setContractAddress(e.target.value)}
                    className="flex-1 h-12 text-base"
                  />
                  <Button
                    onClick={fetchContractSource}
                    disabled={isFetchingContract || !contractAddress}
                    className="px-8 h-12"
                    variant="outline"
                  >
                    {isFetchingContract ? (
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    ) : (
                      <Search className="h-4 w-4 mr-2" />
                    )}
                    Fetch Source
                  </Button>
                </div>

                {contractInfo && (
                  <Card className="bg-green-50 border-green-200">
                    <CardHeader className="pb-3">
                      <CardTitle className="text-lg text-green-800 flex items-center gap-2">
                        <CheckCircle className="h-5 w-5" />
                        Contract Information
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div>
                          <p className="text-sm font-medium text-gray-600">
                            Name
                          </p>
                          <p className="text-base font-semibold">
                            {contractInfo.name}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm font-medium text-gray-600">
                            Compiler
                          </p>
                          <p className="text-base">{contractInfo.compiler}</p>
                        </div>
                        <div>
                          <p className="text-sm font-medium text-gray-600">
                            Optimization
                          </p>
                          <p className="text-base">
                            {contractInfo.optimization ? "Enabled" : "Disabled"}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm font-medium text-gray-600">
                            Runs
                          </p>
                          <p className="text-base">{contractInfo.runs}</p>
                        </div>
                      </div>
                      <Button
                        onClick={handleAnalyze}
                        disabled={isLoading}
                        className="w-full h-12 text-base"
                      >
                        {isLoading ? (
                          <Loader2 className="h-5 w-5 animate-spin mr-2" />
                        ) : (
                          <Shield className="h-5 w-5 mr-2" />
                        )}
                        Analyze Contract Security
                      </Button>
                    </CardContent>
                  </Card>
                )}
              </TabsContent>

              <TabsContent value="code" className="space-y-6">
                <div className="space-y-4">
                  <Textarea
                    placeholder="// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MyContract {
    // Paste your Solidity code here
    mapping(address => uint256) public balances;
    
    function transfer(address to, uint256 amount) public {
        // Your contract logic
    }
}"
                    value={solidityCode}
                    onChange={(e) => setSolidityCode(e.target.value)}
                    className="min-h-[300px] font-mono text-sm resize-none"
                  />
                  <div className="flex justify-between items-center">
                    <p className="text-sm text-gray-500">
                      {solidityCode.split("\n").length} lines,{" "}
                      {solidityCode.length} characters
                    </p>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => copyToClipboard(solidityCode)}
                      disabled={!solidityCode}
                    >
                      <Copy className="h-4 w-4 mr-2" />
                      Copy
                    </Button>
                  </div>
                </div>
                <Button
                  onClick={handleAnalyze}
                  disabled={isLoading}
                  className="w-full h-12 text-base"
                >
                  {isLoading ? (
                    <Loader2 className="h-5 w-5 animate-spin mr-2" />
                  ) : (
                    <Shield className="h-5 w-5 mr-2" />
                  )}
                  Analyze Contract Security
                </Button>
              </TabsContent>
            </Tabs>

            {error && (
              <Alert variant="destructive">
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription className="text-base">
                  {error}
                </AlertDescription>
              </Alert>
            )}
          </CardContent>
        </Card>

        {/* Contract Source Code Display */}
        {contractInfo && inputMethod === "address" && (
          <Card className="mb-8 shadow-lg">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-xl">
                <FileText className="h-6 w-6 text-gray-600" />
                Contract Source Code
              </CardTitle>
              <CardDescription className="text-base">
                Verified source code fetched from blockchain explorer
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <div className="relative">
                <pre className="bg-gray-900 text-gray-100 p-6 rounded-b-lg overflow-x-auto text-sm leading-relaxed">
                  <code>{contractInfo.sourceCode}</code>
                </pre>
                <Button
                  variant="outline"
                  size="sm"
                  className="absolute top-4 right-4 bg-gray-800 border-gray-600 text-gray-200 hover:bg-gray-700"
                  onClick={() => copyToClipboard(contractInfo.sourceCode)}
                >
                  <Copy className="h-4 w-4 mr-2" />
                  Copy Code
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Loading State */}
        {isLoading && (
          <Card className="mb-8 shadow-lg">
            <CardContent className="flex flex-col items-center justify-center py-16">
              <div className="text-center space-y-6 max-w-md">
                <div className="relative">
                  <Shield className="h-16 w-16 mx-auto text-blue-600 animate-pulse" />
                  <div className="absolute inset-0 rounded-full border-4 border-blue-200 border-t-blue-600 animate-spin"></div>
                </div>
                <div className="space-y-3">
                  <h3 className="text-xl font-semibold text-gray-900">
                    Analyzing Contract Security
                  </h3>
                  <p className="text-gray-600">
                    Running comprehensive security analysis...
                  </p>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className="bg-blue-600 h-3 rounded-full transition-all duration-500 ease-out"
                      style={{ width: `${analysisProgress}%` }}
                    ></div>
                  </div>
                  <p className="text-sm text-gray-500">
                    {analysisProgress}% complete
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Analysis Results */}
        {slitherResults && !isLoading && (
          <div className="space-y-8">
            {/* Statistics Overview */}
            <Card className="shadow-lg">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <TrendingUp className="h-6 w-6 text-blue-600" />
                  Analysis Overview
                </CardTitle>
                <CardDescription className="text-base">
                  Security analysis results powered by Slither v
                  {slitherResults.version}
                </CardDescription>
              </CardHeader>
              <CardContent className="p-6">
                <div className="grid grid-cols-2 md:grid-cols-6 gap-6">
                  <div className="text-center">
                    <div className="text-3xl font-bold text-gray-900">
                      {stats.totalIssues}
                    </div>
                    <div className="text-sm text-gray-600">Total Issues</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-red-600">
                      {stats.critical}
                    </div>
                    <div className="text-sm text-gray-600">Critical</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-orange-600">
                      {stats.high}
                    </div>
                    <div className="text-sm text-gray-600">High</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-yellow-600">
                      {stats.medium}
                    </div>
                    <div className="text-sm text-gray-600">Medium</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-blue-600">
                      {stats.low}
                    </div>
                    <div className="text-sm text-gray-600">Low</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-gray-600">
                      {stats.informational}
                    </div>
                    <div className="text-sm text-gray-600">Info</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Detailed Results */}
            <Card className="shadow-lg">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <AlertTriangle className="h-6 w-6 text-red-600" />
                  Security Issues Detected
                </CardTitle>
                <CardDescription className="text-base">
                  {slitherResults.results.detectors.length} issues found by
                  Slither detectors
                </CardDescription>
              </CardHeader>
              <CardContent className="p-6">
                <div className="space-y-4">
                  {slitherResults.results.detectors.map((detector, index) => (
                    <Collapsible key={index}>
                      <CollapsibleTrigger asChild>
                        <div className="flex items-start justify-between p-6 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors">
                          <div className="flex items-start gap-4 flex-1">
                            {getImpactIcon(detector.impact)}
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-3 mb-2">
                                <h3 className="font-semibold text-lg text-gray-900">
                                  {detector.check}
                                </h3>
                                <Badge
                                  className={getImpactColor(detector.impact)}
                                >
                                  {detector.impact.toUpperCase()}
                                </Badge>
                                <Badge
                                  className={getConfidenceColor(
                                    detector.confidence
                                  )}
                                  variant="outline"
                                >
                                  {detector.confidence} Confidence
                                </Badge>
                              </div>
                              <p className="text-gray-600 text-base leading-relaxed">
                                {detector.description.split("\n")[0]}
                              </p>
                              {detector.elements.length > 0 && (
                                <p className="text-sm text-gray-500 mt-2">
                                  📍{" "}
                                  {
                                    detector.elements[0].source_mapping
                                      .filename_relative
                                  }
                                  (Lines:{" "}
                                  {detector.elements[0].source_mapping.lines.join(
                                    ", "
                                  )}
                                  )
                                </p>
                              )}
                            </div>
                          </div>
                          <ChevronDown className="h-5 w-5 text-gray-400 ml-4 flex-shrink-0" />
                        </div>
                      </CollapsibleTrigger>
                      <CollapsibleContent>
                        <div className="px-6 pb-6 space-y-6 border-l-4 border-gray-200 ml-8">
                          <div>
                            <h4 className="font-semibold text-base mb-3 text-gray-900">
                              Detailed Description
                            </h4>
                            <pre className="text-sm text-gray-700 bg-gray-50 p-4 rounded-lg whitespace-pre-wrap font-mono leading-relaxed">
                              {detector.description}
                            </pre>
                          </div>

                          {detector.elements.length > 0 && (
                            <div>
                              <h4 className="font-semibold text-base mb-3 text-gray-900">
                                Affected Elements
                              </h4>
                              <div className="space-y-3">
                                {detector.elements.map((element, elemIndex) => (
                                  <div
                                    key={elemIndex}
                                    className="bg-blue-50 p-4 rounded-lg border border-blue-200"
                                  >
                                    <div className="flex items-center gap-2 mb-2">
                                      <Badge
                                        variant="outline"
                                        className="bg-blue-100 text-blue-800"
                                      >
                                        {element.type}
                                      </Badge>
                                      <span className="font-medium text-blue-900">
                                        {element.name}
                                      </span>
                                    </div>
                                    <div className="text-sm text-blue-700 space-y-1">
                                      <p>
                                        📁 File:{" "}
                                        {
                                          element.source_mapping
                                            .filename_relative
                                        }
                                      </p>
                                      <p>
                                        📍 Lines:{" "}
                                        {element.source_mapping.lines.join(
                                          ", "
                                        )}
                                      </p>
                                      <p>
                                        📏 Position:{" "}
                                        {element.source_mapping.start}-
                                        {element.source_mapping.start +
                                          element.source_mapping.length}
                                      </p>
                                    </div>
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}

                          <div className="flex gap-3">
                            <Button variant="outline" size="sm">
                              <ExternalLink className="h-4 w-4 mr-2" />
                              View Documentation
                            </Button>
                            <Button variant="outline" size="sm">
                              <Copy className="h-4 w-4 mr-2" />
                              Copy Details
                            </Button>
                          </div>
                        </div>
                      </CollapsibleContent>
                    </Collapsible>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Recommendations */}
            <Card className="shadow-lg">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <CheckCircle className="h-6 w-6 text-green-600" />
                  Dynamic Analysis by LLMs
                </CardTitle>
                <CardDescription className="text-base">
                  Actionable steps to improve your contract security
                </CardDescription>
              </CardHeader>
              <CardContent className="p-6">
                <div className="space-y-6">
                  <div className="grid md:grid-cols-2 gap-6">
                    <div className="space-y-4">
                      <h3 className="font-semibold text-lg text-gray-900">
                        🚨 Critical Actions
                      </h3>
                      <ul className="space-y-3">
                        <li className="flex items-start gap-3">
                          <span className="text-gray-700">
                            Implement reentrancy guards for all external calls
                          </span>
                        </li>
                        <li className="flex items-start gap-3">
                          <span className="text-gray-700">
                            Add access control modifiers to sensitive functions
                          </span>
                        </li>
                      </ul>
                    </div>
                    <div className="space-y-4">
                      <h3 className="font-semibold text-lg text-gray-900">
                        ⚡ Optimizations
                      </h3>
                      <ul className="space-y-3">
                        <li className="flex items-start gap-3">
                          <span className="text-gray-700">
                            Replace loops with mappings for gas efficiency
                          </span>
                        </li>
                        <li className="flex items-start gap-3">
                          <span className="text-gray-700">
                            Use events for better transparency
                          </span>
                        </li>
                      </ul>
                    </div>
                  </div>

                  <Separator />

                  <div className="bg-amber-50 p-6 rounded-lg border border-amber-200">
                    <h3 className="font-semibold text-amber-800 mb-3 text-lg">
                      ⚠️ Security Score
                    </h3>
                    <div className="flex items-center gap-4">
                      <div className="flex-1">
                        <Progress value={25} className="h-3" />
                      </div>
                      <span className="font-bold text-amber-800">25/100</span>
                    </div>
                    <p className="text-amber-700 mt-3">
                      <strong>HIGH RISK:</strong> This contract has critical
                      security vulnerabilities that must be addressed before
                      deployment.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
