const withMDX = require("@next/mdx")();

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  rewrites: async () => {
    return [
      {
        source: "/api/v1/:path*",
        destination: "http://localhost:8080/v1/:path*",
      },
    ];
  },
};

module.exports = withMDX(nextConfig);
