package service

import (
	"context"
	"fmt"

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
	//fsiDetails := model.FsiStruct{}
	fmt.Println("REACHED HERE !!!! 222222")

	fsiDetailsPlan, err := service.fsiDetailsDAO.FetchFsiDetailsList(ctx, fsiDetailsValues)
	if err != nil {
		return responseMap, err
	}

	fmt.Println(fsiDetailsPlan)
	// fsiDetailsMap := make(map[string][]model.FsiStruct)
	// for _, detail := range fsiDetailsPlan {
	// 	fsiDetailsMap[detail.Plans.Fsi] = append(fsiDetailsMap[detail.Plans.fsi], detail)
	// }

	// for _, details := range fsiDetailsMap {
	// 	for _, detail := range details {
	// 		compareFsiDetails.FSI = detail.FSI
	// 		compareFsiDetails.Name = detail.Name
	// 		compareFsiDetails.MinDeposit = detail.MinDeposit
	// 		compareFsiDetails.SeniorCitizenBenefit = detail.SeniorCitizenBenefit
	// 		compareFsiDetails.BankAccount = detail.BankAccount
	// 		compareFsiDetails.InsuredAmount = detail.InsuredAmount
	// 		compareFsiDetails.ImageURL = detail.ImageURL
	// 	}

	// 	index := slices.Index(compareFsiValues, compareFsiDetails.FSI)
	// 	key := compareFsiKeys[index]
	// 	responseMap[key] = compareFsiDetails
	// }

	return responseMap, nil
}
