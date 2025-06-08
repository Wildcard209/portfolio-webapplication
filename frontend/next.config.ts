import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
};

module.exports = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_BASE_API_URL: process.env.BASE_API_URL,
    NEXT_PUBLIC_ALLOWED_ORIGIN: process.env.ALLOWED_ORIGIN,
  },
};

export default nextConfig;
