package teamspeak3

type IReasonId int

const (
	ReasonIdIndependentlySwitchChannelsOrEnterServer IReasonId = iota
	ReasonIdUserOrChannelMoved                                 = 1
	ReasonIdTimeout                                            = 3
	ReasonIdChannelKick                                        = 4
	ReasonIdServerKick                                         = 5
	ReasonIdBan                                                = 6
	ReasonIdVoluntarilyLeaveServer                             = 8
	ReasonIdServerOrChannelEdited                              = 10
	ReasonIdServerShutdown                                     = 11
)

type IClientType int

const (
	ClientTypeVoice IClientType = iota
	ClientTypeQuery             = 1
)
