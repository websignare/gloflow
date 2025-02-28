// SPDX-License-Identifier: GPL-2.0
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

package gf_images_jobs

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func Init(pImagesStoreLocalDirPathStr string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pVideoStoreLocalDirPathStr            string,
	pMediaDomainStr                       string,
	pConfig                               *gf_images_core.GFconfig,
	pImageStorage                         *gf_images_storage.GFimageStorage,
	pS3info                               *gf_core.GFs3Info,
	pRuntimeSys                           *gf_core.RuntimeSys) gf_images_jobs_core.JobsMngr {

	lifecycleCallbacks := &gf_images_jobs_core.GF_jobs_lifecycle_callbacks{
		Job_type__transform_imgs__fun: func() *gf_core.GFerror {
			// RUST
			// FIX!! - this just runs Rust job code for testing.
			//         pass in proper job_cmd argument.
			// run_job_rust()

			return nil
		},

		Job_type__uploaded_imgs__fun: func() *gf_core.GFerror {
			// RUST
			// FIX!! - this just runs Rust job code for testing.
			//         pass in proper job_cmd argument.
			// run_job_rust()
			
			return nil	
		},
	}

	jobsMngrCh := gf_images_jobs_core.JobsMngrInit(pImagesStoreLocalDirPathStr,
		pImagesThumbnailsStoreLocalDirPathStr,
		pVideoStoreLocalDirPathStr,
		pMediaDomainStr,
		lifecycleCallbacks,
		pConfig,
		pImageStorage,
		pS3info,
		pRuntimeSys)

	return jobsMngrCh
}