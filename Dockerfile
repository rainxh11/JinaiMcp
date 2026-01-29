# Use Node.js 18 slim image (Debian-based)
FROM node:18-slim

# Install necessary tools and libraries
RUN apt-get update && apt-get install -y \
    chromium \
    libmagic-dev \
    build-essential \
    python3 \
    wget \
    gnupg \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt-get update \
    && apt-get install -y google-chrome-stable \
    && rm -rf /var/lib/apt/lists/*

# Install Python for MCP server
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables for Puppeteer
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/google-chrome-stable

# Set working directory
WORKDIR /app

# Copy Python requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip3 install --no-cache-dir -r requirements.txt

# Copy package files for the reader service
COPY package*.json ./backend/functions/

# Install Node dependencies
RUN cd backend/functions && npm ci

# Copy the reader application code
COPY backend/functions ./backend/functions

# Build the reader application
RUN cd backend/functions && npm run build

# Create local storage directory and set permissions
RUN mkdir -p /app/local-storage && chmod 777 /app/local-storage

# Create directory for MCP server
RUN mkdir -p /app/mcp_server

# Copy MCP server files
COPY mcp_server /app/mcp_server

# Expose the port the reader app runs on
EXPOSE 3000

# Start both the reader app and MCP server using a supervisor script
CMD ["python3", "start.py"]
