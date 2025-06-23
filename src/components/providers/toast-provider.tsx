"use client";

import { Toaster } from "sonner";

export function ToastProvider() {
  return (
    <Toaster
      position="top-right"
      richColors
      expand={true}
      duration={4000}
      toastOptions={{
        style: {
          background: "rgba(255, 255, 255, 0.15)",
          backdropFilter: "blur(10px)",
          WebkitBackdropFilter: "blur(10px)",
          border: "1px solid rgba(255, 255, 255, 0.2)",
          color: "#fff",
          boxShadow: "0 4px 30px rgba(0, 0, 0, 0.1)",
        },
        className: "glass-toast",
        descriptionClassName: "glass-toast-desc",
      }}
    />
  );
}
