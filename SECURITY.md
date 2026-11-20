# Security Policy

## Supported versions

Security fixes go into the latest release of superfile. Please update to the [latest release](https://github.com/yorukot/superfile/releases/latest) before you report a problem.

| Version          | Supported |
| ---------------- | --------- |
| Latest release   | Yes       |
| Older releases   | No        |

## Reporting a vulnerability

Please do not report security problems in public issues, pull requests or discussions.

Report them privately through GitHub instead: [open a security advisory](https://github.com/yorukot/superfile/security/advisories/new).

Include as much of this as you can:

- The superfile version and your operating system. `spf --debug-info` prints both.
- What the problem is and what an attacker could do with it.
- Steps to reproduce it, such as a sample file, archive or config.

We will reply to your report, keep you updated while we work on a fix, and credit you in the advisory unless you ask us not to.

## Scope

Examples of what we treat as a security problem:

- Extracting an archive writes files outside the destination folder.
- Previewing or selecting a file runs code or commands you didn't ask for.
- A file operation changes or deletes files other than the ones you chose.
- The install scripts on superfile.dev can be made to run something unexpected.

Commands you run yourself in the prompt, and programs you set in your config (such as `editor` or `open_with`), are not security problems: superfile runs them on purpose.
