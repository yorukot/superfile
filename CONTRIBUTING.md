# Contributing to superfile

Welcome to **superfile**! This guide will help you get started contributing to the project, whether you're fixing bugs, building features, or just sharing ideas.

There are many ways to contribute:

* Reporting bugs
* Fixing issues
* Adding a theme
* Suggesting and implementing new features
* Sharing ideas or feedback

---

## Issues

### Found a bug?

For a security vulnerability, follow the [security policy](SECURITY.md) to report it privately.

For other bugs, check if there's already an open or closed issue for it. If not, open a new one and describe the problem clearly.

### Want to fix an issue?

1. Fork this repository
2. Create a new branch for the issue you're working on
3. Commit your changes with clear messages
4. Open a pull request (PR) with a description of the problem and your solution

Maintainers may request changes before merging.

---

## Adding a Theme

Before starting, make sure the theme you want to add doesn’t already exist.

1. Copy an existing theme's `.toml` file as a base
2. Customize it to your needs
3. Test it by editing your `~/.config/superfile/config/config.toml`
4. When ready, submit a pull request
5. To ensure the theme looks consistent and functions properly, please include the following screenshots in your PR:
- Full view of superfile (Including sidebar, file previewer, process panel, metadata panel, and clipboard panel )
    - Make sure that file previewer is non empty, process panel has at least one process, and clipboard has at least one entry
- Add a screenshot of these individual panel being focused (To make sure border focus color is good)
    - Sidebar
    - Process panel
- Add a screenshot of help menu (Press ?)
- Add a screenshot of popup that opens when you create a new file (Ctrl+n)
- Add a screenshot of image being preview using your theme.
- Add a screenshot of successful and unsuccessful shell command.
<details>
<summary>Example:</summary>

- Full view of superfile (Including sidebar, file previewer, process panel, metadata panel, and clipboard panel)

  - Make sure that file previewer is non empty, process panel has at least one process, and clipboard has at least one entry

  ![Full view of superfile](website/src/assets/contributing/theme-example/1.png)

- Add a screenshot of these individual panels being focused (To make sure border focus color is good)

  - Sidebar
  - Process panel

  ![Sidebar focused](website/src/assets/contributing/theme-example/2.png)

  ![Process panel focused](website/src/assets/contributing/theme-example/3.png)

- Add a screenshot of help menu (Press `?`)

  ![Help menu](website/src/assets/contributing/theme-example/4.png)

- Add a screenshot of popup that opens when you create a new file (Ctrl+n)

  ![New file popup](website/src/assets/contributing/theme-example/5.png)

- Add a screenshot of image being previewed using your theme

  ![Image preview](website/src/assets/contributing/theme-example/6.png)

- Add a screenshot of successful and unsuccessful shell command

  ![Successful shell command](website/src/assets/contributing/theme-example/7.png)

  ![Failed shell command](website/src/assets/contributing/theme-example/8.png)

</details>


---

## Sharing Ideas

Got a new idea? Awesome!

1. Check if similar ideas exist in Discussions or Issues
2. Open a discussion at: [https://github.com/yorukot/superfile/discussions](https://github.com/yorukot/superfile/discussions)
3. If you want to implement it yourself, follow the PR steps above

---

## Don't Know Where to Start?

Check out GitHub’s official guide:
[https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project](https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project)

Still unsure? Open a discussion — we’re happy to help.

---

## AI policy

Feel free to use AI for the website, documentation, and other work outside the core file manager. The usual [PR checklist](#pull-request-checklist) applies.

If you use AI to change the core file manager, read the relevant Go code first. Before requesting review, read every line the AI changed and test the behavior yourself. Be ready to explain how the change works. For example, if you change how superfile copies files, try copying a file to a folder that already contains one with the same name.

## Pull request checklist

Use the [pull request template](.github/PULL_REQUEST_TEMPLATE/pull_request_template.md) to answer these questions. Add extra headings if they help explain your change.

- **Summary** — What will change for someone using or working on superfile? Keep it under 100 words.
- **Why** — What problem does this solve? Link the issue if there is one.
- **How** — How did you solve the problem, and how did you check that the change works?
- **Other details** — Is there anything else a reviewer should see, such as screenshots or a new dependency? You can leave this section empty.

Before requesting review, check your PR against this list.

- [ ] I have read the [AI policy](#ai-policy).
- [ ] I tested the change and described the result under **How**.
- [ ] I read the full diff and removed leftover debug logs and TODOs.
- [ ] I filled out the template and included any screenshots needed to review the change.
- [ ] My PR title follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

If you changed Go code, also run these checks.

- [ ] I formatted the code with `go fmt ./...`.
- [ ] I ran `golangci-lint run` and fixed the issues it reported.

---

## 🙏 Thank You

Thank you for contributing to superfile! We appreciate every issue, pull request, and idea. Your help makes this project better for everyone.
