/**
* @program: debug_charts
*
* @description:
*
* @author: lemo
*
* @create: 2020-01-06 21:55
**/

package debug_charts

import _ "embed"

//go:embed charts/index.html
var html string

func render() string {
	return html
}
