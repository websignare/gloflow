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
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//---------------------------------------------------
type GFuser struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	UserTypeStr       string                      `bson:"user_type_str"`   // "admin" | "standard"
	User_name_str     gf_identity_core.GFuserName `bson:"user_name_str"`   // set once at the creation of the user
	Screen_name_str   string                      `bson:"screen_name_str"` // changable durring the lifetime of the user
	
	Description_str   string                                 `bson:"description_str"`
	Addresses_eth_lst []gf_identity_core.GF_user_address_eth `bson:"addresses_eth_lst"`

	Email_str            string `bson:"email_str"`
	Email_confirmed_bool bool   `bson:"email_confirmed_bool"` // one-time confirmation on user-creation to validate user
	
	// IMAGES
	Profile_image_url_str string `bson:"profile_image_url_str"`
	Banner_image_url_str  string `bson:"banner_image_url_str"`
}

// ADD!! - provide logic/plugin for storing this record in some alternative store
//         separate from the main DB
type GFuserCreds struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	User_id_str   gf_core.GF_ID                `bson:"user_id_str"`
	User_name_str gf_identity_core.GFuserName  `bson:"user_name_str"`
	Pass_salt_str string                       `bson:"pass_salt_str"`
	Pass_hash_str string                       `bson:"pass_hash_str"`
}

// io_update
type GF_user__input_update struct {
	UserIDstr            gf_core.GF_ID                        `validate:"required"`                 // required - not updated, but for lookup
	User_address_eth_str gf_identity_core.GF_user_address_eth `validate:"omitempty,eth_addr"`       // optional - add an Eth address to the user
	Screen_name_str      *string                              `validate:"omitempty,min=3,max=50"`   // optional
	Email_str            *string                              `validate:"omitempty,email"`          // optional
	Description_str      *string                              `validate:"omitempty,min=1,max=2000"` // optional

	Profile_image_url_str *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
	Banner_image_url_str  *string `validate:"omitempty,min=1,max=100"` // optional // FIX!! - validation
}
type GF_user__output_update struct {
	
}

// io_get
type GF_user__input_get struct {
	UserIDstr gf_core.GF_ID
}

type GF_user__output_get struct {
	User_name_str         gf_identity_core.GFuserName
	Email_str             string
	Description_str       string
	Profile_image_url_str string
	Banner_image_url_str  string
}

//---------------------------------------------------
// PIPELINE__UPDATE
func users__pipeline__update(pInput *GF_user__input_update,
	pServiceInfo *GF_service_info,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (*GF_user__output_update, *gf_core.GF_error) {
	
	//------------------------
	// VALIDATE_INPUT
	gfErr := gf_core.Validate_struct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// USER_NAME
	userNameStr, gfErr := gf_identity_core.DBgetUserNameByID(pInput.UserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// EMAIL
	if pServiceInfo.Enable_email_bool {
		if *pInput.Email_str != "" {

			gfErr = usersEmailPipelineVerify(*pInput.Email_str,
				userNameStr,
				pInput.UserIDstr,
				pServiceInfo.Domain_base_str,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}
		}
	}

	//------------------------

	output := &GF_user__output_update{}

	return output, nil
}

//---------------------------------------------------
// PIPELINE__GET
func usersPipelineGet(pInput *GF_user__input_get,
	pCtx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GF_user__output_get, *gf_core.GF_error) {

	//------------------------
	// VALIDATE
	gfErr := gf_core.Validate_struct(pInput, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	
	user, gfErr := dbUserGetByID(pInput.UserIDstr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	output := &GF_user__output_get{
		User_name_str:         user.User_name_str,
		Email_str:             user.Email_str,
		Description_str:       user.Description_str,
		Profile_image_url_str: user.Profile_image_url_str,
		Banner_image_url_str:  user.Banner_image_url_str,
	}

	return output, nil
}