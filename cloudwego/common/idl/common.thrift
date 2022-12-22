namespace go common

enum Err{
	// inner/out gateway [100000, 200000)
    GateWayBadRequest = 100000,

	// interaction [10000 - 20000)
    InteractionBadRequest = 10000,

	// payment [20000 - 30000)
    PaymentBadRequest = 20000,
    PaymentDbInteralError = 20001,

	// im [30000 - 40000)
    IMBadRequest = 30000,
}




