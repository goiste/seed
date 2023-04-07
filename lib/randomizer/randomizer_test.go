package seed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommandAndArgs(t *testing.T) {
	testCases := []struct {
		name   string
		str    string
		expCmd Command
		expOk  bool
	}{
		{
			"empty",
			"",
			Command{},
			false,
		},
		{
			"no command",
			"test",
			Command{},
			false,
		},
		{
			"invalid command 1",
			"$test!",
			Command{},
			false,
		},
		{
			"invalid command 2",
			"$test :1",
			Command{},
			false,
		},
		{
			"invalid command 3",
			"$test-command:111",
			Command{},
			false,
		},
		{
			"$int",
			"$int",
			Command{Command: "int"},
			true,
		},
		{
			"$int:",
			"$int:",
			Command{Command: "int"},
			true,
		},
		{
			"$int:1,100",
			"$int:1,100",
			Command{Command: "int", Args: []string{"1", "100"}},
			true,
		},
		{
			"$int:1,100",
			"$int:,100",
			Command{Command: "int", Args: []string{"", "100"}},
			true,
		},
		{
			"$test_command1:1,2,3",
			"$test_command1:1,2,3",
			Command{Command: "test_command1", Args: []string{"1", "2", "3"}},
			true,
		},
		{
			"$CamelCommand1:111",
			"$CamelCommand1:111",
			Command{Command: "CamelCommand1", Args: []string{"111"}},
			true,
		},
		{
			"$dot.command:1",
			"$dot.command:1",
			Command{Command: "dot.command", Args: []string{"1"}},
			true,
		},
		{
			"$[]array",
			"$[]array",
			Command{Command: "array", IsArray: true},
			true,
		},
		{
			"$[1]array",
			"$[1]array",
			Command{Command: "array", ArrayFrom: 1, IsArray: true},
			true,
		},
		{
			"$[1,]array",
			"$[1,]array",
			Command{Command: "array", ArrayFrom: 1, IsArray: true},
			true,
		},
		{
			"$[,1]array",
			"$[,1]array",
			Command{Command: "array", ArrayTo: 1, IsArray: true},
			true,
		},
		{
			"$[1,1]array",
			"$[1,1]array",
			Command{Command: "array", ArrayTo: 1, ArrayFrom: 1, IsArray: true},
			true,
		},
		{
			"$[3,1]wrong_array",
			"$[3,1]wrong_array",
			Command{Command: "wrong_array"},
			true,
		},
		{
			"$[a,b]failed_array",
			"$[a,b]failed_array",
			Command{},
			false,
		},
		{
			"$[1]array:,arg",
			"$[1]array:,arg",
			Command{Command: "array", Args: []string{"", "arg"}, ArrayFrom: 1, IsArray: true},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, ok := ParseCommand(tc.str)
			require.Equal(t, tc.expOk, ok)
			require.Equal(t, tc.expCmd, cmd)
		})
	}
}
