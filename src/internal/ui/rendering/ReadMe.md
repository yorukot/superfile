# renderer package
Responsible for rendering

# Dependencies
This package should not not import any other UI package, and should have minimal, ideally zero, dependency on common, utils or any other spf package. Its meant as a utilites to be used by ui components and main model. 
It also should not be even in-directly coupled with any UI components. Assume anything like color, style, border config of any other UI component change. This package should not have any changes.

# To-dos
- [ ] Use rendering package for other models like sort Menu, Help menu, etc.

# Notes
- Can we move this whole thing into a good useful TUI library outside of this repo ?. At least code it in a way that it can be moved
