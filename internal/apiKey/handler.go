package apiKey

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"restapi/internal/config"
	"restapi/internal/handlers"
	"restapi/pkg/logging"
)

type handler struct {
	conf       *config.Config
	repository Repository
	logger     *logging.Logger
}

func processDbResult(bdQueryResult ApiKeyResult, response http.ResponseWriter) http.ResponseWriter {
	if bdQueryResult.ResponseBody == "" {
		response.Header().Add("X-Accel-Redirect", bdQueryResult.Location)
		response.Header().Add("X-API-Account-Name", bdQueryResult.AccountName)
		return response
	} else {
		response.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(response, bdQueryResult.ResponseBody)
		return response
	}
}

func constructHttpResponse(header string, repository Repository, conf *config.Config, response http.ResponseWriter) http.ResponseWriter {
	found := false // Cache.Get(header)
	if found {
		//	return processDbResult(resultCache.(ApiKeyResult), response)
	} else {
		// if we have header hardcoded in config
		if val, ok := conf.FixPlan[header]; ok {
			var res ApiKeyResult
			res.ResponseBody = ""
			res.AccountName = val.AccountName
			res.Location = val.Location
			return processDbResult(res, response)
		}

		resultDB, err := repository.FindeApiKey(context.TODO(), header)
		if err != nil {
			fmt.Println(err) // we have DB error but any way we have the response
			return processDbResult(resultDB, response)
		} else {
			//Cache.Set(header, resultDB, cache.DefaultExpiration)
			return processDbResult(resultDB, response)
		}
	}
	return response
}

func NewHandler(cf *config.Config, repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		conf:       cf,
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	//router.HandlerFunc(http.MethodGet, authorsURL, apperror.Middleware(h.GetList))
	//router.Handle("/",h.parseHeaders)
	router.GET("/", h.parseHeaders)
}

func (h *handler) parseHeaders(response http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var regex *regexp.Regexp
	const regexString = "^BQY[0-9a-zA-Z]{29}$"
	regex, _ = regexp.Compile(regexString)
	header := req.Header.Get(h.conf.NetworkConf.HeaderName)
	if header != "" { // we have empty string when there is no such header
		if regex.MatchString(header) {
			response = constructHttpResponse(header, h.repository, h.conf, response)
			return
		} else {
			response.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(response, "Api key provided has illegal format, check your API Key in the profile page")
			return
		}
	}
	response.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(response, "Api key not provided, get your API Key in the profile page.")
	return
}
