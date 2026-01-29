# Use browserless/chrome image which already has Chrome installed
FROM browserless/chrome:latest

# Switch to root user for all operations (needed for Puppeteer)
USER root

# Install dependencies as root
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install Deno as root
RUN curl -fsSL https://deno.land/install.sh | sh
ENV DENO_INSTALL="/root/.deno"
ENV PATH="${DENO_INSTALL}/bin:${PATH}"

# Install Go
RUN wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz && \
    rm go1.21.5.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum* ./

# Download Go dependencies
RUN go mod download || true

# Copy Go source files
COPY *.go ./

# Build Go MCP server
RUN go build -o mcp-server .

# Copy Deno config and source files
COPY deno /app/deno

# Cache Deno dependencies
RUN deno cache --allow-net --allow-read --allow-write --allow-env --allow-run --allow-sys /app/deno/main.ts || true

# Create local storage directory with proper permissions
RUN mkdir -p /app/local-storage/instant-screenshots && \
    chmod 777 /app/local-storage && \
    chmod 777 /app/local-storage/instant-screenshots

# Expose ports
# 3000 - Deno/Reader service
# 8000 - Go MCP server
EXPOSE 3000 8000

# Start both services
CMD ["sh", "-c", "deno run --allow-net --allow-read --allow-write --allow-env --allow-run --allow-sys /app/deno/main.ts & ./mcp-server"]
