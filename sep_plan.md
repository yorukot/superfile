
- Clipboard seperation
  - Phase 1
    - Separate into a component
    - Add management of copyItems there.
    - Have it store its width and height. Use that while rendering
    - Ensure its dimensions are adjusted on resizing.
  - Phase 2 
    - Unit tests for rendering
    - Critical unit tests for any other important functionality 


# Processbar fix
- Stop taking `footerHeight` as input, and save it in.
- Ensure that it is udpated on model resizes

# sidebar fix
- Stop talking width/height as input in render. save it in
- Ensure that it is udpated on model resizes

# Misc fix
- move `common.FilePanelNoneText` to filepanel const.go and make it un-exported
- List any other global consts that are only used in filepanel