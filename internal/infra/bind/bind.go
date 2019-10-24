package bind

import (
	"encoding/json"
	"encoding/xml"
	"mime"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// TODO: Support bind result validation via JSON schema on Go struct
// TODO: Support bind result validation via optional implementation of a "Validator" interface

func Bind(r *http.Request, w http.ResponseWriter, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		log.Error().Str("type", targetValue.Type().String()).Msg("Bind target must be a struct pointer")
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	targetValue = targetValue.Elem()
	targetType := targetValue.Type()
	for i := 0; i < targetValue.NumField(); i++ {
		fieldValue := targetValue.Field(i)
		fieldType := targetType.Field(i)
		if headerSpec := fieldType.Tag.Get("header"); headerSpec != "" {
			headerName := ""
			required := false
			for _, token := range strings.Split(headerSpec, ",") {
				if token == "required" {
					required = true
				} else if headerName != "" {
					log.Error().
						Str("fieldName", fieldType.Name).
						Str("tag", headerSpec).
						Msg("Invalid tag")
					w.WriteHeader(http.StatusInternalServerError)
					return false
				} else {
					headerName = token
				}
			}
			if statusCode, err := bindRequestHeader(r, headerName, required, fieldValue); statusCode > 0 {
				if err != nil {
					log.Error().
						Err(err).
						Str("fieldName", fieldType.Name).
						Msg("Header binding failed")
				}
				w.WriteHeader(statusCode)
				return false
			}
		} else if querySpec := fieldType.Tag.Get("query"); querySpec != "" {
			parameterName := ""
			required := false
			for _, token := range strings.Split(querySpec, ",") {
				if token == "required" {
					required = true
				} else if parameterName != "" {
					log.Error().
						Str("fieldName", fieldType.Name).
						Str("tag", querySpec).
						Msg("Invalid tag")
					w.WriteHeader(http.StatusInternalServerError)
					return false
				} else {
					parameterName = token
				}
			}
			if statusCode, err := bindRequestQuery(r, parameterName, required, fieldValue); statusCode > 0 {
				if err != nil {
					log.Error().
						Err(err).
						Str("fieldName", fieldType.Name).
						Msg("Query binding failed")
				}
				w.WriteHeader(statusCode)
				return false
			}
		} else if cookieSpec := fieldType.Tag.Get("cookie"); cookieSpec != "" {
			cookieName := ""
			required := false
			for _, token := range strings.Split(cookieSpec, ",") {
				if token == "required" {
					required = true
				} else if cookieName != "" {
					log.Error().
						Str("fieldName", fieldType.Name).
						Str("tag", querySpec).
						Msg("Invalid tag")
					w.WriteHeader(http.StatusInternalServerError)
					return false
				} else {
					cookieName = token
				}
			}
			if statusCode, err := bindRequestCookie(r, cookieName, required, fieldValue); statusCode > 0 {
				if err != nil {
					log.Error().
						Err(err).
						Str("fieldName", fieldType.Name).
						Msg("Query binding failed")
				}
				w.WriteHeader(statusCode)
				return false
			}
		} else if bodySpec, ok := fieldType.Tag.Lookup("body"); ok {
			if bodySpec != "" {
				log.Error().
					Str("fieldName", fieldType.Name).
					Str("bodyTag", bodySpec).
					Msg("Invalid body tag")
				w.WriteHeader(http.StatusInternalServerError)
				return false
			}
			if statusCode, err := bindRequestBody(r, fieldValue); statusCode > 0 {
				if err != nil {
					log.Error().
						Err(err).
						Str("fieldName", fieldType.Name).
						Msg("Body payload binding failed")
				}
				w.WriteHeader(statusCode)
				return false
			}
		}
	}
	return true
}

func bindRequestHeader(r *http.Request, headerName string, required bool, v reflect.Value) (statusCode int, err error) {
	if statusCode, err := bindValues(r.Header[headerName], v, required); err != nil {
		return statusCode, err
	}
	return 0, nil
}

func bindRequestQuery(r *http.Request, queryParamName string, required bool, v reflect.Value) (statusCode int, err error) {
	// TODO: implement bindRequestQuery
	panic("query tag not implemented")
}

func bindRequestCookie(r *http.Request, cookieName string, required bool, v reflect.Value) (statusCode int, err error) {
	// TODO: implement bindRequestCookie
	panic("cookie tag not implemented")
}

func bindRequestBody(r *http.Request, fv reflect.Value) (statusCode int, err error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return http.StatusUnsupportedMediaType, nil
	}
	switch mediaType {
	case "application/x-yaml", "application/yaml", "text/yaml":
		decoder := yaml.NewDecoder(r.Body)
		decoder.SetStrict(false)
		if err := decoder.Decode(fv); err != nil {
			return http.StatusBadRequest, errors.Wrap(err, "decoding failed")
		}

	case "application/json", "text/json":
		decoder := json.NewDecoder(r.Body)
		decoder.UseNumber()
		decoder.DisallowUnknownFields()
		vptr := reflect.New(fv.Type())
		if err := decoder.Decode(vptr.Interface()); err != nil {
			return http.StatusBadRequest, errors.Wrap(err, "decoding failed")
		} else {
			fv.Set(vptr.Elem())
		}

	case "application/xml", "text/xml":
		decoder := xml.NewDecoder(r.Body)
		if err := decoder.Decode(fv); err != nil {
			return http.StatusBadRequest, errors.Wrap(err, "decoding failed")
		}

	default:
		return http.StatusBadRequest, nil
	}
	return 0, nil
}

func bindValues(values []string, v reflect.Value, required bool) (statusCode int, err error) {
	switch v.Kind() {
	case reflect.Bool:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if b, err := strconv.ParseBool(values[0]); err != nil {
			return http.StatusBadRequest, nil
		} else {
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if i, err := strconv.ParseInt(values[0], 10, 64); err != nil {
			return http.StatusBadRequest, nil
		} else {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if i, err := strconv.ParseUint(values[0], 10, 64); err != nil {
			return http.StatusBadRequest, nil
		} else {
			v.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else if f, err := strconv.ParseFloat(values[0], 64); err != nil {
			return http.StatusBadRequest, nil
		} else {
			v.SetFloat(f)
		}
	case reflect.Ptr:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else {
			targetType := v.Type().Elem()
			targetValue := reflect.New(targetType)
			if statusCode, err := bindValues(values, targetValue, true); err != nil || statusCode > 0 {
				return statusCode, err
			} else {
				v.Set(targetValue)
			}
		}
	case reflect.String, reflect.Interface:
		if len(values) == 0 {
			if required {
				return http.StatusBadRequest, nil
			} else {
				v.Set(reflect.Zero(v.Type()))
			}
		} else {
			v.SetString(values[0])
		}
	case reflect.Slice:
		v.SetLen(len(values))
		v.SetCap(len(values))
		for i, value := range values {
			vi := v.Index(i)
			if statusCode, err := bindValues([]string{value}, vi, true); err != nil || statusCode > 0 {
				return statusCode, err
			}
		}
	default:
		panic("not implemented")
	}
	return 0, nil
}
