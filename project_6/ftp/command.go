package ftp

type CommandName struct {
	Name string
}

type Command struct {
	CommandName
	Desc string
}

type ConnectCommand struct {
	CommandName
}

type PutCommand struct {
	CommandName
	Filepath string
}

type GetCommand struct {
	CommandName
	Remotepath string
}
