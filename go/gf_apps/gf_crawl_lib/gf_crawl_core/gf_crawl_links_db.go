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
	"net/url"
	"context"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
func Link__db_index__init(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.Link__db_index__init()")

	indexes_keys_lst := [][]string{
		[]string{"t", "crawler_name_str"}, // all stat queries first match on "t"
		[]string{"t", "id_str"},
		[]string{"t", "hash_str"},
		[]string{"t", "hash_str", "valid_for_crawl_bool", "fetched_bool", "error_type_str", "error_id_str"}, // Link__get_unresolved()
		[]string{"t", "hash_str", "valid_for_crawl_bool"}, // Link__mark_as_resolved()
	}

	_, gf_err := gf_core.Mongo__ensure_index(indexes_keys_lst, "gf_crawl", pRuntimeSys)
	return gf_err
}

//--------------------------------------------------
// LINK_DB_GET_UNRESOLVED
func Link__db_get_unresolved(p_crawler_name_str string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFcrawlerPageOutgoingLink, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.Link__get_unresolved()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	black  := color.New(color.FgBlack).Add(color.BgWhite).SprintFunc()

	fmt.Println("INFO", cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))
	fmt.Println("INFO", black("GET__UNRESOLVED_LINK")+" - for crawler - "+yellow(p_crawler_name_str))
	fmt.Println("INFO", cyan(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ---------------------------------------"))

	ctx := context.Background()
	var unresolved_link GFcrawlerPageOutgoingLink
	err := pRuntimeSys.Mongo_db.Collection("gf_crawl").FindOne(ctx, bson.M{
			"t":                    "crawler_page_outgoing_link",
			"crawler_name_str":     p_crawler_name_str, //get links that were discovered by this crawler
			"valid_for_crawl_bool": true,
			"fetched_bool":         false,

			// IMPORTANT!! - get all unresolved links that also dont have any errors associated
			//               with them. this way repeated processing of unresolved links that always cause 
			//               an error is avoided (wasted resources)
			"error_type_str": bson.M{"$exists": false,},
			"error_id_str":   bson.M{"$exists": false,},

			// //-------------------
			// // IMPORTANT!! - this gets all unresolved links that come from the domain 
			// //               that the crawler is assigned to
			// //"origin_domain_str"   :p_crawler_domain_str,
			// "$or":domains_query_lst,
			// //-------------------
		}).Decode(&unresolved_link)

	/*query := pRuntimeSys.Mongodb_db.C("gf_crawl").Find(bson.M{
			"t":                    "crawler_page_outgoing_link",
			"crawler_name_str":     p_crawler_name_str, //get links that were discovered by this crawler
			"valid_for_crawl_bool": true,
			"fetched_bool":         false,

			// IMPORTANT!! - get all unresolved links that also dont have any errors associated
			//               with them. this way repeated processing of unresolved links that always cause 
			//               an error is avoided (wasted resources)
			"error_type_str": bson.M{"$exists": false,},
			"error_id_str":   bson.M{"$exists": false,},

			// //-------------------
			// // IMPORTANT!! - this gets all unresolved links that come from the domain 
			// //               that the crawler is assigned to
			// //"origin_domain_str"   :p_crawler_domain_str,
			// "$or":domains_query_lst,
			// //-------------------
		})


	var unresolved_link GFcrawlerPageOutgoingLink
	err := query.One(&unresolved_link)*/


	// NO_DOCUMENTS
	// IMPORTANT!! - link not being found in the DB is actually expected state, and should not throw an error.
	//               instead a nil value is returned for the link without error.
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get unresolved_link from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"crawler_name_str": p_crawler_name_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}

	//-------------------
	// IMPORTANT!! - some unresolved links in the DB might possibly be urlescaped,
	//               so for proper usage it is unescaped here and stored back in the unresolved_link struct.
	unescaped_unresolved_link_url_str, err := url.QueryUnescape(unresolved_link.A_href_str)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get unresolved_link from mongodb", "url_unescape_error",
			map[string]interface{}{
				"crawler_name_str":        p_crawler_name_str,
				"unresolved_link_url_str": unresolved_link.A_href_str,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}
	unresolved_link.A_href_str = unescaped_unresolved_link_url_str

	//-------------------

	fmt.Printf("unresolved_link URL - %s\n", unresolved_link.A_href_str)
	return &unresolved_link, nil
}

//--------------------------------------------------
func Link__db_get(p_link_id_str string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFcrawlerPageOutgoingLink, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.Link__db_get()")

	ctx := context.Background()
	var unresolved_link GFcrawlerPageOutgoingLink
	err := pRuntimeSys.Mongo_db.Collection("gf_crawl").FindOne(ctx, bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link_id_str,
		}).Decode(&unresolved_link)
		
	/*err := pRuntimeSys.Mongodb_db.C("gf_crawl").Find(bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link_id_str,
		}).One(&unresolved_link)*/

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get crawler_page_outgoing_link by ID from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"link_id_str": p_link_id_str,},
			err, "gf_crawl_core", pRuntimeSys)
		return nil, gf_err
	}

	return &unresolved_link, nil	
}

