/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_admin_lib

import (
	"net/http"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
)

//-------------------------------------------------
type GF_service_info struct {

	Name_str string

	// ADMIN_EMAIL - what the default admin email is (for auth)
	Admin_email_str string

	// EVENTS_APP - enable sending of app events from various functions
	Enable_events_app_bool bool

	// enable storage of user_creds in a secret store
	Enable_user_creds_in_secrets_store_bool bool

	// enable sending of emails for any function that needs it
	Enable_email_bool bool
}

//-------------------------------------------------
func InitNewService(pTemplatesPathsMap map[string]string,
	pServiceInfo         *GF_service_info,
	pIdentityServiceInfo *gf_identity_lib.GF_service_info,
	pHTTPmux             *http.ServeMux,
	pLocalHub            *sentry.Hub,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// STATIC FILES SERVING
	staticFilesURLbaseStr := "/v1/admin"
	localDirPathStr       := "./static"

	gf_core.HTTP__init_static_serving_with_mux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)
		
	//------------------------
	// IDENTITY_HANDLERS

	gfErr := gf_identity_lib.InitService(pHTTPmux,
		pIdentityServiceInfo,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	//------------------------
	// ADMIN_HANDLERS
	
	gfErr = init_handlers(pTemplatesPathsMap,
		pHTTPmux,
		pServiceInfo,
		pIdentityServiceInfo,
		pLocalHub,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	gfErr = init_handlers__users(pHTTPmux,
		pServiceInfo,
		pIdentityServiceInfo,
		pLocalHub,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}