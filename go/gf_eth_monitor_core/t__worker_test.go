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

package gf_eth_monitor_core

import (
	"os"
	"fmt"
	"testing"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__worker(p_test *testing.T) {

	fmt.Println("TEST__WORKER ==============================================")
	
	ctx := context.Background()

	block_int := 4634748
	host_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST")
	worker_inspector__port_int := 2000


	//--------------------
	// RUNTIME_SYS
	log_fun     := gf_core.Init_log_fun()
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_eth_monitor_core__tests",
		Log_fun:          log_fun,
		
		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,
	}


	config := &GF_config{
		Mongodb_host_str:    "localhost:27017",
		Mongodb_db_name_str: "gf_test",
	}
	runtime, err := Runtime__get(config, runtime_sys)
	if err != nil {
		p_test.Fail()
	}





	

	
	// METRICS
	metrics_port := 9110
	metrics, gf_err := Metrics__init(metrics_port)
	if gf_err != nil {
		p_test.Fail()
	}


	//--------------------

	// GET_BLOCK__FROM_WORKER
	gf_block, gf_err := eth_blocks__worker_inspector__get_block(uint64(block_int),
		host_str,
		uint(worker_inspector__port_int),
		ctx,
		runtime_sys)

	if gf_err != nil {

		p_test.Fail()


	}




	spew.Dump(gf_block)



	abis_map := t__get_abis()

	gf_err = eth_tx__enrich_from_block(gf_block,
		abis_map,
		ctx,
		metrics,
		runtime)
	if gf_err != nil {

		p_test.Fail()
	}




}