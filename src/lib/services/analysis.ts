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
import { LOADING_DELAYS, API_BASE_URL } from "../constants";
import { generateSourceHash } from "../utils";

export class AnalysisService {
  /**
   * Performs static analysis with GET-first caching approach
   */
  static async performStaticAnalysis(
    sourceCode: string
  ): Promise<StaticAnalysisResponse> {
    try {
      // Calculate source hash for caching
      const sourceHash = await generateSourceHash(sourceCode);

      // Try GET first to check if analysis already exists
      const getUrl = `${API_BASE_URL}/static-analysis/${sourceHash}`;

      const getCacheResponse = await fetch(getUrl, {
        method: "GET",
        headers: {
          accept: "application/json",
        },
      });

      if (getCacheResponse.ok) {
        const cachedData: StaticAnalysisResponse =
          await getCacheResponse.json();
        return cachedData;
      }

      // If not cached, perform new analysis with POST
      const postUrl = `${API_BASE_URL}/static-analysis/`;
      const requestBody = {
        source_code: sourceCode,
      };

      const postResponse = await fetch(postUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          accept: "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (!postResponse.ok) {
        throw new Error(`HTTP error! status: ${postResponse.status}`);
      }

      const data: StaticAnalysisResponse = await postResponse.json();

      if (!data.success) {
        throw new Error("Static analysis failed");
      }

      return data;
    } catch (error) {
      console.error("Static Analysis API Error:", error);
      await new Promise((resolve) =>
        setTimeout(resolve, LOADING_DELAYS.STATIC_ANALYSIS)
      );
      return mockStaticAnalysisResult;
    }
  }

  /**
   * Performs AI analysis with GET-first caching approach
   */
  static async performAIAnalysis(
    sourceCode: string
  ): Promise<AIAnalysisResponse> {
    try {
      // Calculate source hash for caching
      const sourceHash = await generateSourceHash(sourceCode);

      // Try GET first to check if analysis already exists
      const getUrl = `${API_BASE_URL}/dynamic-analysis/${sourceHash}`;

      const getCacheResponse = await fetch(getUrl, {
        method: "GET",
        headers: {
          accept: "application/json",
        },
      });

      if (getCacheResponse.ok) {
        const cachedData: AIAnalysisResponse = await getCacheResponse.json();
        return cachedData;
      }

      // If not cached, perform new analysis with POST
      const postUrl = `${API_BASE_URL}/dynamic-analysis/`;
      const requestBody = {
        source_code: sourceCode,
      };

      const postResponse = await fetch(postUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          accept: "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (!postResponse.ok) {
        throw new Error(`HTTP error! status: ${postResponse.status}`);
      }

      const data: AIAnalysisResponse = await postResponse.json();

      if (!data.success) {
        throw new Error("Analysis failed");
      }

      return data;
    } catch (error) {
      console.error("AI Analysis API Error:", error);
      await new Promise((resolve) =>
        setTimeout(resolve, LOADING_DELAYS.AI_ANALYSIS)
      );
      return mockAIAnalysisResult;
    }
  }

  /**
   * Generates test cases with GET-first caching approach
   */
  static async generateTestCases(
    sourceCode: string,
    framework: string,
    language: string
  ): Promise<TestCaseResponse> {
    try {
      // Calculate source hash for caching
      const sourceHash = await generateSourceHash(sourceCode);

      // Try GET first to check if test cases already exist
      const getUrl = `${API_BASE_URL}/test-cases/${sourceHash}`;

      const getCacheResponse = await fetch(getUrl, {
        method: "GET",
        headers: {
          accept: "application/json",
        },
      });

      if (getCacheResponse.ok) {
        const cachedData: TestCaseResponse = await getCacheResponse.json();
        return cachedData;
      }

      // If not cached, generate new test cases with POST
      const postUrl = `${API_BASE_URL}/test-cases/generate`;
      const requestBody = {
        source_code: sourceCode,
        test_framework: framework,
        test_language: language,
      };

      const postResponse = await fetch(postUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          accept: "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (!postResponse.ok) {
        throw new Error(`HTTP error! status: ${postResponse.status}`);
      }

      const data: TestCaseResponse = await postResponse.json();

      if (!data.success) {
        throw new Error("Test case generation failed");
      }

      return data;
    } catch (error) {
      console.error("Test Generation API Error:", error);
      await new Promise((resolve) =>
        setTimeout(resolve, LOADING_DELAYS.TEST_GENERATION)
      );
      return generateMockTestResult(framework, language);
    }
  }

  /**
   * Get cached static analysis result by source hash
   */
  static async getCachedStaticAnalysis(
    sourceHash: string
  ): Promise<StaticAnalysisResponse> {
    const url = `${API_BASE_URL}/static-analysis/${sourceHash}`;

    const response = await fetch(url, {
      method: "GET",
      headers: {
        accept: "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: StaticAnalysisResponse = await response.json();

    if (!data.success) {
      throw new Error("Cached static analysis not found");
    }

    return data;
  }

  /**
   * Get cached AI analysis result by source hash
   */
  static async getCachedAIAnalysis(
    sourceHash: string
  ): Promise<AIAnalysisResponse> {
    const url = `${API_BASE_URL}/dynamic-analysis/${sourceHash}`;

    const response = await fetch(url, {
      method: "GET",
      headers: {
        accept: "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: AIAnalysisResponse = await response.json();

    if (!data.success) {
      throw new Error("Cached AI analysis not found");
    }

    return data;
  }

  /**
   * Get cached test cases result by source hash
   */
  static async getCachedTestCases(
    sourceHash: string
  ): Promise<TestCaseResponse> {
    const url = `${API_BASE_URL}/test-cases/${sourceHash}`;

    const response = await fetch(url, {
      method: "GET",
      headers: {
        accept: "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: TestCaseResponse = await response.json();

    if (!data.success) {
      throw new Error("Cached test cases not found");
    }

    return data;
  }
}
