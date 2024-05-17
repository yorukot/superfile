---
title: Overview
description: An overview of why we built this starter, including its features, the libraries used, and more.
head:
  - tag: title
    content: Overview | React Native / Expo Starter
---

As a team of experienced developers at Obytes Mobile Tribe, we have spent years refining our approach to building high-quality React Native applications. Our starter kit is based on the best practices and tools that we have found to be most effective in our own projects.

This starter kit has been thoroughly tested and proven successful in multiple projects over the past four years. It is regularly used by our team on a daily basis and has helped us deliver great results for our clients.

While our starter kit is heavily influenced by our own opinions and experiences, we have carefully selected the included solutions to address common challenges and meet the needs of the majority of use cases. We believe it offers a streamlined and efficient approach to building React Native apps, and we are confident that it can help others achieve their project goals as well.

## 🚀 Motivation

Our goal with this starter kit was to streamline the process of building React Native apps, both for our own team and for our clients. We wanted to create a resource that would allow us to create high-quality apps faster and with less effort, while ensuring that all of our projects adhere to the same code standards and architectural principles.

The benefits of using this starter kit are numerous. It helps our team easily switch between projects, as we can rely on a consistent foundation of code. It also allows us to focus on the business logic of each project rather than getting bogged down in boilerplate code. And, because it promotes consistency across projects, it makes it easier to maintain and scale our apps, as well as share code between teams.

Overall, our starter kit is designed to facilitate efficient and effective app development, helping us to bring the best possible products to our clients

## ✍️ Philosophy

When creating this starter kit, we had several guiding principles in mind::

- **🚀 Production-ready**: We wanted to ensure that this starter was ready for real-world use, providing a solid foundation for building production-grade apps.
- **🥷 Developer experience and productivity**: Our focus was on creating a starter that would enhance the developer experience and increase productivity.
- **🧩 Minimal code and dependencies**: We aimed to keep the codebase and dependencies as small as possible.
- **💪 Well-maintained third-party libraries**: We included only well-maintained and reliable third-party libraries, to provide stability and support for our projects.

## ⭐ Key Features

