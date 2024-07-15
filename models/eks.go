package models

type EKSCommonOutput struct {
	Value string `json:"value"`
}

type EKSApplyOutput struct {
	EKSId       string `json:"eks_id"`
	EKSArn      string `json:"eks_arn"`
	EKSEndpoint string `json:"eks_endpoint"`
}
