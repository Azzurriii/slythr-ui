export const DEFAULT_CONTRACT = `// SPDX-License-Identifier: MIT
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

export const NOTIFICATION_TIMEOUTS = {
  SUCCESS: 5000,
  ERROR: 8000,
} as const;

export const LOADING_DELAYS = {
  STATIC_ANALYSIS: 2000,
  AI_ANALYSIS: 3000,
  TEST_GENERATION: 2500,
} as const;

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
