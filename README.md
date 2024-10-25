# ecresolve
Resolves AWS ECR images with prioritized tag-based lookup

## Installation

You can download the binary from the [releases page](https://github.com/ebi-yade/ecresolve/releases/) or install it using Go:

```shell
go install github.com/ebi-yade/ecresolve@latest
```

## Usage

Let's say you have a repository named `<your-repository>` with only two tags: `latest` and `stable`.

If you run a query with three tags like this:

```shell
ecresolve foo latest stable --repository-name <your-repository-name>
```

The `latest` image will be returned since it's the **first** matching tag **found** in the list of provided candidates:

```json
{
  "imageId": {
    "imageDigest": "sha256:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "imageTag": "latest"
  },
  "imageManifest": " <abbreviated escaped JSON string>",
  "imageManifestMediaType": "application/vnd.oci.image.index.v1+json",
  "registryId": "<your-aws-account-id>",
  "repositoryName": "<your-repository-name>"
}
```
Note: The response is compatible with the [Image object](https://docs.aws.amazon.com/AmazonECR/latest/APIReference/API_Image.html) from the AWS ECR API.

If none of the specified tags exist, the command will fail with the following error and exit with code 1:

```shell
ecresolve foo bar --repository-name <your-repository-name>
# 2024/10/24 18:15:02 ERROR error Resolve: no matching images found
```

## Contributing

Feel free to open an issue or a pull request if you have any suggestions or improvements in mind.

日本語でのIssueやPRも歓迎です!
