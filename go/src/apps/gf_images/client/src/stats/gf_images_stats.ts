
///<reference path="../d/d3.d.ts" />
///<reference path="../d/nvd3.d.ts" />

namespace gf_images_stats {
//-------------------------------------------------
export function stat_view__completed_image_jobs_runtime_infos(p_completed_image_jobs_runtime_infos_lst, p_log_fun) {
	p_log_fun("FUN_ENTER","gf_images_stats.stat_view__completed_image_jobs_runtime_infos()");

	const plot = $(
		'<div id="stat_view__completed_image_jobs_runtime_infos" class="plot">'+
			'<div class="close_plot_btn">x</div>'+
			'<div class="title">completed image jobs runtime (s)</div>'+
			'<svg width="800" height="400"></svg>'+
		'</div>');
	$('#plots').append(plot);

	$(plot).find('.close_plot_btn').click(function(e){
		$(plot).remove()
	});
	//------------------------------
	//FORMAT DATA

	const d3_data_lst = [];
	$.each(p_completed_image_jobs_runtime_infos_lst,function(p_i,p_data_map){

		const runtime_duration_sec_f = p_data_map['runtime_duration_sec_f'];
		d3_data_lst.push({
			x:p_i,
			y:runtime_duration_sec_f
		});
	});
	//------------------------------

	nv.addGraph(function(){
		const chart = nv.models.lineChart()
						.margin({left: 100})           //give x-axis some space
						.useInteractiveGuideline(true) //tooltips
						//.transitionDuration(350)       
						.showLegend(true)              
						.showYAxis(true)              
						.showXAxis(true);             

		//X-axis settings
		chart.xAxis
			.axisLabel('image_job index')
			.tickFormat(function(d){
				return d3.time.format('%x')(new Date(d));
			});

		//Y-axis settings
		chart.yAxis
			.axisLabel('image_job runtime duration (seconds)');
			//.tickFormat(d3.format('.02f'));

		d3.select('#stat_view__completed_image_jobs_runtime svg')

			//populate plot with data
			.datum([{
				values:d3_data_lst,
				//key:"chronological_sorted_counts",
				//color:'#2ca02c',
				//area: false //dont fill the space under the curve
			}])
			.transition().duration(500)
			.call(chart); //render

		//update the chart when window resizes
		nv.utils.windowResize(function(){chart.update()});
		return chart;
	});
}
//-------------------------------------------------
export function stat_view__image_jobs_errors(p_image_jobs_errors_lst, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_images_stats.stat_view__image_jobs_errors()');
	
}
//-------------------------------------------------
//HTTP
//-------------------------------------------------
export function http__get_stat_data(p_stat_name_str,
	p_onComplete_fun,
	p_log_fun) {
	p_log_fun('FUN_ENTER','gf_images_stats.http__get_stat_data()');

	const url_str  = '/images/stats';
	const data_str = JSON.stringify({
		'stat_name_str':p_stat_name_str
	});

	$.ajax({
		'url'        :url_str,
		'type'       :'POST',
		'data'       :data_str,
		'contentType':'application/json',
		'success'    :function(p_response_str){
			console.log('+++++++  START -- HTTP RESPONSE ---------------');
			//console.log(p_response_str);

			const response_map = JSON.parse(p_response_str);
			const status_str   = response_map['status_str'];
			const data_map     = response_map['data'];

			//console.log(data_map);

			switch(status_str) {
				case 'OK':
					p_onComplete_fun(data_map);
					break;
				//-------------------
				//ERROR HANDLING
				
				case 'ERROR':
					break;
				//-------------------
			}
		}
	});
}
//-------------------------------------------------
}