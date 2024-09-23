package gowebserver

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

type AddFunc func(key string, value string)

type AddHeaders struct {
	Regex   string            `json:"Regex"`
	Headers map[string]string `json:"Headers"`
}

type Header struct {
	RequestHeaders  []AddHeaders `json:"RequestHeaders"`
	ResponseHeaders []AddHeaders `json:"ResponseHeaders"`
}

type Skipper []string
type Headers []Header

type IndexHeaders struct {
	RequestHeaders  map[string]string `json:"RequestHeaders"`
	ResponseHeaders map[string]string `json:"ResponseHeaders"`
}

type Proxy struct {
	Path    string  `json:"Path"`
	Target  string  `json:"Target"`
	Skipper Skipper `json:"Skipper"`
	Headers Headers `json:"Headers"`
}

type Router struct {
	Path         string       `json:"Path"`
	FilePath     string       `json:"FilePath"`
	Headers      Headers      `json:"Headers"`
	IndexHeaders IndexHeaders `json:"IndexHeaders"`
}

type Proxys []Proxy
type Routers []Router

type ServerConfig struct {
	Port    int     `json:"Port"`
	Proxys  Proxys  `json:"Proxys"`
	Routers Routers `json:"Routers"`
}

func (a *AddHeaders) Add(path string, addFunc AddFunc) {
	if a.Regex != "" {
		re := regexp.MustCompile(a.Regex)
		if !re.MatchString(path) {
			return // 如果路径不匹配，则跳出
		}
	}

	for key, value := range a.Headers {
		addFunc(key, value)
	}
}

// var goServerHeaderKey = "Go-Server-CustomKey"

func (h *Header) Add(r *http.Request, w *echo.Response) {

	if len(h.RequestHeaders) > 0 {
		for _, rh := range h.RequestHeaders {
			rh.Add(r.URL.Path, func(key string, value string) {
				r.Header.Set(key, value)
				// keys := r.Header.Get(goServerHeaderKey)
				// if keys == "" {
				// 	keys = keys + key
				// } else {
				// 	keys = keys + "," + key
				// }
				// r.Header.Set(goServerHeaderKey, keys)
			})
		}
	}

	if len(h.ResponseHeaders) > 0 {
		for _, rh := range h.ResponseHeaders {
			rh.Add(r.URL.Path, func(key string, value string) {
				w.Header().Set(key, value)
				// keys := r.Header.Get(goServerHeaderKey)
				// if keys == "" {
				// 	keys = keys + key
				// } else {
				// 	keys = keys + "," + key
				// }
				// r.Header.Set(goServerHeaderKey, keys)
			})
		}
	}
}

func (s Skipper) Test(path string) bool {
	if len(s) > 0 {
		for _, r := range s {
			regex := regexp.MustCompile(r)
			if regex.MatchString(path) {
				return true
			}
		}
	}
	return false
}

func (h Headers) Add(r *http.Request, w *echo.Response) {
	for _, header := range h {
		header.Add(r, w)
	}
}

func (i IndexHeaders) Add(r *http.Request, w *echo.Response) {
	for key, value := range i.RequestHeaders {
		r.Header.Set(key, value)
	}

	for key, value := range i.ResponseHeaders {
		w.Header().Set(key, value)
	}
}

func (p Proxys) Add(e *echo.Echo) {
	if len(p) > 0 {
		fwd := newForward()

		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				path := c.Request().URL.Path

				if path == healthPath {
					return next(c)
				}

			Loop_Proxy:
				for _, proxy := range p {
					if !strings.HasPrefix(path, proxy.Path) {
						continue
					}

					if proxy.Skipper.Test(path) {
						break Loop_Proxy
					}

					target, err := url.Parse(proxy.Target)
					if err != nil {
						logger.Printf("url.Parse Error: %s\n", proxy.Target)
						continue
					}

					logger.Printf("proxy: %s, target: %s\n", path, proxy.Target)

					proxy.Headers.Add(c.Request(), c.Response())

					c.Request().URL = target
					fwd.ServeHTTP(c.Response(), c.Request())

					return nil
				}

				return next(c)
			}
		})
	}
}

func (r Routers) Add(e *echo.Echo) {
	if len(r) > 0 {
		for _, router := range r {
			e.Static(router.Path, router.FilePath)
			e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					path := c.Request().URL.Path

					if path == healthPath {
						return next(c)
					}

					if !strings.HasPrefix(path, router.Path) {
						return next(c)
					}

					router.Headers.Add(c.Request(), c.Response())

					if path == router.Path {
						router.IndexHeaders.Add(c.Request(), c.Response())
					}
					err := next(c)

					if err != nil {
						if he, ok := err.(*echo.HTTPError); ok {
							if he.Code == http.StatusNotFound {
								// 重定向到index.html
								router.IndexHeaders.Add(c.Request(), c.Response())
								return c.File(router.FilePath + "/index.html")
							}
						}
					}
					return err
				}
			})
		}
	}
}
