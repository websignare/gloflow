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

package gf_eth_indexer

import (
	"os"
	"fmt"
	"testing"
	"context"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_contract"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__indexer(p_test *testing.T) {

	fmt.Println("TEST__INDEXER ==============================================")

	worker__host_port_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST_PORT")

	ctx := context.Background()
	runtime, _ := gf_eth_core.T__get_runtime(p_test)

	
	// ABI_DEFS
	abis_defs_map, gf_err := gf_eth_contract.Eth_abi__get_defs(ctx, nil, runtime)
	if gf_err != nil {
		p_test.Fail()
	}



	gf_errs_lst := index__range(2_000_000, 2_000_001,
		func(p_ctx context.Context, p_runtime *gf_eth_core.GF_runtime) []string {
			return []string{worker__host_port_str,}
		},
		abis_defs_map,
		ctx,
		nil,
		nil,
		runtime)

	if len(gf_errs_lst) > 0 {
		p_test.Fail()
	}
}