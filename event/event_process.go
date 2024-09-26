package event

import (
	"context"
	"github.com/dapr-platform/common"
	"log"
	"strings"
	"time"
)

func ConstructAndSendEvent(ctx context.Context, eventType int32, eventSubType int32, eventTitle string, eventDescription string, eventStatus int32, eventLevel int32, eventTime time.Time, objectID string, objectName string, location string) {
	dn := objectID + "_" + eventTitle
	dn = strings.Replace(dn, " ", "_", -1)
	event := common.Event{
		Dn:          dn,
		Title:       eventTitle,
		Type:        eventType,
		Description: eventDescription,
		Status:      eventStatus,
		Level:       eventLevel,
		EventTime:   common.LocalTime(eventTime),
		CreateAt:    common.LocalTime(time.Now()),
		UpdatedAt:   common.LocalTime(time.Now()),
		ObjectID:    objectID,
		ObjectName:  objectName,
		Location:    location,
	}
	err := publishInternalEvent(ctx, &event)
	if err != nil {
		log.Println("publishInternalEvent err:", err)
	}
}

func publishInternalEvent(ctx context.Context, event *common.Event) error {

	return common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, common.EVENT_TOPIC_NAME, event)
}

func PublishInternalMessage(ctx context.Context, msg *common.InternalMessage) error {

	return common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, common.INTERNAL_MESSAGE_TOPIC_NAME, msg)
}

func PublishCommonMessage(ctx context.Context, msg *common.CommonMessage) error {
	return common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, common.WEB_MESSAGE_TOPIC_NAME, msg)
}
