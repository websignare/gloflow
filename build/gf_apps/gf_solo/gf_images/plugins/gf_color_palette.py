# GloFlow application and media management/publishing platform
# Copyright (C) 2021 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import argparse
import json
from colored import fg, bg, attr

import numpy as np
from PIL import Image, ImageDraw

#----------------------------------------------
def main(p_test_bool=True):

	# TEST
	if p_test_bool:
		test_image_paths_lst = [
			'./../../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg',
			'./../../gf_ml_worker/test/data/input/canvas.png',
			"./../../gf_ml_worker/test/data/input/4b14ca75070ac78323cf2ddef077ae92_thumb_medium.jpeg",
			"./../../gf_ml_worker/test/data/input/4838df39722bc2d681b67bf739f29357_thumb_small.jpeg",
			"./../../gf_ml_worker/test/data/input/1234cd19517b939d3eb726c817985fe4_thumb_medium.jpeg",
			"./../../gf_ml_worker/test/data/input/3a61e7d68fb17198e8dc0476cc862ddd_thumb_small.jpeg",
		]


		for i in range(0, len(test_image_paths_lst)):
			img = Image.open(test_image_paths_lst[i])

			# pix = img.load()
			width_int, height_int = img.size

			o_str = f"./output/out_{i}.png"

			# RUN
			run(img, o_str, width_int, height_int)

	# PRODUCTION
	else:

		args_map = parse_args()


		input_images_local_file_paths_lst = args_map["input_images_local_file_paths_str"].split(",")
		
		# RUN
		run_multiple(input_images_local_file_paths_lst)


		out_map = {}
		print(f"GF_OUT:{json.dumps(out_map)}")


#----------------------------------------------
def run_multiple(p_input_images_local_file_paths_lst):

	print(f"RUNNING {fg('green')}GF_COLOR_PALETTE{attr(0)} PLUGIN ")
	# VERIFY
	for f in p_input_images_local_file_paths_lst:
		assert os.path.isfile(f)
		print(f"input - {fg('yellow')}{f}{attr(0)}")

	# RUN
	for f in p_input_images_local_file_paths_lst:
		img = Image.open(f)
		width_int, height_int = img.size

		run(img, o_str, width_int, height_int)

