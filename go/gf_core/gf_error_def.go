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

package gf_core

//-------------------------------------------------
type Error_def struct {
	Descr_str string
}

//-------------------------------------------------
func error__get_defs() map[string]Error_def {

	error_defs_map := map[string]Error_def{

		//---------------
		"panic_error": Error_def{
			Descr_str: "a golang panic was caught with recover()",
		},

		//---------------
		"int_parse_error": Error_def{
			Descr_str: "failed to parse an integer string",
		},
		"url_parse_error": Error_def{
			Descr_str: "failed to parse a url with url.Parse()",
		},
		"url_unescape_error": Error_def{
			Descr_str: "failed to unescape a url with url.QueryUnescape()",
		},
		"io_reader_error": Error_def{
			Descr_str: "failed to read bytes using the io.Reader.ReadBytes()",
		},

		//---------------
		// DATA_VERIFICATION
		"verify__invalid_value_error": Error_def{
			Descr_str: "data failed verification, not an expected value",
		},
		"verify__value_not_integer_error": Error_def{
			Descr_str: "data failed verification, the supplied value is not an integer",
		},
		"verify__value_too_many_error": Error_def{
			Descr_str: "data failed verification, the supplied too many values",
		},
		"verify__missing_key_error": Error_def{
			Descr_str: "data failed verification, the needed key is missing",
		},
		"verify__invalid_key_value_error": Error_def{
			Descr_str: "data failed verification, the key does not have the expected value",
		},

		// length
		"verify__string_too_short_error": Error_def{
			Descr_str: "data failed verification, the string is too short",
		},
		"verify__string_too_long_error": Error_def{
			Descr_str: "data failed verification, the string is too long",
		},
		"verify__string_not_correct_length_error": Error_def{
			Descr_str: "data failed verification, the string is too long",
		},

		"verify__invalid_image_extension_error": Error_def{
			Descr_str: "an unsupported image file extension was encountered",
		},
		"verify__invalid_query_string_encoding_error": Error_def{
			Descr_str: "string is not a valid query-string encoding",
		},
		"verify__invalid_image_nsfv_error": Error_def{
			Descr_str: "image NSFV verification failed",
		},

		//---------------
		// FILESYSTEM
		"file_open_error": Error_def{
			Descr_str: "os.Create() failed to create a file - package (os)",
		},
		"file_create_error": Error_def{
			Descr_str: "os.Open() failed to open a file - package (os)",
		},
		"file_read_error": Error_def{
			Descr_str: "f.Read()/ioutil.ReadFile() failed to read file - package (os/ioutil)",
		},
		"file_remove_error": Error_def{
			Descr_str: "os.Remove() failed to remove a file - package (os)",
		},
		"file_write_error": Error_def{
			Descr_str: "f.Write() failed to write to a file - package (os)",
		},
		"file_sync_error": Error_def{
			Descr_str: "f.Sync() failed to sync a file to the FS - package (os)",
		},
		"file_missing_error": Error_def{
			Descr_str: "file doesnt exist in the FS",
		},
		"file_buffer_copy_error": Error_def{
			Descr_str: "using a file as a source/target of a buffer copy failed - (io.Copy(),etc.)",
		},
		"file_stat_error": Error_def{
			Descr_str: "getting info on a file via a stat() system call (golang API or CLI) failed - (os.Stat())",
		},
		"dir_list_error": Error_def{
			Descr_str: "failed to list contents of a dir in the FS",
		},

		//---------------
		// CLI
		"cli_run_error": Error_def{
			Descr_str: "failed to run a CLI command from Go",
		},

		//---------------
		// ENCODE/DECODE
		// JSON
		"json_decode_error": Error_def{
			Descr_str: "json.Unmarshal() failed to parse byte array - package (encoding/json)",
		},
		"json_encode_error": Error_def{
			Descr_str: "json.Marshal() failed to parse byte array - package (encoding/json)",
		},
		// YAML
		"yaml_decode_error": Error_def{
			Descr_str: "yaml.Unmarshal() failed to parse byte array - package (gopkg.in/yaml.v2)",
		},
		// BASE64
		"base64_decoding_error": Error_def{
			Descr_str: "base64.StdEncoding.DecodeString() failed - package (encoding/base64)",
		},
		// HEX
		"decode_hex": Error_def{
			Descr_str: "failed to decode hex string",
		},
		// MAPSTRUCT
		"mapstruct__decode": Error_def{
			Descr_str: "failed to decode a map into a struct using mapstructure lib",
		},

		//---------------
		// HTTP
		"http_client_req_error": Error_def{
			Descr_str:"failed to execute a http_client request",
		},
		"http_client_req_status_error": Error_def{
			Descr_str:"http_client received a non 2xx/3xx HTTP status code",
		},
		"http_server_flusher_not_supported_error": Error_def{
			Descr_str:"http_server not supporting http.Flusher (probably for SSE support,etc.)",
		},
		"http_client_gf_status_error": Error_def{
			Descr_str:"http_client received a non-OK GF error",
		},
		
		//---------------
		// WEBSOCKETS
		"ws_connection_init_error": Error_def{
			Descr_str: "websocket client failed to connect to a url",
		},
		
		//---------------
		// IMAGES
		"image_decoding_error": Error_def{
			Descr_str: "image.Decode() failed to decode image data - package (image)",
		},
		"image_decoding_config_error": Error_def{
			Descr_str: "image.DecodeConfig() failed to decode image data - package (image,image/png,image/jpeg,etc.)",
		},
		"jpeg_decoding_error": Error_def{
			Descr_str: "jpeg.Decode() failed to decode image data - package (image/jpeg)",
		},
		"png_decoding_error": Error_def{
			Descr_str: "png.Decode() failed to decode image data - package (image/png)",
		},
		"png_encoding_error": Error_def{
			Descr_str: "png.Encode() failed to encode image data - package (image/png)",
		},
		"gif_decoding_frames_error": Error_def{
			Descr_str: "gif.DecodeAll() failed to decode GIF frames - package (image/gif)",
		},

		//---------------
		// MONGODB
		"mongodb_connect_error": Error_def{
			Descr_str: "failed to connect to a mongodb host - package (go.mongodb.org/mongo-driver)",
		},
		"mongodb_ping_error": Error_def{
			Descr_str: "failed to ping a mongodb host - package (go.mongodb.org/mongo-driver)",
		},		
		"mongodb_find_error": Error_def{
			Descr_str:"c.Find() failed to find a mongodb document",
		},
		"mongodb_count_error": Error_def{
			Descr_str:"Count of documents failed in mongodb",
		},
		"mongodb_not_found_error": Error_def{
			Descr_str:"target document not found in mongodb",
		},
		"mongodb_insert_error": Error_def{
			Descr_str:"c.Insert() failed to insert/create new mongodb document",
		},
		"mongodb_write_bulk_error": Error_def{
			Descr_str:"c.BulkWrite() failed to write new mongodb documents",
		},
		"mongodb_update_error": Error_def{
			Descr_str:"c.Update() failed to update a mongodb document",
		},
		"mongodb_delete_error": Error_def{
			Descr_str:"c.Update() failed to update a mongodb document",
		},
		"mongodb_aggregation_error": Error_def{
			Descr_str:"failed to run a aggregation pipeline in mongodb",
		},
		"mongodb_ensure_index_error": Error_def{
			Descr_str:"c.EnsureIndex() failed to create a mongodb index",
		},
		"mongodb_cursor_decode": Error_def{
			Descr_str:"failed to decode value from the mongodb results Cursor",
		},

		//---------------
		// ELASTICSEARCH
		"elasticsearch_get_client": Error_def{
			Descr_str:"c.NewClient() failed to get elasticsearch client - package (elastic)",
		},
		"elasticsearch_ping": Error_def{
			Descr_str:"c.Ping() failed to ping elasticsearch server from client - package (elastic)",
		},
		"elasticsearch_add_to_index": Error_def{
			Descr_str:"c.Index() failed to add a record to the index - package (elastic)",
		},
		"elasticsearch_query_index": Error_def{
			Descr_str:"c.Search() failed issue a query - package (elastic)",
		},

		//---------------
		// TEMPLATES
		"template_create_error": Error_def{
			Descr_str:"template.New() failed to create/load a template - package (text/template)",
		},
		"template_render_error": Error_def{
			Descr_str:"template.Execute() failed to render a template - package (text/template)",
		},
		
		//---------------
		// AWS
		"aws_general_error": Error_def{
			Descr_str: "AWS general error",
		},
		"aws_ec2_instances_describe_error": Error_def{
			Descr_str: "failed to describe EC2 instances",
		},
		"aws_ecs_service_update_error": Error_def{
			Descr_str: "failed to update an AWS ECS service",
		},

		//---------------
		// S3
		"s3_credentials_error": Error_def{
			Descr_str: "S3 credentials operation failed",
		},
		"s3_file_upload_error": Error_def{
			Descr_str: "failed to upload a file to S3 bucket",
		},
		"s3_file_upload_url_presign_error": Error_def{
			Descr_str: "failed to get a presigned URL for uploading a file to S3 bucket",
		},
		"s3_file_copy_error": Error_def{
			Descr_str: "failed to copy a file within S3",
		},
		"s3_file_download_error": Error_def{
			Descr_str: "failed to download a file from S3 to a local FS",
		},

		//---------------
		// HTML_PARSING
		"html_parse_error": Error_def{
			Descr_str: "parsing of a HTML document failed",
		},

		//---------------
		// LIBRARY_ERROR
		"library_error": Error_def{
			Descr_str: "third-party library has failed",
		},

		//---------------
	}
	return error_defs_map
}