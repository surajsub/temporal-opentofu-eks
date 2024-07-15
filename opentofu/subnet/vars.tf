
variable "vpc_id" {
  description = "The VPC to create this subnet in"
}


variable "private_subnet_cdir" {
  description = "Private Subnet CDIR Block"
  default = "10.0.0.0/20"
}

variable "private_subnet_zone" {
  description = "Private Subnet Zone"
  default ="us-west-2a"
}

variable "public_subnet_cdir" {
  description = "Public Subnet CDIR Block"
  default = "10.0.32.0/20"
}

variable "public_subnet_zone" {
  description = "Public Subnet Zone"
  default = "us-west-2b"
}