//--------------------------------------------------
func link__db_create(p_link *GFcrawlerPageOutgoingLink,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	//-------------
	// IMPORTANT!! - REXAMINE!! - to conserve on storage for potentially large savings (should be checked empirically?), links are persisted
	//                            in the DB only if their hash is unique. Hashes are composed of origin page URL and target URL hashed, so multiple links coming from the 
	//                            same origin page URL, and targeting the same URL, are only stored once.
	//                            this is a potentially loss of information, for pages that have a lot of these duplicate links. having this information 
	//                            on pages could maybe prove useful for some kind of analysis or algo. 
	//                            - so maybe store links even if their hashes are duplicates?
	//                            - add some kind of tracking where these duplicates are counted for pages.
	link_exists_bool, gf_err := link__db_exists(p_link.Hash_str, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//-------------

	// crawler_page_outgoing_link already exists, from previous crawls, so ignore it
	if link_exists_bool {
		fmt.Println(">> "+yellow(">>>>>>>> - DB PAGE_LINK ALREADY EXISTS >-- ")+cyan(fmt.Sprint(p_link.A_href_str)))
		return nil
	} else {

		ctx           := context.Background()
		coll_name_str := "gf_crawl"
		gf_err        := gf_core.Mongo__insert(p_link,
			coll_name_str,
			map[string]interface{}{
				"link_a_href_str":    p_link.A_href_str,
				"caller_err_msg_str": "failed to insert a crawler_page_outgoing_link into the DB",
			},
			ctx,
			pRuntimeSys)

		if gf_err != nil {
			return gf_err
		}
		
		/*err := pRuntimeSys.Mongodb_db.C("gf_crawl").Insert(p_link)
		if err != nil {

			gf_err := gf_core.Mongo__handle_error("failed to insert a crawler_page_outgoing_link in mongodb",
				"mongodb_insert_error",
				map[string]interface{}{
					"link_a_href_str": p_link.A_href_str,
				},
				err, "gf_crawl_core", pRuntimeSys)
			return gf_err
		}*/
	}

	return nil
}

//--------------------------------------------------
func link__db_exists(pLinkHashStr string,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	ctx := context.Background()
	c, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").CountDocuments(ctx,
		bson.M{
			"t":        "crawler_page_outgoing_link",
			"hash_str": pLinkHashStr,
		})

	/*c, err := pRuntimeSys.Mongodb_db.C("gf_crawl").Find(bson.M{
		"t":        "crawler_page_outgoing_link",
		"hash_str": pLinkHashStr,
		}).Count()*/

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to count crawler_page_outgoing_link by its hash",
			"mongodb_find_error",
			map[string]interface{}{"hash_str": pLinkHashStr,},
			err, "gf_crawl_core", pRuntimeSys)
		return false, gf_err
	}

	// crawler_page_outgoing_link already exists, from previous crawls, so ignore it
	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
	return false, nil
}

//--------------------------------------------------
func Link__db_mark_import_in_progress(pStatusBool bool,
	p_unix_time_f float64,
	p_link        *GFcrawlerPageOutgoingLink,
	p_runtime     *GFcrawlerRuntime,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.Link__db_mark_import_in_progress()")

	//----------------
	update_map := bson.M{
		"import__in_progress_bool": pStatusBool,
	}
	if pStatusBool {
		update_map["import__start_time_f"] = p_unix_time_f
	} else {
		update_map["import__end_time_f"] = p_unix_time_f
	}

	//----------------
	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link.Id_str,
		},
		bson.M{"$set": update_map,})

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as in_progress/complete",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str": p_link.Id_str,
				"status_bool": pStatusBool,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}
	return nil
}

//--------------------------------------------------
func Link__db_mark_as_resolved(p_link *GFcrawlerPageOutgoingLink,
	p_fetch_id_str          string,
	p_fetch_creation_time_f float64,
	pRuntimeSys             *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.Link__db_mark_as_resolved()")
	
	ctx := context.Background()

	p_link.Fetched_bool = true
	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
				"t":                    "crawler_page_outgoing_link",
				"id_str":               p_link.Id_str,
				"valid_for_crawl_bool": true,
			},
			bson.M{"$set": bson.M{
				"fetched_bool":      true,
				"fetch_last_id_str": p_fetch_id_str,
				"fetch_last_time_f": p_fetch_creation_time_f,
			},
		})
	
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as resolved/fetched",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str":  p_link.Id_str,
				"fetch_id_str": p_fetch_id_str,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}

	return nil
}

//--------------------------------------------------
func link__db_mark_as_failed(p_error *Gf_crawler_error,
	p_link      *GFcrawlerPageOutgoingLink,
	p_runtime   *GFcrawlerRuntime,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_crawl_links_db.link__mark_as_failed()")

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_db.Collection("gf_crawl").UpdateMany(ctx, bson.M{
			"t":      "crawler_page_outgoing_link",
			"id_str": p_link.Id_str,
		},
		bson.M{"$set": bson.M{
				"error_id_str":   p_error.Id_str,
				"error_type_str": p_error.Type_str,
			},
		})

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a crawler_page_outgoing_link in mongodb as failed",
			"mongodb_update_error",
			map[string]interface{}{
				"link_id_str":    p_link.Id_str,
				"error_id_str":   p_error.Id_str,
				"error_type_str": p_error.Type_str,
			},
			err, "gf_crawl_core", pRuntimeSys)
		return gf_err
	}

	return nil
}