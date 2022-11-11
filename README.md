# 项目说明
根据官方的文档说明在我看文档的时候已经有2个版本的v1 v2的使用方法：https://aws.github.io/aws-sdk-go-v2/docs/migrating/
官方建议是使用v2版本的sdk方式，并且给出了v1版本下的时候方式怎么对应迁移到v2方式使用，我这里就直接用v2版本的方式。

# 目录说明
- config-file 配置文件方式配置s3的访问密钥
  > https://docs.aws.amazon.com/zh_cn/sdkref/latest/guide/file-format.html#file-format-creds
- static-credential 静态认证方配置s3的访问密钥
  > https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#static-credentials
- bucket-operate 对bucket的操作
  > https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/example_s3_CreateBucket_section.html
  > https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/example_s3_DeleteBucket_section.html
- upload-file 单个文件传
  > https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/example_s3_PutObject_section.html
- muti-upload-file 大文件分段上传
  > https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/mpuoverview.html
- gen-file-sign-url 生成带有效期的访问签名
 > https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/example_s3_Scenario_PresignedUrl_section.html
