package request

import (
	"akcom/internal/app/service/util"
	"fmt"
)

type ResetPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Pagination struct {
	Limit      *int   `json:"limit,omitempty" form:"limit"`
	Page       *int   `json:"page,omitempty" form:"page"`
	Offset     int    `json:"offset,omitempty" form:"offset"`
	Sort       string `json:"sort,omitempty" form:"sort"`
	Order      string `json:"order,omitempty" form:"order"`
	Query      string `json:"query,omitempty" form:"query"`
	GetAllData bool   `json:"get_all_data,omitempty" form:"get_all_data"`
	Total      int    `json:"total" form:"total"`
	TotalPage  int    `json:"total_page" form:"total_page"`
}

func (r *Pagination) Validate() error {
	if !r.GetAllData {
		if r.Limit == nil || r.Page == nil {
			r.Limit = util.Int(10)
			r.Page = util.Int(1)
		}
		if r.Limit != nil {
			if *r.Limit == 0 {
				r.Limit = util.Int(10)
			}
		}
		if r.Page != nil {
			if *r.Page == 0 {
				r.Page = util.Int(1)
			}
		}
		r.Offset = *r.Limit * (*r.Page - 1)
	}

	if r.Sort == "" {
		r.Sort = "DESC"
	}
	if r.Order == "" {
		r.Order = "created_at"
	}
	return nil
}

func (r ResetPassword) ValidateRequest() error {
	if r.Email == "" || r.Password == "" {
		return fmt.Errorf("invalid request body, require email/otp/password")
	}
	return nil
}
