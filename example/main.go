/**
* @program: debug_charts
*
* @description:
*
* @author: lemo
*
* @create: 2020-01-07 20:36
**/

package main

import (
	"os"

	"github.com/Lemo-yxk/lemo/console"
	"github.com/Lemo-yxk/lemo/utils"

	_ "github.com/Lemo-yxk/debug_charts"
)

func main() {

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Info(sig)
	})
}