- ✅ The latest version of Expo SDK, Along with the Custom Dev client, you can leverage the best of the Expo ecosystem and maintain full control over your app.
- 🎉 [TypeScript](https://www.typescriptlang.org/) for type checking, to help you catch bugs and improve code quality.
- 💅 A minimal UI kit built with [TailwindCSS](https://www.nativewind.dev/), With the most common components you should have in your app.
- ⚙️ Support for multiple environments builds, including Production, Staging, and Development, using Expo configuration.
- 🦊 Husky for Git Hooks, to automate your git hooks and enforce code standards.
- 💡 A clean project structure with Absolute Imports, to make it easier to navigate and manage your code.
- 🚫 Lint-staged for running Eslint and TypeScript checks on Git staged files, to ensure that your code is always up to standards.
- 🗂 VSCode recommended extensions, settings, and snippets to enhance the developer experience.
- ☂️ Pre-installed [Expo Router](https://docs.expo.dev/router/introduction/) with examples, to provide a comprehensive navigation solution for your app.
- 💫 An auth flow with [Zustand](https://github.com/pmndrs/zustand) and [react-native-mmkv](https://github.com/mrousavy/react-native-mmkv) as a storage solution to save sensitive data.
- 🛠 +10 workflows for building, releasing, testing and distributing your app using [Github action](https://github.com/features/actions).
- 🔥 [React Query](https://react-query.tanstack.com/) and [axios](https://github.com/axios/axios) for fetching data, to help you build efficient and performant apps.
- 🧵 A good approach for handling forms with [react-hook-form](https://react-hook-form.com/) and [zod](https://github.com/colinhacks/zod) for validation + keyboard handling.
- 🎯 Localization with [i18next](https://www.i18next.com/), along with Eslint for validation.
- Unit testing with [Jest](https://jestjs.io/) and [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/) setup to help you write tests for your app.

## 😉 Why Expo?

Expo is a powerful tool for building React Native apps, offering a range of features and benefits that can help developers create high-quality apps more efficiently. One question we often receive from the community is why we choose to use Expo instead of the React Native CLI.

In the past, our team used the React Native CLI for our starter kit. However, we found that using Expo presented several advantages. In particular, the introduction of the Custom dev client feature allowed us to take advantage of the Expo ecosystem and utilize native libraries without the need for ejecting the app. This has greatly simplified our development process and enabled us to focus on the business logic of our projects.

Additionally, we have found that using Expo has made it easier to upgrade our apps to new versions, eliminating the issues we previously encountered when using the React Native CLI.

Overall, we believe that Expo offers numerous benefits for building React Native apps and is a valuable tool for any developer. The real question may be, **why not use Expo?**

## 🤔 Is this starter for you?

If you are planning to build a React Native app and are looking for a strong foundation, well-designed architecture, and a positive developer experience, then this starter kit is an excellent resource to consider. It offers a comprehensive set of best practices and tools that have been tested and proven effective in multiple projects.

Even if you are not sure that using a starter kit is the right choice for your project, you can still benefit from this resource. You can explore the starter kit and draw inspiration from the solutions it provides for common challenges faced by React Native developers. This can be a helpful way to discover good practices and find effective solutions for your own app development process.

Overall, whether you choose to use this starter kit as is or simply take some ideas from it, we believe it offers valuable insights and resources for anyone looking to build a high-quality React Native app.

## 🧑‍💻 Stay up to date

We are committed to continually improving our starter kit and providing the best possible resources for building React Native apps. To that end, we regularly add new features and fix any bugs that are discovered.

If you want to stay up to date with the latest developments in our starter kit, you can either watch the repository or hit the "star" button. This will allow you to receive notifications whenever new updates are available.

We value the feedback and contributions of our users, and we encourage you to let us know if you have any suggestions for improving our starter kit. We are always looking for ways to make it even more effective and useful for our community. So, please do not hesitate to reach out and share your thoughts with us.

<!-- add a gif image here  -->

## 💎 Libraries used

- [Expo](https://docs.expo.io/)
- [Expo Router](https://docs.expo.dev/router/introduction/)
- [Nativewind](https://www.nativewind.dev/v4/overview)
- [Flash list](https://github.com/Shopify/flash-list)
- [React Query](https://tanstack.com/query/v4)
- [Axios](https://axios-http.com/docs/intro)
- [React Hook Form](https://react-hook-form.com/)
- [i18next](https://www.i18next.com/)
- [zustand](https://github.com/pmndrs/zustand)
- [React Native MMKV](https://github.com/mrousavy/react-native-mmkv)
- [React Native Gesture Handler](https://docs.swmansion.com/react-native-gesture-handler/docs/)
- [React Native Reanimated](https://docs.swmansion.com/react-native-reanimated/docs/)
- [React Native Svg](https://github.com/software-mansion/react-native-svg)
- [ React Error Boundaries](https://github.com/bvaughn/react-error-boundary)
- [Expo Image](https://docs.expo.dev/versions/unversioned/sdk/image/)

## Contributors

This starter is maintained by [Obytes mobile tribe team](https://www.obytes.com/team) and we welcome new contributors to join us in improving it. If you are interested in getting involved in the project, please don't hesitate to open an issue or submit a pull request.

In addition to maintaining this starter kit, we are also available to work on custom projects and help you build your dream app. If you are looking for experienced and reliable developers to bring your app vision to life, please visit our website at [obytes.com/contact](https://www.obytes.com/contact) to get in touch with us. We would be happy to discuss your project in more detail and explore how we can help you achieve your goals.

## ❓ FAQ

If you have any questions about the starter and want answers, please check out the [Discussions](https://github.com/obytes/react-native-template-obytes/discussions) page.
