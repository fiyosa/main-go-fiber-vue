package cronrunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
)

type Runner struct {
	client *http.Client
}

func NewRunner() *Runner {
	return &Runner{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute runs the scheduled JS script.
func (r *Runner) Execute(ctx context.Context, script string, event map[string]interface{}, env map[string]string) (output string, resultErr string, err error) {
	loop := eventloop.NewEventLoop()
	loop.Start()
	defer loop.Stop()

	var buf bytes.Buffer
	var mu sync.Mutex

	done := make(chan struct{})
	var execErr error
	var scriptRej string

	loop.RunOnLoop(func(vm *goja.Runtime) {
		vm.Set("console", map[string]interface{}{
			"log": func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				for _, arg := range call.Arguments {
					fmt.Fprintf(&buf, "%v ", arg.Export())
				}
				buf.WriteString("\n")
				mu.Unlock()
				return goja.Undefined()
			},
			"error": func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				fmt.Fprintf(&buf, "[error] ")
				for _, arg := range call.Arguments {
					fmt.Fprintf(&buf, "%v ", arg.Export())
				}
				buf.WriteString("\n")
				mu.Unlock()
				return goja.Undefined()
			},
			"warn": func(call goja.FunctionCall) goja.Value {
				mu.Lock()
				fmt.Fprintf(&buf, "[warn] ")
				for _, arg := range call.Arguments {
					fmt.Fprintf(&buf, "%v ", arg.Export())
				}
				buf.WriteString("\n")
				mu.Unlock()
				return goja.Undefined()
			},
		})

		// Setup fetch
		vm.Set("fetch", func(call goja.FunctionCall) goja.Value {
			url := call.Argument(0).String()
			var method = "GET"
			var body io.Reader
			headers := make(http.Header)

			if len(call.Arguments) > 1 && !goja.IsUndefined(call.Argument(1)) && !goja.IsNull(call.Argument(1)) {
				opts := call.Argument(1).ToObject(vm)
				if m := opts.Get("method"); m != nil && !goja.IsUndefined(m) {
					method = m.String()
				}
				if b := opts.Get("body"); b != nil && !goja.IsUndefined(b) {
					body = strings.NewReader(b.String())
				}
				if h := opts.Get("headers"); h != nil && !goja.IsUndefined(h) {
					hObj := h.ToObject(vm)
					for _, key := range hObj.Keys() {
						headers.Set(key, hObj.Get(key).String())
					}
				}
			}

			promise, resolve, reject := vm.NewPromise()

			go func() {
				req, err := http.NewRequestWithContext(ctx, method, url, body)
				if err != nil {
					loop.RunOnLoop(func(vm *goja.Runtime) {
						reject(vm.NewGoError(err))
					})
					return
				}
				req.Header = headers

				resp, err := r.client.Do(req)
				if err != nil {
					loop.RunOnLoop(func(vm *goja.Runtime) {
						reject(vm.NewGoError(err))
					})
					return
				}
				defer resp.Body.Close()

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					loop.RunOnLoop(func(vm *goja.Runtime) {
						reject(vm.NewGoError(err))
					})
					return
				}

				loop.RunOnLoop(func(vm *goja.Runtime) {
					responseObj := vm.NewObject()
					responseObj.Set("status", resp.StatusCode)
					responseObj.Set("statusText", http.StatusText(resp.StatusCode))
					responseObj.Set("ok", resp.StatusCode >= 200 && resp.StatusCode < 300)
					responseObj.Set("url", resp.Request.URL.String())

					headersObj := vm.NewObject()
					for k, v := range resp.Header {
						if len(v) > 0 {
							headersObj.Set(k, v[0])
						}
					}
					responseObj.Set("headers", headersObj)

					responseObj.Set("text", func(call goja.FunctionCall) goja.Value {
						p, res, _ := vm.NewPromise()
						go func() {
							loop.RunOnLoop(func(vm *goja.Runtime) {
								res(vm.ToValue(string(respBody)))
							})
						}()
						return vm.ToValue(p)
					})

					responseObj.Set("json", func(call goja.FunctionCall) goja.Value {
						p, res, rej := vm.NewPromise()
						go func() {
							loop.RunOnLoop(func(vm *goja.Runtime) {
								val, err := vm.RunString("(" + string(respBody) + ")")
								if err != nil {
									rej(vm.NewGoError(err))
									return
								}
								res(val)
							})
						}()
						return vm.ToValue(p)
					})

					resolve(responseObj)
				})
			}()

			return vm.ToValue(promise)
		})

		// Run script definition
		_, err := vm.RunString(script)
		if err != nil {
			execErr = err
			close(done)
			return
		}

		vm.Set("event", event)
		vm.Set("env", env)

		vm.Set("__go_resolve", func(val goja.Value) {
			if val != nil && !goja.IsUndefined(val) {
				mu.Lock()
				fmt.Fprintf(&buf, "[return] %v\n", val.Export())
				mu.Unlock()
			}
			close(done)
		})

		vm.Set("__go_reject", func(errVal goja.Value) {
			scriptRej = fmt.Sprintf("%v", errVal.Export())
			close(done)
		})

		// Execute scheduled(event, env) wrapped in Promise.resolve
		_, err = vm.RunString(`
			try {
				Promise.resolve(scheduled(event, env)).then(__go_resolve).catch(__go_reject);
			} catch (e) {
				__go_reject(e);
			}
		`)
		if err != nil {
			execErr = err
			close(done)
		}
	})

	select {
	case <-done:
		mu.Lock()
		out := buf.String()
		mu.Unlock()
		if execErr != nil {
			return out, "", execErr
		}
		if scriptRej != "" {
			return out, scriptRej, nil
		}
		return out, "", nil
	case <-ctx.Done():
		return buf.String(), "", ctx.Err()
	case <-time.After(35 * time.Second):
		return buf.String(), "", fmt.Errorf("execution timeout")
	}
}
