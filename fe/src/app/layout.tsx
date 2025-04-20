import { Toaster } from "@/components/ui/sonner";
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Link from "next/link";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Changi Airport Weather Report System",
  description:
    "A system for generating and comparing weather reports for Changi Airport",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="flex min-h-screen flex-col">
          <header className="sticky top-0 z-50 w-full border-b bg-background mb-4">
            <div className="container mx-auto flex justify-between items-center h-16 px-4 py-2">
              <h1 className="text-xl font-bold text-blue-500">
                Weather Report System
              </h1>
              <nav className="flex space-x-4 text-blue-500">
                <Link href="/" className="text-sm font-medium hover:underline">
                  Generate Report
                </Link>
                <Link
                  href="/history"
                  className="text-sm font-medium hover:underline"
                >
                  History
                </Link>
              </nav>
            </div>
          </header>
          <main className="flex-1 py-6 px-4">{children}</main>
          <footer className="border-t py-4">
            <div className="text-center text-sm text-muted-foreground">
              &copy; {new Date().getFullYear()} Changi Airport Weather Report
              System
            </div>
          </footer>
        </div>
        <Toaster />
      </body>
    </html>
  );
}
