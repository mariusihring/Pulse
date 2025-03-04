import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'cryptologos.cc',
                port: '',
                pathname: '/**',
            },
        ],
    },
};

export default nextConfig;
