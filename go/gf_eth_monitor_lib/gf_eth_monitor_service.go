/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_lib

import (
	"fmt"
	"net/http"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_service_info struct {
	Port_str           string
	Port_metrics_str   string
	SQS_queue_name_str string
	Influxdb_client    *influxdb2.Client
}

//-------------------------------------------------
func Run_service(p_service_info *GF_service_info,
	p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_eth_monitor_service.Run_service()")



	//-------------
	// METRICS
	metrics_port_int := 9110

	metrics, gf_err := metrics__init(metrics_port_int)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// QUEUE
	queue_name_str  := p_service_info.SQS_queue_name_str
	queue_info, err := Event__init_queue(queue_name_str)
	if err != nil {
		fmt.Println("failed to initialize event queue")
		panic(err)
	}

	
	// QUEUE_START_CONSUMING
	event__start_sqs_consumer(queue_info, metrics)

	//-------------
	// HANDLERS
	gf_err = init_handlers(queue_info,
		p_service_info,
		metrics,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------

	p_runtime_sys.Log_fun("INFO", fmt.Sprintf("STARTING HTTP SERVER - PORT - %s", p_service_info.Port_str))
	http_err := http.ListenAndServe(fmt.Sprintf(":%s", p_service_info.Port_str), nil)

	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - ", p_service_info.Port_str)
		p_runtime_sys.Log_fun("ERROR", msg_str)
		p_runtime_sys.Log_fun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}