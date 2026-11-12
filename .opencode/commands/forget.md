---
description: Drop prior conversation context; treat the next message as a new chat
---

Context reset boundary.

Discard everything before this message: earlier requests, plans, decisions, findings, file contents, tool output, and assumptions. Do not summarize, recap, or reference them. Do not continue any earlier task unless the next user message asks for it again.

From here on:

1. Treat the next user message as the first message of a new conversation.
2. Re-derive anything you need from the current repository state, not from memory of this chat.
3. Do not mention this reset unless asked.
4. Reply with exactly one line: `Context cleared. Ready for new task.`

Only context carried forward, if provided:

$ARGUMENTS
