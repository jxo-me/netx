package bot

const (
	CodeStart  = "```"
	CodeEnd    = "```"
	CodeMsgTpl = `
%s
%s
%s
%s
`
	CodeTpl = `
%s
%s
%s
`
	ErrInvalid  = "object invalid"
	ErrDup      = "object duplicated"
	ErrCreate   = "object creation failed"
	ErrNotFound = "object not found"
	ErrSave     = "save config failed"
)
