package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
)

type FsiDetailService interface {
	GetFsiDetails(ctx context.Context, fsiDetailsKeys []string, fsiDetailsValues []string) (map[string]model.FsiStruct, error)
}

type fsiDetailServiceImpl struct {
	fsiDetailsDAO dao.FsiDetailsDAO
}

func DefaultFsiDetailService() FsiDetailService {
	return &fsiDetailServiceImpl{fsiDetailsDAO: dao.DefaultFsiDetailsDAO()}
}

func (service *fsiDetailServiceImpl) GetFsiDetails(ctx context.Context, fsiDetailsKeys []string, fsiDetailsValues []string) (map[string]model.FsiStruct, error) {
	responseMap := map[string]model.FsiStruct{}
	fsiDetails := model.FsiStruct{}
	fmt.Println("REACHED HERE !!!! 222222")

	fsiDetailsPlan, err := service.fsiDetailsDAO.FetchFsiDetailsList(ctx, fsiDetailsValues)
	if err != nil {
		return responseMap, err
	}

	fmt.Println("WE HAVEE REEAACHHEDDD HHEERRREEEE")
	fsiDetailsMap := make(map[string][]model.FsiStruct)
	for _, details := range fsiDetailsPlan {
		for _, detail := range details.Plans {
			fsiDetailsMap[detail.Fsi] = append(fsiDetailsMap[detail.Fsi], details)
		}
	}

	for _, details := range fsiDetailsMap {
		for _, detail := range details {
			fsiDetails.Plans = detail.Plans
			fsiDetails.About = detail.About
			fsiDetails.Calculator = detail.Calculator
			fsiDetails.FAQs = detail.FAQs

			for _, fsi := range detail.Plans {
				index := slices.Index(fsiDetailsValues, fsi.Fsi)
				key := fsiDetailsKeys[index]
				responseMap[key] = fsiDetails
			}
		}
	}

	return responseMap, nil
}
