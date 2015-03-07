package zerver_rest

import "net/http"

type StatusResponse interface {
	ReportContinue()           // 100
	ReportSwitchingProtocols() // 101

	ReportOK()                   // 200
	ReportCreated()              // 201
	ReportAccepted()             // 202
	ReportNonAuthoritativeInfo() // 203
	ReportNoContent()            // 204
	ReportResetContent()         // 205
	ReportPartialContent()       // 206

	ReportMultipleChoices()   // 300
	ReportMovedPermanently()  // 301
	ReportFound()             // 302
	ReportSeeOther()          // 303
	ReportNotModified()       // 304
	ReportUseProxy()          // 305
	ReportTemporaryRedirect() // 307

	ReportBadRequest()                   // 400
	ReportUnauthorized()                 // 401
	ReportPaymentRequired()              // 402
	ReportForbidden()                    // 403
	ReportNotFound()                     // 404
	ReportMethodNotAllowed()             // 405
	ReportNotAcceptable()                // 406
	ReportProxyAuthRequired()            // 407
	ReportRequestTimeout()               // 408
	ReportConflict()                     // 409
	ReportGone()                         // 410
	ReportLengthRequired()               // 411
	ReportPreconditionFailed()           // 412
	ReportRequestEntityTooLarge()        // 413
	ReportRequestURITooLong()            // 414
	ReportUnsupportedMediaType()         // 415
	ReportRequestedRangeNotSatisfiable() // 416
	ReportExpectationFailed()            // 417
	ReportTeapot()                       // 418

	ReportInternalServerError()     // 500
	ReportNotImplemented()          // 501
	ReportBadGateway()              // 502
	ReportServiceUnavailable()      // 503
	ReportGatewayTimeout()          // 504
	ReportHTTPVersionNotSupported() // 505
}

//100
func (resp *response) ReportContinue() {
	resp.ReportStatus(http.StatusContinue)
}

//101
func (resp *response) ReportSwitchingProtocols() {
	resp.ReportStatus(http.StatusSwitchingProtocols)
}

//200
func (resp *response) ReportOK() {
	resp.ReportStatus(http.StatusOK)
}

//201
func (resp *response) ReportCreated() {
	resp.ReportStatus(http.StatusCreated)
}

//202
func (resp *response) ReportAccepted() {
	resp.ReportStatus(http.StatusAccepted)
}

//203
func (resp *response) ReportNonAuthoritativeInfo() {
	resp.ReportStatus(http.StatusNonAuthoritativeInfo)
}

//204
func (resp *response) ReportNoContent() {
	resp.ReportStatus(http.StatusNoContent)
}

//205
func (resp *response) ReportResetContent() {
	resp.ReportStatus(http.StatusResetContent)
}

//206
func (resp *response) ReportPartialContent() {
	resp.ReportStatus(http.StatusPartialContent)
}

//300
func (resp *response) ReportMultipleChoices() {
	resp.ReportStatus(http.StatusMultipleChoices)
}

//301
func (resp *response) ReportMovedPermanently() {
	resp.ReportStatus(http.StatusMovedPermanently)
}

//302
func (resp *response) ReportFound() {
	resp.ReportStatus(http.StatusFound)
}

//303
func (resp *response) ReportSeeOther() {
	resp.ReportStatus(http.StatusSeeOther)
}

//304
func (resp *response) ReportNotModified() {
	resp.ReportStatus(http.StatusNotModified)
}

//305
func (resp *response) ReportUseProxy() {
	resp.ReportStatus(http.StatusUseProxy)
}

//307
func (resp *response) ReportTemporaryRedirect() {
	resp.ReportStatus(http.StatusTemporaryRedirect)
}

//400
func (resp *response) ReportBadRequest() {
	resp.ReportStatus(http.StatusBadRequest)
}

//401
func (resp *response) ReportUnauthorized() {
	resp.ReportStatus(http.StatusUnauthorized)
}

//402
func (resp *response) ReportPaymentRequired() {
	resp.ReportStatus(http.StatusPaymentRequired)
}

//403
func (resp *response) ReportForbidden() {
	resp.ReportStatus(http.StatusForbidden)
}

//404
func (resp *response) ReportNotFound() {
	resp.ReportStatus(http.StatusNotFound)
}

//405
func (resp *response) ReportMethodNotAllowed() {
	resp.ReportStatus(http.StatusMethodNotAllowed)
}

//406
func (resp *response) ReportNotAcceptable() {
	resp.ReportStatus(http.StatusNotAcceptable)
}

//407
func (resp *response) ReportProxyAuthRequired() {
	resp.ReportStatus(http.StatusProxyAuthRequired)
}

//408
func (resp *response) ReportRequestTimeout() {
	resp.ReportStatus(http.StatusRequestTimeout)
}

//409
func (resp *response) ReportConflict() {
	resp.ReportStatus(http.StatusConflict)
}

//410
func (resp *response) ReportGone() {
	resp.ReportStatus(http.StatusGone)
}

//411
func (resp *response) ReportLengthRequired() {
	resp.ReportStatus(http.StatusLengthRequired)
}

//412
func (resp *response) ReportPreconditionFailed() {
	resp.ReportStatus(http.StatusPreconditionFailed)
}

//413
func (resp *response) ReportRequestEntityTooLarge() {
	resp.ReportStatus(http.StatusRequestEntityTooLarge)
}

//414
func (resp *response) ReportRequestURITooLong() {
	resp.ReportStatus(http.StatusRequestURITooLong)
}

//415
func (resp *response) ReportUnsupportedMediaType() {
	resp.ReportStatus(http.StatusUnsupportedMediaType)
}

//416
func (resp *response) ReportRequestedRangeNotSatisfiable() {
	resp.ReportStatus(http.StatusRequestedRangeNotSatisfiable)
}

//417
func (resp *response) ReportExpectationFailed() {
	resp.ReportStatus(http.StatusExpectationFailed)
}

//418
func (resp *response) ReportTeapot() {
	resp.ReportStatus(http.StatusTeapot)
}

//500
func (resp *response) ReportInternalServerError() {
	resp.ReportStatus(http.StatusInternalServerError)
}

//501
func (resp *response) ReportNotImplemented() {
	resp.ReportStatus(http.StatusNotImplemented)
}

//502
func (resp *response) ReportBadGateway() {
	resp.ReportStatus(http.StatusBadGateway)
}

//503
func (resp *response) ReportServiceUnavailable() {
	resp.ReportStatus(http.StatusServiceUnavailable)
}

//504
func (resp *response) ReportGatewayTimeout() {
	resp.ReportStatus(http.StatusGatewayTimeout)
}

//505
func (resp *response) ReportHTTPVersionNotSupported() {
	resp.ReportStatus(http.StatusHTTPVersionNotSupported)
}
