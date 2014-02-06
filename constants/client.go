package constants

type ClientVersion uint16

const (
	ClientVersion1109 ClientVersion = 1109
	ClientVersion1110 ClientVersion = 1110
	ClientVersion1111 ClientVersion = 1111
	ClientVersion1112 ClientVersion = 1112
	ClientVersion1113 ClientVersion = 1113
	ClientVersion1114 ClientVersion = 1114
)

const (
	ClientVersionMinimum ClientVersion = ClientVersion1109
	ClientVersionMaximum ClientVersion = ClientVersion1114
)

type LoginError uint8

const (
	LoginErrorWrongPassword                       LoginError = 0x01
	LoginErrorAccountInvalid                      LoginError = 0x02
	LoginErrorAuthorizationServerUnavailable      LoginError = 0x03
	LoginErrorClientVersionTooLow                 LoginError = 0x05
	LoginErrorCannotAccessUserAccount             LoginError = 0x06
	LoginErrorAccountNotFound                     LoginError = 0x07
	LoginErrorAccountNoAccessAnyGame              LoginError = 0x08
	LoginErrorAccountNoAccessThisGame             LoginError = 0x09
	LoginErrorAccountClosed                       LoginError = 0x0A
	LoginErrorAccountAlreadyLoggedIn              LoginError = 0x0B
	LoginErrorTooManyPlayersLoggedIn              LoginError = 0x0C
	LoginErrorGameCurrentlyClosed                 LoginError = 0x0D
	LoginErrorAccountAlreadyLoggedIntoOtherServer LoginError = 0x10
	LoginErrorAccountIsInLogoutProcedure          LoginError = 0x11
	LoginErrorExpansionPacketNotAllowed           LoginError = 0x12
	LoginErrorAccountIsBannedFromThisServerType   LoginError = 0x13
	LoginErrorCafeIsOutOfPlayingTime              LoginError = 0x14
	LoginErrorPersonalAccountIsOutOfTime          LoginError = 0x15
	LoginErrorCafesAccountIsSuspended             LoginError = 0x16
	LoginErrorNotAuthorizedToUseExpansionVersion  LoginError = 0x17
	LoginErrorServiceNotAvailable                 LoginError = 0xAA
)
