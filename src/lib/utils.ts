import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isValidAddress(address: string): boolean {
  return /^0x[a-fA-F0-9]{40}$/.test(address);
}

/**
 * Generate a SHA-256 hash of a source code string, equivalent to:
 * Go: sha256.Sum256 + hex.EncodeToString
 */
export async function generateSourceHash(sourceCode: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(sourceCode);
  const hashBuffer = await crypto.subtle.digest("SHA-256", data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}
