In src/internal/ui/bulk_rename/model.go

There is 

	m.findInput = common.GenerateBulkRenameTextInput("Find text")
	m.replaceInput = common.GenerateBulkRenameTextInput("Replace with")
	m.prefixInput = common.GenerateBulkRenameTextInput("Add prefix")
	m.suffixInput = common.GenerateBulkRenameTextInput("Add suffix")


Only generate them once on model init.

List down