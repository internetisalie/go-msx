// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restfulcontext

import "github.com/emicklei/go-restful"

type RouteBuilderFunc func(*restful.RouteBuilder)
