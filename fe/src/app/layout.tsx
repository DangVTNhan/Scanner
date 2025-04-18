import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Link from "next/link";
import { Toaster } from "sonner";
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
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="flex min-h-screen flex-col">
          <header className="sticky top-0 z-50 w-full border-b bg-background">
            <div className="container flex h-16 items-center">
              <h1 className="text-xl font-bold">
                Changi Airport Weather Report System
              </h1>
              <nav className="ml-auto flex gap-4">
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
          <main className="flex-1 container py-6">{children}</main>
          <footer className="border-t py-4">
            <div className="container text-center text-sm text-muted-foreground">
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
