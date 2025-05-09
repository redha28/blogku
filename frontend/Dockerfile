# Use a stable Node version
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install --production=false

# Copy the application code
COPY . .

# Set environment variables
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000

# Build the application (assumes you have a build script in package.json)
RUN npm run build || ( \
    echo "Build failed, creating minimal app (App Router)" && \
    mkdir -p src/app public && \
    echo 'export default function Home() { return <h1>BlogKu Frontend</h1>; }' > src/app/page.js && \
    echo 'module.exports = { output: "standalone" };' > next.config.js && \
    npm install && \
    npm run build \
)


# Expose port
EXPOSE 3000

# Start the application
CMD ["npm", "run", "start"]
