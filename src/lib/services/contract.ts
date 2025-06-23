import type { ContractSourceResponse, NetworkType } from "../types";
import { API_BASE_URL } from "../constants";

export class ContractService {
  static async fetchSourceCode(
    address: string,
    network: NetworkType
  ): Promise<ContractSourceResponse> {
    try {
      const response = await fetch(
        `${API_BASE_URL}/contracts/${address}/source-code?network=${network}`,
        {
          method: "GET",
          headers: {
            accept: "application/json",
          },
        }
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(
          errorData.message ||
            `Failed to fetch contract source code: ${response.status} ${response.statusText}`
        );
      }

      const data: ContractSourceResponse = await response.json();
      return data;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error(
        "An unexpected error occurred while fetching contract source code"
      );
    }
  }
}
