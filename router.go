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
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type group struct {
	pattern  string
	handlers []Handler
}

// Router represents a Macaron router layer.
type Router struct {
	m *Macaron

	groups              []group
	internalServerError func(*Context, error)

	// handlerWrapper is used to wrap arbitrary function from Handler to inject.FastInvoker.
	handlerWrapper func(Handler) Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Handle(method, pattern string, handlers []Handler) {
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handler, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}

		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}

	handlers = validateAndWrapHandlers(handlers, r.handlerWrapper)
	r.m.httprouter.Handle(method, pattern, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := r.m.createContext(w, req)
		c.params = params
		c.handlers = make([]Handler, 0, len(r.m.handlers)+len(handlers))
		c.handlers = append(c.handlers, r.m.handlers...)
		c.handlers = append(c.handlers, handlers...)
		c.run()
	})
}

func (r *Router) Group(pattern string, fn func(), h ...Handler) {
	r.groups = append(r.groups, group{pattern, h})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

// Get is a shortcut for r.Handle("GET", pattern, handlers)
func (r *Router) Get(pattern string, h ...Handler) {
	r.Handle("GET", pattern, h)
}

// Patch is a shortcut for r.Handle("PATCH", pattern, handlers)
func (r *Router) Patch(pattern string, h ...Handler) {
	r.Handle("PATCH", pattern, h)
}

// Post is a shortcut for r.Handle("POST", pattern, handlers)
func (r *Router) Post(pattern string, h ...Handler) {
	r.Handle("POST", pattern, h)
}

// Put is a shortcut for r.Handle("PUT", pattern, handlers)
func (r *Router) Put(pattern string, h ...Handler) {
	r.Handle("PUT", pattern, h)
}

// Delete is a shortcut for r.Handle("DELETE", pattern, handlers)
func (r *Router) Delete(pattern string, h ...Handler) {
	r.Handle("DELETE", pattern, h)
}

// Options is a shortcut for r.Handle("OPTIONS", pattern, handlers)
func (r *Router) Options(pattern string, h ...Handler) {
	r.Handle("OPTIONS", pattern, h)
}

// Head is a shortcut for r.Handle("HEAD", pattern, handlers)
func (r *Router) Head(pattern string, h ...Handler) {
	r.Handle("HEAD", pattern, h)
}

// Route is a shortcut for same handlers but different HTTP methods.
//
// Example:
// 		m.Route("/", "GET,POST", h)
func (r *Router) Route(pattern, methods string, h ...Handler) {
	for _, m := range strings.Split(methods, ",") {
		r.Handle(strings.TrimSpace(m), pattern, h)
	}
}

// NotFound configurates http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(handlers ...Handler) {
	handlers = validateAndWrapHandlers(handlers)
	r.m.httprouter.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		c := r.m.createContext(rw, req)
		c.handlers = make([]Handler, 0, len(r.m.handlers)+len(handlers))
		c.handlers = append(c.handlers, r.m.handlers...)
		c.handlers = append(c.handlers, handlers...)
		c.run()
	})
}

// InternalServerError configurates handler which is called when route handler returns
// error. If it is not set, default handler is used.
// Be sure to set 500 response code in your handler.
func (r *Router) InternalServerError(handlers ...Handler) {
	handlers = validateAndWrapHandlers(handlers)
	r.internalServerError = func(c *Context, err error) {
		c.index = 0
		c.handlers = handlers
		c.Map(err)
		c.run()
	}
}
