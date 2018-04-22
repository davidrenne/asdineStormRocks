package controllers

import (
	"github.com/davidrenne/asdineStormRocks/constants"
	"github.com/davidrenne/asdineStormRocks/queries"
	"github.com/davidrenne/asdineStormRocks/sessionFunctions"
	"github.com/davidrenne/asdineStormRocks/viewModel"
)

func (self *TransactionsController) RollBackFromSnackbar(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.AppViewModel
	vm.Parse(state)

	err := queries.Transactions.Rollback(context, "Snackbar Rollback", vm.SnackBarUndoTransactionId)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_ROLLBACK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_RERENDER, constants.TRANSACTION_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}
