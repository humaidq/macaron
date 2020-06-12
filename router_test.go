// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package macaron

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Router_Handle(t *testing.T) {
	test_Router_Handle(t, false)
}
func Test_Router_FastInvoker_Handle(t *testing.T) {
	test_Router_Handle(t, true)
}

func test_Router_Handle(t *testing.T, isFast bool) {
	Convey("Register all HTTP methods routes", t, func() {
		m := New()

		m.Get("/get", func() string {
			return "GET"
		})
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/get", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "GET")

		m.Patch("/patch", func() string {
			return "PATCH"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("PATCH", "/patch", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "PATCH")

		m.Post("/post", func() string {
			return "POST"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("POST", "/post", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "POST")

		m.Put("/put", func() string {
			return "PUT"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("PUT", "/put", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "PUT")

		m.Delete("/delete", func() string {
			return "DELETE"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("DELETE", "/delete", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "DELETE")

		m.Options("/options", func() string {
			return "OPTIONS"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("OPTIONS", "/options", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "OPTIONS")

		m.Head("/head", func() string {
			return "HEAD"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("HEAD", "/head", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldHaveLength, 0)

		m.Route("/route", "GET,POST", func() string {
			return "ROUTE"
		})
		resp = httptest.NewRecorder()
		req, err = http.NewRequest("POST", "/route", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "ROUTE")
	})
}

func Test_Router_Group(t *testing.T) {
	Convey("Register route group", t, func() {
		m := New()
		m.Group("/api", func() {
			m.Group("/v1", func() {
				m.Get("/list", func() string {
					return "Well done!"
				})
			})
		})
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/list", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "Well done!")
	})
}

func Test_Router_NotFound(t *testing.T) {
	Convey("Custom not found handler", t, func() {
		m := New()
		m.Get("/", func() {})
		m.NotFound(func() string {
			return "Custom not found"
		})
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/404", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "Custom not found")
	})
}

func Test_Router_InternalServerError(t *testing.T) {
	Convey("Custom internal server error handler", t, func() {
		m := New()
		m.Get("/", func() error {
			return errors.New("Custom internal server error")
		})
		m.InternalServerError(func(rw http.ResponseWriter, err error) {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
		})
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, 500)
		So(resp.Body.String(), ShouldEqual, "Custom internal server error")
	})
}

func Test_Router_splat(t *testing.T) {
	Convey("Register router with glob", t, func() {
		m := New()
		m.Get("/*glob", func(ctx *Context) string {
			return ctx.Params("glob")
		})
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/hahaha", nil)
		So(err, ShouldBeNil)
		m.ServeHTTP(resp, req)
		So(resp.Body.String(), ShouldEqual, "/hahaha")
	})
}
