/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_core

import (
	"fmt"
	"time"
	// "strings"
	"errors"
	"runtime"
	"runtime/debug"
	"github.com/fatih/color"
	"github.com/globalsign/mgo/bson"
	"github.com/getsentry/sentry-go"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type Gf_error struct {
	Id                   bson.ObjectId          `bson:"_id,omitempty"`
	Id_str               string                 `bson:"id_str"` 
	T_str                string                 `bson:"t"`                    //"gf_error"
	Creation_unix_time_f float64                `bson:"creation_unix_time_f"`
	Type_str             string                 `bson:"type_str"`
	User_msg_str         string                 `bson:"user_msg_str"`
	Data_map             map[string]interface{} `bson:"data_map"`
	Descr_str            string                 `bson:"descr_str"`
	Error                error                  `bson:"error"`
	Service_name_str     string                 `bson:"service_name_str"`
	Subsystem_name_str   string                 `bson:"subsystem_name_str"`   //major portion of functionality, a particular package, or a logical group of functions
	Stack_trace_str      string                 `bson:"stack_trace_str"`
	Function_name_str    string                 `bson:"func_name_str"`
	File_str             string                 `bson:"file_str"`
	Line_num_int         int                    `bson:"line_num_int"`
}


//-------------------------------------------------
func Panic__check_and_handle(p_user_msg_str string,
	p_panic_data_map     map[string]interface{},
	p_oncomplete_fn      func(),
	p_subsystem_name_str string,
	p_runtime_sys        *Runtime_sys) {

	// call to recover stops the unwinding and returns the argument passed to panic
	// If the goroutine is not panicking, recover returns nil.
	if panic_info := recover(); panic_info != nil {

		err := errors.New(fmt.Sprint(panic_info))

		//--------------------
		// SENTRY
		if p_runtime_sys.Errors_send_to_sentry_bool {
			
			/*sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetExtra("gf_error.service_name",   gf_error.Service_name_str)
				scope.SetExtra("gf_error.subsystem_name", gf_error.Subsystem_name_str)
				scope.SetExtra("gf_error.type",           gf_error.Type_str)
			})*/

			sentry.WithScope(func(scope *sentry.Scope) {


				scope.SetTag("gf_panic.service_name",   p_runtime_sys.Service_name_str)
				scope.SetTag("gf_panic.subsystem_name", p_subsystem_name_str)
				scope.SetTag("gf_panic.type",           "panic_error")

				for k, v := range p_panic_data_map {
					scope.SetTag(fmt.Sprintf("gf_panic.%s", k),
						fmt.Sprint(v))
				}

				// scope.SetLevel(sentry.LevelWarning);

				sentry.CaptureException(err)
			})

			// FLUSH
			defer sentry.Flush(2 * time.Second)
		}

		//--------------------

		if p_oncomplete_fn != nil {
			p_oncomplete_fn()
		}
	}
}

//-------------------------------------------------
func Error__create_with_hook(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_hook_fun           func(*Gf_error) map[string]interface{},
	p_runtime_sys        *Runtime_sys) *Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_error.Error__create_with_hook()")

	gf_error := Error__create(p_user_msg_str,
		p_error_type_str,
		p_error_data_map,
		p_error,
		p_subsystem_name_str,
		p_runtime_sys)

	p_hook_fun(gf_error)
	return gf_error
}


//-------------------------------------------------
func Error__create(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_runtime_sys        *Runtime_sys) *Gf_error {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_error.Error__create()")

	error_defs_map := error__get_defs()
	
	gf_err := Error__create_with_defs(p_user_msg_str,
		p_error_type_str,
		p_error_data_map,
		p_error,
		p_subsystem_name_str,
		error_defs_map,
		p_runtime_sys)

	return gf_err
}

