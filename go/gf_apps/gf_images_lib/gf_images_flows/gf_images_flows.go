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

package gf_images_flows

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_policy"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//-------------------------------------------------
// IMPORTANT!! - image_flow's are ordered sequences of images, that the user creates and then
//               over time adds images to it... 

type GFflow struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`
	NameStr           string             `bson:"name_str"`
}

type GFimageExistsCheck struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty"`
	IDstr               gf_core.GF_ID      `bson:"id_str"`
	Tstr                string             `bson:"t"`
	CreationUNIXtimeF   float64            `bson:"creation_unix_time_f"`
	ImagesExternURLsLst []string           `bson:"images_extern_urls_lst"`
}

// //-------------------------------------------------
// // GET_MAPPING_TO_S3_BUCKETS
// func flows__get_mapping_to_s3_buckets() map[string]string {
// 	// FLOW_TO_S3_BUCKET_MAPPING - maps which image flows are going to use which S3 buckets
// 	//                             to store their images.
// 	flow_to_s3_bucket_map := map[string]string{
// 		"general":    "gf--img",
// 		"discovered": "gf--discovered--img", // "gf--img--discover"
// 	}
//
// 	return flow_to_s3_bucket_map
// }

//-------------------------------------------------
func pipelineGetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GF_error) {

	resultsLst, gfErr := DBgetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	allFlowsLst := []map[string]interface{}{}
	for _, flowInfoMap := range resultsLst {
		flowNameStr      := flowInfoMap["_id"].(string)
		flowImgsCountInt := flowInfoMap["count_int"].(int32)

		allFlowsLst = append(allFlowsLst, map[string]interface{}{
			"flow_name_str":       flowNameStr,
			"flow_imgs_count_int": flowImgsCountInt,
		})
	}
	return allFlowsLst, nil
}

//-------------------------------------------------
// GET_PAGE__PIPELINE
func pipelineGetPage(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) ([]*gf_images_core.GF_image, *gf_core.GF_error) {

	//--------------------
	// INPUT

	qs_map := p_req.URL.Query()

	flow_name_str := "general" // default
	if a_lst,ok := qs_map["fname"]; ok {
		flow_name_str = a_lst[0]
	}

	var err error
	page_index_int := 0 // default
	if a_lst, ok := qs_map["pg_index"]; ok {
		pg_index           := a_lst[0]
		page_index_int, err = strconv.Atoi(pg_index) // user supplied value
		
		if err != nil {
			gf_err := gf_core.ErrorCreate("failed to parse integer pg_index query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_index": pg_index,},
				err, "gf_images_lib", p_runtime_sys)
			return nil, gf_err
		}
	}

	page_size_int := 10 // default
	if a_lst, ok := qs_map["pg_size"]; ok {
		pg_size          := a_lst[0]
		page_size_int,err = strconv.Atoi(pg_size) // user supplied value
		if err != nil {
			gf_err := gf_core.ErrorCreate("failed to parse integer pg_size query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_size": pg_size,},
				err, "gf_images_lib", p_runtime_sys)
			return nil, gf_err
		}
	}

	p_runtime_sys.LogFun("INFO",fmt.Sprintf("flow_name_str  - %s", flow_name_str))
	p_runtime_sys.LogFun("INFO",fmt.Sprintf("page_index_int - %d", page_index_int))
	p_runtime_sys.LogFun("INFO",fmt.Sprintf("page_size_int  - %d", page_size_int))

	//--------------------

	//--------------------
	// GET_PAGES
	cursor_start_position_int := page_index_int*page_size_int
	pages_lst, gf_err := flows_db__get_page(flow_name_str,  // "general", //p_flow_name_str
		cursor_start_position_int, // p_cursor_start_position_int
		page_size_int,             // p_elements_num_int
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------
	return pages_lst, nil
}

