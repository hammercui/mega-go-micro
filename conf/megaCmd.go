/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/19
 * Time: 14:15
 * Mail: hammercui@163.com
 *
 */
package conf

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
)

//生成通用cmd
func GenCmd() cmd.Cmd {
	customCmd := cmd.NewCmd()
	customFlag := customCmd.App().Flags
	customFlag = append(customFlag, &cli.StringFlag{
		Name:    "test.v",
		EnvVars: []string{"test"},
		Usage:   "mega cmd: test.v",
	},
		&cli.StringFlag{
			Name:    "configs",
			EnvVars: []string{"configs"},
			Usage:   "mega cmd: configs",
		},
		&cli.StringFlag{
			Name:    "logout",
			EnvVars: []string{"configs"},
			Usage:   "mega cmd: configs",
		})

	customCmd.App().Flags = customFlag
	return customCmd
}