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

package gf_crawl_core

import (
	"fmt"
	"testing"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__img_add_to_flow(p_test *testing.T) {

	//cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//IMPORTANT!! - in this test there is no downloading of a file. a gf_page_img__pipeline_info reference is created manually
	//              with a local image file path set manually. this local path is the path of the test image (test__local_image_file_path_str).
	//              crawler image ADT's are manually created first 

	//-------------------
	//INIT 

	test__crawler_name_str                  := "test-crawler"
	test__cycle_run_id_str                  := "test__cycle_run_id"
	test__image_flows_names_lst             := []string{"test_flow",}
	test__img_src_url_str                   := "/some/origin/test_image_wasp_small.jpeg"
	test__origin_page_url_str               := "/some/origin/url.html"
	test__local_image_file_path_str         := "../test_data/test_image_wasp_small.jpeg"
	test__images_store_local_dir_path_str   := "../test_data/processed_images" //image tmp thumbnails, or downloaded gif's and their frames
	test__crawled_images_s3_bucket_name_str := "gf--test--discovered--img"
	test__gf_images_s3_bucket_name_str      := "gf--test--img"




	

	runtime_sys, crawler_runtime := T__init(p_test)
	if runtime_sys == nil || crawler_runtime == nil {
		return
	}

	t__cleanup__test_page_imgs(test__crawler_name_str, runtime_sys)



	test__local_gf_image_file_path_str := t__create_test_gf_image_named_image_file(p_test,
		test__img_src_url_str,
		test__local_image_file_path_str,
		runtime_sys)
	if test__local_gf_image_file_path_str == "" {
		return
	}
	


	test__crawled_image, test__crawled_image_ref := t__create_test_image_ADTs(p_test, 
		test__crawler_name_str,
		test__cycle_run_id_str,
		test__img_src_url_str,
		test__origin_page_url_str,
		crawler_runtime,
		runtime_sys)
	if test__crawled_image == nil || test__crawled_image_ref == nil {
		return
	}
	//-------------------
	//PIPELINE_STAGE__PROCESS_IMAGES - apply image transformations, create thumbnails, etc.

	//GF_PAGE_IMAGE_LINK
	page_img_link := &gf_page_img_link{
		img_src_str:         test__img_src_url_str,
		origin_page_url_str: test__origin_page_url_str,
	}

	//GF_PAGE_IMAGE__PIPELINE_INFO - this is the struct thats passed through the crawler image processing pipeline, 
	//                               from stage to stage. here we're createing it manually and populating with test values. 
	page_img__pipeline_info := &gf_page_img__pipeline_info{
		link:         page_img_link,
		page_img:     test__crawled_image,
		page_img_ref: test__crawled_image_ref,
		exists_bool:  false, //artificially set test image to be declared as not existing already, in order to be fully processed

		//-------------------
		//IMPORTANT!! - this is critical, that the gf_image file_path is used, not the unprocessed/original file_path (test__local_image_file_path_str).
		//              this is because gf_crawl uses the gf_images naming scheme for image file_names that is based on the gf_image ID. 
		//              (this ID is used for file naming because on a lot of URL/domains some generic/common image file names are used, even though the 
		//              contents might be different)
		local_file_path_str: test__local_gf_image_file_path_str,
		//-------------------

		nsfv_bool: false,
		thumbs:    nil,
	}

	page_imgs__pinfos_lst := []*gf_page_img__pipeline_info{
		page_img__pipeline_info,
	}

	page_imgs__pinfos_with_thumbs_lst := images__stage__process_images(test__crawler_name_str,
		page_imgs__pinfos_lst,
		test__images_store_local_dir_path_str,
		test__origin_page_url_str,
		test__crawled_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)

	fmt.Println("   STAGE_COMPLETE --------------")

	assert.Equal(p_test, len(page_imgs__pinfos_lst), len(page_imgs__pinfos_with_thumbs_lst), "more page_imgs pipeline_info's returned from images__stage__process_images() then inputed")
	assert.Equal(p_test, len(page_imgs__pinfos_lst), len(page_imgs__pinfos_with_thumbs_lst), "more page_imgs pipeline_info's returned from images__stage__process_images() then inputed")

	for _, page_img__pinfo := range page_imgs__pinfos_with_thumbs_lst {

		fmt.Printf("  ------- page_img__pinfo")
		spew.Dump(page_img__pinfo)

		assert.Equal(p_test, page_img__pinfo.page_img.S3_stored_bool, false)

		if page_img__pinfo.thumbs == nil {
			p_test.Errorf("page_img.thumbs has not been set to a gf_images_utils.Gf_image_thumbs instance pointer")
			return
		}
	}

	//------------------
	//PIPELINE_STAGE__S3_STORE_IMAGES

	page_imgs__pinfos_with_s3_lst := images_s3__stage__store_images(test__crawler_name_str,
		page_imgs__pinfos_with_thumbs_lst,
		test__origin_page_url_str,
		test__crawled_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)
	
	//CHECK!! - Downloaded_bool flag doesnt seem to be set at this point, so Im setting it here 
	//          directly for consistency. check if it makes sense for this flag to be set someplace
	//          within images__stage__process_images() for example, so that here it tests id doesnt
	//          have to be managed. 
	test__crawled_image.Downloaded_bool = true

	fmt.Println("   STAGE_COMPLETE --------------")
		
	spew.Dump(page_imgs__pinfos_with_s3_lst)

	for _, page_img__pinfo := range page_imgs__pinfos_with_s3_lst {

		spew.Dump(page_img__pinfo)

		assert.Equal(p_test, page_img__pinfo.page_img.S3_stored_bool, true)
	}
	//-------------------
	//FLOWS__ADD_EXTERN_IMAGE - copying files from one FS location to another (S3 bucket to another)



	fmt.Printf("+++++++++++++++++++++++++++++++++++++\n\n")
	fmt.Printf("%s\n", yellow("TEST_CRAWLED_IMAGE"))
	spew.Dump(test__crawled_image)
	

	gf_err := Flows__add_extern_image(test__crawled_image.Id_str,
		test__image_flows_names_lst,
		test__crawled_images_s3_bucket_name_str,
		test__gf_images_s3_bucket_name_str,
		crawler_runtime,
		runtime_sys)

	if gf_err != nil {
		p_test.Errorf("failed to add external image with ID [%s] to a flows [%s] from S3 bucket [%s] to [%s]",
			test__crawled_image.Id_str, 
			fmt.Sprint(test__image_flows_names_lst),
			test__crawled_images_s3_bucket_name_str,
			test__gf_images_s3_bucket_name_str)
		return
	}
	//-------------------
}