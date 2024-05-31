//go:build windows

package device

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetBaseBoardID() (string, error) {
	var boardID string
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return boardID, err
	}

	boardID = string(b)
	boardID = strings.ReplaceAll(boardID, "\n", "")
	fmt.Println(boardID)

	return boardID, nil
}