//-------------------------------------------------
func Error__create_with_defs(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_err_defs_map       map[string]Error_def,
	p_runtime_sys        *Runtime_sys) *Gf_error {

	

	creation_unix_time_f := float64(time.Now().UnixNano()) / 1000000000.0
	id_str               := fmt.Sprintf("%s:%f", p_error_type_str, creation_unix_time_f)
	stack_trace_str      := string(debug.Stack())

	// IMPORTANT!! - number of stack frames to skip before recording. without skipping 
	//               we would get info on this function, not its caller which is where
	//               the error occured.
	skip_stack_frames_num_int := 1

	// https://golang.org/pkg/runtime/#Caller
	program_counter, file_str, line_num_int,_ := runtime.Caller(skip_stack_frames_num_int)

	// FuncForPC - returns a *Func describing the function that contains the given program counter address
	function          := runtime.FuncForPC(program_counter)
	function_name_str := function.Name()

	//--------------------
	// ERROR_DEF

	error_def, ok := p_err_defs_map[p_error_type_str]
	if !ok {
		panic(fmt.Sprintf("unknown gf_error type encountered - %s", p_error_type_str))
	}

	//--------------------

	gf_error := Gf_error{
		Id_str:               id_str,
		T_str:                "gf_error",
		Creation_unix_time_f: creation_unix_time_f,
		Type_str:             p_error_type_str,
		User_msg_str:         p_user_msg_str,
		Data_map:             p_error_data_map,
		Descr_str:            error_def.Descr_str,
		Error:                p_error,
		Service_name_str:     p_runtime_sys.Service_name_str,
		Subsystem_name_str:   p_subsystem_name_str,
		Stack_trace_str:      stack_trace_str,
		Function_name_str:    function_name_str,
		File_str:             file_str,
		Line_num_int:         line_num_int,
	}


	red      := color.New(color.FgRed).SprintFunc()
	cyan     := color.New(color.FgCyan, color.BgWhite).SprintFunc()
	yellow   := color.New(color.FgYellow).SprintFunc()
	yellowBg := color.New(color.FgBlack, color.BgYellow).SprintFunc()
	green    := color.New(color.FgBlack, color.BgGreen).SprintFunc()



	//--------------------
	// VIEW
	// fmt.Printf("\n\n  %s ------------- %s\n\n\n", red("FAILED FUNCTION CALL"), yellow(function_name_str))

	fmt.Printf("GF_ERROR:\n")
	fmt.Printf("file           - %s\n", yellowBg(gf_error.File_str))
	fmt.Printf("line_num       - %s\n", yellowBg(gf_error.Line_num_int))
	fmt.Printf("user_msg       - %s\n", yellowBg(gf_error.User_msg_str))
	fmt.Printf("id             - %s\n", yellow(gf_error.Id_str))
	fmt.Printf("type           - %s\n", yellow(gf_error.Type_str))
	fmt.Printf("service_name   - %s\n", yellow(gf_error.Service_name_str))
	fmt.Printf("subsystem_name - %s\n", yellow(gf_error.Subsystem_name_str))
	fmt.Printf("function_name  - %s\n", yellow(gf_error.Function_name_str))
	fmt.Printf("data           - %s\n", yellow(gf_error.Data_map))
	fmt.Printf("error          - %s\n", red(p_error))
	fmt.Printf("%s:\n%s\n", cyan("STACK TRACE"), green(gf_error.Stack_trace_str))
	
	p_runtime_sys.Log_fun("ERROR", fmt.Sprintf("gf_error created - type:%s - service:%s - subsystem:%s - func:%s - usr_msg:%s",
		p_error_type_str,
		p_runtime_sys.Service_name_str,
		p_subsystem_name_str,
		function_name_str,
		p_user_msg_str))

	fmt.Printf("\n\n")

	//--------------------
	// DB_PERSIST
	if p_runtime_sys.Errors_send_to_mongodb_bool {
		err := p_runtime_sys.Mongodb_coll.Insert(gf_error)
		if err != nil {

		}
	}
	
	//--------------------
	// SENTRY
	if p_runtime_sys.Errors_send_to_sentry_bool {
		
		/*sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetExtra("gf_error.service_name",   gf_error.Service_name_str)
			scope.SetExtra("gf_error.subsystem_name", gf_error.Subsystem_name_str)
			scope.SetExtra("gf_error.type",           gf_error.Type_str)
		})*/

		sentry.WithScope(func(scope *sentry.Scope) {


			scope.SetTag("gf_error.service_name",   gf_error.Service_name_str)
			scope.SetTag("gf_error.subsystem_name", gf_error.Subsystem_name_str)
			scope.SetTag("gf_error.type",           gf_error.Type_str)

			for k, v := range gf_error.Data_map {
				scope.SetTag(fmt.Sprintf("gf_error.%s", k),
					fmt.Sprint(v))
			}

			// scope.SetLevel(sentry.LevelWarning);

			if p_error != nil {
				sentry.CaptureException(p_error)
			} else {

				// IMPORTANT!! - in case the GF_error doesnt have a correspoting
				//               golang error. this is for GF error conditions that are 
				//               not caused by a golang error.
				err := errors.New(fmt.Sprintf("GF error - %s", gf_error.Type_str))
				sentry.CaptureException(err)
			}
		})

		// FLUSH
		defer sentry.Flush(2 * time.Second)
	}
	
	//--------------------

	return &gf_error
}

//-------------------------------------------------
func Error__init_sentry(p_sentry_endpoint_str string,
	p_transactions__to_trace_map map[string]bool,
	p_sample_rate_f              float64) error {



	err := sentry.Init(sentry.ClientOptions{
		Dsn: p_sentry_endpoint_str,

		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,

		// TRACING
		// TracesSampleRate: p_sample_rate_f, // 1.0,

		TracesSampler: sentry.TracesSamplerFunc(func(p_ctx sentry.SamplingContext) sentry.Sampled {


			
			hub                  := sentry.GetHubFromContext(p_ctx.Span.Context())
			transaction_name_str := hub.Scope().Transaction()
			// handler_path_str     := strings.Split(transaction_name_str, " ")[1]



			fmt.Printf("SENTRY TX - %s\n", transaction_name_str)



			// exclude traces of all transactions that are not for expected handlers
			if _, ok := p_transactions__to_trace_map[transaction_name_str]; !ok {
				
				// exclude

				fmt.Println("SENTRY TX EXCLUDE")
				return sentry.SampledFalse
			}

			fmt.Println("INCLUDE")
			return sentry.SampledTrue // sentry.UniformTracesSampler(p_sample_rate_f).Sample(p_ctx)

			// // Sample all other transactions for testing. On
			// // production, use TracesSampleRate with a rate adequate
			// // for your traffic, or use the SamplingContext to
			// // customize sampling per-transaction.
			// return sentry.SampledTrue
		}),
	})
	if err != nil {
		return err
	}
	
	return nil
}