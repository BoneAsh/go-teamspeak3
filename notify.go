package teamspeak3

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type NotifyEventType string

const (
	ClientEnterView NotifyEventType = "notifycliententerview"
	ClientLeftView                  = "notifyclientleftview"
)

var notifyEventTypeMap = map[NotifyEventType]Notify{
	ClientEnterView: &NotifyClientEnterView{},
	ClientLeftView:  &NotifyClientLeftView{},
}

func newNotifyEvent(s string) (n Notify, err error) {
	if event, ok := notifyEventTypeMap[NotifyEventType(s)]; ok {
		return event, nil
	} else {
		return nil, errors.New(fmt.Sprintf("event type(%s) is not support", s))
	}
}

type Notify interface {
	Type() NotifyEventType
	Decode(content map[string]interface{}) error
}

func NewNotify(content string) (n Notify, err error) {
	contentSplits := strings.SplitN(content, " ", 2)
	n, err = newNotifyEvent(contentSplits[0])
	if err != nil {
		return nil, err
	}
	err = n.Decode(DecodeResponse(contentSplits[1]))
	if err != nil {
		return nil, err
	}
	return n, nil
}

type NotifyClientEnterView struct {
	// Cfid source channel; "0" when entering the server
	Cfid int `mapstructure:"cfid"`
	// Ctid target channel
	Ctid                                 int         `mapstructure:"ctid"`
	ReasonId                             IReasonId   `mapstructure:"reasonid"`
	Clid                                 int         `mapstructure:"clid"`
	ClientUniqueIdentifier               string      `mapstructure:"client_unique_identifier"`
	ClientNickname                       string      `mapstructure:"client_nickname"`
	ClientInputMuted                     int         `mapstructure:"client_input_muted"`
	ClientOuputMuted                     int         `mapstructure:"client_ouput_muted"`
	ClientOutputOnlyMuted                int         `mapstructure:"client_outputonly_muted"`
	ClientInputHardware                  int         `mapstructure:"client_input_hardware"`
	ClientOutputHardWare                 int         `mapstructure:"client_output_hardware"`
	ClientMetaData                       string      `mapstructure:"client_meta_data"`
	ClientIsRecording                    int         `mapstructure:"client_is_recording"`
	ClientDatabaseId                     int         `mapstructure:"client_database_id"`
	clientChannelGroupId                 int         `mapstructure:"client_channel_group_id"`
	ClientServerGroups                   int         `mapstructure:"client_server_groups"`
	ClientAway                           int         `mapstructure:"client_away"`
	ClientAwayMessage                    string      `mapstructure:"client_away_message"`
	ClientType                           IClientType `mapstructure:"client_type"`
	ClientFlagAvatar                     string      `mapstructure:"client_flag_avatar"`
	ClientTalkPower                      int         `mapstructure:"client_talk_power"`
	ClientTalkRequest                    int         `mapstructure:"client_talk_request"`
	ClientTalkRequestMsg                 string      `mapstructure:"client_talk_request_msg"`
	ClientDescription                    string      `mapstructure:"client_description"`
	ClientIsTalker                       int         `mapstructure:"client_is_talker"`
	ClientIsPrioritySpeaker              int         `mapstructure:"client_is_priority_speaker"`
	ClientUnreadMessages                 int         `mapstructure:"client_unread_messages"`
	ClientNicknamePhonetic               string      `mapstructure:"client_nickname_phonetic"`
	ClientNeededServerQueryViewPower     int         `mapstructure:"client_needed_serverquery_view_power"`
	ClientIconId                         int         `mapstructure:"client_icon_id"`
	ClientIsChannelCommander             int         `mapstructure:"client_is_channel_commander"`
	ClientCountry                        string      `mapstructure:"client_country"`
	ClientChannelGroupInheritedChannelId int         `mapstructure:"client_channel_group_inherited_channel_id"`
	ClientBadges                         string      `mapstructure:"client_badges"`
	ClientMyteamspeakId                  string      `mapstructure:"client_myteamspeak_id"`
	ClientIntegrations                   string      `mapstructure:"client_integrations"`
	ClientMyteamspeakAvatar              string      `mapstructure:"client_myteamspeak_avatar"`
	ClientSignedBadges                   string      `mapstructure:"client_signed_badges"`
}

func (n *NotifyClientEnterView) Type() NotifyEventType {
	return ClientEnterView
}

func (n *NotifyClientEnterView) Decode(content map[string]interface{}) (err error) {
	return mapstructure.Decode(content, &n)
}

type NotifyClientLeftView struct {
	Cfid        int       `mapstructure:"cfid"`
	Ctid        int       `mapstructure:"ctid"`
	ReasonId    IReasonId `mapstructure:"reasonid"`
	ReasonMsg   string    `mapstructure:"reasonmsg"`
	Clid        int       `mapstructure:"clid"`
	InvokerId   string    `mapstructure:"invokerid"`
	InvokerName string    `mapstructure:"invokername"`
	InvokerUid  int       `mapstructure:"invokeruid"`
	Bantime     int       `mapstructure:"bantime"`
}

func (n *NotifyClientLeftView) Type() NotifyEventType {
	return ClientLeftView
}

func (n *NotifyClientLeftView) Decode(content map[string]interface{}) (err error) {
	return mapstructure.Decode(content, &n)
}
