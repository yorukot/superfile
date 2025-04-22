# renderer package
Responsible for rendering

# Dependencies
This package should not not import any other UI package. Its meant as a utilites to be used by ui components and main model. 
It also should not be even in-directly coupled with any UI components. Assume anything like color, style, border config of any other UI component change. This package should not have any changes.

# To-dos
- Rename to rendering package