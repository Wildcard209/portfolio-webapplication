import type { Metadata } from "next";
import Navbar from "@/app/components/Navbar/Navbar";
import Footer from "@/app/components/Footer/Footer";
import ClientWrapper from "@/app/components/ClientWrapper";
import localFont from "next/font/local";
import "./globals.scss";
import Script from "next/script";

const peanutButter = localFont({
  src: [{ path: "./fonts/Peanut Butter/Peanut-Butter.woff2" }],
  variable: "--font-peanut-butter",
});

const roboto = localFont({
  src: [{ path: "./fonts/Roboto/Roboto-Regular.ttf" }],
  variable: "--font-roboto",
});

const googleAnalyticsId = process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID ?? "";

export const metadata: Metadata = {
  title: "Jessica Wylde",
  description: "Jessica Wyldes portfolio page",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${roboto.variable} ${peanutButter.variable}`}>
<head>

  {/* Google Analytics */}
        <Script
          src="https://www.googletagmanager.com/gtag/js?id=G-36VT3MYPY5"
          strategy="afterInteractive"
        />
        <Script id="gtag-init" strategy="afterInteractive" dangerouslySetInnerHTML={{
          __html: `
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', 'G-36VT3MYPY5');
          `,
        }} />

  <link
    rel="icon"
    type="image/png"
    sizes="32x32"
    href="/assets/favicons/favicon-32x32.png"
  />
  <link
    rel="icon"
    type="image/png"
    sizes="16x16"
    href="/assets/favicons/favicon-16x16.png"
  />
  <meta
    name="viewport"
    content="width=device-width, initial-scale=1, shrink-to-fit=no"
  />
  <meta name="theme-color" content="#ffffff" />
  <meta name="description" content="Jessica Wylde" />
  <meta name="author" content="Jessica Wylde" />
  <link
    rel="apple-touch-icon"
    sizes="180x180"
    href="/assets/favicons/apple-touch-icon.png"
  />
  <link rel="manifest" href="/assets/favicons/site.webmanifest" />
        <title>{`${metadata.title}`}</title>
      </head>
      <body>
        <Navbar />
        <main className="main-content">
          <ClientWrapper>{children}</ClientWrapper>
        </main>
        <Footer />
      </body>
    </html>
  );
}
