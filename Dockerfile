# Use browserless/chrome image which already has Chrome installed
FROM browserless/chrome:latest

# Switch to root user for all operations (needed for Puppeteer)
USER root

# Install dependencies as root
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Install Deno as root
RUN curl -fsSL https://deno.land/install.sh | sh
ENV DENO_INSTALL="/root/.deno"
ENV PATH="${DENO_INSTALL}/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy Python requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip3 install --no-cache-dir -r requirements.txt

# Copy Deno config and source files
COPY deno /app/deno

# Cache Deno dependencies by running the type check
# This downloads all npm packages without executing the code
RUN deno cache --allow-net --allow-read --allow-write --allow-env --allow-run --allow-sys /app/deno/main.ts || true

# Copy MCP server files
COPY mcp_server /app/mcp_server

# Copy startup script
COPY start.py .

# Create local storage directory with proper permissions
RUN mkdir -p /app/local-storage/instant-screenshots && \
    chmod 777 /app/local-storage && \
    chmod 777 /app/local-storage/instant-screenshots

# Expose ports
# 3000 - Deno/Reader service
# 8000 - MCP server
EXPOSE 3000 8000

# Start both services as root (required for Puppeteer)
CMD ["python3", "start.py"]
