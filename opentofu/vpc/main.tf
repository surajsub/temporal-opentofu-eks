

resource "aws_vpc" "temporal" {

 cidr_block = var.cidr_block

 tags = {

   Name = "Temporal-Automated-VPC"

 }

}



output "vpc_id" {
  value = aws_vpc.temporal.id
}

output "vpc_cidr_block" {
  value = aws_vpc.temporal.cidr_block
}