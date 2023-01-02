namespace go common

enum Err{
	// out/inner gateway [100000, 200000)
    GateWayBadRequest = 100000,
    GateWayMethodNotFound = 100001,
    GateWayHttpNewRequestError = 100002,
    GateWayServerHandlerError = 100003,
    GateWayServerInnerCallError = 100004,

	// interaction [10000 - 20000)
    InteractionBadRequest = 10000,

	// payment [20000 - 30000)
    PaymentBadRequest = 20000,
    PaymentDbInteralError = 20001,
    PaymentStationInteralError = 20002,

	// im [30000 - 40000)
    IMBadRequest = 30000,
}




