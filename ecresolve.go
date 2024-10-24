package ecresolve

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/pkg/errors"
)

type Input struct {
	Tags           []string `arg:"" help:"Image tags to search for, in order of preference (e.g. :latest :dev :main)" required:""`
	RepositoryName string   `help:"The repository that contains the images to describe." required:""`
	RegistryId     string   `help:"The AWS account ID associated with the registry containing the images." optional:""`
	Region         string   `help:"AWS region to use." optional:""`
}

var ErrNoMatchingImages = errors.New("no matching images found")

func Resolve(ctx context.Context, input Input) (*types.Image, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	if input.Region != "" {
		cfg.Region = input.Region
	}
	client := ecr.NewFromConfig(cfg)

	// Build the input for BatchGetImage API
	imageIds := make([]types.ImageIdentifier, 0, len(input.Tags))
	for _, rev := range input.Tags {
		imageIds = append(imageIds, types.ImageIdentifier{
			ImageTag: &rev,
		})
	}
	bgiInput := &ecr.BatchGetImageInput{
		RepositoryName: &input.RepositoryName,
		ImageIds:       imageIds,
	}
	if input.RegistryId != "" {
		bgiInput.RegistryId = &input.RegistryId
	}

	output, err := client.BatchGetImage(ctx, bgiInput)
	if err != nil {
		return nil, errors.Wrap(err, "error BatchGetImage")
	}
	foundImage := findFirstMatchingImage(output.Images, input.Tags)
	if foundImage == nil {
		return nil, ErrNoMatchingImages
	}

	return foundImage, nil
}

func findFirstMatchingImage(images []types.Image, tags []string) *types.Image {
	// Create a map of tag to index for quick lookup
	index := make(map[string]int)
	for i, tag := range tags {
		index[tag] = i
	}

	// Find the image with the earliest revision in our list
	var bestImage *types.Image
	bestIndex := len(tags)

	for i := range images {
		image := &images[i]
		if image.ImageId.ImageTag == nil {
			continue
		}

		if idx, exists := index[*image.ImageId.ImageTag]; exists {
			if idx < bestIndex {
				bestIndex = idx
				bestImage = image
			}
		}
	}

	return bestImage
}
