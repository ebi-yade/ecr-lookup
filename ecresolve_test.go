package ecresolve

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stretchr/testify/assert"
)

func TestFindFirstMatchingImage(t *testing.T) {
	digest1 := "sha256:111"
	digest2 := "sha256:222"
	digest3 := "sha256:333"
	tag1 := "foo"
	tag2 := "bar"
	tag3 := "baz"

	tests := []struct {
		name      string
		images    []types.Image
		revisions []string
		want      *types.Image
	}{
		{
			name: "returns first matching image in revision order",
			images: []types.Image{
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest3,
						ImageTag:    &tag3,
					},
				},
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest2,
						ImageTag:    &tag2,
					},
				},
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest1,
						ImageTag:    &tag1,
					},
				},
			},
			revisions: []string{"foo", "bar", "baz"},
			want: &types.Image{
				ImageId: &types.ImageIdentifier{
					ImageDigest: &digest1,
					ImageTag:    &tag1,
				},
			},
		},
		{
			name: "returns second revision when first not found",
			images: []types.Image{
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest3,
						ImageTag:    &tag3,
					},
				},
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest2,
						ImageTag:    &tag2,
					},
				},
			},
			revisions: []string{"foo", "bar", "baz"},
			want: &types.Image{
				ImageId: &types.ImageIdentifier{
					ImageDigest: &digest2,
					ImageTag:    &tag2,
				},
			},
		},
		{
			name: "handles images without tags",
			images: []types.Image{
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest1,
					},
				},
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest2,
						ImageTag:    &tag2,
					},
				},
			},
			revisions: []string{"foo", "bar"},
			want: &types.Image{
				ImageId: &types.ImageIdentifier{
					ImageDigest: &digest2,
					ImageTag:    &tag2,
				},
			},
		},
		{
			name: "returns nil when no matches found",
			images: []types.Image{
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest1,
						ImageTag:    &tag3,
					},
				},
			},
			revisions: []string{"foo", "bar"},
			want:      nil,
		},
		{
			name:      "handles empty image list",
			images:    []types.Image{},
			revisions: []string{"foo", "bar"},
			want:      nil,
		},
		{
			name: "handles empty revisions list",
			images: []types.Image{
				{
					ImageId: &types.ImageIdentifier{
						ImageDigest: &digest1,
						ImageTag:    &tag1,
					},
				},
			},
			revisions: []string{},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findFirstMatchingImage(tt.images, tt.revisions)
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}
			if assert.NotNil(t, got) && assert.NotNil(t, got.ImageId) {
				assert.Equal(t, *tt.want.ImageId.ImageDigest, *got.ImageId.ImageDigest)
				assert.Equal(t, *tt.want.ImageId.ImageTag, *got.ImageId.ImageTag)
			}
		})
	}
}