#----------------------------------------------
def run(p_img, p_palette__output_file_path_str, p_img_width_int, p_img_height_int):

	img_pixels_arr = np.array(p_img.getdata())
	
	#----------------------------------------------
	def process_img__with_px_coords(p_img_pixels_with_global_index_arr, p_level_int, p_levels_max_int=6):
		
		# p_img_pixels_arr
		# shape - (px_num, 3) - 3 element tuple per pixel
		# [[r1, g1, b1], [r2, g2, b2], ...]
		
		# TERMINATION
		if p_level_int >= p_levels_max_int:
			

			img_pixels_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 0]

  
			# for all rows, only include last 3 elements (rgb pixel values), since thats what
			# we want to calculate the average with.
			img_pixels_no_global_indexes_arr = p_img_pixels_with_global_index_arr[:, 1:4]

			# get average pixel, by taking averages of all pixels on each of the channels.
			# this is the last image subdivision (dividing by pixel channel with maximum range), 
			# and so calculate average pixel value for that subdivision and return.
			# 
			# img_pixels_no_global_indexes_arr.T - average is calculated on a transpose of the pixels matrix,
			#                                      so that each channel is an array, and average is
			#                                      calculated on each channel. 
			average_pixel_arr = np.average(img_pixels_no_global_indexes_arr.T, axis=1)

			# round and cast as integer
			a = np.round(average_pixel_arr).astype(int)
			return [(img_pixels_global_indexes_arr, a[0], a[1], a[2])]

		# p_img_pixels_with_global_index_arr.T output:
		# [[     0      1      2 ... 115997 115998 115999]
		# [   162    163    164 ...    142    140    140]
		# [   200    201    202 ...    193    191    191]
		# [   203    204    205 ...    194    195    195]]
		# 
		# [1:] - exclude the the first element/row of the transposed matrix,
		#        since its the pixel_global_index and we dont need ranges on that.
		# 
		# output:
		# [   162    163    164 ...    142    140    140]
		# [   200    201    202 ...    193    191    191]
		# [   203    204    205 ...    194    195    195]]
		img_by_channel_arr = p_img_pixels_with_global_index_arr.T[1:]

		#----------
		# PICK_CHANNEL_WITH_MAX_RANGE
		# get the range (max_val-min_val) of each of the r/g/b channels.
		ranges_per_channel_lst           = img_by_channel_arr.ptp(axis=1)
		channel_with_max_range_index_int = np.argmax(ranges_per_channel_lst)

		# values of pixels in a particular channel which has the maximum range of values.
		#
		# img_by_channel_arr+1 - to account for the fact that the first element is global_index,
		#                        so the color index has to be offset by 1.
		img_max_range_channel_vals_arr = p_img_pixels_with_global_index_arr[:, channel_with_max_range_index_int+1]

		#----------
		# SORT_BY_CHANNEL
		# argsort() - sorts values and returns their indexes
		img_max_range_channel_vals_sorted_indexes_arr = img_max_range_channel_vals_arr.argsort()

		# split sorted list of pixels (on channel with max range by pixel value) indicies
		# into half/median. 
		# np.array_split() - wont error if the image doesnt split in 2 equal parts.
		upper_half__px_indicies_arr, lower_half__px_indicies_arr = np.array_split(img_max_range_channel_vals_sorted_indexes_arr, 2)



		# index into image pixels by sorted indexes of the channel with biggest range/variance
		upper_half__px_arr = p_img_pixels_with_global_index_arr[upper_half__px_indicies_arr]
		lower_half__px_arr = p_img_pixels_with_global_index_arr[lower_half__px_indicies_arr]


		upper_half__avrg_pixels_arr = process_img__with_px_coords(upper_half__px_arr, p_level_int+1)
		lower_half__avrg_pixels_arr = process_img__with_px_coords(lower_half__px_arr, p_level_int+1)

		r=[]
		r.extend(upper_half__avrg_pixels_arr)
		r.extend(lower_half__avrg_pixels_arr)
		return r

	#----------------------------------------------

	# ADD_PIXEL_GLOBAL_INDEX - index of the pixel relative to the whole image.
	#                          this is used to be able to track where in the image particular pixels come from.
	# [[r0, g0, b0], [r1, g1, b1], [r2, g2, b2], ...] -> [[0, r0, g0, b0], [1, r1, g1, b1], [2, r2, g2, b2], ...]
	# 
	# np.column_stack() - Stack 1-D arrays as columns into a 2-D array.
	img_pixels_with_global_index_arr = np.column_stack((np.arange(len(img_pixels_arr)), img_pixels_arr))



	r = process_img__with_px_coords(img_pixels_with_global_index_arr, 0)


	# IMAGE_PALETTE
	r_img = Image.new('RGB', (len(r), 1))

	palette_pixels_only_arr = np.array(r)[:, 1:4]

	print(palette_pixels_only_arr)
	print(palette_pixels_only_arr.shape)

	# r_img.putdata(palette_pixels_only_arr)
	r_img = Image.fromarray(palette_pixels_only_arr.reshape(palette_pixels_only_arr.shape[0], 1, 3).astype(np.uint8))
	r_img.save(p_palette__output_file_path_str)



	#----------
	




	# colorize original image with quantized color palette,
	# using the global_indexes of every average pixel color calcuated
	# to set those indexes to those exact values.
	# ADD!! - figure out how to vectorize this operation, withyout needing
	#         to iterate through pixels individually.
	colors_lst = []
	for c in r:
		print(c)

		global_indexes_arr = c[0]
		color_arr          = c[1:4]
		colors_lst.append(color_arr)

		for i in global_indexes_arr:
			img_pixels_arr[i] = color_arr

	# height/width/3 - height has to go before width
	img_pixels_3d_arr = img_pixels_arr.reshape(p_img_height_int, p_img_width_int, 3)



	r_img = Image.fromarray(img_pixels_3d_arr.astype(np.uint8))

	#----------------------------------------------
	# def draw_sectors():
	#     draw = ImageDraw.Draw(r_img)
	#     draw.rectangle(((0, 00), (100, 100)), outline='black', width=1)

	#----------------------------------------------
	# draw_sectors()


	r_img.save(f"{p_palette__output_file_path_str}__sectors.png")

	#----------


	print("K MEANS ------------")
	kmeans(colors_lst)


	exit()


#----------------------------------------------

