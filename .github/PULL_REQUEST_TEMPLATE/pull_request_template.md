## Summary

<!-- What changes for users or contributors? Keep it under 100 words. -->

## Why

<!-- What problem does this solve? If it fixes an open issue, link it with “Fixes #123”. -->

## How

<!-- How did you solve it? Describe a test you ran or the steps you tried, and what happened. -->

## Other details

<!-- Add screenshots, new dependencies, or other details a reviewer needs. Leave this empty if it doesn’t apply. Extra headings are welcome. -->

## Pre-submission checklist

<!-- Complete the applicable checks before requesting review. You may remove this checklist afterward. -->

- [ ] I have read the [AI policy](https://superfile.dev/contribute/ai-policy/).
- [ ] I tested the change and described the result under **How**.
- [ ] I read the full diff and removed leftover debug logs and TODOs.
- [ ] I filled out the template and included any screenshots needed to review the change.
- [ ] My PR title follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

If you changed Go code, also run these checks.

- [ ] I formatted the code with `go fmt ./...`.
- [ ] I ran `golangci-lint run` and fixed the issues it reported.
