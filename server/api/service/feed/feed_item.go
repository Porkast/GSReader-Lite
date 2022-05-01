package feed

import (
	"context"
	"gs-reader-lite/server/component"
	"gs-reader-lite/server/model"
	"gs-reader-lite/server/model/biz"

	"github.com/gogf/gf/v2/os/gtime"
)

func GetFeedItemByChannelId(ctx context.Context, start, size int, channelId, userId string) (itemList []biz.RssFeedItemData) {

	if userId != "" {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql+", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked, usc.status as sub").
			Joins("left join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Joins("left join user_sub_channel usc on usc.channel_id=rfi.channel_id and usc.user_id="+"'"+userId+"'").
			Joins("left join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id and umfi.user_id="+"'"+userId+"'").
			Where("rfi.channel_id", channelId).
			Order("rfi.input_date desc").
			Limit(size).
			Offset(start).
			Scan(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	} else {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql+", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl").
			Joins("left join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Where("rfi.channel_id", channelId).
			Order("rfi.input_date desc").
			Limit(size).
			Offset(start).
			Find(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	}

	return
}

func GetFeedItemListByUserId(ctx context.Context, userId string, start, size int) (itemList []biz.RssFeedItemData) {
	if err := component.GetDatabase().Table("rss_feed_item rfi").
		Select(model.RFIWithoutContentFieldSql+", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked, usc.status as sub").
		Joins("inner join user_sub_channel usc on usc.channel_id=rfi.channel_id").
		Joins("left join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id").
		Joins("inner join rss_feed_channel rfc on usc.channel_id=rfc.id").
		Where("usc.user_id = ? and usc.status = 1", userId).
		Order("rfi.input_date desc").
		Limit(size).
		Offset(start).
		Find(&itemList); err != nil {
		component.Logger().Error(ctx, err)
	}

	return
}

func GetLatestFeedItem(ctx context.Context, userId string, start, size int) (itemList []biz.RssFeedItemData) {
	if userId != "" {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql + ", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked").
			Joins("inner join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Joins("left join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id and umfi.user_id=" + "'" + userId + "'").
			Group("rfc.id").
			Order("rfi.input_date desc").
			Limit(size).
			Offset(start).
			Find(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	} else {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql + ", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl").
			Joins("inner join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Group("rfc.id").
			Order("rfi.input_date desc").
			Limit(size).
			Offset(start).
			Find(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	}

	return
}

func GetRandomFeedItem(ctx context.Context, start, size int, userId string) (itemList []biz.RssFeedItemData) {
	threeDayBefore := gtime.Now().AddDate(0, 0, -3).Format("Y-m-d")
	if userId != "" {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql+", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked, usc.status as sub").
			Joins("inner join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Joins("left join user_sub_channel usc on usc.channel_id=rfi.channel_id and usc.user_id="+"'"+userId+"'").
			Joins("left join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id and umfi.user_id="+"'"+userId+"'").
			Order("RAND()").
			Where("rfi.input_date>=?", threeDayBefore).
			Limit(size).
			Offset(start).
			Find(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	} else {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Select(model.RFIWithoutContentFieldSql+", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl").
			Joins("inner join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Order("RAND()").
			Where("rfi.input_date>=?", threeDayBefore).
			Limit(size).
			Offset(start).
			Find(&itemList); err != nil {
			component.Logger().Error(ctx, err)
		}
	}

	return
}

func MarkFeedItem(ctx context.Context, userId, feedItemId string) (err error) {
	markItem := model.UserMarkFeedItem{}
	result := component.GetDatabase().Table("user_mark_feed_item").Where("user_id = ? and channel_item_id = ?", userId, feedItemId).Find(&markItem)

	if err == nil && markItem.ChannelItemId != "" {
		if markItem.Status == 1 {
			markItem.Status = 0
		} else {
			markItem.Status = 1
		}
		result := component.GetDatabase().Table("user_mark_feed_item").Updates(model.UserMarkFeedItem{Status: markItem.Status}).Where("user_id = ? and channel_item_id = ?", userId, feedItemId)
		return result.Error
	}

	nowTime := gtime.Now().Format("Y-m-d H:i:s.u")
	markItem = model.UserMarkFeedItem{
		UserId:        userId,
		ChannelItemId: feedItemId,
		InputTime:     gtime.NewFromStr(nowTime),
		Status:        1,
	}

	result = component.GetDatabase().Create(&markItem)

	return result.Error
}

func GetMarkedFeedItemListByUserId(ctx context.Context, userId string, start, size int) (itemList []biz.RssFeedItemData) {
	if err := component.GetDatabase().Table("rss_feed_item rfi").
		Select(model.RFIWithoutContentFieldSql + ", rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked, usc.status as sub").
		Joins("inner join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id and umfi.user_id=" + "'" + userId + "'").
		Joins("left join user_sub_channel usc on usc.channel_id=rfi.channel_id and usc.user_id=" + "'" + userId + "'").
		Joins("inner join rss_feed_channel rfc on rfi.channel_id=rfc.id").
		Order("umfi.input_time desc").
		Where("umfi.status = 1").
		Limit(size).
		Offset(start).
		Find(&itemList); err != nil {
		component.Logger().Error(ctx, err)
	}

	return
}

func GetFeedItemByItemId(ctx context.Context, itemId, userId string) (item biz.RssFeedItemData) {

	if userId != "" {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Joins("left join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Joins("left join user_sub_channel usc on usc.channel_id=rfi.channel_id and usc.user_id="+"'"+userId+"'").
			Joins("left join user_mark_feed_item umfi on umfi.channel_item_id=rfi.id and umfi.user_id="+"'"+userId+"'").
			Select("rfi.*, rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl, umfi.status as marked, usc.status as sub").
			Where("rfi.id", itemId).
			Find(&item); err != nil {
			component.Logger().Error(ctx, err)
		}
	} else {
		if err := component.GetDatabase().Table("rss_feed_item rfi").
			Joins("left join rss_feed_channel rfc on rfi.channel_id=rfc.id").
			Select("rfi.*, rfc.rsshub_link as rsshubLink, rfc.title as channelTitle, rfc.image_url as channelImageUrl").
			Where("rfi.id", itemId).
			Find(&item); err != nil {
			component.Logger().Error(ctx, err)
		}
	}

	return
}