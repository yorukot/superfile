# Filepanel dimension fix
## Phase 1
- Have filePanel store its height and width and use it everywhere
- Mind that mainPanelHeight is height - BorderSize(2)
- Make panelElementHeight reciever, have it do height - padding - border
- have a minWidth, minHeight, have an update function that does validations, like prompt model


## Phase 2
- Replace all functions to not take hieght and width arguements.
- Fix unit tests and make sure they pass too
- List out all the places where filePanel's height/width may change and use the update function 
