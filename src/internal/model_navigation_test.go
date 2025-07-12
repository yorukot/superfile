package internal

import "testing"

func TestFilePanelNavigation(t *testing.T) {
	/*
	We want to test
	(1) Switching to parent directory 
	(2) Switching to parent on being at root "/"
	(3) Entering current directory
	(4) Entering via cd / command
	
	Make sure to validate
	- Search bar is cleared
	- The cursor and render values are restored correctly
	*/
}