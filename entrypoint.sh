#!/bin/bash
set -e

# Fix ownership of mounted volumes
chown -R claude:claude /home/claude/.claude-config 2>/dev/null || true
chown -R claude:claude /home/claude/.config 2>/dev/null || true
chown claude:claude /workspace 2>/dev/null || true

# Set git safe directory for workspace (as claude user)
su - claude -c "git config --global --add safe.directory /workspace" 2>/dev/null || true

# Set git identity from environment if provided
if [ -n "$GIT_AUTHOR_NAME" ]; then
    su - claude -c "git config --global user.name '$GIT_AUTHOR_NAME'"
fi
if [ -n "$GIT_AUTHOR_EMAIL" ]; then
    su - claude -c "git config --global user.email '$GIT_AUTHOR_EMAIL'"
fi

# Initialize beads hooks if this is a beads project
if [ -f "/workspace/.beads/issues.jsonl" ]; then
    su - claude -c "cd /workspace && bd doctor --fix" 2>/dev/null || true
fi

# Create MOTD for first-run instructions
cat > /etc/motd << 'EOF'
==================================
FIRST RUN SETUP (one-time only)
==================================
1. Run 'claude' and authenticate
2. Install plugins (inside claude):
   /plugin marketplace add steveyegge/beads
   /plugin install beads@beads-marketplace
   /plugin marketplace add obra/superpowers-marketplace
   /plugin install superpowers@superpowers-marketplace
==================================

EOF

# Execute the command (sshd by default)
exec "$@"
