package infra

/*func AdaptHandler(handler interface{}) (gin.HandlerFunc, error) {

	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil, errors.Errorf("not a function: %+v", handler)
	} else if handlerType.IsVariadic() {
		return nil, errors.Errorf("function must not be variadic")
	} else if handlerType.NumIn() != 2 {
		return nil, errors.Errorf("function must have two arguments")
	}

	contextPtrParamType := handlerType.In(0)
	if contextPtrParamType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("function 1st parameter must be of type '*gin.Context'")
	}
	contextParamType := contextPtrParamType.Elem()
	if contextParamType.Kind() != reflect.Struct {
		return nil, errors.Errorf("function 1st parameter must be of type '*gin.Context'")
	} else if contextParamType.PkgPath() != "github.com/gin-gonic/gin" {
		return nil, errors.Errorf("function 1st parameter must be of type '*gin.Context'")
	} else if contextParamType.Name() != "Context" {
		return nil, errors.Errorf("function 1st parameter must be of type '*gin.Context'")
	}

	bindingParamType := handlerType.In(1)
	if bindingParamType.Kind() != reflect.Struct {
		return nil, errors.Errorf("function 2nd parameter must be a binding-target struct")
	}

	if handlerType.NumOut() != 2 {
		return nil, errors.Errorf("function must return have two return values (result struct, and an error)")
	}
	resultParamPtrType := handlerType.Out(0)
	if resultParamPtrType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("function 1st return value must be a pointer to a struct")
	}
	resultParamType := resultParamPtrType.Elem()
	if resultParamType.Kind() != reflect.Struct {
		return nil, errors.Errorf("function 1st return value must be a struct")
	}
	errorParamType := handlerType.Out(1)
	if errorParamType.Kind() != reflect.Interface {
		return nil, errors.Errorf("function 2nd return value must be error")
	} else if errorParamType.PkgPath() != "" {
		return nil, errors.Errorf("function 2nd return value must be error")
	} else if errorParamType.Name() != "error" {
		return nil, errors.Errorf("function 2nd return value must be error")
	}

	handlerFunc := reflect.ValueOf(handler)
	return func(c *gin.Context) {
		acceptedMimeType := c.NegotiateFormat(supportedContentTypes...)
		if acceptedMimeType == "" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		inputParamValue := reflect.New(bindingParamType)
		if err := c.Bind(inputParamValue.Interface()); err != nil {
			// this will reach the error handler middleware to be handled there, so just terminate early
			return
		}

		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(c)
		in[1] = inputParamValue.Elem()

		out := handlerFunc.Call(in)

		errOut := out[1]
		if !errOut.IsNil() {
			c.Error(errOut.Interface().(error))
			return
		}

		resultOut := out[0]
		statusCode := http.StatusOK

		if !resultOut.IsNil() {

			switch c.Request.Method {
			case http.MethodGet:

			case http.MethodDelete:
			case http.MethodHead:
			case http.MethodOptions:
			case http.MethodPatch:
			case http.MethodPost:
			case http.MethodPut:
			case http.MethodTrace:
			default:
			}
			switch acceptedMimeType {
			case gin.MIMEJSON:
				c.JSON(statusCode, map[string]interface{}{"errors": errors})

			case gin.MIMEHTML:
				c.Header("Content-Type", gin.MIMEHTML)
				c.Status(http.StatusOK) // send HTTP 200 since this is most probably a browser
				if err := errorsTemplate.Execute(c.Writer, errors); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				}

			case gin.MIMEXML, gin.MIMEXML2:
				c.XML(statusCode, map[string]interface{}{"errors": errors})

			case gin.MIMEPlain:
				c.Header("Content-Type", gin.MIMEPlain)
				c.Status(statusCode)
				for _, err := range errors {
					_, _ = c.Writer.WriteString(err.Error() + "\n")
				}

			case "text/yaml", "application/yaml", gin.MIMEYAML:
				c.YAML(statusCode, map[string]interface{}{"errors": errors})

			default:
				c.AbortWithStatus(http.StatusNotAcceptable)
			}
		}

	}, nil
}
*/
