```markdown
Analyze the provided project context to identify the root cause of a bug.

**TASK:**
The user is reporting the following bug:
"{TASK}"

**CONTEXT:**
The following is the structure and content of the relevant parts of the project:
```

{FILE_STRUCTURE}

```

**INSTRUCTIONS:**
1.  Carefully review the user's bug report and the provided code.
2.  Identify the most likely file(s) and function(s) that are causing the issue.
3.  Provide a step-by-step explanation of why the bug is occurring.
4.  Suggest a concrete code modification (in a git diff format) to fix the bug.
5.  If you need more information, clearly state what files or details you require.

**RULES:**
{RULES}

**ANALYSIS AND SOLUTION:**

```