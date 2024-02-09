# cloudwatch-dashboard-mover

## これはなに？

あるAWSアカウントにある CloudWatch Dashboard (以下 Dashboard)を、
別のAWSアカウントに移行するためのサポートツールです。
Dashboard の定義に含まれるリソースIDを置換します。

前提条件として以下の2つを満たしている必要があります。

- 両アカウントのAWSリソースがterraformで管理されている
- 両アカウント用のterraform上のリソース名が一致している
- 両アカウントのAWSリソース構成が同一である

## 使い方

```console
$ AWS_PROFILE=account-from aws cloudwatch get-dashboard \
    --dashboard-name "dashboard" \
    | jq -r '.DashboardBody' \
    | jq -S '.' \
    > definition.json

$ go run cmd/cloudwatch-dashboard-mover/main.go \
    -dashboard-definition definition.json \
    -tfstate-from s3://account-from-terraform/terraform.tfstate \
    -tfstate-to s3://account-from-terraform/terraform.tfstate \
    -aws-profile-from account-from \
    -aws-profile-to account-to \
    > definition-replaced.json

$ AWS_PROFILE=account-to aws cloudwatch put-dashboard \
    --dashboard-name "dashboard" \
    --dashboard-body file://definition-replaced.json
```

## 注意

- Dashboard 定義に含まれるidのようなものを置換しているだけです。
idとリソースの対応を正確に紐づけているわけではありません。
異なる種類のリソース間でidが重複している場合などは、正しく動作しません。
- idのようなもの以外は置換しません。
識別子として任意の文字列を指定するようなリソース(`aws_db_cluster`)などは対象外です。
自分のユースケースでは命名ルールが決まっており、単純な置換で対応できたためです。

## 免責

個人的な必要性にかられて作成した限定的な用途のツールであり、すべての環境で正常に動作するとは限りません。
このツールを使用したことによるいかなる損害に対しても、保証いたしかねます。

## 作者

[@handlename](https://github.com/handlename)
