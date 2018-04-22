package controllers

import (
	_ "github.com/davidrenne/asdineStormRocks/constants"
	"github.com/davidrenne/asdineStormRocks/models/v1/model"
	"github.com/davidrenne/asdineStormRocks/queries"
	"github.com/davidrenne/asdineStormRocks/viewModel"

	"github.com/davidrenne/asdineStormRocks/constants"
	"github.com/davidrenne/asdineStormRocks/sessionFunctions"
)

func (self *ServerSettingsModifyController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	controller := ServerSettingsModifyController{}
	controller.Load(context, uriParams, respond)
}

func (self *ServerSettingsModifyController) PopulateVm(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse, callRespond bool) (vm viewModel.ServerSettingsModifyViewModel) {
	vm.LoadDefaultState()

	selectedTab, ok := uriParams["Tab"]
	if ok {
		vm.SelectedTab = selectedTab
	}

	var retrievedLockout model.ServerSetting
	retrievedLockout, err := queries.ServerSettings.LoginAttempts()
	if err != nil {
		if callRespond {
			respond(PARAM_REDIRECT_NONE, "ERROR WITH QUERY FOR SERVER SETTING GATEWAY", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		}
		return
	}

	vm.LockoutSettings.Lockout = model.ServerSetting{
		Category: "users",
		Key:      "lockoutAttempts",
	}
	if retrievedLockout.Key == "lockoutAttempts" {
		vm.LockoutSettings.Lockout = retrievedLockout
	}

	vm.TimeZones = model.TimeZoneLocations
	vm.TimeZone, _ = queries.ServerSettings.ById(constants.SERVER_SETTING_TIMEZONE)

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
	return vm
}

func (self *ServerSettingsModifyController) Load(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	controller := ServerSettingsModifyController{}
	controller.PopulateVm(context, uriParams, respond, true)
}
