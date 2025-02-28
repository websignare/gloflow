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

package gf_eth_contract

import (
	// "os"
	"fmt"
	"testing"
	// "github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func Test__plugins(pTest *testing.T) {

	fmt.Println("TEST__PLUGINS ==============================================")
	
	//--------------------
	// RUNTIME_SYS
	logFun, _   := gf_core.InitLogs()
	runtime_sys := &gf_core.RuntimeSys{
		Service_name_str: "gf_eth_monitor_core__tests",
		LogFun:           logFun,
		
		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,
	}

	//--------------------

	new_contract_addr_str := "0xTestContractAddr"
	plugins_info := &GF_py_plugins{
		Base_dir_path_str: "./../../py/plugins",
	}
	gf_err := py__run_plugin__get_contract_info(new_contract_addr_str,
		plugins_info,
		runtime_sys)
	if gf_err != nil {
		pTest.Fail()
	}


}