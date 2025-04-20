"use client"

import { useTheme } from "next-themes";
import { Toaster as Sonner, ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "light" } = useTheme();

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      className="toaster group"
      position="top-right"
      expand={false}
      richColors
      closeButton
      style={
        {
          // Base styles
          "--normal-bg": "white",
          "--normal-text": "#334155",
          "--normal-border": "#e2e8f0",

          // Success toast styles
          "--success-bg": "#e6f7ef",
          "--success-text": "#0a7d5c",
          "--success-border": "#bae7d7",

          // Error toast styles
          "--error-bg": "#fef2f2",
          "--error-text": "#b91c1c",
          "--error-border": "#fecaca",

          // Info toast styles
          "--info-bg": "#eff6ff",
          "--info-text": "#1e40af",
          "--info-border": "#bfdbfe",

          // Warning toast styles
          "--warning-bg": "#fffbeb",
          "--warning-text": "#b45309",
          "--warning-border": "#fef3c7",

          // Toast dimensions and styling
          "--border-radius": "0.5rem",
          "--font-family": "var(--font-geist-sans)",
          "--z-index": "9999",
        } as React.CSSProperties
      }
      toastOptions={{
        classNames: {
          toast: "group rounded-lg border shadow-lg",
          title: "font-medium text-sm",
          description: "text-xs text-muted-foreground",
          actionButton:
            "bg-primary text-primary-foreground text-xs px-3 py-1 rounded-md",
          cancelButton:
            "bg-muted text-muted-foreground text-xs px-3 py-1 rounded-md",
          closeButton: "text-foreground/50 hover:text-foreground",
          success:
            "bg-[--success-bg] text-[--success-text] border-[--success-border]",
          error: "bg-[--error-bg] text-[--error-text] border-[--error-border]",
          info: "bg-[--info-bg] text-[--info-text] border-[--info-border]",
          warning:
            "bg-[--warning-bg] text-[--warning-text] border-[--warning-border]",
        },
      }}
      {...props}
    />
  );
};

export { Toaster };

