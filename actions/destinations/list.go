package destinations

import (
	"encoding/json"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// list will fetch a list of destinations
func (a *Action) list(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)

	// todo GetJSON should be able to marshal into a model
	metadataReq := params.GetJSON(bux.ModelMetadata.String())
	var metadata *bux.Metadata
	if len(metadataReq) > 0 {
		// marshal the metadata into the Metadata model
		metaJSON, _ := json.Marshal(metadataReq)
		_ = json.Unmarshal(metaJSON, &metadata)
	}

	// Record a new transaction (get the hex from parameters)a
	var err error
	var destinations []*bux.Destination
	if destinations, err = a.Services.Bux.GetDestinations(
		req.Context(),
		reqXPub,
		metadata,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(destinations))
}
