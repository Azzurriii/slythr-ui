"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, Shield, Zap, TestTube, Download } from "lucide-react";
import type { NetworkType } from "@/lib/types";
import { isValidAddress } from "@/lib/utils";

interface SlythrActionsProps {
  isLoading: boolean;
  loadingType?: "static" | "ai" | "tests" | "fetch";
  onAnalyze: () => void;
  onAIAnalyze: () => void;
  onGenerateTests: (framework: string, language: string) => void;
  onFetchContract: (address: string, network: NetworkType) => void;
}

export function SlythrActions({
  isLoading,
  loadingType,
  onAnalyze,
  onAIAnalyze,
  onGenerateTests,
  onFetchContract,
}: SlythrActionsProps) {
  const [testFramework, setTestFramework] = useState("hardhat");
  const [testLanguage, setTestLanguage] = useState("javascript");
  const [contractAddress, setContractAddress] = useState("");
  const [network, setNetwork] = useState<NetworkType>("ethereum");

  const handleGenerateTests = () => {
    onGenerateTests(testFramework, testLanguage);
  };

  const handleFetchContract = () => {
    if (contractAddress && isValidAddress(contractAddress)) {
      onFetchContract(contractAddress, network);
    }
  };

  const isValidContractAddress =
    contractAddress && isValidAddress(contractAddress);

  return (
    <div className="space-y-6">
      {/* Primary Analysis Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-primary" />
            Security Analysis
          </CardTitle>
          <CardDescription>
            Run comprehensive security checks on your smart contract
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <Button
            onClick={onAnalyze}
            disabled={isLoading}
            className="w-full"
            size="lg"
          >
            {isLoading && loadingType === "static" ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Analyzing...
              </>
            ) : (
              <>
                <Shield className="mr-2 h-4 w-4" />
                Static Analysis
              </>
            )}
          </Button>

          <Button
            onClick={onAIAnalyze}
            disabled={isLoading}
            variant="secondary"
            className="w-full"
            size="lg"
          >
            {isLoading && loadingType === "ai" ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                AI Analyzing...
              </>
            ) : (
              <>
                <Zap className="mr-2 h-4 w-4" />
                AI Audit
              </>
            )}
          </Button>
        </CardContent>
      </Card>

      {/* Test Generation */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TestTube className="h-5 w-5 text-primary" />
            Generate Test Cases
          </CardTitle>
          <CardDescription>
            Automatically generate comprehensive test suites
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="framework">Test Framework</Label>
            <Select value={testFramework} onValueChange={setTestFramework}>
              <SelectTrigger>
                <SelectValue placeholder="Select framework" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="hardhat">Hardhat</SelectItem>
                <SelectItem value="foundry">Foundry</SelectItem>
                <SelectItem value="truffle">Truffle</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="language">Test Language</Label>
            <Select value={testLanguage} onValueChange={setTestLanguage}>
              <SelectTrigger>
                <SelectValue placeholder="Select language" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="javascript">JavaScript</SelectItem>
                <SelectItem value="typescript">TypeScript</SelectItem>
                <SelectItem value="solidity">Solidity</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Button
            onClick={handleGenerateTests}
            disabled={isLoading}
            variant="secondary"
            className="w-full"
          >
            {isLoading && loadingType === "tests" ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Generating...
              </>
            ) : (
              <>
                <TestTube className="mr-2 h-4 w-4" />
                Generate Tests
              </>
            )}
          </Button>
        </CardContent>
      </Card>

      {/* Fetch from Blockchain */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Download className="h-5 w-5 text-primary" />
            Fetch from Blockchain
          </CardTitle>
          <CardDescription>
            Import contract source code from verified contracts
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="address">Contract Address</Label>
            <Input
              id="address"
              placeholder="0x..."
              value={contractAddress}
              onChange={(e) => setContractAddress(e.target.value)}
              className={
                !isValidContractAddress && contractAddress
                  ? "border-destructive"
                  : ""
              }
            />
            {contractAddress && !isValidContractAddress && (
              <p className="text-sm text-destructive">
                Please enter a valid contract address
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="network">Network</Label>
            <Select
              value={network}
              onValueChange={(value: NetworkType) => setNetwork(value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select network" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ethereum">Ethereum Mainnet</SelectItem>
                <SelectItem value="polygon">Polygon</SelectItem>
                <SelectItem value="bsc">BSC (Binance Smart Chain)</SelectItem>
                <SelectItem value="base">Base</SelectItem>
                <SelectItem value="arbitrum">Arbitrum</SelectItem>
                <SelectItem value="avalanche">Avalanche</SelectItem>
                <SelectItem value="optimism">Optimism</SelectItem>
                <SelectItem value="gnosis">Gnosis Chain</SelectItem>
                <SelectItem value="fantom">Fantom</SelectItem>
                <SelectItem value="celo">Celo</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Button
            onClick={handleFetchContract}
            disabled={!isValidContractAddress || isLoading}
            variant="outline"
            className="w-full"
          >
            {isLoading && loadingType === "fetch" ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Fetching...
              </>
            ) : (
              <>
                <Download className="mr-2 h-4 w-4" />
                Fetch Contract
              </>
            )}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
