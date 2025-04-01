package common

func NoPromptAction() PromptAction {
	return PromptAction{
		Action: NoAction,
	}
}

func (p PromptAction) IsNoAction() bool {
	return p.Action == NoAction
}
