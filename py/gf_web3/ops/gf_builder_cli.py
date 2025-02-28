# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
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

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import argparse
import subprocess
import urllib.parse
import json
import requests

from colored import fg, bg, attr
import delegator

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli
import gf_ops_container

sys.path.append("%s/../test"%(modd_str))
import gf_test

#--------------------------------------------------
def main():

	args_map = parse_args()

	docker_user_str            = "glofloworg"
	service_cont_image_tag_str = "latest"

	services_map = {
		"gf_web3_monitor": {
			"service_name_str":                 "gf_web3_monitor",
			"service_dir_path_str":             "%s/../../../go/gf_web3"%(modd_str),
			"service_bin_output_path_str":      "%s/../../../build/gf_apps/gf_web3_monitor/gf_web3_monitor"%(modd_str),
			"service_cont_image_name_str":      f"glofloworg/gf_web3_monitor:{service_cont_image_tag_str}",
			"service_cont_dockerfile_path_str": f"{modd_str}/../../../build/gf_apps/gf_web3_monitor/Dockerfile",
		},
		"gf_eth_monitor_worker_inspector": {
			"service_name_str":                 "gf_eth_monitor_worker_inspector",
			"service_dir_path_str":             "%s/../../../go/gf_web3/gf_eth_monitor_worker_inspector"%(modd_str),
			"service_bin_output_path_str":      "%s/../../../build/gf_apps/gf_web3_monitor/gf_eth_monitor_worker_inspector"%(modd_str),
			"service_cont_image_name_str":      f"glofloworg/gf_eth_monitor_worker_inspector:{service_cont_image_tag_str}",
			"service_cont_dockerfile_path_str": f"{modd_str}/../../../build/gf_apps/gf_web3_monitor/Dockerfile__worker_inspector",
		}
	}
	
	#------------------------
	# TEST_GO

	if args_map["run"] == "test_go":
		
		test_ci_bool = args_map["test_ci_bool"]
		gf_test.run_go(test_ci_bool)

	#------------------------
	# TEST_PY

	elif args_map["run"] == "test_py":

		test_ci_bool = args_map["test_ci_bool"]
		gf_test.run_py(test_ci_bool)

	#------------------------
	# BUILD
	elif args_map["run"] == "build":

		for service_name_str, v in services_map.items():
			
			build_go(service_name_str,
				v["service_dir_path_str"],
				v["service_bin_output_path_str"],
				p_static_bool = args_map["static_bool"])
				
	#------------------------
	# BUILD_CONTAINER
	elif args_map["run"] == "build_containers":
		
		for service_name_str, v in services_map.items():
			gf_ops_container.build(v["service_cont_image_name_str"],
				v["service_cont_dockerfile_path_str"],
				p_docker_sudo_bool=args_map["docker_sudo_bool"])

	#------------------------
	# PUBLISH_CONTAINER
	elif args_map["run"] == "publish_containers":
		docker_pass_str = args_map["gf_docker_pass_str"]
		assert not docker_pass_str == None

		for service_name_str, v in services_map.items():
			gf_ops_container.publish(v["service_cont_image_name_str"],
				docker_user_str,
				docker_pass_str,
				p_docker_sudo_bool=args_map["docker_sudo_bool"])

	#------------------------
	# NOTIFY_COMPLETION
	elif args_map["run"] == "notify_completion":

		gf_notify_completion_url_str = args_map["gf_notify_completion_url_str"]
		assert not gf_notify_completion_url_str == None

		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

		for service_name_str, v in services_map.items():

			notify_completion(gf_notify_completion_url_str,
				service_name_str,
				p_git_commit_hash_str = git_commit_hash_str)

	#------------------------

#--------------------------------------------------
# NOTIFY_COMPLETION
def notify_completion(p_gf_notify_completion_url_str,
	p_service_name_str,
	p_git_commit_hash_str = None):
	
	url_str = None

	# add git_commit_hash as a querystring argument to the notify_completion URL.
	# the entity thats receiving the completion notification needs to know what the tag
	# is of the newly created container.
	if not p_git_commit_hash_str == None:
		url = urllib.parse.urlparse(p_gf_notify_completion_url_str)
		
		# QUERY_STRING
		qs_lst = urllib.parse.parse_qsl(url.query)
		qs_lst.append(("git_commit", p_git_commit_hash_str)) # .parse_qs() places all values in lists

		qs_str = "&".join(["%s=%s"%(k, v) for k, v in qs_lst])

		# _replace() - "url" is of type ParseResult which is a subclass of namedtuple;
		#              _replace is a namedtuple method that:
		#              "returns a new instance of the named tuple replacing
		#              specified fields with new values".
		url_new = url._replace(query=qs_str)
		url_str = url_new.geturl()
	else:
		url_str = p_gf_notify_completion_url_str

	print("NOTIFY_COMPLETION - HTTP REQUEST - %s"%(url_str))

	# HTTP_GET
	data_map = {
		"app_name": p_service_name_str # "gf_eth_monitor"
	}
	r = requests.post(url_str, data=json.dumps(data_map))
	print(r.text)

	if not r.status_code == 200:
		print("notify_completion http request failed")
		exit(1)


	r_map = json.loads(r.text)

	if r_map["status"] == "ERROR":
		print("notify_completion failed")
		exit(1)


