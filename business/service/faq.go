package service

import (
	"context"
	"encoding/json"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
	"github.com/angel-one/fd-core/constants"
)

type FAQService interface {
	GetFAQDetails(ctx context.Context, tag string) (model.FAQResponse, error)
}

type FAQServiceImpl struct {
	faqDAO dao.FAQDAO
}

func DefaultFAQService() FAQService {
	return &FAQServiceImpl{faqDAO: dao.DefaultFAQDAO()}
}

func (service *FAQServiceImpl) GetFAQDetails(ctx context.Context, tag string) (model.FAQResponse, error) {
	response := model.FAQResponse{}
	faqs := []model.FAQ{}
	faqDetails, err := service.faqDAO.FetchFAQDetails(ctx, tag)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(faqDetails, &faqs)
	if err != nil {
		return response, err
	}

	return model.FAQResponse{
		Status: constants.StatusSuccess,
		FAQs:   faqs,
	}, nil
}
