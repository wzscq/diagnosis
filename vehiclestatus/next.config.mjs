/** @type {import('next').NextConfig} */
const nextConfig = {
    basePath: '/nextservice',
    output: 'standalone',
    env: {
        PORT:"8400",
    },
};

export default nextConfig;
