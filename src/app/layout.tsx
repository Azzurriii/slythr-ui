import type React from "react";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import { ToastProvider } from "@/components/providers/toast-provider";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Slythr - Smart Contract Security Analysis",
  description:
    "Advanced security analysis platform for Solidity smart contracts",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body
        className={cn(
          inter.className,
          "min-h-screen bg-background antialiased"
        )}
      >
        <div className="flex flex-col min-h-screen">
          <header className="border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="container flex h-14 items-center">
              <div className="flex items-center space-x-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                  <span className="text-sm font-bold text-primary-foreground">
                    S
                  </span>
                </div>
                <span className="text-xl font-bold">Slythr</span>
              </div>
              <div className="ml-auto flex items-center space-x-2">
                <span className="text-sm text-muted-foreground">
                  Smart Contract Security Analysis
                </span>
              </div>
            </div>
          </header>
          <main className="flex-1">{children}</main>
        </div>
        <ToastProvider />
      </body>
    </html>
  );
}
