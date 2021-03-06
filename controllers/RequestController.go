package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"time"

	"github.com/astaxie/beego"
)

type RequestController struct {
	BaseController
}

//Request Ad
func (this *RequestController) RequestAd() {

	t1 := time.Now().UnixNano()

	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug(this.Ctx.Input.RequestBody)
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else if ValidRequest(&adRequest) != true {
		adResponse.StatusCode = lib.ERROR_REQUIRED_FIELD_MISSING
	} else {
		adRequest.Did = lib.GenerateBid(adRequest.AdspaceKey)

		adRequest.RequestTime = time.Now().Unix()
		tmp := engine.InvokeDemand(&adRequest)

		if tmp == nil {
			adResponse.StatusCode = lib.ERROR_NO_DEMAND_ERROR
			adResponse.Bid = adRequest.Bid
			adResponse.Did = adRequest.Did
			adResponse.AdspaceKey = adRequest.AdspaceKey
			adResponse.ResponseTime = time.Now().Unix()
		} else {
			adResponse = tmp
		}

		//only running pmp adspace need track request log
		if adResponse.StatusCode != lib.ERROR_NO_PMP_ADSPACE_ERROR {
			adRequest.StatusCode = adResponse.StatusCode

			//这里添加代码
			t2 := time.Now().UnixNano()
			adRequest.ProcessDuration = (t2 - t1) / 1000000
			engine.SendRequestLog(&adRequest, 1)
		}

	}

	commonResponse := GetCommonResponse(adResponse)

	this.Data["json"] = commonResponse
	this.ServeJSON()

}
