import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
};

module.exports = {
  reactStrictMode: true,
  env: {
    API_BASE_URL: "http://localhost:8080/api",
  },
};

export default nextConfig;
