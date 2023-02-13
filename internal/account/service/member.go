package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
)

type MemberSrv interface {
	Push(ctx context.Context, members []*model.Member, syncTime int64) error

	Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Member, error)
}

type memberService struct {
	memberStore store.MemberStore
}

func NewMemberSrv() MemberSrv {
	return &memberService{
		memberStore: store.NewMemberStore(),
	}
}

func (m *memberService) Push(ctx context.Context, members []*model.Member, syncTime int64) error {
	return WithTransaction(ctx, func(ctx context.Context) error {
		for _, member := range members {
			member.SyncTime = syncTime
			if member.SyncState == constant.SYNC_ADD {
				member.SyncState = constant.SYNC_SUCCESS
				if e := m.memberStore.Add(ctx, member); e != nil {
					return e
				}
			} else {
				member.SyncState = constant.SYNC_SUCCESS
				if e := m.memberStore.Update(ctx, member); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (m *memberService) Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Member, error) {
	return m.memberStore.QueryBySyncTime(ctx, bookId, lastSyncTime)
}
