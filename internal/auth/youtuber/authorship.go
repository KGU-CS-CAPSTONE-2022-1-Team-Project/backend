package youtuber

import (
	"github.com/pkg/errors"
	"google.golang.org/api/youtube/v3"
)

const (
	MinSubscriber   = 1_000
	MinViewerCount  = 10_000
	MinCommentCount = 1_000
	MinVideoCount   = 20
)

func ValidateChannel(channel *youtube.Channel) error {
	auditDetails := channel.AuditDetails
	statistic := channel.Statistics
	if !(auditDetails.CopyrightStrikesGoodStanding &&
		auditDetails.ContentIdClaimsGoodStanding &&
		auditDetails.CommunityGuidelinesGoodStanding) {
		return errors.New("유튜브 정책 위반")
	}
	if channel.ContentOwnerDetails == nil {
		return errors.New("파트너가 아님")
	}
	if statistic == nil {
		return errors.New(" 권한 부재")
	} else if statistic.SubscriberCount < MinSubscriber {
		return errors.New("구독자수 부족")
	} else if statistic.ViewCount < MinViewerCount {
		return errors.New("조회수 부족")
	} else if statistic.CommentCount < MinCommentCount {
		return errors.New("댓글수 부족")
	} else if statistic.VideoCount < MinVideoCount {
		return errors.New("영상 개수 부족")
	}
	return nil
}
