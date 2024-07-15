package models

type SubnetCommonOutput struct {
	Value string `json:"value"`
}

type SubnetApplyOutput struct {
	SubnetId   string `json:"subnet_id"`
	SubnetArn  string `json:"subnet_arn"`
	SubnetCIDR string `json:"subnet_cidr"`
}
