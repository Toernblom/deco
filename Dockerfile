FROM golang:1.23-bookworm

# Install tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    jq \
    nodejs \
    npm \
    sudo \
    openssh-server \
    && rm -rf /var/lib/apt/lists/*

# Install Claude Code globally
RUN npm install -g @anthropic-ai/claude-code

# Install beads CLI
RUN curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash

# Create non-root user with password for SSH
RUN useradd -m -s /bin/bash claude && \
    echo "claude:claude" | chpasswd && \
    echo "claude ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Setup SSH server with agent forwarding
RUN mkdir -p /var/run/sshd && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PermitUserEnvironment no/PermitUserEnvironment yes/' /etc/ssh/sshd_config && \
    sed -i 's/#AllowAgentForwarding yes/AllowAgentForwarding yes/' /etc/ssh/sshd_config

# Create necessary directories with proper ownership
RUN mkdir -p /home/claude/.claude-config /home/claude/.config/git /home/claude/.ssh && \
    chown -R claude:claude /home/claude

# Set Claude config directory (used by Claude Code for all config/auth)
ENV CLAUDE_CONFIG_DIR=/home/claude/.claude-config

# Switch to claude user for remaining setup
USER claude
WORKDIR /home/claude

# Set up git credential caching (uses XDG path, persisted via volume)
RUN git config --global credential.helper 'store --file=/home/claude/.config/git/credentials' && \
    git config --global init.defaultBranch main && \
    git config --global core.autocrlf input

# Note: Beads plugin is installed via /plugin commands after first auth
# Run: /plugin marketplace add steveyegge/beads
# Then: /plugin install beads@beads-marketplace

# Add convenience aliases and environment to bashrc
RUN echo 'alias c="claude --dangerously-skip-permissions"' >> /home/claude/.bashrc && \
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> /home/claude/.bashrc && \
    echo 'cd /workspace 2>/dev/null || true' >> /home/claude/.bashrc

# Also add to .profile for SSH sessions
RUN echo 'alias c="claude --dangerously-skip-permissions"' >> /home/claude/.profile && \
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> /home/claude/.profile && \
    echo 'cd /workspace' >> /home/claude/.profile

# Switch back to root for entrypoint
USER root

# Set workspace
WORKDIR /workspace

# Copy entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 22

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/sbin/sshd", "-D"]