#--------------------------------------------------
# BUILD_GO
def build_go(p_name_str,
	p_go_dir_path_str,
	p_go_output_path_str,
	p_static_bool       = False,
	p_exit_on_fail_bool = True):
	assert isinstance(p_static_bool, bool)
	assert os.path.isdir(p_go_dir_path_str)

	print(p_go_output_path_str)
	assert os.path.isdir(os.path.dirname(p_go_output_path_str))

	print("=============================")
	if p_static_bool:
		print(" -- %sSTATIC BINARY BUILD%s"%(fg("yellow"), attr(0)))
		
	print(" -- build %s%s%s service"%(fg("green"), p_name_str, attr(0)))
	print(" -- go_dir_path    - %s%s%s"%(fg("green"), p_go_dir_path_str, attr(0)))  
	print(" -- go_output_path - %s%s%s"%(fg("green"), p_go_output_path_str, attr(0)))  

	cwd_str = os.getcwd()
	os.chdir(p_go_dir_path_str) # change into the target main package dir

	# GO_GET
	print("go get")
	_, _, exit_code_int = gf_core_cli.run("go get -u")
	print("")
	print("")

	# STATIC_LINKING - when deploying to containers it is not always guaranteed that all
	#                  required libraries are present. so its safest to compile to a statically
	#                  linked lib.
	#                  build time a few times larger then regular, so slow for dev.
	# "-ldflags '-s'" - omit the symbol table and debug information.
	# "-a" - forces all packages to be rebuilt
	if p_static_bool:
		
		# https://golang.org/cmd/link/
		# IMPORTANT!! - "CGO_ENABLED=0" and "-installsuffix cgo" no longer necessary since golang 1.10.
		#               "CGO_ENABLED=0" we also dont want to disable since Rust libs are used in Go via CGO.
		# "-extldflags flags" - Set space-separated flags to pass to the external linker
		args_lst = [
			"CGO_ENABLED=0",
			"GOOS=linux",
			"go build",
			"-a",
			# "-installsuffix cgo",

			# LINKER_FLAGS
			# "-ldl" - "-l" provides lib path. links in  /usr/lib/libdl.so/.a
			#          this is needed to prevent Rust .a lib errors relating
			#          to undefined references to "dlsym","dladdr"
			('''-ldflags '-s -extldflags "-static -ldl"' ''').strip(),
			
			
			"-o %s"%(p_go_output_path_str),
		]
		c_str = " ".join(args_lst)

	# DYNAMIC_LINKING - fast build for dev.
	else:
		c_str = "go build -o %s"%(p_go_output_path_str)

	print(c_str)
	_, _, exit_code_int = gf_core_cli.run(c_str)

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
	#               want to fail with a non-zero exit code as well - this way other CI 
	#               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	os.chdir(cwd_str) # return to initial dir



#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#-------------
	# RUN
	arg_parser.add_argument("-run", action = "store", default = "build",
		help = '''
- '''+fg('yellow')+'test_go'+attr(0)+'''            - run app Go code tests
- '''+fg('yellow')+'test_py'+attr(0)+'''            - run app Py code tests
- '''+fg('yellow')+'build'+attr(0)+'''              - build app golang/web code
- '''+fg('yellow')+'build_containers'+attr(0)+'''   - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+''' - publish app Docker containers
- '''+fg('yellow')+'notify_completion'+attr(0)+'''  - notify remote HTTP endpoint of build completion
		''')

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument("-docker_sudo", action = "store_true", default=False,
		help = "specify if certain Docker CLI commands are to run with 'sudo'")

	#----------------------------
	# STATIC - boolean flag
	arg_parser.add_argument("-static", action = "store_true", default=False,
		help = "compile binaries with static linking")

	#----------------------------
	# TEST_CI - boolean flag
	# IMPORTANT!! - flag to indicate if tests are run in a CI env (non-local on a remote machine or in a CI environment) or localy on a dev machine.
	#               affects on how supporting services for tests are started up, if local dir paths should be used or if a CI layout is expected.
	arg_parser.add_argument("-test_ci", action = "store_true", default=False,
		help = "run tests in a CI environment")

	#-------------
	# ENV_VARS
	drone_commit_sha_str         = os.environ.get("DRONE_COMMIT_SHA", None) # Drone defined ENV var
	gf_docker_user_str           = os.environ.get("GF_DOCKER_USER", None)
	gf_docker_pass_str           = os.environ.get("GF_DOCKER_PASS", None)
	gf_notify_completion_url_str = os.environ.get("GF_NOTIFY_COMPLETION_URL", None)

	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"run":                          args_namespace.run,
		"drone_commit_sha":             drone_commit_sha_str,
		"gf_docker_user_str":           gf_docker_user_str,
		"gf_docker_pass_str":           gf_docker_pass_str,
		"gf_notify_completion_url_str": gf_notify_completion_url_str,
		"docker_sudo_bool":             args_namespace.docker_sudo,
		"static_bool":                  args_namespace.static,
		"test_ci_bool":                 args_namespace.test_ci,
	}

#--------------------------------------------------
main()