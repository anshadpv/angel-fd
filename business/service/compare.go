package service

import (
	"context"
	"fmt"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
	"golang.org/x/exp/slices"
)

type CompareService interface {
	GetCompareList(ctx context.Context) (model.FsiList, error)
	GetCompareFsiDetails(ctx context.Context, compareFsiKeys []string, compareFsiValues []string) (map[string]model.CompareFSIDetails, error)
}

type compareServiceImpl struct {
	compareDAO dao.CompareDAO
}

func DefaultCompareService() CompareService {
	return &compareServiceImpl{compareDAO: dao.DefaultCompareDAO()}
}

func (service *compareServiceImpl) GetCompareList(ctx context.Context) (model.FsiList, error) {
	response := model.FsiList{}
	fsiList, err := service.compareDAO.FetchCompareList(ctx)
	if err != nil {
		return response, err
	}
	response.FsiList = fsiList

	return response, nil
}

func (service *compareServiceImpl) GetCompareFsiDetails(ctx context.Context, compareFsiKeys []string, compareFsiValues []string) (map[string]model.CompareFSIDetails, error) {
	responseMap := map[string]model.CompareFSIDetails{}
	compareFsiDetails := model.CompareFSIDetails{}

	compareFsiDBDetails, err := service.compareDAO.FetchCompareFsiDetails(ctx, compareFsiValues)
	if err != nil {
		return responseMap, err
	}

	fsiDetailsMap := make(map[string][]model.CompareFSIDBDetails)
	for _, detail := range compareFsiDBDetails {
		fsiDetailsMap[detail.FSI] = append(fsiDetailsMap[detail.FSI], detail)
	}
	fmt.Println(fsiDetailsMap)

	for _, details := range fsiDetailsMap {
		yearlyInterestRate := map[int]float64{}
		for _, detail := range details {
			compareFsiDetails.FSI = detail.FSI
			compareFsiDetails.Name = detail.Name
			compareFsiDetails.MinDeposit = detail.MinDeposit
			compareFsiDetails.SeniorCitizenBenefit = detail.SeniorCitizenBenefit
			compareFsiDetails.BankAccount = detail.BankAccount
			compareFsiDetails.InsuredAmount = detail.InsuredAmount
			compareFsiDetails.ImageURL = detail.ImageURL
			if detail.TenureYears != 0 && detail.TenureMonths == 0 && detail.TenureDays == 0 {
				if _, exists := yearlyInterestRate[detail.TenureYears-1]; !exists {
					yearlyInterestRate[detail.TenureYears-1] = 0.0
				}
				if yearlyInterestRate[detail.TenureYears-1] < detail.InterestRate {
					yearlyInterestRate[detail.TenureYears-1] = detail.InterestRate
				}
			} else {
				if _, exists := yearlyInterestRate[detail.TenureYears]; !exists {
					yearlyInterestRate[detail.TenureYears] = 0.0
				}
				if yearlyInterestRate[detail.TenureYears] < detail.InterestRate {
					yearlyInterestRate[detail.TenureYears] = detail.InterestRate
				}
			}
		}
		compareFsiDetails.YearlyInterestRate.ZeroToOne = yearlyInterestRate[0]
		compareFsiDetails.YearlyInterestRate.OneToTwo = yearlyInterestRate[1]
		compareFsiDetails.YearlyInterestRate.TwoToThree = yearlyInterestRate[2]
		compareFsiDetails.YearlyInterestRate.ThreeToFour = yearlyInterestRate[3]
		compareFsiDetails.YearlyInterestRate.FourToFive = yearlyInterestRate[4]
		compareFsiDetails.YearlyInterestRate.FiveToSix = yearlyInterestRate[5]

		index := slices.Index(compareFsiValues, compareFsiDetails.FSI)
		key := compareFsiKeys[index]
		responseMap[key] = compareFsiDetails
	}

	return responseMap, nil
}
