/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_images_flows

import (
	"os"
	"fmt"
	"testing"
	// "github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

var log_fun func(string,string)
var cli_args_map map[string]interface{}

//---------------------------------------------------
func TestMain(m *testing.M) {
	log_fun = gf_core.Init_log_fun()
	cli_args_map = gf_images_core.CLI__parse_args(log_fun)
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func Test__templates(p_test *testing.T) {

}