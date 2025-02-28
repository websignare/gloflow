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

package gf_identity_lib

import (
	"fmt"
	"time"
	"net/http"
	"context"
	"testing"
	"github.com/parnurzeal/gorequest"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
var logFun func(p_g string, p_m string)
var cliArgsMap map[string]interface{}

//-------------------------------------------------
func TestCreateAndLoginNewUser(pTest *testing.T,
	pHTTPagent              *gorequest.SuperAgent,
	pIdentityServicePortInt int,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) {

	testUserNameStr := "ivan_t"
	testUserPassStr := "pass_lksjds;lkdj"
	testEmailStr    := "ivan_t@gloflow.com"

	//---------------------------------
	// CLEANUP
	TestDBcleanup(pCtx, pRuntimeSys)
	
	//---------------------------------
	// ADD_TO_INVITE_LIST
	gfErr := DBuserAddToInviteList(testEmailStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		fmt.Println(gfErr.Error)
		pTest.FailNow()
	}

	//---------------------------------
	// GF_IDENTITY_INIT
	TestUserHTTPcreate(testUserNameStr,
		testUserPassStr,
		testEmailStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)

	TestUserHTTPlogin(testUserNameStr,
		testUserPassStr,
		pHTTPagent,
		pIdentityServicePortInt,
		pTest)
		
	//---------------------------------
}

//-------------------------------------------------
func TestStartService(pPortInt int,
	pRuntimeSys *gf_core.RuntimeSys) {

	// testPortInt := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		serviceInfo := &GF_service_info{

			// IMPORTANT!! - durring testing dont send emails
			Enable_email_bool: false,
		}
		InitService(HTTPmux, serviceInfo, pRuntimeSys)
		gf_rpc_lib.Server__init_with_mux(pPortInt, HTTPmux)
	}()
	time.Sleep(2*time.Second) // let server startup
}

//-------------------------------------------------
func T__init() *gf_core.RuntimeSys {

	test__mongodb_host_str    := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	test__mongodb_db_name_str := "gf_tests"
	test__mongodb_url_str := fmt.Sprintf("mongodb://%s", test__mongodb_host_str)


	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_identity_tests",
		LogFun:           logFun,
		Validator:        gf_core.Validate__init(),
	}




	mongo_db, _, gf_err := gf_core.Mongo__connect_new(test__mongodb_url_str,
		test__mongodb_db_name_str,
		nil,
		runtimeSys)
	if gf_err != nil {
		panic(-1)
	}


	mongo_coll := mongo_db.Collection("data_symphony")
	runtimeSys.Mongo_db   = mongo_db
	runtimeSys.Mongo_coll = mongo_coll




	return runtimeSys
}

//-------------------------------------------------
func TestDBcleanup(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	// CLEANUP
	collNameStr := "gf_users"
	gf_core.Mongo__delete(bson.M{}, collNameStr, 
		map[string]interface{}{
			"caller_err_msg_str": "failed to cleanup test user DB",
		},
		pCtx, pRuntimeSys)
}