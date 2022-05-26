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
	"github.com/lemonyxk/debug_charts"
	_ "github.com/lemonyxk/debug_charts"
)

func main() {
	debug_charts.Start()

	select {}
}