//-------------------------------------------------
// IMAGES_EXIST_CHECK
func flowsImagesExistCheck(pImagesExternURLsLst []string,
	pFlowNameStr   string,
	pClientTypeStr string,
	p_runtime_sys  *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GF_error) {
	
	existing_images_lst, gf_err := flows_db__images_exist(pImagesExternURLsLst,
		pFlowNameStr,
		pClientTypeStr,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//-------------------------
	// PERSIST IMAGE_EXISTS_CHECK

	go func() {
		creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
		idStr             := gf_core.GF_ID(fmt.Sprintf("img_exists_check:%f", creationUNIXtimeF))
		
		check := GFimageExistsCheck{
			IDstr:               idStr,
			Tstr:                "img_exists_check",
			CreationUNIXtimeF:   creationUNIXtimeF,
			ImagesExternURLsLst: pImagesExternURLsLst,
		}

		ctx           := context.Background()
		coll_name_str := "gf_flows_img_exists_check" // p_runtime_sys.Mongo_coll.Name()
		_              = gf_core.Mongo__insert(check,
			coll_name_str,
			map[string]interface{}{
				"images_extern_urls_lst": pImagesExternURLsLst,
				"flow_name_str":          pFlowNameStr,
				"client_type_str":        pClientTypeStr,
				"caller_err_msg_str":     "failed to insert a img_exists_check into the DB",
			},
			ctx,
			p_runtime_sys)
	}()

	//-------------------------

	return existing_images_lst, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGE_WITH_POLICY
func FlowsAddExternImageWithPolicy(pImageExternURLstr string,
	pImageOriginPageURLstr string,
	pFlowsNamesLst         []string,
	pClientTypeStr         string,
	pJobsMngrCh            chan gf_images_jobs_core.JobMsg,
	pUserIDstr             gf_core.GF_ID,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (*string, *string, gf_images_core.GF_image_id, *gf_core.GF_error) {

	//-------------------------
	// POLICY_VERIFY - raises error if policy rejects the op
	opStr := gf_policy.GF_POLICY_OP__FLOW_ADD_IMG
	gfErr := flowsVerifyPolicy(opStr,
		pFlowsNamesLst,
		pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gf_images_core.GF_image_id(""), gfErr
	}

	//-------------------------

	runningJobIDstr, thumbnailSmallRelativeURLstr, imageIDstr, gfErr := FlowsAddExternImage(pImageExternURLstr,
		pImageOriginPageURLstr,
		pFlowsNamesLst,
		pClientTypeStr,
		pJobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gf_images_core.GF_image_id(""), gfErr
	}

	return runningJobIDstr, thumbnailSmallRelativeURLstr, imageIDstr, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGES - BATCH
func FlowsAddExternImages(pImagesExternURLsLst []string,
	pImagesOriginPagesURLsStr []string,
	pFlowsNamesLst            []string,
	pClientTypeStr            string,
	pJobsMngrCh               chan gf_images_jobs_core.JobMsg,
	pRuntimeSys               *gf_core.RuntimeSys) (*string, []*string, []gf_images_core.GFimageID, *gf_core.GFerror) {

	//------------------
	imagesURLsToProcessLst := []gf_images_jobs_core.GF_image_extern_to_process{}
	for i := 0; i < len(pImagesExternURLsLst); i++ {

		imageExternURLstr := pImagesExternURLsLst[i]
		imageOriginPageURLstr := pImagesOriginPagesURLsStr[i]
		
		image := gf_images_jobs_core.GF_image_extern_to_process{
			Source_url_str:      imageExternURLstr,
			Origin_page_url_str: imageOriginPageURLstr,
		}
		imagesURLsToProcessLst = append(imagesURLsToProcessLst, image)
	}
	
	// GF_IMAGES_JOBS_CLIENT
	runningJob, jobExpectedOutputsLst, gfErr := gf_images_jobs_client.RunExternImages(pClientTypeStr,
		imagesURLsToProcessLst,
		pFlowsNamesLst,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, nil, nil, gfErr
	}

	//------------------

	imagesIDsLst                   := []gf_images_core.GFimageID{}
	imagesThumbSmallRelativeURLlst := []*string{}

	for i:=0; i < len(jobExpectedOutputsLst); i++ {	
		imageIDstr               := gf_images_core.GFimageID(jobExpectedOutputsLst[i].Image_id_str)
		thumbSmallRelativeURLstr := jobExpectedOutputsLst[i].Thumbnail_small_relative_url_str


		imagesIDsLst = append(imagesIDsLst, imageIDstr)
		imagesThumbSmallRelativeURLlst = append(imagesThumbSmallRelativeURLlst, &thumbSmallRelativeURLstr)
	}

	return &runningJob.Id_str, imagesThumbSmallRelativeURLlst, imagesIDsLst, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGE
func FlowsAddExternImage(pImageExternURLstr string,
	pImageOriginPageURLstr string,
	pFlowsNamesLst         []string,
	pClientTypeStr         string,
	pJobsMngrCh            chan gf_images_jobs_core.JobMsg,
	pRuntimeSys            *gf_core.RuntimeSys) (*string, *string, gf_images_core.GF_image_id, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_images_flows.FlowsAddExternImage()")

	//------------------
	imagesURLsToProcessLst := []gf_images_jobs_core.GF_image_extern_to_process{
			{
				Source_url_str:      pImageExternURLstr,
				Origin_page_url_str: pImageOriginPageURLstr,
			},
		}
	
	// GF_IMAGES_JOBS_CLIENT
	runningJob, jobExpectedOutputsLst, gfErr := gf_images_jobs_client.RunExternImages(pClientTypeStr,
		imagesURLsToProcessLst,
		pFlowsNamesLst,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, nil, gf_images_core.GFimageID(""), gfErr
	}

	//------------------

	imageIDstr                       := gf_images_core.GF_image_id(jobExpectedOutputsLst[0].Image_id_str)
	thumbnail_small_relative_url_str := jobExpectedOutputsLst[0].Thumbnail_small_relative_url_str

	return &runningJob.Id_str, &thumbnail_small_relative_url_str, imageIDstr, nil
}

//-------------------------------------------------
// CREATE
func flowsCreate(pFlowNameStr string,
	pOwnerUserIDstr gf_core.GF_ID,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFflow, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	idStr             := gf_core.GF_ID(fmt.Sprintf("img_flow:%f", creationUNIXtimeF))
	

	flow := &GFflow{
		Vstr:              "0",
		IDstr:             idStr,
		NameStr:           pFlowNameStr,
		CreationUNIXtimeF: creationUNIXtimeF,
	}

	// DB
	coll_name_str := pRuntimeSys.Mongo_coll.Name()
	gfErr := gf_core.Mongo__insert(flow,
		coll_name_str,
		map[string]interface{}{
			"images_flow_name_str": pFlowNameStr,
			"caller_err_msg_str":   "failed to insert a image Flow into the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	

	// POLICY_CREATE
	gfErr = gf_policy.PipelineCreate(idStr, pOwnerUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return flow, nil
}