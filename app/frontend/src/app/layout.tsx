import type { Metadata } from "next";
import { Toaster } from "@/components/ui/toaster";
import "./globals.css";

export const metadata: Metadata = {
  title: "CMS",
  description: "git-based CMS",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <main className="flex min-h-screen flex-col items-center max-w-4xl mx-auto">
          {children}
        </main>
        <Toaster />
      </body>
    </html>
  );
}
