# renderer package
Responsible for rendering

# Dependencies
This package should not not import any other UI package, and should have minimal, ideally zero, dependency on common, utils or any other spf package. Its meant as a utilites to be used by ui components and main model. 
It also should not be even in-directly coupled with any UI components. Assume anything like color, style, border config of any other UI component change. This package should not have any changes.

# To-dos
x Rename to rendering package
x Sectionization
x Dynamic Height / Height truncation
- Complete unit tests ?
x Use in filepanel
x Use in prompt
- Fix empty border issue

- Fix expensive debug logs
- Finish Code todos
- One self review
- sanity testing 
- Windows testing ?
- Rabbit AI comments


# Notes
- Prefer Test Driven development
- Can we move this whole thing into a good useful TUI library outside of this repo ?. At least code it in a way that it can be moved
