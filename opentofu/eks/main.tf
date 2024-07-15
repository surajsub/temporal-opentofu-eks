# Resource: aws_iam_role
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role
resource "aws_iam_role" "temporal" {
  name = "eks-cluster-temporal"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

# Resource: aws_iam_role_policy_attachment
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment
resource "aws_iam_role_policy_attachment" "demo-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.temporal.name
}

resource "aws_eks_cluster" "temporal" {
  name     = "temporal"
  role_arn = aws_iam_role.temporal.arn

  vpc_config {
    subnet_ids = [

      var.private_subnet_id,
      var.public_subnet_id

    ]
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
  depends_on = [aws_iam_role_policy_attachment.demo-AmazonEKSClusterPolicy]
}


output "eks_id" {
  value = aws_eks_cluster.temporal.id
}

output "eks_arn" {
  value  = aws_eks_cluster.temporal.arn
}

output "eks_endpoint" {
  value = aws_eks_cluster.temporal.endpoint
}