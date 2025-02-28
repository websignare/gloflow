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

///<reference path="../../../../d/jquery.d.ts" />

//---------------------------------------------------
export async function init(p_log_fun) {

    const all_flows_lst = await http__get_all_flows(p_log_fun) as {}[];


    // <div id="flows_experimental_label">experimental:</div>
    const all_flows_container = $(`
        <div id="flows_picker">
            <div id="flows">
            </div>

            
            <div id="flows_experimental">
            </div>
        </div>`);
    $('body').append(all_flows_container);


    const experimental_flows_lst = [
        "discovered",
        "gifs"
    ];
    for (const flow_map of all_flows_lst ) {
        const flow_name_str = flow_map["flow_name_str"];

        // FIX!! - allow access to these flows only to logged in users, ton of content there
        //         but not filtered yet for NSFW.
        if (flow_name_str == "discovered" || flow_name_str == "gifs") {
            continue;
        }

        const flow_imgs_count_int = flow_map["flow_imgs_count_int"];
        const flow_url_str = `/images/flows/browser?fname=${flow_name_str}`;

        var target_container_id_str :string;
        if (experimental_flows_lst.includes(flow_name_str)) {
            target_container_id_str = "flows_experimental";
        } else {
            target_container_id_str = "flows";
        }

        $(all_flows_container).find(`#${target_container_id_str}`).append(`
            <div id="flow_info">
                <div class="flow_imgs_count">${flow_imgs_count_int}</div>
                <div class="flow_name">
                    <a href="${flow_url_str}">${flow_name_str}</a>
                </div>
            </div>
        `);
    }
}

//---------------------------------------------------
async function http__get_all_flows(p_log_fun) {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const url_str = `/v1/images/flows/all`;
        p_log_fun('INFO', 'url_str - '+url_str);

        //-------------------------
        // HTTP AJAX
        $.get(url_str,
            function(p_data_map) {
                if (p_data_map["status"] == 'OK') {
                    const all_flows_lst = p_data_map['data']['all_flows_lst'];
                    p_resolve_fun(all_flows_lst);
                }
                else {
                    p_reject_fun(p_data_map["data"]);
                }
            });

        //-------------------------	
    });
    return p;
}