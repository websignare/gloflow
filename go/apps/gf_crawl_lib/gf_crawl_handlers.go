package gf_crawl_lib

import (
	"time"
	"net/http"
	"gf_rpc_lib"
	"gf_core"
	"apps/gf_crawl_lib/gf_crawl_core"
)
//-------------------------------------------------
func init_handlers(p_crawled_images_s3_bucket_name_str string,
		p_gf_images_s3_bucket_name_str string,
		p_runtime                      *gf_crawl_core.Crawler_runtime,
		p_runtime_sys                  *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_handlers.init_handlers()")

	//----------------
	http.HandleFunc("/a/crawl/image/recent",func(p_resp http.ResponseWriter,
										p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /a/crawl/image/recent ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------------
			recent_images_lst,gf_err := gf_crawl_core.Images__get_recent(p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/crawl/image/recent",
								"failed to get recently crawled images", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
			//------------------
			//OUTPUT
			data_map := map[string]interface{}{
				"recent_images_lst":recent_images_lst,
			}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/crawl/image/recent",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//----------------
	http.HandleFunc("/a/crawl/image/add_to_flow",func(p_resp http.ResponseWriter,
										p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /a/crawl/image/add_to_flow ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------------
			//INPUT
			i,gf_err := gf_rpc_lib.Get_http_input("/a/crawl/image/add_to_flow",
											p_resp,
											p_req,
											p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/crawl/image/add_to_flow",
								"failed to get input for adding a crawled image to a flow", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}

			crawler_page__gf_image_id_str := i["crawler_page_image_id_str"].(string)

			//flow_name_str := "general" //i["flow_name_str"].(string) //DEPRECATED
			flows_names_lst := []string{}
			for _,s := range i["flows_names_lst"].([]interface{}) {
				flows_names_lst = append(flows_names_lst,s.(string))
			}
			//--------------------------

			gf_err = gf_crawl_core.Flows__add_extern_image(crawler_page__gf_image_id_str,
												flows_names_lst,
												p_crawled_images_s3_bucket_name_str,
												p_gf_images_s3_bucket_name_str,
												p_runtime,
												p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/crawl/image/add_to_flow",
								"failed to add a crawled image to a flow", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
			//------------------
			//OUTPUT
			data_map := map[string]interface{}{}
			gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
			//------------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/crawl/image/add_to_flow",
									start_time__unix_f,
									end_time__unix_f,
									p_runtime_sys)
			}()
		}
	})
	//----------------
	http.HandleFunc("/a/crawl/search",func(p_resp http.ResponseWriter,
										p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /a/crawl/search ----------")

		query_term_str := p_req.URL.Query()["term"][0]
		p_runtime_sys.Log_fun("INFO","query_term_str - "+query_term_str)

		//IMPORTANT!! - only query if the indexer is enabled
		if p_runtime.Esearch_client != nil {
			gf_err := gf_crawl_core.Index__query(query_term_str,p_runtime,p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/crawl/search",
								"failed to query the crawled index", //p_user_msg_str
								gf_err,p_resp,p_runtime_sys)
				return
			}
		}
		//------------------
		//OUTPUT
		data_map := map[string]interface{}{}
		gf_rpc_lib.Http_Respond(data_map,"OK",p_resp,p_runtime_sys)
		//------------------
	})
	//----------------

	return nil
}