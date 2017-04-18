//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with p work for additional information
// regarding copyright ownership.  The ASF licenses p file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use p file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
// @project fatima
// @author DeockJin Chung (jin.freestyle@gmail.com)
// @date 2017. 3. 6. PM 2:56
//

package builder


const (
	BUILTIN_VARIABLE_HOME            = "${var.builtin.user.home}"
	BUILTIN_VARIABLE_FATIMA_HOME     = "${var.builtin.fatima.home}"
	BUILTIN_VARIABLE_LOCAL_IPADDRESS = "${var.builtin.local.ipaddress}"
	BUILTIN_VARIABLE_YYYYMM          = "${var.builtin.date.yyyymm}"
	BUILTIN_VARIABLE_YYYYMMDD        = "${var.builtin.date.yyyymmdd}"
	BUILTIN_VARIABLE_APP_NAME        = "${var.builtin.app.name}"
	BUILTIN_VARIABLE_APP_FOLDER_DATA = "${var.builtin.app.folder.data}"

	GLOBAL_DEFINE_PACKAGE_HOSTNAME = "var.global.package.hostname"
	GLOBAL_DEFINE_PACKAGE_GROUPNAME = "var.global.package.groupname"
	GLOBAL_DEFINE_PACKAGE_NAME = "var.global.package.name"
)

const (
	GOFATIMA_PROP_PPROF_ADDRESS		= "gofatima.pprof.address"		// e.g :6060, localhost:6060
)