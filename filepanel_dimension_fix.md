# Filepanel dimension fix
## Phase 1
- Have filePanel store its height and width and use it everywhere
- Mind that mainPanelHeight is height - BorderSize(2)
- Make panelElementHeight reciever, have it do height - padding - border
- have a minWidth, minHeight constants like prompt model, have an update function that does validations, like prompt model
- Have constructor initialize width and height to min
- Replace all functions to not take hight and width arguements.
  - Ex: `scrollToCursor(cursor int, mainPanelHeight int)` should change to `scrollToCursor(cursor int)`. 
  - Ex: `renderTopBar(r *rendering.Renderer, filePanelWidth int)` should change to `renderTopBar(r *rendering.Renderer)` 
  - Ex: `ItemSelectUp(mainPanelHeight int)` should change to ` ItemSelectUp()`
  - Ex: `renderCount := panelElementHeight(mainPanelHeight)` in `scrollToCursor` should change to `renderCount := m.panelElementHeight()`
  - Similiar changes in all places to not depend on mainPanelHeight and get the values from the model.

## Phase 2
- Fix unit tests and make sure they pass too
- List out all the places where filePanel's height/width may change and use the update function 