def kmeans(p_colors_lst, p_k = 5):

	import scipy.spatial.distance
	import matplotlib.pyplot as plt
	import sklearn.cluster

	# array([[ 25, 236, 102],
	#       [129, 147, 103],
	#       [226,  35,  64],
	#       [ 94,  15,  33],
	#       [ 56, 104, 115]], dtype=uint64)
	centroids_arr = (np.random.rand(p_k, 3)*255).round().astype(np.uint)

	print("------------------------")
	print(p_colors_lst)
	print("==")
	print(centroids_arr)

	





	

	#----------------------------------------------
	def algo():
		for i in range(0, 10):

			# ASSIGN_COLORS_TO_CENTROIDS
			colors_centroids_lst = []
			for color_tpl in p_colors_lst:
				color_to_centroids_distances_arr = scipy.spatial.distance.cdist(centroids_arr, [color_tpl], 'euclidean')

				# print("result:")
				# print(color_to_centroids_distances_arr)

				# get the index/label of the centroid that the color is closest to
				closest_centroid_index_int = color_to_centroids_distances_arr.argmin()
				colors_centroids_lst.append(closest_centroid_index_int)
			
			# GROUP_COLORS_BY_CENTROIDS
			centroids_colors_map = {}
			for i in range(0, len(colors_centroids_lst)):

				color_arr          = p_colors_lst[i]
				centroid_label_int = colors_centroids_lst[i]
				if centroid_label_int in centroids_colors_map.keys():
					centroids_colors_map[centroid_label_int].append(color_arr)
				else:
					centroids_colors_map[centroid_label_int] = [color_arr]



			for centroid_label_int, colors_lst in centroids_colors_map.items():

				colors_by_channel_arr = np.array(colors_lst).T
				new_r_f = np.average(colors_by_channel_arr[0])
				new_g_f = np.average(colors_by_channel_arr[1])
				new_b_f = np.average(colors_by_channel_arr[2])

				centroid_new_coords_arr = np.array([new_r_f, new_g_f, new_b_f])

				# print(centroid_new_coords_arr)

				centroids_arr[centroid_label_int] = centroid_new_coords_arr
		
		return colors_centroids_lst

	#----------------------------------------------
	def plot_colors(p_colors_centroids_lst):
		

		fig = plt.figure()
		ax1 = fig.add_subplot(1, 2, 1, projection="3d") # plt.axes(projection='3d')
		ax2 = fig.add_subplot(1, 2, 2, projection="3d")

		color_map_lst = [
			"red",
			"green",
			"blue",
			"yellow",
			"orange"
		]

		#----------------
		# CLUSTERS
		colors_rgb_coords_arr = np.array(p_colors_lst).T
		colors_x_arr = colors_rgb_coords_arr[0]
		colors_y_arr = colors_rgb_coords_arr[1]
		colors_z_arr = colors_rgb_coords_arr[2]

		colors_1range_arr = np.array(p_colors_lst)/255
		# ax.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_1range_arr)

		colors_lst = [color_map_lst[l] for l in colors_centroids_lst]
		ax1.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_lst)

		#----------------

		# CENTROIDS
		centroid_coords_arr = centroids_arr.T
		centroids_x_arr = centroid_coords_arr[0]
		centroids_y_arr = centroid_coords_arr[1]
		centroids_z_arr = centroid_coords_arr[2]
		ax1.scatter3D(centroids_x_arr, centroids_y_arr, centroids_z_arr, c="black")

		# plt.show()

		#----------------
		# KMEANS_SKLEARN		
		clusters_labels_lst = sklearn.cluster.KMeans(n_clusters=p_k).fit_predict(p_colors_lst)
		print(clusters_labels_lst)

		
		colors_lst = [color_map_lst[l] for l in clusters_labels_lst]
		ax2.scatter3D(colors_x_arr, colors_y_arr, colors_z_arr, c=colors_lst) # colors_1range_arr)

		#----------------

		plt.show()

	#----------------------------------------------
	colors_centroids_lst = algo()
	plot_colors(colors_centroids_lst)

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	
	#----------------------------
	# INPUT_IMAGES_LOCAL_FILE_PATHS
	arg_parser.add_argument("-input_images_local_file_paths", action = "store", default=None,
		help = "list of image file paths (',' delimited) to process")

	#----------------------------
	# OUTPUT_DIR_PATH
	arg_parser.add_argument("-output_dir_path", action = "store", default=None,
		help = "dir path of the output images")

	#----------------------------
	# MEDIAN_CUT_LEVELS_NUM
	arg_parser.add_argument("-median_cut_levels_num", action = "store", default=None,
		help = "number of levels to use for the median-cut algo subdivisions")

	#----------------------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"input_images_local_file_paths_str": args_namespace.input_images_local_file_paths,
		"output_dir_path_str":               args_namespace.output_dir_path,
		"median_cut_levels_num":             int(args_namespace.median_cut_levels_num)
	}

#----------------------------------------------
if __name__ == "__main__":
	main()