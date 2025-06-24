import type React from "react";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Image from "next/image";
import "./globals.css";
import { cn } from "@/lib/utils";
import { ToastProvider } from "@/components/providers/toast-provider";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Slythr - Smart Contract Security Analysis",
  description:
    "Advanced security analysis platform for Solidity smart contracts",
  icons: {
    icon: "/favicon.svg",
    shortcut: "/favicon.svg",
    apple: "/favicon.svg",
  },
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
          "min-h-screen antialiased overflow-hidden"
        )}
      >
        {/* Animated Background */}
        <div className="fixed inset-0 -z-10">
          {/* Base gradient */}
          <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-purple-900/20 to-slate-900" />

          {/* Animated orbs */}
          <div className="absolute top-0 -left-4 w-72 h-72 bg-purple-500/30 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob" />
          <div className="absolute top-0 -right-4 w-72 h-72 bg-cyan-500/30 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000" />
          <div className="absolute -bottom-8 left-20 w-72 h-72 bg-pink-500/30 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000" />

          {/* Grid pattern */}
          <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:50px_50px]" />

          {/* Noise texture */}
          <div className='absolute inset-0 opacity-[0.015] bg-[url(&apos;data:image/svg+xml,%3Csvg width="60" height="60" viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg"%3E%3Cg fill="none" fillRule="evenodd"%3E%3Cg fill="%23ffffff" fillOpacity="1"%3E%3Ccircle cx="7" cy="7" r="1"/%3E%3Ccircle cx="27" cy="7" r="1"/%3E%3Ccircle cx="47" cy="7" r="1"/%3E%3Ccircle cx="7" cy="27" r="1"/%3E%3Ccircle cx="27" cy="27" r="1"/%3E%3Ccircle cx="47" cy="27" r="1"/%3E%3Ccircle cx="7" cy="47" r="1"/%3E%3Ccircle cx="27" cy="47" r="1"/%3E%3Ccircle cx="47" cy="47" r="1"/%3E%3C/g%3E%3C/g%3E%3C/svg%3E&apos;)]' />
        </div>

        <div className="flex flex-col min-h-screen relative">
          {/* Glassmorphism Header */}
          <header className="relative border-b border-white/10 bg-white/5 backdrop-blur-xl supports-[backdrop-filter]:bg-white/5">
            {/* Header glow effect */}
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-500/10 via-transparent to-purple-500/10" />

            <div className="container flex h-16 items-center justify-between px-6 relative z-10">
              <div className="flex items-center space-x-3">
                {/* Logo with glass effect */}
                <div className="relative flex items-center justify-center p-1">
                  <div className="absolute inset-0 bg-gradient-to-br from-cyan-400/20 to-purple-400/20 rounded-xl blur-sm" />
                  <div className="relative bg-white/10 backdrop-blur-sm border border-white/20 rounded-xl p-2">
                    <Image
                      src="/favicon.svg"
                      alt="Slythr Logo"
                      width={32}
                      height={32}
                      className="h-8 w-8 drop-shadow-lg"
                    />
                  </div>
                </div>

                <div className="flex flex-col">
                  <span className="text-2xl font-bold bg-gradient-to-r from-[#6AD7E5] to-[#9D7AEA] bg-clip-text text-transparent tracking-tight">
                    Slythr
                  </span>
                  <span className="text-xs text-white/60 -mt-1">
                    Smart Contract Security
                  </span>
                </div>
              </div>
            </div>

            {/* Bottom border glow */}
            <div className="absolute bottom-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-cyan-500/50 to-transparent" />
          </header>

          {/* Main content with glass container */}
          <main className="flex-1 relative">
            <div className="absolute inset-0 bg-black/20 backdrop-blur-sm" />
            <div className="relative z-10 h-full">{children}</div>
          </main>
        </div>

        <ToastProvider />
      </body>
    </html>
  );
}
