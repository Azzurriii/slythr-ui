import type {
  StaticAnalysisResponse,
  AIAnalysisResponse,
  TestCaseResponse,
} from "../types";
import {
  mockStaticAnalysisResult,
  mockAIAnalysisResult,
  generateMockTestResult,
} from "../mock-data";
import { LOADING_DELAYS } from "../constants";

export class AnalysisService {
  /**
   * Performs static analysis on the provided source code
   */
  static async performStaticAnalysis(
    sourceCode: string
  ): Promise<StaticAnalysisResponse> {
    // Simulate API call delay
    await new Promise((resolve) =>
      setTimeout(resolve, LOADING_DELAYS.STATIC_ANALYSIS)
    );

    // In a real implementation, this would make an API call
    // return await fetch('/api/analysis/static', { method: 'POST', body: JSON.stringify({ sourceCode }) })

    return mockStaticAnalysisResult;
  }

  /**
   * Performs AI-powered analysis on the provided source code
   */
  static async performAIAnalysis(
    sourceCode: string
  ): Promise<AIAnalysisResponse> {
    // Simulate API call delay
    await new Promise((resolve) =>
      setTimeout(resolve, LOADING_DELAYS.AI_ANALYSIS)
    );

    // In a real implementation, this would make an API call
    // return await fetch('/api/analysis/ai', { method: 'POST', body: JSON.stringify({ sourceCode }) })

    return mockAIAnalysisResult;
  }

  /**
   * Generates test cases for the provided source code
   */
  static async generateTestCases(
    sourceCode: string,
    framework: string,
    language: string
  ): Promise<TestCaseResponse> {
    // Simulate API call delay
    await new Promise((resolve) =>
      setTimeout(resolve, LOADING_DELAYS.TEST_GENERATION)
    );

    // In a real implementation, this would make an API call
    // return await fetch('/api/analysis/tests', {
    //   method: 'POST',
    //   body: JSON.stringify({ sourceCode, framework, language })
    // })

    return generateMockTestResult(framework, language);
  }
}
