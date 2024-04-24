package gobd2

type Commander struct {
	connector Connector
}

func NewCommander(connector Connector) *Commander {
	return &Commander{connector}
}

func (cmd *Commander) ExecuteCommand(command CommandCode) (string, error) {
	response, err := cmd.connector.SendCommand(command)
	if err != nil {
		return "", err
	}

	return response, nil
}
