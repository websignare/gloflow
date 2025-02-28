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

package gf_publisher_lib

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------
// TAGS
//------------------------------------------
func Add_tags_to_post_in_db(p_post_title_str string,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.RuntimeSys) (*gf_publisher_core.Gf_post, *gf_core.Gf_error) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post_utils.Add_tags_to_post_in_db()")
	
	post, gf_err := gf_publisher_core.DB__get_post(p_post_title_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	add_tags_to_post(post, p_tags_lst, p_runtime_sys)
	fmt.Println(fmt.Sprintf(" --------- post tags - %s", post.Tags_lst))

	gf_err = gf_publisher_core.DB__update_post(post, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return post, nil
}

//------------------------------------------
func add_tags_to_post(p_post *gf_publisher_core.Gf_post,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.RuntimeSys) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post_utils.add_tags_to_post()")
	
	if len(p_tags_lst) > 0 {
		p_post.Tags_lst = append(p_post.Tags_lst, p_tags_lst...)

		//---------------
		// eliminate duplicates from the list, in case 
		// some of the tags just added already exist in the list of all tags

		encountered_map   := map[string]bool{}
		no_dupliactes_lst := []string{}

		for _, t_str := range p_post.Tags_lst {
			if encountered_map[t_str] {
				// tuplicate exists
			} else {
				encountered_map[t_str] = true
				no_dupliactes_lst      = append(no_dupliactes_lst, t_str)
 			}
		}

		//---------------
		
		p_post.Tags_lst = no_dupliactes_lst
	} else {
		p_post.Tags_lst = append(p_post.Tags_lst, p_tags_lst...)
	}
}

//---------------------------------------------------
func get_posts_small_thumbnails_urls(p_posts_lst []*gf_publisher_core.Gf_post, p_runtime_sys *gf_core.RuntimeSys) map[string][]string {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post_utils.get_posts_small_thumbnails_urls()")
	
	posts_small_thumbnails_urls_map := map[string][]string{}
	for _, post := range p_posts_lst {

		post_small_thumbnails_urls_lst := []string{}
		for _,post_element := range post.Post_elements_lst {

			thumb_url_str                           := post_element.Img_thumbnail_small_url_str
			post_element.Img_thumbnail_small_url_str = thumb_url_str
		}
		posts_small_thumbnails_urls_map[post.Title_str] = post_small_thumbnails_urls_lst
	}
	return posts_small_thumbnails_urls_map
}

/*//------------------------------------------
//VARIOUS
//---------------------------------------------------
//->:List<:Tuple<:String(post_day_month_year_str),:List<:Post_ADT>>>(posts_by_day_groups_sorted_lst)
List group_posts_by_creation_date(List<gf_post.Post_ADT> p_posts_adts_lst,
								  Function               pLogFun) {
	pLogFun("FUN_ENTER","gf_post_utils.group_posts_by_creation_date()");
	
	//--------------------
	//SORTING_1
	//sort the all posts by creation_datetime
	//newest dates first (hence reverse = True)
	
	//ADD!! - reverse the sorted list, so that newest dates are first
	
	final List<gf_post.Post_ADT> sorted_by_creation_datetime_lst = 
		p_posts_adts_lst.sort((gf_post.Post_ADT p_post_1_adt,
							   gf_post.Post_ADT p_post_2_adt) {
			return p_post_1_adt.creation_datetime - p_post_2_adt.creation_datetime;
		});
	//--------------------
	//GROUP POSTS BY DATE
	
	final Map posts_by_date_groups_dict = {};
	sorted_by_creation_datetime_lst.forEach((gf_post.Post_ADT p_post_adt) {

		final String key_str = p_post_adt.creation_datetime.toString();

		if posts_by_date_groups_dict.containsKey(key_str) {
			posts_by_date_groups_dict[key_str].add(p_post_adt);
		}
		else {
			posts_by_date_groups_dict[key_str] = [p_post_adt];
		}
	});
	//--------------------	
	//SORTING_2
	//dict by definition is unsorted, 
	//so an extra sorting is done on the day_groups as a whole. SORTING_1 gurantees that 
	//posts are sorted properly in the day_groups
	
	//:List<:Tuple<:datetime.date,:List>>
	
	final List posts_by_date_groups_lst        = posts_by_date_groups_dict.forEach
	final List posts_by_date_groups_sorted_lst = posts_by_date_groups_lst.sort((x,y) {
														
													});


	posts_by_date_groups_sorted_lst = sorted(posts_by_date_groups_dict.items(),
											 reverse:true)
	//--------------------
	final List<List> posts_by_date_string_groups_sorted_lst = [];
  
  	posts_by_date_groups_sorted_lst.forEach((List p) {

  	});

	for posts_group_date,posts_in_group_lst in posts_by_date_groups_sorted_lst:
		assert isinstance(posts_group_date,datetime.date)
  	
		post_day_month_year_str = "%s#%s#%s"%(posts_group_date.day,
			                                  posts_group_date.month,
			                                  posts_group_date.year)
		
		posts_by_date_string_groups_sorted_lst.append([post_day_month_year_str,posts_in_group_lst])
  
  
	return posts_by_date_string_groups_sorted_lst;
}*/
//---------------------------------------------------
//IMAGES
//---------------------------------------------------
/*//DEPRECEATED!!
String get_post_thumbnail_url_str(gf_post.Post_ADT p_post_adt,
								  Function         pLogFun) {
	pLogFun("FUN_ENTER","gf_post_utils.get_post_thumbnail_url_str()");

	//HACK!!
	//IMPORTANT!! - old versions of code counlnt depend on post_adt.thumbnail_url_str, so they had
	//              to either get it manually or use this get_post_thumbnail_url_str().
	//              thats why Im hacking it here and testting thumbnail_url_str for null, and if not null
	//              then using the opportunity to initialize the new SYMPHONY thumbnail_url_str property
	if (p_post_adt.thumbnail_url_str == null) {

		List<gf_post_element.PostElement_ADT> image_post_elements_lst = gf_post_element.get_post_elements_of_type(p_post_adt, "image", pLogFun);

		//some posts dont have any image elements, and since the first image element of the post is used as its thumbnail,
		//this will fail... 
		if (image_post_elements_lst.length > 0) {
			final gf_post_element.PostElement_ADT main_image_post_element_adt = image_post_elements_lst[0];
			final String                          post_thumbnail_url_str      = main_image_post_element_adt.img_thumbnail_medium_url_str;

			//----------------
			//HACK!! - setting it here, even though post_adt might not get persisted. 
			//         what really should be done is for thumbnail_url_str to get added to all posts in a batch
			//         DB job run over the entire posts collection. 
			p_post_adt.thumbnail_url_str = post_thumbnail_url_str;
			//----------------
			
			return post_thumbnail_url_str;
		}
		else {
			return null;
		}
	}
	else {
		return p_post_adt.thumbnail_url_str;
	}
}*/