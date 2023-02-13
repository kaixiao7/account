package member

import "kaixiao7/account/internal/account/service"

type MemberController struct {
	memberSrv service.MemberSrv
}

func NewMemberController() *MemberController {
	return &MemberController{
		memberSrv: service.NewMemberSrv(),
	}
}
